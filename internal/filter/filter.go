package filter

import (
	"strings"

	"github.com/Equationzhao/g/internal/entry"
	"github.com/Equationzhao/g/internal/git"
	"github.com/Equationzhao/g/internal/request"
	"github.com/gobwas/glob"
)

func Apply(ents []entry.Entry, req request.Request, repo git.RepoStatus) []entry.Entry {
	var pats []glob.Glob
	for _, s := range req.Ignore {
		if g, err := glob.Compile(s); err == nil {
			pats = append(pats, g)
		}
	}
	out := ents[:0:0]
	for _, e := range ents {
		if keep(e, req, pats, repo) {
			out = append(out, e)
		}
	}
	return out
}

func keep(e entry.Entry, req request.Request, pats []glob.Glob, repo git.RepoStatus) bool {
	if e.IsRootArg {
		return true
	}
	if e.Name == "." || e.Name == ".." {
		return req.Visibility == request.VisAll
	}
	if e.Hidden && req.Visibility == request.VisHidden {
		return false
	}
	if req.IgnoreBackups && strings.HasSuffix(e.Name, "~") {
		return false
	}
	for _, g := range pats {
		if g.Match(e.Name) {
			return false
		}
	}
	switch req.KindFilter {
	case request.KindDirsOnly:
		if !e.IsDir() {
			return false
		}
	case request.KindFilesOnly:
		if e.IsDir() {
			return false
		}
	}
	if req.GitIgnore && repo.OK {
		if git.Ignored(repo, relToRepo(e, repo)) {
			return false
		}
	}
	return true
}

func relToRepo(e entry.Entry, repo git.RepoStatus) string {
	p := e.Path
	root := repo.Root
	if strings.HasPrefix(p, root) {
		rel := strings.TrimPrefix(p, root)
		return strings.TrimPrefix(rel, "/")
	}
	return e.Name
}
