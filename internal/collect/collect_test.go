package collect

import (
	"strings"
	"testing"
	"time"

	"github.com/Equationzhao/g/internal/entry"
	gfs "github.com/Equationzhao/g/internal/fs"
	"github.com/Equationzhao/g/internal/request"
)

func TestWalkRecurseEmitsPerDirectoryRoots(t *testing.T) {
	m := gfs.NewMem()
	m.Now = time.Date(2026, 1, 2, 15, 4, 0, 0, time.UTC)
	m.AddDir("/fix")
	m.AddDir("/fix/a")
	m.AddFile("/fix/a/b", 10)
	m.AddFile("/fix/readme", 2048)
	m.AddDir("/fix/.secret")
	m.AddFile("/fix/.secret/leaked", 1)

	req := request.Default()
	req.Paths = []string{"/fix"}
	req.Recurse = true

	roots := Walk(m, req)
	if len(roots) != 2 {
		t.Fatalf("want 2 roots (parent + a/), got %d: %v", len(roots), rootPaths(roots))
	}
	if roots[0].Path != "/fix" {
		t.Fatalf("first root %q", roots[0].Path)
	}
	if roots[1].Path != "/fix/a" {
		t.Fatalf("second root %q, want /fix/a", roots[1].Path)
	}
	if namesContain(roots[0].Entries, "b") {
		t.Fatalf("flattened: b is in the parent root: %v", entryNames(roots[0].Entries))
	}
	if !namesContain(roots[1].Entries, "b") {
		t.Fatalf("a/ missing b: %v", entryNames(roots[1].Entries))
	}
	for _, r := range roots {
		if strings.Contains(r.Path, ".secret") {
			t.Fatalf("recursed into hidden dir: %v", rootPaths(roots))
		}
		if namesContain(r.Entries, "leaked") {
			t.Fatalf("collected hidden-dir child: path=%s names=%v", r.Path, entryNames(r.Entries))
		}
	}

	req.Visibility = request.VisAlmostAll
	shown := Walk(m, req)
	var sawLeaked bool
	for _, r := range shown {
		if namesContain(r.Entries, "leaked") {
			sawLeaked = true
		}
	}
	if !sawLeaked {
		t.Fatalf("-A -R should enter hidden dirs, roots=%v", rootPaths(shown))
	}
}

func rootPaths(roots []Root) []string {
	s := make([]string, len(roots))
	for i, r := range roots {
		s[i] = r.Path
	}
	return s
}

func entryNames(ents []entry.Entry) []string {
	s := make([]string, len(ents))
	for i, e := range ents {
		s[i] = e.Name
	}
	return s
}

func namesContain(ents []entry.Entry, name string) bool {
	for _, e := range ents {
		if e.Name == name {
			return true
		}
	}
	return false
}
