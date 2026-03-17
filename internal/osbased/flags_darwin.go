package osbased

import (
	"os"
	"slices"
	"syscall"

	"github.com/Equationzhao/g/internal/item"
	"golang.org/x/sys/unix"
)

var flags = map[int]string{
	unix.UF_APPEND:     "uappnd",
	unix.UF_COMPRESSED: "compressed",
	unix.UF_HIDDEN:     "hidden",
	unix.UF_IMMUTABLE:  "uchg",
	unix.UF_NODUMP:     "nodump",
	unix.UF_OPAQUE:     "opaque",

	// unix.UF_SETTABLE:  "UF_SETTABLE",
	// unix.UF_TRACKED:   "UF_TRACKED",
	// unix.UF_DATAVAULT: "UF_DATAVAULT",

	unix.SF_APPEND:     "sappnd",
	unix.SF_ARCHIVED:   "arch",
	unix.SF_DATALESS:   "dataless",
	unix.SF_IMMUTABLE:  "schg",
	unix.SF_RESTRICTED: "restricted",

	// unix.SF_SETTABLE:  "SF_SETTABLE",
	// unix.SF_SUPPORTED: "SF_SUPPORTED",
	// unix.SF_SYNTHETIC: "SF_SYNTHETIC",
	// unix.SF_FIRMLINK:  "SF_FIRMLINK",
}

func getFlags(filename string) uint32 {
	file, err := os.Open(filename)
	if err != nil {
		return 0
	}
	defer func() {
		_ = file.Close()
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return 0
	}

	stat := fileInfo.Sys().(*syscall.Stat_t)
	return stat.Flags
}

func CheckFlags(i *item.FileInfo) []string {
	res := make([]string, 0, 8)
	f := getFlags(i.FullPath)
	for key, val := range flags {
		if f&uint32(key) != 0 {
			res = append(res, val)
		}
	}
	slices.Sort(res)
	return res
}
