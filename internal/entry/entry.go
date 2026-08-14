package entry

import "time"

type Kind uint8

const (
	KindFile Kind = iota
	KindDir
	KindSymlink
	KindBrokenSymlink
	KindExec
	KindPipe
	KindSocket
	KindChar
	KindBlock
)

type Entry struct {
	Name       string
	Path       string
	Parent     string
	Depth      int
	Kind       Kind
	Mode       uint32
	Size       int64
	Nlink      uint64
	Inode      string
	Blocks     int64
	UID, GID   string
	User       string
	Group      string
	ModTime    time.Time
	AccTime    time.Time
	ChangeTime time.Time
	Birth      time.Time
	HasBirth   bool
	Target     string
	TargetOK   bool
	Hidden     bool
	Dev, Ino   uint64
	HasDevIno  bool
	ReadOrder  int
	IsRootArg  bool
	Git        string // two chars, empty if not requested
}

func (e Entry) IsDir() bool { return e.Kind == KindDir }
