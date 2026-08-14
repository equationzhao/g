//go:build linux

package fs

import (
	"syscall"
	"time"
)

func atime(st *syscall.Stat_t) time.Time {
	return time.Unix(st.Atim.Sec, st.Atim.Nsec)
}

func ctime(st *syscall.Stat_t) time.Time {
	return time.Unix(st.Ctim.Sec, st.Ctim.Nsec)
}
