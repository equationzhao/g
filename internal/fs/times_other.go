//go:build unix && !linux && !darwin

package fs

import (
	"syscall"
	"time"
)

func atime(st *syscall.Stat_t) time.Time { _ = st; return time.Unix(0, 0) }
func ctime(st *syscall.Stat_t) time.Time { _ = st; return time.Unix(0, 0) }
