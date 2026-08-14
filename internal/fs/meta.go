package fs

import "time"

type meta struct {
	dev, ino  uint64
	hasDevIno bool
	nlink     uint64
	blocks    int64
	uid, gid  uint32
	acc       time.Time
	ctime     time.Time
	birth     time.Time
	hasBirth  bool
	hidden    bool
}
