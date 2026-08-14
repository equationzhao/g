package print

import (
	"strings"
	"testing"
	"time"

	"github.com/Equationzhao/g/internal/entry"
	"github.com/Equationzhao/g/internal/request"
)

func TestLongGitNotInDisplayName(t *testing.T) {
	e := entry.Entry{Name: "readme", Git: "M-", Kind: entry.KindFile, Mode: 0o644, Nlink: 1, Size: 7, ModTime: time.Date(2026, 1, 2, 15, 4, 0, 0, time.UTC)}
	req := request.Default()
	req.Long = true
	req.Git = true
	name := displayName(e, req)
	if strings.Contains(name, "M-") {
		t.Fatalf("displayName must not prefix git when Long: %q", name)
	}
	line := longLine(e, req, e.ModTime)
	if strings.Count(line, "M-") != 1 {
		t.Fatalf("want git cell once, got %q", line)
	}
	if !strings.HasSuffix(line, " readme") && !strings.HasSuffix(line, "readme") {
		t.Fatalf("name should be last field: %q", line)
	}
}
