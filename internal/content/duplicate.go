package content

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/Equationzhao/g/internal/cached"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/util"
)

type (
	filenameList = *util.Slice[string]
	hashStr      = string
)

type DuplicateDetect struct {
	IsThrough bool
	hashTb    *cached.Map[hashStr, filenameList]
}

type DOption func(d *DuplicateDetect)

const defaultTbSize = 200

func NewDuplicateDetect(options ...DOption) *DuplicateDetect {
	d := &DuplicateDetect{}

	for _, option := range options {
		option(d)
	}

	if d.hashTb == nil {
		d.hashTb = cached.NewCacheMap[hashStr, filenameList](defaultTbSize)
	}

	return d
}

func DuplicateWithTbSize(size int) DOption {
	return func(d *DuplicateDetect) {
		d.hashTb = cached.NewCacheMap[hashStr, filenameList](size)
	}
}

func DetectorFallthrough(d *DuplicateDetect) {
	d.IsThrough = true
}

func (d *DuplicateDetect) Enable() NoOutputOption {
	job := make(chan *item.FileInfo, 100)
	isJobFinished := cached.NewCacheMap[string, *chan struct{}](1000)
	go func() {
		for info := range job {
			func() {
				c, _ := isJobFinished.Get(info.FullPath)
				defer func() {
					*c <- struct{}{}
				}()
				afterHash, err := fileHash(info, d.IsThrough)
				if err != nil {
					return
				}
				actual, _ := d.hashTb.GetOrCompute(
					afterHash, func() filenameList {
						return util.NewSlice[string](10)
					},
				)
				actual.AppendTo(info.Name())
			}()
		}
	}()
	return func(info *item.FileInfo) {
		c := make(chan struct{}, 1)
		isJobFinished.Set(info.FullPath, &c)
		job <- info
		<-c
	}
}

type Duplicate struct {
	Filenames []string
}

func (d *DuplicateDetect) Result() []Duplicate {
	list := d.hashTb.Values()
	res := make([]Duplicate, 0, len(list))
	for _, i := range list {
		if l := i.Len(); l > 1 {
			f := i.GetCopy()
			slices.SortStableFunc(f, func(a, b string) int {
				return strings.Compare(a, b)
			})
			res = append(res, Duplicate{Filenames: f})
		}
	}
	return res
}

func (d *DuplicateDetect) Reset() {
	d.hashTb.ForEach(func(k string, v *util.Slice[string]) bool {
		v.Clear()
		return true
	})
}

func (d *DuplicateDetect) Fprint(w io.Writer) {
	r := d.Result()
	if len(r) != 0 {
		_, _ = fmt.Fprintln(w, "Duplicates:")
		for _, i := range r {
			for _, filename := range i.Filenames {
				_, _ = fmt.Fprint(w, "    ", filename)
			}
			_, _ = fmt.Fprintln(w)
		}
	}
}

var thresholdFileSize = int64(16 * KiB)

// fileHash calculates the hash of the file provided.
// If isThorough is true, then it uses SHA256 of the entire file.
// Otherwise, it uses CRC32 of "crucial bytes" of the file.
func fileHash(fileInfo *item.FileInfo, isThorough bool) (string, error) {
	if !fileInfo.Mode().IsRegular() {
		return "", fmt.Errorf("can't compute hash of non-regular file")
	}
	var prefix string
	var bytes []byte
	var fileReadErr error
	if isThorough || fileInfo.Size() <= thresholdFileSize {
		bytes, fileReadErr = os.ReadFile(fileInfo.FullPath)
		if fileReadErr != nil {
			return "", fmt.Errorf("couldn't read file: %w", fileReadErr)
		}
		if fileInfo.Size() <= thresholdFileSize {
			prefix = "f"
		}
	} else {
		prefix = "s"
		bytes, fileReadErr = readCrucialBytes(fileInfo.FullPath, fileInfo.Size())
		if fileReadErr != nil {
			return "", fmt.Errorf("couldn't calculate hash: %w", fileReadErr)
		}
	}
	var h hash.Hash
	if isThorough {
		h = sha256.New()
	} else {
		h = crc32.NewIEEE()
	}
	_, hashErr := h.Write(bytes)
	if hashErr != nil {
		return "", fmt.Errorf("error while computing hash: %w", hashErr)
	}
	hashBytes := h.Sum(nil)
	return prefix + hex.EncodeToString(hashBytes), nil
}

// readCrucialBytes reads the first few bytes, middle bytes and last few bytes of the file
func readCrucialBytes(filePath string, fileSize int64) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	firstBytes := make([]byte, thresholdFileSize/2)
	_, fErr := file.ReadAt(firstBytes, 0)
	if fErr != nil {
		return nil, fmt.Errorf("couldn't read first few bytes (maybe file is corrupted?): %w", fErr)
	}
	middleBytes := make([]byte, thresholdFileSize/4)
	_, mErr := file.ReadAt(middleBytes, fileSize/2)
	if mErr != nil {
		return nil, fmt.Errorf("couldn't read middle bytes (maybe file is corrupted?): %w", mErr)
	}
	lastBytes := make([]byte, thresholdFileSize/4)
	_, lErr := file.ReadAt(lastBytes, fileSize-thresholdFileSize/4)
	if lErr != nil {
		return nil, fmt.Errorf("couldn't read end bytes (maybe file is corrupted?): %w", lErr)
	}
	bytes := append(append(firstBytes, middleBytes...), lastBytes...)
	return bytes, nil
}
