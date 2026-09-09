package config

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestParseValid(t *testing.T) {
	tests := []struct {
		name    string
		content string
		strs    map[string]string
		arrays  map[string][]string
		sects   map[string]map[string]string
	}{
		{
			name:    "blank lines and comments are ignored",
			content: "\n# a comment\n\nprofile = \"rust-systems\"\n# trailing comment\n",
			strs:    map[string]string{"profile": "rust-systems"},
		},
		{
			name:    "inline trailing comment after a string",
			content: `name = "spec" # this is the skill name`,
			strs:    map[string]string{"name": "spec"},
		},
		{
			name:    "single-line array",
			content: `skills = ["spec", "adr", "onboard"]`,
			arrays:  map[string][]string{"skills": {"spec", "adr", "onboard"}},
		},
		{
			name: "multi-line array with trailing comma",
			content: "skills = [\n" +
				"  \"spec\", \"adr\", \"onboard\", \"scout\", \"handoff\",\n" +
				"  \"rust-verify\", \"rust-systems\",\n" +
				"]\n",
			arrays: map[string][]string{
				"skills": {"spec", "adr", "onboard", "scout", "handoff", "rust-verify", "rust-systems"},
			},
		},
		{
			name:    "empty array",
			content: `templates = []`,
			arrays:  map[string][]string{"templates": nil},
		},
		{
			name:    "section with a string field",
			content: "[verify]\nsummary = \"cargo test\"\n",
			sects:   map[string]map[string]string{"verify": {"summary": "cargo test"}},
		},
		{
			name:    "keys may contain hyphens",
			content: `toolbox-version = "0.1.0"`,
			strs:    map[string]string{"toolbox-version": "0.1.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parse("test.toml", tt.content)
			if err != nil {
				t.Fatalf("parse() error = %v, want nil", err)
			}
			if tt.strs != nil && !reflect.DeepEqual(doc.strings, tt.strs) {
				t.Errorf("strings = %#v, want %#v", doc.strings, tt.strs)
			}
			if tt.arrays != nil {
				for k, want := range tt.arrays {
					if got := doc.arrays[k]; !reflect.DeepEqual(got, want) {
						t.Errorf("arrays[%q] = %#v, want %#v", k, got, want)
					}
				}
			}
			if tt.sects != nil && !reflect.DeepEqual(doc.sections, tt.sects) {
				t.Errorf("sections = %#v, want %#v", doc.sections, tt.sects)
			}
		})
	}
}

func TestParseUnsupportedSyntax(t *testing.T) {
	tests := []struct {
		name    string
		content string
		line    int
	}{
		{
			name:    "unquoted value",
			content: "profile = rust-systems",
			line:    1,
		},
		{
			name:    "number value",
			content: "count = 3",
			line:    1,
		},
		{
			name:    "inline table",
			content: `verify = { summary = "x" }`,
			line:    1,
		},
		{
			name:    "nested array",
			content: `skills = [["spec"], ["adr"]]`,
			line:    1,
		},
		{
			name:    "no '=' on a non-comment, non-section line",
			content: "not a valid line at all",
			line:    1,
		},
		{
			name:    "array inside a section",
			content: "[verify]\nsteps = [\"a\", \"b\"]\n",
			line:    2,
		},
		{
			name:    "unclosed array",
			content: "skills = [\n  \"spec\",\n",
			line:    1,
		},
		{
			name:    "error is reported on the second line of a file",
			content: "profile = \"rust-systems\"\nbad line\n",
			line:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parse("test.toml", tt.content)
			if err == nil {
				t.Fatal("parse() error = nil, want an unsupported-syntax error")
			}
			wantPrefix := "test.toml:" + strconv.Itoa(tt.line) + ":"
			if !strings.HasPrefix(err.Error(), wantPrefix) {
				t.Errorf("parse() error = %q, want prefix %q", err.Error(), wantPrefix)
			}
		})
	}
}
