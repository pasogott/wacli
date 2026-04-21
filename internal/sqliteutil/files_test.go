package sqliteutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestChmodFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	for _, p := range []string{path, path + "-wal"} {
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatalf("WriteFile %s: %v", p, err)
		}
	}

	if err := ChmodFiles(path, 0o600); err != nil {
		t.Fatalf("ChmodFiles: %v", err)
	}

	for _, p := range []string{path, path + "-wal"} {
		info, err := os.Stat(p)
		if err != nil {
			t.Fatalf("Stat %s: %v", p, err)
		}
		if got := info.Mode().Perm(); got != 0o600 {
			t.Fatalf("%s mode = %04o, want 0600", filepath.Base(p), got)
		}
	}
}

func TestChmodFilesIgnoresMissingSidecars(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := ChmodFiles(path, 0o600); err != nil {
		t.Fatalf("ChmodFiles: %v", err)
	}
}
