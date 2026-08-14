package fs

import (
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

// Mem is an in-memory Filesystem for tests.
type Mem struct {
	Now   time.Time
	nodes map[string]*memNode
}

type memNode struct {
	name     string
	mode     fs.FileMode
	size     int64
	mod      time.Time
	acc      time.Time
	ctime    time.Time
	birth    time.Time
	hasBirth bool
	target   string
	hidden   bool
	uid, gid uint32
	user     string
	group    string
	nlink    uint64
	ino      uint64
	dev      uint64
	blocks   int64
}

func NewMem() *Mem {
	now := time.Now()
	m := &Mem{Now: now, nodes: map[string]*memNode{}}
	m.nodes["/"] = &memNode{name: "/", mode: fs.ModeDir | 0o755, mod: now, acc: now, ctime: now, nlink: 2, ino: 1, dev: 1}
	return m
}

func (m *Mem) AddDir(p string) *memNode {
	p = clean(p)
	n := &memNode{name: path.Base(p), mode: fs.ModeDir | 0o755, mod: m.Now, acc: m.Now, ctime: m.Now, nlink: 2, ino: uint64(len(m.nodes) + 1), dev: 1}
	if strings.HasPrefix(n.name, ".") {
		n.hidden = true
	}
	m.nodes[p] = n
	return n
}

func (m *Mem) AddFile(p string, size int64) *memNode {
	p = clean(p)
	n := &memNode{name: path.Base(p), mode: 0o644, size: size, mod: m.Now, acc: m.Now, ctime: m.Now, nlink: 1, ino: uint64(len(m.nodes) + 1), dev: 1, blocks: (size + 511) / 512}
	if strings.HasPrefix(n.name, ".") {
		n.hidden = true
	}
	m.nodes[p] = n
	return n
}

func (m *Mem) AddSymlink(p, target string) *memNode {
	p = clean(p)
	n := &memNode{name: path.Base(p), mode: fs.ModeSymlink | 0o777, target: target, mod: m.Now, acc: m.Now, ctime: m.Now, nlink: 1, ino: uint64(len(m.nodes) + 1), dev: 1}
	m.nodes[p] = n
	return n
}

func (m *Mem) Lstat(name string) (FileInfo, error) {
	n, ok := m.nodes[clean(name)]
	if !ok {
		return nil, &os.PathError{Op: "lstat", Path: name, Err: os.ErrNotExist}
	}
	return n, nil
}

func (m *Mem) Stat(name string) (FileInfo, error) {
	seen := map[string]bool{}
	cur := clean(name)
	for i := 0; i < 40; i++ {
		n, ok := m.nodes[cur]
		if !ok {
			return nil, &os.PathError{Op: "stat", Path: name, Err: os.ErrNotExist}
		}
		if n.mode&fs.ModeSymlink == 0 {
			return n, nil
		}
		if seen[cur] {
			return nil, &os.PathError{Op: "stat", Path: name, Err: os.ErrInvalid}
		}
		seen[cur] = true
		if path.IsAbs(n.target) {
			cur = clean(n.target)
		} else {
			cur = clean(path.Join(path.Dir(cur), n.target))
		}
	}
	return nil, &os.PathError{Op: "stat", Path: name, Err: os.ErrInvalid}
}

func (m *Mem) ReadDir(name string) ([]DirEntry, error) {
	dir := clean(name)
	n, ok := m.nodes[dir]
	if !ok {
		return nil, &os.PathError{Op: "readdir", Path: name, Err: os.ErrNotExist}
	}
	if !n.mode.IsDir() {
		return nil, &os.PathError{Op: "readdir", Path: name, Err: os.ErrInvalid}
	}
	var out []DirEntry
	prefix := dir
	if prefix != "/" {
		prefix += "/"
	}
	for p, child := range m.nodes {
		if p == dir {
			continue
		}
		if !strings.HasPrefix(p, prefix) {
			continue
		}
		rest := strings.TrimPrefix(p, prefix)
		if rest == "" || strings.Contains(rest, "/") {
			continue
		}
		out = append(out, memDir{n: child, path: p})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name() < out[j].Name() })
	return out, nil
}

func (m *Mem) Readlink(name string) (string, error) {
	n, ok := m.nodes[clean(name)]
	if !ok {
		return "", &os.PathError{Op: "readlink", Path: name, Err: os.ErrNotExist}
	}
	if n.mode&fs.ModeSymlink == 0 {
		return "", &os.PathError{Op: "readlink", Path: name, Err: os.ErrInvalid}
	}
	return n.target, nil
}

func (m *Mem) Abs(name string) (string, error) {
	if path.IsAbs(name) || strings.HasPrefix(name, "/") {
		return clean(name), nil
	}
	return clean("/" + name), nil
}

func clean(p string) string {
	p = strings.ReplaceAll(p, "\\", "/")
	if p == "" {
		return "/"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return path.Clean(p)
}

func (n *memNode) Name() string       { return n.name }
func (n *memNode) Size() int64        { return n.size }
func (n *memNode) Mode() fs.FileMode  { return n.mode }
func (n *memNode) ModTime() time.Time { return n.mod }
func (n *memNode) IsDir() bool        { return n.mode.IsDir() }
func (n *memNode) Sys() any           { return n }
func (n *memNode) DevIno() (uint64, uint64, bool) {
	return n.dev, n.ino, true
}
func (n *memNode) Hidden() bool                 { return n.hidden || strings.HasPrefix(n.name, ".") }
func (n *memNode) Nlink() uint64                { return n.nlink }
func (n *memNode) Inode() uint64                { return n.ino }
func (n *memNode) Blocks() int64                { return n.blocks }
func (n *memNode) UID() uint32                  { return n.uid }
func (n *memNode) GID() uint32                  { return n.gid }
func (n *memNode) AccTime() time.Time           { return n.acc }
func (n *memNode) ChangeTime() time.Time        { return n.ctime }
func (n *memNode) BirthTime() (time.Time, bool) { return n.birth, n.hasBirth }

type memDir struct {
	n    *memNode
	path string
}

func (d memDir) Name() string      { return d.n.name }
func (d memDir) IsDir() bool       { return d.n.IsDir() }
func (d memDir) Type() fs.FileMode { return d.n.mode.Type() }
func (d memDir) Info() (FileInfo, error) {
	return d.n, nil
}
