package filter

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Equationzhao/g/internal/git"
	constval "github.com/Equationzhao/g/internal/global"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/osbased"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gobwas/glob"
)

const (
	keep   = true
	remove = false
)

type ItemFilter struct {
	tfs []*ItemFilterFunc
}

func (tf *ItemFilter) AppendTo(typeFunc ...*ItemFilterFunc) {
	tf.tfs = append(tf.tfs, typeFunc...)
}

func NewItemFilter(tfs ...*ItemFilterFunc) *ItemFilter {
	return &ItemFilter{tfs: tfs}
}

func (tf *ItemFilter) Match(e *item.FileInfo) bool {
	ok := keep
	for _, funcPtr := range tf.tfs {
		if !(*funcPtr)(e) {
			ok = remove
			break
		}
	}
	return ok
}

func (tf *ItemFilter) Filter(e ...*item.FileInfo) (res []*item.FileInfo) {
	for _, entry := range e {
		ok := keep
		for _, funcPtr := range tf.tfs {
			if !(*funcPtr)(entry) {
				ok = remove
				break
			}
		}
		if ok {
			res = append(res, entry)
		}
	}
	return res
}

// ItemFilterFunc return true -> Keep
// return false -> remove
type ItemFilterFunc = func(e *item.FileInfo) bool

var RemoveDir = func(e *item.FileInfo) bool {
	return !e.IsDir()
}

var DirOnly = func(e *item.FileInfo) bool {
	return e.IsDir()
}

// RemoveByExt
//
//	eg:
//		a.go b.c c.rs d.cxx dir
//		RemoveByExt([]string{"go", "cxx"})
//	result:
//		b.c c.rs dir
func RemoveByExt(ext ...string) ItemFilterFunc {
	return func(e *item.FileInfo) bool {
		for _, extI := range ext {
			if strings.HasSuffix(e.Name(), "."+extI) {
				return remove
			}
		}
		return keep
	}
}

func ExtOnly(ext ...string) ItemFilterFunc {
	return func(e *item.FileInfo) bool {
		for _, extI := range ext {
			if strings.HasSuffix(e.Name(), "."+extI) {
				return keep
			}
		}
		return remove
	}
}

// RemoveGlob if all pattern complied successfully, return a func and nil error,
// if match any one, the fn will return remove, else return keep
// if error occurred, return nil func and error
func RemoveGlob(globPattern ...string) (ItemFilterFunc, error) {
	compiled := make([]glob.Glob, 0, len(globPattern))
	for _, v := range globPattern {
		compile, err := glob.Compile(v)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, compile)
	}

	return func(e *item.FileInfo) bool {
		for _, r := range compiled {
			if r.Match(e.Name()) {
				return remove
			}
		}
		return keep
	}, nil
}

// GlobOnly if all pattern complied successfully, return a func and nil error,
// if match any one, the fn will return keep, else return remove
// if error occurred, return nil func and error
func GlobOnly(globPattern ...string) (ItemFilterFunc, error) {
	compiled := make([]glob.Glob, 0, len(globPattern))
	for _, v := range globPattern {
		compile, err := glob.Compile(v)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, compile)
	}

	return func(e *item.FileInfo) bool {
		for _, r := range compiled {
			if r.Match(e.Name()) {
				return keep
			}
		}
		return remove
	}, nil
}

var RemoveHidden = func(e *item.FileInfo) bool {
	return !strings.HasPrefix(e.Name(), ".")
}

var HiddenOnly = func(e *item.FileInfo) bool {
	return strings.HasPrefix(e.Name(), ".")
}

var RemoveBackups = func(e *item.FileInfo) bool {
	return !strings.HasSuffix(e.Name(), "~")
}

func RemoveGitIgnore(repoPath git.RepoPath) ItemFilterFunc {
	isOrIsParentOf := func(parent, child string) bool {
		if parent == child {
			return true
		}
		if strings.HasPrefix(child, parent+string(filepath.Separator)) { // should not use filepath.Separator
			return true
		}
		return false
	}
	ignoredCache := git.GetCache()

	return func(e *item.FileInfo) (ok bool) {
		actual, _ := ignoredCache.GetOrCompute(repoPath, git.DefaultInit(repoPath))
		ok = true
		topLevel, err := git.GetTopLevel(repoPath)
		if err != nil {
			return keep
		}
		rel, err := filepath.Rel(topLevel, e.FullPath)
		if err != nil {
			return keep
		}
		for _, fileGit := range *actual {
			if fileGit.X == git.Ignored || fileGit.Y == git.Ignored {
				if isOrIsParentOf(fileGit.Name, rel) {
					ok = remove
					break
				}
			}
		}
		return ok
	}
}

func isOrIsSonOf(a, b string) bool {
	if a == b {
		return true
	}
	if strings.HasPrefix(a, b+"/") {
		return true
	}
	return false
}

const MimeTypeName = constval.NameOfMIME

func MimeTypeOnly(fileTypes ...string) ItemFilterFunc {
	return func(e *item.FileInfo) bool {
		if e.IsDir() {
			return keep
		}
		file, err := os.Open(e.Name())
		if err != nil {
			return keep
		}
		mtype, err := mimetype.DetectReader(file)
		if err != nil {
			return keep
		}
		s := mtype.String()
		for i := range fileTypes {
			if strings.Contains(s, ";") {
				// remove charset
				s = strings.SplitN(s, ";", 2)[0]
			}
			if isOrIsSonOf(s, fileTypes[i]) {
				e.Cache[MimeTypeName] = []byte(s)
				return keep
			}
		}
		return remove
	}
}

func RemoveMimeType(fileTypes ...string) ItemFilterFunc {
	return func(e *item.FileInfo) bool {
		file, err := os.Open(e.Name())
		if err != nil {
			return keep
		}
		mtype, err := mimetype.DetectReader(file)
		if err != nil {
			return keep
		}

		s := mtype.String()
		for i := range fileTypes {
			if strings.Contains(s, ";") {
				// remove charset
				s = strings.SplitN(s, ";", 2)[0]
			}
			if fileTypes[i] == mtype.String() {
				return remove
			}
		}
		return keep
	}
}

func BeforeTime(t time.Time, timeFunc func(os.FileInfo) time.Time) ItemFilterFunc {
	return func(e *item.FileInfo) bool {
		return timeFunc(e).Before(t)
	}
}

func AfterTime(t time.Time, timeFunc func(os.FileInfo) time.Time) ItemFilterFunc {
	return func(e *item.FileInfo) bool {
		return timeFunc(e).After(t)
	}
}

func WhichTimeFiled(mod string) (t func(os.FileInfo) time.Time) {
	switch strings.ToLower(mod) {
	case "mod", "modified":
		t = osbased.ModTime
	case "create", "cr":
		t = osbased.CreateTime
	case "access", "ac":
		t = osbased.AccessTime
	case "birth":
		t = osbased.BirthTime
	default:
		t = osbased.ModTime
	}
	return t
}
