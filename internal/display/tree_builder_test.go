package display

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/util"
)

func TestTreeBuilderBuildsTreeWithOrder(t *testing.T) {
	base := filepath.Join(string(filepath.Separator), "root")
	now := time.Now()

	root := newTreeFileInfo(t, filepath.Join(base), "", 0, true, now)
	dirA := newTreeFileInfo(t, filepath.Join(base, "a"), base, 1, true, now)
	dirB := newTreeFileInfo(t, filepath.Join(base, "b"), base, 1, true, now)
	fileA1 := newTreeFileInfo(t, filepath.Join(base, "a", "a1.txt"), dirA.FullPath, 2, false, now)
	fileA2 := newTreeFileInfo(t, filepath.Join(base, "a", "a2.txt"), dirA.FullPath, 2, false, now)
	fileB1 := newTreeFileInfo(t, filepath.Join(base, "b", "b1.txt"), dirB.FullPath, 2, false, now)

	items := []*item.FileInfo{root, dirA, dirB, fileA1, fileA2, fileB1}
	buildTree := NewTreeBuilder().Build(items)

	if buildTree.Root.Meta != root {
		t.Fatalf("expected root meta %s, got %v", root.FullPath, buildTree.Root.Meta)
	}
	if buildTree.Root.Level != 0 {
		t.Fatalf("expected root level 0, got %d", buildTree.Root.Level)
	}
	if len(buildTree.Root.Child) != 2 {
		t.Fatalf("expected 2 root children, got %d", len(buildTree.Root.Child))
	}
	if buildTree.Root.Child[0].Meta != dirA || buildTree.Root.Child[1].Meta != dirB {
		t.Fatalf("root children order mismatch")
	}
	if buildTree.Root.Child[0].Level != 1 || buildTree.Root.Child[1].Level != 1 {
		t.Fatalf("expected child level 1")
	}

	dirAChildren := buildTree.Root.Child[0].Child
	if len(dirAChildren) != 2 {
		t.Fatalf("expected 2 children under dirA, got %d", len(dirAChildren))
	}
	if dirAChildren[0].Meta != fileA1 || dirAChildren[1].Meta != fileA2 {
		t.Fatalf("dirA children order mismatch")
	}

	dirBChildren := buildTree.Root.Child[1].Child
	if len(dirBChildren) != 1 {
		t.Fatalf("expected 1 child under dirB, got %d", len(dirBChildren))
	}
	if dirBChildren[0].Meta != fileB1 {
		t.Fatalf("dirB child mismatch")
	}
	if dirBChildren[0].Level != 2 {
		t.Fatalf("expected dirB child level 2, got %d", dirBChildren[0].Level)
	}
}

func newTreeFileInfo(t *testing.T, fullPath, parent string, level int, isDir bool, modTime time.Time) *item.FileInfo {
	t.Helper()
	mode := os.FileMode(0)
	if isDir {
		mode = os.ModeDir
	}
	info, err := item.NewFileInfoWithOption(
		item.WithAbsPath(fullPath),
		item.WithFileInfo(util.NewMockFileInfo(0, isDir, filepath.Base(fullPath), mode, modTime)),
	)
	if err != nil {
		t.Fatalf("new file info: %v", err)
	}
	info.ParentPath = parent
	info.Level = level
	return info
}
