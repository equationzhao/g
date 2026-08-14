//go:build unix && !darwin && !linux

package fs

import (
	"syscall"
	"time"
)

func birthTimespec(st *syscall.Stat_t) (time.Time, bool) {
	_ = st
	return time.Time{}, false
}
