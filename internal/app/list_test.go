package app

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	gfs "github.com/Equationzhao/g/internal/fs"
	"github.com/Equationzhao/g/internal/git"
	"github.com/Equationzhao/g/internal/parse"
)

func fixture() *gfs.Mem {
	m := gfs.NewMem()
	m.Now = time.Date(2026, 1, 2, 15, 4, 0, 0, time.UTC)
	m.AddDir("/fix")
	m.AddDir("/fix/a")
	m.AddFile("/fix/a/b", 10)
	m.AddFile("/fix/readme", 2048)
	m.AddFile("/fix/.hidden", 3)
	m.AddFile("/fix/zzz", 100)
	return m
}

func run(t *testing.T, mem *gfs.Mem, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	var out, errb bytes.Buffer
	d := Deps{
		FS:         mem,
		Git:        git.Fake{Err: errors.New("no git")},
		IDs:        NewIdentCache(),
		Stdout:     &out,
		Stderr:     &errb,
		Now:        func() time.Time { return mem.Now },
		IsTerminal: func() bool { return true },
		TermWidth:  func() int { return 80 },
	}
	code = Run(args, []string{"COLUMNS=80"}, d)
	return out.String(), errb.String(), code
}

func TestListDefaultLongTree(t *testing.T) {
	mem := fixture()
	out, errb, code := run(t, mem, "/fix")
	if code != 0 {
		t.Fatalf("code=%d err=%q out=%q", code, errb, out)
	}
	if !strings.Contains(out, "readme") || !strings.Contains(out, "a") {
		t.Fatalf("default listing missing names: %q", out)
	}
	if strings.Contains(out, "rewrite/v1 stub") {
		t.Fatal("stub leaked")
	}

	lout, _, code := run(t, mem, "-l", "/fix")
	if code != 0 {
		t.Fatal(code)
	}
	if !strings.Contains(lout, "readme") || !strings.Contains(lout, "2048") {
		t.Fatalf("long missing size/name: %q", lout)
	}
	if !strings.Contains(lout, "-rw") && !strings.Contains(lout, "Jan") && !strings.Contains(lout, "2026") {
		// mode string from FileMode plus a date
		if !strings.Contains(lout, "Jan") {
			t.Fatalf("long missing time-like field: %q", lout)
		}
	}

	tout, _, code := run(t, mem, "-T", "/fix")
	if code != 0 {
		t.Fatal(code)
	}
	if !strings.Contains(tout, "a") || !strings.Contains(tout, "b") {
		t.Fatalf("tree missing a/b: %q", tout)
	}
	if !strings.Contains(tout, "├") && !strings.Contains(tout, "|--") && !strings.Contains(tout, "└") {
		t.Fatalf("tree missing branches: %q", tout)
	}
}

func TestGNUShortsListing(t *testing.T) {
	mem := fixture()

	out, _, _ := run(t, mem, "-a", "/fix")
	if !strings.Contains(out, ".hidden") || !strings.Contains(out, ".") {
		t.Fatalf("-a should include hidden and dot: %q", out)
	}

	out, _, _ = run(t, mem, "-A", "/fix")
	if !strings.Contains(out, ".hidden") {
		t.Fatalf("-A missing hidden: %q", out)
	}
	// -A must not list . and .. as names on their own lines in a way that is only dots
	for _, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) == "." || strings.TrimSpace(line) == ".." {
			t.Fatalf("-A listed %q", line)
		}
	}

	out, _, _ = run(t, mem, "-1", "/fix")
	lines := nonempty(out)
	if len(lines) < 2 {
		t.Fatalf("-1 expected multiple lines: %q", out)
	}
	for _, ln := range lines {
		if strings.Contains(strings.TrimSpace(ln), "  ") && strings.Contains(ln, "readme") && strings.Contains(ln, "zzz") {
			t.Fatalf("-1 should be one name per line: %q", ln)
		}
	}

	out, _, _ = run(t, mem, "-d", "/fix")
	if !strings.Contains(out, "fix") {
		t.Fatalf("-d should list the directory itself: %q", out)
	}
	if strings.Contains(out, "readme") {
		t.Fatalf("-d should not list children: %q", out)
	}

	out, _, _ = run(t, mem, "-F", "--classify=always", "/fix")
	if !strings.Contains(out, "a/") {
		t.Fatalf("-F should classify dir: %q", out)
	}

	lout, _, _ := run(t, mem, "-lh", "/fix")
	if !strings.Contains(lout, "K") && !strings.Contains(lout, "2.0K") && !strings.Contains(lout, "2K") {
		t.Fatalf("-h should humanize 2048: %q", lout)
	}

	out, _, _ = run(t, mem, "-1", "-S", "/fix")
	names := namesOnly(out)
	if idx(names, "readme") > idx(names, "zzz") {
		// size 2048 > 100 so readme first
		t.Fatalf("-S largest first, got %v", names)
	}

	out, _, _ = run(t, mem, "-1", "-S", "-r", "/fix")
	names = namesOnly(out)
	if idx(names, "readme") < idx(names, "zzz") && len(names) > 1 {
		t.Fatalf("-Sr should reverse size, got %v", names)
	}

	out, _, _ = run(t, mem, "-R", "/fix")
	if !strings.Contains(out, "b") {
		t.Fatalf("-R should recurse to a/b: %q", out)
	}

	out, _, _ = run(t, mem, "-1", "-t", "/fix")
	if !strings.Contains(out, "readme") {
		t.Fatalf("-t listing failed: %q", out)
	}
}

func TestGitDegrade(t *testing.T) {
	mem := fixture()
	out, errb, code := run(t, mem, "--git", "-1", "/fix")
	if code != 0 {
		t.Fatalf("git missing must not fail listing: code=%d err=%q", code, errb)
	}
	if !strings.Contains(out, "readme") {
		t.Fatalf("listing missing: %q", out)
	}
	if !strings.Contains(out, "--") {
		t.Fatalf("expected degraded git cells: %q", out)
	}
}

func TestSpecsBudget(t *testing.T) {
	if n := len(parse.Specs()); n != 40 {
		t.Fatalf("Specs()=%d", n)
	}
}

func nonempty(s string) []string {
	var o []string
	for _, ln := range strings.Split(strings.TrimRight(s, "\n"), "\n") {
		if strings.TrimSpace(ln) != "" {
			o = append(o, ln)
		}
	}
	return o
}

func namesOnly(s string) []string {
	var o []string
	for _, ln := range nonempty(s) {
		o = append(o, strings.TrimSpace(ln))
	}
	return o
}

func idx(a []string, name string) int {
	for i, s := range a {
		if s == name || strings.HasSuffix(s, name) || strings.Contains(s, name) {
			return i
		}
	}
	return 99
}
