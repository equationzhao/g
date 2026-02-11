package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Equationzhao/g/internal/filter"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/util"
)

func TestDiveSetsParentAndLevel(t *testing.T) {
	root := t.TempDir()
	dir1 := filepath.Join(root, "dir1")
	dir2 := filepath.Join(dir1, "dir2")

	if err := os.Mkdir(dir1, 0o755); err != nil {
		t.Fatalf("mkdir dir1: %v", err)
	}
	if err := os.Mkdir(dir2, 0o755); err != nil {
		t.Fatalf("mkdir dir2: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("a"), 0o644); err != nil {
		t.Fatalf("write a.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir1, "b.txt"), []byte("b"), 0o644); err != nil {
		t.Fatalf("write b.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir2, "c.txt"), []byte("c"), 0o644); err != nil {
		t.Fatalf("write c.txt: %v", err)
	}

	infoSlice := util.NewSlice[*item.FileInfo](10)
	errSlice := util.NewSlice[error](10)
	dive(root, 1, -1, infoSlice, errSlice, filter.NewItemFilter())

	for _, err := range *errSlice.GetRaw() {
		if err != nil {
			t.Fatalf("dive error: %v", err)
		}
	}

	expected := map[string]struct {
		parent string
		level  int
	}{
		filepath.Join(root, "a.txt"): {parent: root, level: 1},
		dir1:                         {parent: root, level: 1},
		filepath.Join(dir1, "b.txt"): {parent: dir1, level: 2},
		dir2:                         {parent: dir1, level: 2},
		filepath.Join(dir2, "c.txt"): {parent: dir2, level: 3},
	}

	got := *infoSlice.GetRaw()
	if len(got) != len(expected) {
		t.Fatalf("expected %d entries, got %d", len(expected), len(got))
	}

	for _, info := range got {
		exp, ok := expected[info.FullPath]
		if !ok {
			t.Fatalf("unexpected entry: %s", info.FullPath)
		}
		if info.ParentPath != exp.parent {
			t.Fatalf("parent for %s: expected %s, got %s", info.FullPath, exp.parent, info.ParentPath)
		}
		if info.Level != exp.level {
			t.Fatalf("level for %s: expected %d, got %d", info.FullPath, exp.level, info.Level)
		}
	}
}
