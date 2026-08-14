package git

import (
	"bytes"
	"context"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Status byte

const (
	StatusNone Status = iota
	StatusModified
	StatusAdded
	StatusDeleted
	StatusRenamed
	StatusCopied
	StatusUntracked
	StatusIgnored
	StatusTypeChanged
	StatusUnmerged
)

var rank = [...]int{
	StatusNone:        0,
	StatusIgnored:     1,
	StatusUntracked:   2,
	StatusAdded:       3,
	StatusModified:    4,
	StatusTypeChanged: 5,
	StatusCopied:      6,
	StatusRenamed:     7,
	StatusDeleted:     8,
	StatusUnmerged:    9,
}

type FileStatus struct {
	RelPath string
	X, Y    Status
}

type RepoStatus struct {
	Root  string
	OK    bool
	Err   error
	Files []FileStatus
}

type Client interface {
	Status(ctx context.Context, dir string) (RepoStatus, error)
}

type Exec struct{}

func (Exec) Status(ctx context.Context, dir string) (RepoStatus, error) {
	nullf := "/dev/null"
	if runtime.GOOS == "windows" {
		nullf = "NUL"
	}
	cmd := exec.CommandContext(ctx, "git", "-C", dir, "status", "--porcelain=v1", "--ignored", "--untracked-files=normal")
	cmd.Env = append(cmd.Environ(),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL="+nullf,
		"GIT_CONFIG_SYSTEM="+nullf,
	)
	out, err := cmd.Output()
	if err != nil {
		return RepoStatus{Root: dir, Err: err}, err
	}
	rs := RepoStatus{Root: dir, OK: true}
	for _, line := range bytes.Split(out, []byte{'\n'}) {
		if len(line) < 3 {
			continue
		}
		x, y := ParseShort(line[0], line[1])
		rel := string(line[3:])
		if i := strings.Index(rel, " -> "); i >= 0 {
			rel = rel[i+4:]
		}
		rel = strings.Trim(rel, "\"")
		rs.Files = append(rs.Files, FileStatus{RelPath: filepath.ToSlash(rel), X: x, Y: y})
	}
	return rs, nil
}

func ParseShort(xb, yb byte) (Status, Status) {
	return byte2(xb), byte2(yb)
}

func byte2(b byte) Status {
	switch b {
	case ' ', 0:
		return StatusNone
	case 'M':
		return StatusModified
	case 'A':
		return StatusAdded
	case 'D':
		return StatusDeleted
	case 'R':
		return StatusRenamed
	case 'C':
		return StatusCopied
	case '?':
		return StatusUntracked
	case '!':
		return StatusIgnored
	case 'T':
		return StatusTypeChanged
	case 'U':
		return StatusUnmerged
	default:
		return StatusNone
	}
}

func letter(s Status) byte {
	switch s {
	case StatusModified:
		return 'M'
	case StatusAdded:
		return 'A'
	case StatusDeleted:
		return 'D'
	case StatusRenamed:
		return 'R'
	case StatusCopied:
		return 'C'
	case StatusUntracked:
		return '?'
	case StatusIgnored:
		return '!'
	case StatusTypeChanged:
		return 'T'
	case StatusUnmerged:
		return 'U'
	default:
		return '-'
	}
}

func Lookup(rs RepoStatus, rel string, isDir bool) string {
	if !rs.OK {
		return "--"
	}
	rel = filepath.ToSlash(rel)
	if !isDir {
		for _, f := range rs.Files {
			if f.RelPath == rel {
				return string([]byte{letter(f.X), letter(f.Y)})
			}
		}
		return "--"
	}
	var mx, my Status
	prefix := rel
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	found := rel == ""
	for _, f := range rs.Files {
		if f.RelPath == rel || (prefix != "" && strings.HasPrefix(f.RelPath, prefix)) || rel == "." {
			found = true
			if rank[f.X] > rank[mx] {
				mx = f.X
			}
			if rank[f.Y] > rank[my] {
				my = f.Y
			}
		}
	}
	if !found {
		return "--"
	}
	return string([]byte{letter(mx), letter(my)})
}

func Ignored(rs RepoStatus, rel string) bool {
	rel = filepath.ToSlash(rel)
	for _, f := range rs.Files {
		if f.RelPath == rel && (f.X == StatusIgnored || f.Y == StatusIgnored) {
			return true
		}
	}
	return false
}

func WithTimeout(c Client, dir string, d time.Duration) RepoStatus {
	if d <= 0 {
		d = 2 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	rs, err := c.Status(ctx, dir)
	if err != nil {
		return RepoStatus{Root: dir, OK: false, Err: err}
	}
	return rs
}

type Fake struct {
	Repos map[string]RepoStatus
	Err   error
}

func (f Fake) Status(_ context.Context, dir string) (RepoStatus, error) {
	if f.Err != nil {
		return RepoStatus{Root: dir, Err: f.Err}, f.Err
	}
	if rs, ok := f.Repos[dir]; ok {
		rs.OK = true
		return rs, nil
	}
	return RepoStatus{Root: dir, OK: true}, nil
}
