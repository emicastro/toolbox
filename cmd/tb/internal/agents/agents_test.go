package agents

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectIn(t *testing.T) {
	noneOnPath := func(string) bool { return false }

	t.Run("neither on PATH nor marker file present", func(t *testing.T) {
		home := t.TempDir()
		claude, grok := DetectIn(noneOnPath, home)
		if claude || grok {
			t.Errorf("DetectIn() = (%v, %v), want (false, false)", claude, grok)
		}
	})

	t.Run("both present via PATH", func(t *testing.T) {
		onPath := func(binary string) bool { return binary == "claude" || binary == "grok" }
		claude, grok := DetectIn(onPath, t.TempDir())
		if !claude || !grok {
			t.Errorf("DetectIn() = (%v, %v), want (true, true)", claude, grok)
		}
	})

	t.Run("both present via marker file, absent from PATH", func(t *testing.T) {
		home := t.TempDir()
		mustWriteFile(t, filepath.Join(home, ".claude", "settings.json"))
		mustWriteFile(t, filepath.Join(home, ".grok", "config.toml"))
		claude, grok := DetectIn(noneOnPath, home)
		if !claude || !grok {
			t.Errorf("DetectIn() = (%v, %v), want (true, true)", claude, grok)
		}
	})

	t.Run("only claude present", func(t *testing.T) {
		home := t.TempDir()
		mustWriteFile(t, filepath.Join(home, ".claude", "settings.json"))
		claude, grok := DetectIn(noneOnPath, home)
		if !claude || grok {
			t.Errorf("DetectIn() = (%v, %v), want (true, false)", claude, grok)
		}
	})

	t.Run("only grok present, via PATH", func(t *testing.T) {
		onPath := func(binary string) bool { return binary == "grok" }
		claude, grok := DetectIn(onPath, t.TempDir())
		if claude || !grok {
			t.Errorf("DetectIn() = (%v, %v), want (false, true)", claude, grok)
		}
	})

	t.Run("empty home never crashes and reports absent", func(t *testing.T) {
		claude, grok := DetectIn(noneOnPath, "")
		if claude || grok {
			t.Errorf("DetectIn() = (%v, %v), want (false, false)", claude, grok)
		}
	})

	// Regression test for docs/adr/0004-agent-detection-marker-file.md:
	// the bare config directory existing (with nothing inside it) must
	// NOT count as present. `tb install` creates exactly this directory
	// shape (an empty skills/ subdirectory, no settings.json/config.toml)
	// as a side effect of linking skills — if this test ever fails, that
	// self-inflicted false positive is back.
	t.Run("a bare config directory with no marker file inside does not count", func(t *testing.T) {
		home := t.TempDir()
		mustMkdir(t, filepath.Join(home, ".claude", "skills")) // what `tb install` alone creates
		mustMkdir(t, filepath.Join(home, ".grok", "skills"))
		claude, grok := DetectIn(noneOnPath, home)
		if claude || grok {
			t.Errorf("DetectIn() = (%v, %v), want (false, false): a bare config dir must not read as present", claude, grok)
		}
	})

	t.Run("a marker path that is a directory, not a file, still counts (Stat, not a file-type check)", func(t *testing.T) {
		// Documents the actual behavior: present() only calls os.Stat,
		// it does not require the marker path to be a regular file. This
		// is intentionally permissive rather than a claim that a
		// directory there is expected.
		home := t.TempDir()
		mustMkdir(t, filepath.Join(home, ".claude", "settings.json"))
		claude, _ := DetectIn(noneOnPath, home)
		if !claude {
			t.Error("DetectIn() claude = false, want true: present() only Stats the marker path")
		}
	})
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}

func mustWriteFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
}
