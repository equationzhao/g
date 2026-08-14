//go:build windows

package fs

import (
	"os"
	"syscall"
	"time"
)

func inspect(fi os.FileInfo, _ string) meta {
	m := meta{
		acc:   fi.ModTime(),
		ctime: fi.ModTime(),
	}
	st, ok := fi.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return m
	}
	m.hidden = st.FileAttributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0
	m.acc = time.Unix(0, st.LastAccessTime.Nanoseconds())
	m.ctime = fi.ModTime()
	m.birth = time.Unix(0, st.CreationTime.Nanoseconds())
	m.hasBirth = true
	return m
}
