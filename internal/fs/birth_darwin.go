//go:build darwin

package fs

import (
	"syscall"
	"time"
)

func birthTimespec(st *syscall.Stat_t) (time.Time, bool) {
	return time.Unix(st.Birthtimespec.Sec, st.Birthtimespec.Nsec), true
}
