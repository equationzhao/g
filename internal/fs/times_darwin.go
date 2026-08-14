//go:build darwin

package fs

import (
	"syscall"
	"time"
)

func atime(st *syscall.Stat_t) time.Time {
	return time.Unix(st.Atimespec.Sec, st.Atimespec.Nsec)
}

func ctime(st *syscall.Stat_t) time.Time {
	return time.Unix(st.Ctimespec.Sec, st.Ctimespec.Nsec)
}
