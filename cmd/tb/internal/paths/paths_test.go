package paths

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolve(t *testing.T) {
	t.Run("TOOLBOX_HOME set to an existing directory", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("TOOLBOX_HOME", dir)

		got, err := Resolve()
		if err != nil {
			t.Fatalf("Resolve() error = %v, want nil", err)
		}
		if got != filepath.Clean(dir) {
			t.Errorf("Resolve() = %q, want %q", got, dir)
		}
	})

	t.Run("TOOLBOX_HOME unset, falls back to $HOME/toolbox", func(t *testing.T) {
		home := t.TempDir()
		toolboxHome := filepath.Join(home, "toolbox")
		if err := os.Mkdir(toolboxHome, 0o755); err != nil {
			t.Fatal(err)
		}
		t.Setenv("TOOLBOX_HOME", "")
		t.Setenv("HOME", home)

		got, err := Resolve()
		if err != nil {
			t.Fatalf("Resolve() error = %v, want nil", err)
		}
		if got != toolboxHome {
			t.Errorf("Resolve() = %q, want %q", got, toolboxHome)
		}
	})

	t.Run("TOOLBOX_HOME set but empty, falls back like unset", func(t *testing.T) {
		home := t.TempDir()
		toolboxHome := filepath.Join(home, "toolbox")
		if err := os.Mkdir(toolboxHome, 0o755); err != nil {
			t.Fatal(err)
		}
		t.Setenv("TOOLBOX_HOME", "")
		t.Setenv("HOME", home)

		got, err := Resolve()
		if err != nil {
			t.Fatalf("Resolve() error = %v, want nil", err)
		}
		if got != toolboxHome {
			t.Errorf("Resolve() = %q, want %q", got, toolboxHome)
		}
	})

	t.Run("TOOLBOX_HOME set to a missing directory fails", func(t *testing.T) {
		missing := filepath.Join(t.TempDir(), "does-not-exist")
		t.Setenv("TOOLBOX_HOME", missing)

		_, err := Resolve()
		if err == nil {
			t.Fatal("Resolve() error = nil, want an error for a missing directory")
		}
		if !strings.Contains(err.Error(), missing) {
			t.Errorf("Resolve() error = %q, want it to mention %q", err.Error(), missing)
		}
		if !strings.Contains(err.Error(), "TOOLBOX_HOME") {
			t.Errorf("Resolve() error = %q, want it to mention TOOLBOX_HOME", err.Error())
		}
	})

	t.Run("TOOLBOX_HOME pointing at a file, not a directory, fails", func(t *testing.T) {
		file := filepath.Join(t.TempDir(), "not-a-dir")
		if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Setenv("TOOLBOX_HOME", file)

		_, err := Resolve()
		if err == nil {
			t.Fatal("Resolve() error = nil, want an error for a non-directory path")
		}
	})
}

func TestResolveWith(t *testing.T) {
	t.Run("env value wins and is cleaned", func(t *testing.T) {
		got, err := resolveWith("/some//path/../path", func() (string, error) { return "/should/not/be/used", nil })
		if err != nil {
			t.Fatalf("resolveWith() error = %v, want nil", err)
		}
		if want := filepath.Clean("/some//path/../path"); got != want {
			t.Errorf("resolveWith() = %q, want %q", got, want)
		}
	})

	t.Run("empty env value falls back to userHomeDir + toolbox", func(t *testing.T) {
		got, err := resolveWith("", func() (string, error) { return "/home/someone", nil })
		if err != nil {
			t.Fatalf("resolveWith() error = %v, want nil", err)
		}
		if want := filepath.Join("/home/someone", "toolbox"); got != want {
			t.Errorf("resolveWith() = %q, want %q", got, want)
		}
	})
}
