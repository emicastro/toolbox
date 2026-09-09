package agents

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectIn(t *testing.T) {
	noneOnPath := func(string) bool { return false }

	t.Run("neither on PATH nor config dir present", func(t *testing.T) {
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

	t.Run("both present via config dir, absent from PATH", func(t *testing.T) {
		home := t.TempDir()
		mustMkdir(t, filepath.Join(home, ".claude"))
		mustMkdir(t, filepath.Join(home, ".grok"))
		claude, grok := DetectIn(noneOnPath, home)
		if !claude || !grok {
			t.Errorf("DetectIn() = (%v, %v), want (true, true)", claude, grok)
		}
	})

	t.Run("only claude present", func(t *testing.T) {
		home := t.TempDir()
		mustMkdir(t, filepath.Join(home, ".claude"))
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

	t.Run("a config dir that is a file, not a directory, does not count", func(t *testing.T) {
		home := t.TempDir()
		if err := os.WriteFile(filepath.Join(home, ".claude"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		claude, _ := DetectIn(noneOnPath, home)
		if claude {
			t.Error("DetectIn() claude = true, want false for a non-directory .claude")
		}
	})

	t.Run("empty home never crashes and reports absent", func(t *testing.T) {
		claude, grok := DetectIn(noneOnPath, "")
		if claude || grok {
			t.Errorf("DetectIn() = (%v, %v), want (false, false)", claude, grok)
		}
	})
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}
