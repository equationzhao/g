package collect

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Equationzhao/g/internal/entry"
	gfs "github.com/Equationzhao/g/internal/fs"
	"github.com/Equationzhao/g/internal/request"
)

type Root struct {
	Path    string
	Abs     string
	Entries []entry.Entry
	Err     error
	Code    int
}

func Walk(fsys gfs.Filesystem, req request.Request) []Root {
	out := make([]Root, 0, len(req.Paths))
	for _, p := range req.Paths {
		out = append(out, walkRoot(fsys, req, p))
	}
	return out
}

func walkRoot(fsys gfs.Filesystem, req request.Request, p string) Root {
	abs, err := fsys.Abs(p)
	if err != nil {
		abs = p
	}
	r := Root{Path: p, Abs: abs}
	fi, err := lstatOrStat(fsys, req, p)
	if err != nil {
		r.Err = err
		r.Code = 2
		return r
	}
	self := toEntry(fsys, req, p, abs, filepath.Dir(abs), 0, fi, true)
	deep := req.Recurse || req.Format == request.FormatTree
	if req.DirSelf || !fi.IsDir() || (req.HasDepth && req.Depth == 0 && !deep) {
		r.Entries = []entry.Entry{self}
		return r
	}
	if req.HasDepth && req.Depth == 0 && req.Recurse && req.Format != request.FormatTree {
		r.Entries = []entry.Entry{self}
		return r
	}
	seen := map[[2]uint64]bool{}
	if self.HasDevIno {
		seen[[2]uint64{self.Dev, self.Ino}] = true
	}
	kids, code := walkDir(fsys, req, p, abs, 1, seen)
	r.Code = code
	if req.Format == request.FormatTree {
		r.Entries = append([]entry.Entry{self}, kids...)
		return r
	}
	r.Entries = kids
	return r
}

func walkDir(fsys gfs.Filesystem, req request.Request, display, abs string, depth int, seen map[[2]uint64]bool) ([]entry.Entry, int) {
	if req.HasDepth && depth > req.Depth {
		return nil, 0
	}
	ents, err := fsys.ReadDir(abs)
	if err != nil {
		return nil, 1
	}
	out := make([]entry.Entry, 0, len(ents)+2)
	if depth == 1 && req.Visibility == request.VisAll {
		out = append(out, dotEntry(abs, ".", depth), dotEntry(abs, "..", depth))
	}
	code := 0
	order := 0
	for _, de := range ents {
		fi, err := de.Info()
		if err != nil {
			code = 1
			continue
		}
		childAbs := filepath.Join(abs, de.Name())
		childDisp := filepath.Join(display, de.Name())
		ent := toEntry(fsys, req, childDisp, childAbs, abs, depth, fi, false)
		ent.ReadOrder = order
		order++
		out = append(out, ent)
		descend := (req.Recurse || req.Format == request.FormatTree) && fi.IsDir()
		if !descend {
			continue
		}
		if req.HasDepth && depth >= req.Depth {
			continue
		}
		if ent.HasDevIno {
			key := [2]uint64{ent.Dev, ent.Ino}
			if seen[key] {
				continue
			}
			seen[key] = true
		}
		nested, c := walkDir(fsys, req, childDisp, childAbs, depth+1, seen)
		if c > code {
			code = c
		}
		out = append(out, nested...)
	}
	return out, code
}

func lstatOrStat(fsys gfs.Filesystem, req request.Request, p string) (gfs.FileInfo, error) {
	if req.Dereference {
		return fsys.Stat(p)
	}
	return fsys.Lstat(p)
}

func toEntry(fsys gfs.Filesystem, req request.Request, display, abs, parent string, depth int, fi gfs.FileInfo, isRoot bool) entry.Entry {
	name := fi.Name()
	if base := filepath.Base(display); base != "" && base != "." && base != string(filepath.Separator) {
		name = base
	}
	if display == "." || display == ".." {
		name = display
	}
	e := entry.Entry{
		Name:       name,
		Path:       abs,
		Parent:     parent,
		Depth:      depth,
		Size:       fi.Size(),
		Mode:       uint32(fi.Mode()),
		ModTime:    fi.ModTime(),
		AccTime:    fi.AccTime(),
		ChangeTime: fi.ChangeTime(),
		Hidden:     fi.Hidden() || strings.HasPrefix(name, "."),
		IsRootArg:  isRoot,
		Nlink:      fi.Nlink(),
		Blocks:     fi.Blocks(),
	}
	if ino := fi.Inode(); ino != 0 {
		e.Inode = fmt.Sprintf("%d", ino)
	} else {
		e.Inode = "-"
	}
	e.Dev, e.Ino, e.HasDevIno = fi.DevIno()
	if b, ok := fi.BirthTime(); ok {
		e.Birth, e.HasBirth = b, true
	}
	e.UID = fmt.Sprintf("%d", fi.UID())
	e.GID = fmt.Sprintf("%d", fi.GID())
	e.Kind = kindOf(fi)
	if fi.Mode()&os.ModeSymlink != 0 {
		e.Kind = entry.KindSymlink
		if t, err := fsys.Readlink(abs); err == nil {
			e.Target = t
			if _, err := fsys.Stat(abs); err != nil {
				e.Kind = entry.KindBrokenSymlink
			} else {
				e.TargetOK = true
			}
		}
	}
	if !fi.IsDir() && fi.Mode()&0o111 != 0 && e.Kind == entry.KindFile {
		e.Kind = entry.KindExec
	}
	return e
}

func kindOf(fi gfs.FileInfo) entry.Kind {
	m := fi.Mode()
	switch {
	case m.IsDir():
		return entry.KindDir
	case m&fs.ModeSymlink != 0:
		return entry.KindSymlink
	case m&fs.ModeNamedPipe != 0:
		return entry.KindPipe
	case m&fs.ModeSocket != 0:
		return entry.KindSocket
	case m&fs.ModeDevice != 0 && m&fs.ModeCharDevice != 0:
		return entry.KindChar
	case m&fs.ModeDevice != 0:
		return entry.KindBlock
	default:
		return entry.KindFile
	}
}

func dotEntry(abs, name string, depth int) entry.Entry {
	p := abs
	if name == ".." {
		p = filepath.Dir(abs)
	}
	return entry.Entry{
		Name: name, Path: p, Parent: abs, Depth: depth,
		Kind: entry.KindDir, Mode: uint32(os.ModeDir | 0o755),
	}
}
