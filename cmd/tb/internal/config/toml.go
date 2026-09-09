// Package config reads tb's own closed TOML subset (docs/design.md §3.1,
// docs/adr/0001-tb-dependency-policy.md): comments, blank lines, bare
// keys, double-quoted string values, flat double-quoted string arrays
// (which may span multiple lines up to the closing bracket — see the
// note on parseArray), and a single level of [section] tables holding
// their own string fields. Anything outside that grammar is a parse
// error naming the file and line, by design: tb only ever needs to read
// the files it or a human following this schema writes.
package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	beginArray = "["
	endArray   = "]"
)

var bareKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_-]*$`)

// document is the generic parse result of one file, before it is shaped
// into a typed Pointer or Profile.
type document struct {
	strings  map[string]string
	arrays   map[string][]string
	sections map[string]map[string]string
}

// parseFile reads and parses path.
func parseFile(path string) (*document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("tb: %w", err)
	}
	return parse(path, string(data))
}

func parse(path, content string) (*document, error) {
	doc := &document{
		strings:  map[string]string{},
		arrays:   map[string][]string{},
		sections: map[string]map[string]string{},
	}

	lines := strings.Split(content, "\n")
	section := "" // "" means top level, not inside a [section]

	for i := 0; i < len(lines); {
		lineNo := i + 1
		trimmed := strings.TrimSpace(stripComment(lines[i]))
		i++

		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "[") {
			name, err := parseSectionHeader(trimmed)
			if err != nil {
				return nil, fmt.Errorf("%s:%d: unsupported syntax: %s", path, lineNo, trimmed)
			}
			section = name
			if _, ok := doc.sections[section]; !ok {
				doc.sections[section] = map[string]string{}
			}
			continue
		}

		key, rawValue, err := splitKeyValue(trimmed)
		if err != nil {
			return nil, fmt.Errorf("%s:%d: unsupported syntax: %s", path, lineNo, trimmed)
		}
		rawValue = strings.TrimSpace(rawValue)

		switch {
		case strings.HasPrefix(rawValue, `"`):
			val, err := parseQuotedString(rawValue)
			if err != nil {
				return nil, fmt.Errorf("%s:%d: unsupported syntax: %s", path, lineNo, trimmed)
			}
			target := doc.strings
			if section != "" {
				target = doc.sections[section]
			}
			target[key] = val

		case strings.HasPrefix(rawValue, beginArray):
			if section != "" {
				return nil, fmt.Errorf("%s:%d: unsupported syntax: array inside [%s] (arrays are top-level only)", path, lineNo, section)
			}
			values, next, err := parseArray(path, lines, i-1)
			if err != nil {
				return nil, err
			}
			doc.arrays[key] = values
			i = next

		default:
			return nil, fmt.Errorf("%s:%d: unsupported syntax: %s", path, lineNo, trimmed)
		}
	}

	return doc, nil
}

// stripComment cuts a line at the first '#' that is not inside a
// double-quoted string.
func stripComment(line string) string {
	inString := false
	for idx, r := range line {
		switch r {
		case '"':
			inString = !inString
		case '#':
			if !inString {
				return line[:idx]
			}
		}
	}
	return line
}

func parseSectionHeader(trimmed string) (string, error) {
	if !strings.HasSuffix(trimmed, "]") {
		return "", fmt.Errorf("malformed section header")
	}
	name := strings.TrimSpace(trimmed[1 : len(trimmed)-1])
	if name == "" || strings.ContainsAny(name, "[]") {
		return "", fmt.Errorf("malformed section header")
	}
	return name, nil
}

func splitKeyValue(trimmed string) (key, rawValue string, err error) {
	idx := strings.Index(trimmed, "=")
	if idx < 0 {
		return "", "", fmt.Errorf("no '=' in line")
	}
	key = strings.TrimSpace(trimmed[:idx])
	if !bareKeyPattern.MatchString(key) {
		return "", "", fmt.Errorf("invalid key %q", key)
	}
	return key, trimmed[idx+1:], nil
}

// parseQuotedString requires rest to be exactly one double-quoted string
// with nothing else on the line (comments are already stripped). This
// subset has no escape sequences: a value containing '"' is unsupported.
func parseQuotedString(rest string) (string, error) {
	if len(rest) < 2 || rest[0] != '"' || rest[len(rest)-1] != '"' {
		return "", fmt.Errorf("not a quoted string")
	}
	inner := rest[1 : len(rest)-1]
	if strings.Contains(inner, `"`) {
		return "", fmt.Errorf("unsupported escaping")
	}
	return inner, nil
}

// parseArray parses a flat double-quoted string array whose opening line
// is lines[startIdx]. Real TOML treats newlines inside "[ ... ]" as
// insignificant, and profiles/*.toml (docs/design.md §3.3) relies on
// that to keep long skill lists readable, so this reader follows suit:
// it keeps consuming lines (each with its own comment stripped) until one
// contains the closing bracket. Nested arrays and inline tables are not
// supported — the first "]" found is always taken as the close.
//
// It returns the parsed values and the index of the line after the
// closing bracket, so the caller can resume line-by-line parsing there.
func parseArray(path string, lines []string, startIdx int) ([]string, int, error) {
	var buf strings.Builder
	i := startIdx
	closed := false
	for i < len(lines) {
		clean := stripComment(lines[i])
		buf.WriteString(clean)
		buf.WriteByte('\n')
		i++
		if strings.Contains(clean, endArray) {
			closed = true
			break
		}
	}
	if !closed {
		return nil, 0, fmt.Errorf("%s:%d: unsupported syntax: array not closed", path, startIdx+1)
	}

	content := buf.String()
	open := strings.Index(content, beginArray)
	closeIdx := strings.LastIndex(content, endArray)
	if open < 0 || closeIdx < open {
		return nil, 0, fmt.Errorf("%s:%d: unsupported syntax: malformed array", path, startIdx+1)
	}

	var values []string
	for _, part := range strings.Split(content[open+1:closeIdx], ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue // tolerate a trailing comma and blank lines
		}
		val, err := parseQuotedString(part)
		if err != nil {
			return nil, 0, fmt.Errorf("%s:%d: unsupported syntax: array element %q", path, startIdx+1, part)
		}
		values = append(values, val)
	}
	return values, i, nil
}
