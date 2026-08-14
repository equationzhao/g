package fs

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileInfo interface {
	fs.FileInfo
	DevIno() (dev, ino uint64, ok bool)
	Hidden() bool
	Nlink() uint64
	Inode() uint64
	Blocks() int64
	UID() uint32
	GID() uint32
	AccTime() time.Time
	ChangeTime() time.Time
	BirthTime() (time.Time, bool)
}

type DirEntry interface {
	Name() string
	IsDir() bool
	Type() fs.FileMode
	Info() (FileInfo, error)
}

type Filesystem interface {
	Lstat(name string) (FileInfo, error)
	Stat(name string) (FileInfo, error)
	ReadDir(name string) ([]DirEntry, error)
	Readlink(name string) (string, error)
	Abs(name string) (string, error)
}

type OS struct{}

func (OS) Lstat(name string) (FileInfo, error) {
	fi, err := os.Lstat(name)
	if err != nil {
		return nil, err
	}
	return wrap(name, fi), nil
}

func (OS) Stat(name string) (FileInfo, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	return wrap(name, fi), nil
}

func (OS) ReadDir(name string) ([]DirEntry, error) {
	ents, err := os.ReadDir(name)
	if err != nil {
		return nil, err
	}
	out := make([]DirEntry, len(ents))
	for i, e := range ents {
		out[i] = osDir{dir: name, e: e}
	}
	return out, nil
}

func (OS) Readlink(name string) (string, error) { return os.Readlink(name) }

func (OS) Abs(name string) (string, error) { return filepath.Abs(name) }

type osFile struct {
	os.FileInfo
	path string
	meta
}

func wrap(path string, fi os.FileInfo) *osFile {
	o := &osFile{FileInfo: fi, path: path}
	o.meta = inspect(fi, path)
	return o
}

func (o *osFile) DevIno() (uint64, uint64, bool) { return o.dev, o.ino, o.hasDevIno }
func (o *osFile) Hidden() bool                   { return o.hidden || strings.HasPrefix(o.Name(), ".") }
func (o *osFile) Nlink() uint64                  { return o.nlink }
func (o *osFile) Inode() uint64                  { return o.ino }
func (o *osFile) Blocks() int64                  { return o.blocks }
func (o *osFile) UID() uint32                    { return o.uid }
func (o *osFile) GID() uint32                    { return o.gid }
func (o *osFile) AccTime() time.Time             { return o.acc }
func (o *osFile) ChangeTime() time.Time          { return o.ctime }
func (o *osFile) BirthTime() (time.Time, bool)   { return o.birth, o.hasBirth }

type osDir struct {
	dir string
	e   os.DirEntry
}

func (d osDir) Name() string      { return d.e.Name() }
func (d osDir) IsDir() bool       { return d.e.IsDir() }
func (d osDir) Type() fs.FileMode { return d.e.Type() }
func (d osDir) Info() (FileInfo, error) {
	fi, err := d.e.Info()
	if err != nil {
		return nil, err
	}
	return wrap(filepath.Join(d.dir, d.e.Name()), fi), nil
}
