//go:build unix

package fs

import (
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func inspect(fi os.FileInfo, _ string) meta {
	m := meta{
		acc:   fi.ModTime(),
		ctime: fi.ModTime(),
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return m
	}
	m.dev = uint64(st.Dev)
	m.ino = uint64(st.Ino)
	m.hasDevIno = true
	m.nlink = uint64(st.Nlink)
	m.blocks = int64(st.Blocks)
	m.uid = st.Uid
	m.gid = st.Gid
	m.acc = atime(st)
	m.ctime = ctime(st)
	if b, ok := birthTimespec(st); ok {
		m.birth = b
		m.hasBirth = true
	}
	_ = unix.Getpagesize()
	return m
}
