package app

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Equationzhao/g/internal/collect"
	"github.com/Equationzhao/g/internal/entry"
	"github.com/Equationzhao/g/internal/filter"
	gfs "github.com/Equationzhao/g/internal/fs"
	"github.com/Equationzhao/g/internal/git"
	"github.com/Equationzhao/g/internal/parse"
	"github.com/Equationzhao/g/internal/print"
	"github.com/Equationzhao/g/internal/request"
	"github.com/Equationzhao/g/internal/sort"
	"golang.org/x/term"
)

type Deps struct {
	FS         gfs.Filesystem
	Git        git.Client
	IDs        *IdentCache
	Stdout     io.Writer
	Stderr     io.Writer
	Now        func() time.Time
	LookupEnv  func(string) (string, bool)
	IsTerminal func() bool
	TermWidth  func() int
	GitTimeout time.Duration
}

func OSDeps() Deps {
	return Deps{
		FS:        gfs.OS{},
		Git:       git.Exec{},
		IDs:       NewIdentCache(),
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		Now:       time.Now,
		LookupEnv: os.LookupEnv,
		IsTerminal: func() bool {
			return term.IsTerminal(int(os.Stdout.Fd()))
		},
		TermWidth: func() int {
			w, _, err := term.GetSize(int(os.Stdout.Fd()))
			if err != nil {
				return 0
			}
			return w
		},
	}
}

// Run is the single listing entry point used by main and tests.
func Run(args, env []string, d Deps) int {
	if d.FS == nil {
		d.FS = gfs.OS{}
	}
	if d.Git == nil {
		d.Git = git.Exec{}
	}
	if d.Stdout == nil {
		d.Stdout = os.Stdout
	}
	if d.Stderr == nil {
		d.Stderr = os.Stderr
	}
	if d.Now == nil {
		d.Now = time.Now
	}
	req, err := parse.Parse(args, env)
	if err != nil {
		fmt.Fprintf(d.Stderr, "g: %s\n", err.Error())
		if e, ok := err.(*parse.Error); ok && e.Code != 0 {
			return e.Code
		}
		return 2
	}
	if req.Help {
		fmt.Fprint(d.Stdout, parse.Usage())
		return 0
	}
	if req.Version {
		fmt.Fprintf(d.Stdout, "g %s\n", parse.Version)
		return 0
	}
	req, err = parse.ApplyConfigFile(req)
	if err != nil {
		fmt.Fprintf(d.Stderr, "g: %s\n", err.Error())
		if e, ok := err.(*parse.Error); ok && e.Code != 0 {
			return e.Code
		}
		return 2
	}
	rt := request.Runtime{
		StdoutTTY:  d.IsTerminal != nil && d.IsTerminal(),
		WidthIOCTL: 0,
		COLUMNS:    parse.EnvCOLUMNS(env),
		NOCOLOR:    parse.EnvNOCOLOR(env),
	}
	if d.TermWidth != nil {
		rt.WidthIOCTL = d.TermWidth()
	}
	req = request.Resolve(req, rt)

	roots := collect.Walk(d.FS, req)
	code := 0
	var jobRoots []print.Root
	for _, r := range roots {
		if r.Code > code {
			code = r.Code
		}
		if r.Err != nil {
			fmt.Fprintf(d.Stderr, "g: %s: %s\n", r.Path, r.Err.Error())
			if r.Code == 0 {
				code = 2
			}
			continue
		}
		var repo git.RepoStatus
		if req.Git || req.GitIgnore {
			repo = git.WithTimeout(d.Git, r.Abs, d.GitTimeout)
			if !repo.OK {
				fmt.Fprintf(d.Stderr, "g: %s: git unavailable\n", r.Path)
			}
		}
		ents := filter.Apply(r.Entries, req, repo)
		if req.Git {
			for i := range ents {
				rel := ents[i].Name
				if r.Abs != "" && len(ents[i].Path) >= len(r.Abs) {
					rel = ents[i].Path[len(r.Abs):]
					if len(rel) > 0 && (rel[0] == '/' || rel[0] == '\\') {
						rel = rel[1:]
					}
				}
				cell := git.Lookup(repo, rel, ents[i].IsDir())
				ents[i].Git = cell
			}
		}
		fillOwners(ents, d.IDs)
		sort.Apply(ents, req)
		jobRoots = append(jobRoots, print.Root{Path: r.Path, Abs: r.Abs, Entries: ents})
	}
	if len(jobRoots) == 0 && code != 0 {
		return code
	}
	err = print.Print(d.Stdout, print.Job{
		Roots: jobRoots,
		Req:   req,
		Width: req.Width,
		Now:   d.Now(),
	})
	if err != nil {
		fmt.Fprintf(d.Stderr, "g: %s\n", err.Error())
		return 2
	}
	return code
}

func fillOwners(ents []entry.Entry, ids *IdentCache) {
	if ids == nil {
		return
	}
	for i := range ents {
		if ents[i].UID != "" && ents[i].User == "" {
			ents[i].User = ids.User(ents[i].UID)
		}
		if ents[i].GID != "" && ents[i].Group == "" {
			ents[i].Group = ids.Group(ents[i].GID)
		}
	}
}
