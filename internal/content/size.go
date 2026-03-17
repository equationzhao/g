package content

import (
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/Equationzhao/g/internal/align"

	constval "github.com/Equationzhao/g/internal/global"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/osbased"
	"github.com/Equationzhao/g/internal/render"
	"github.com/Equationzhao/g/internal/util"
)

const Unknown SizeUnit = -1

const (
	Auto SizeUnit = iota
	Bit  SizeUnit = 1.0 << (10 * iota)
	Byte
	KiB
	MiB
	GiB
	TiB
)

// SI units
const (
	KB = 1000 * Byte
	MB = 1000 * KB
	GB = 1000 * MB
	TB = 1000 * GB
)

const SizeName = constval.NameOfSize

type Size struct {
	Bytes uint64
}

func CountBytes(size SizeUnit) uint64 {
	return uint64(size / Byte)
}

func ParseSize(size string) (Size, error) {
	sizeFloat, unit := util.SplitNumberAndUnit(size)
	if sizeFloat < 0 {
		return Size{}, fmt.Errorf("size can't be negative")
	}
	sizeUnit := string2SizeUnit(unit)
	if sizeUnit <= Bit {
		return Size{}, fmt.Errorf("invalid size unit: %s", unit)
	}
	return Size{Bytes: uint64(sizeFloat * float64(CountBytes(sizeUnit)))}, nil
}

type SizeUnit float64

func sizeStringSets(size SizeUnit) []string {
	switch size {
	case Unknown:
		return []string{"unknown"}
	case Auto:
		return []string{""}
	case Bit:
		return []string{"bit", "Bit"}
	case Byte:
		return []string{"B", "b", "byte", "Byte"}
	case KB:
		return []string{"KB", "kb", "K", "k"}
	case KiB:
		return []string{"KiB"}
	case MB:
		return []string{"MB", "mb", "M", "m"}
	case MiB:
		return []string{"MiB"}
	case GB:
		return []string{"GB", "gb", "G", "g"}
	case GiB:
		return []string{"GiB"}
	case TB:
		return []string{"TB", "tb", "T", "t"}
	case TiB:
		return []string{"TiB"}
	default:
		return []string{"unknown"}
	}
}

var string2SizeUnitMap = map[string]SizeUnit{
	"bit":  Bit,
	"Bit":  Bit,
	"byte": Byte,
	"Byte": Byte,
	"KB":   KB,
	"KiB":  KiB,
	"kb":   KB,
	"K":    KB,
	"k":    KB,
	"MB":   MB,
	"MiB":  MiB,
	"mb":   MB,
	"M":    MB,
	"m":    MB,
	"GB":   GB,
	"GiB":  GiB,
	"gb":   GB,
	"G":    GB,
	"g":    GB,
	"TB":   TB,
	"TiB":  TiB,
	"tb":   TB,
	"T":    TB,
	"t":    TB,
}

func string2SizeUnit(size string) SizeUnit {
	if unit, ok := string2SizeUnitMap[size]; ok {
		return unit
	}
	return Unknown
}

func Convert2SizeString(size SizeUnit) string {
	return sizeStringSets(size)[0]
}

func ConvertFromSizeString(size string, isSI bool) SizeUnit {
	switch size {
	case "bit", "Bit", "BIT":
		return Bit
	case "B", "b", "byte", "Byte", "BYTE":
		return Byte
	case "KB", "kb", "Kb", "k":
		if isSI {
			return KB
		}
		return KiB
	case "MB", "mb", "Mb", "M", "m":
		if isSI {
			return MB
		}
		return MiB
	case "GB", "gb", "Gb", "G", "g":
		if isSI {
			return GB
		}
		return GiB
	case "TB", "tb", "Tb", "T", "t":
		if isSI {
			return TB
		}
		return TiB
	default:
		return Unknown
	}
}

type SizeEnabler struct {
	total       atomic.Int64
	enableTotal bool
	sizeUint    SizeUnit
	recursive   *SizeRecursive
	isSi        bool
}

func (s *SizeEnabler) Recursive() *SizeRecursive {
	return s.recursive
}

func (s *SizeEnabler) SetRecursive(sr *SizeRecursive) {
	s.recursive = sr
}

func (s *SizeEnabler) SetSI() *SizeEnabler {
	s.isSi = true
	return s
}

type SizeRecursive struct {
	depth int
}

func NewSizeRecursive(depth int) *SizeRecursive {
	return &SizeRecursive{depth: depth}
}

func NewSizeEnabler() *SizeEnabler {
	return &SizeEnabler{
		total:       atomic.Int64{},
		enableTotal: false,
		sizeUint:    Auto,
		recursive:   nil,
	}
}

func (s *SizeEnabler) SizeUint() SizeUnit {
	return s.sizeUint
}

func (s *SizeEnabler) SetEnableTotal() {
	s.enableTotal = true
}

func (s *SizeEnabler) DisableTotal() {
	s.enableTotal = false
}

func (s *SizeEnabler) Total() (size int64, ok bool) {
	if s.enableTotal {
		return s.total.Load(), s.enableTotal
	}
	return 0, false
}

func (s *SizeEnabler) Reset() {
	if s.enableTotal {
		s.total.Store(0)
	}
}

func (s *SizeEnabler) Size2String(b int64) (string, SizeUnit) {
	var res string
	actualUnit := s.sizeUint
	v := float64(b)
	switch s.sizeUint {
	case Bit:
		res = strconv.FormatInt(int64(v*8), 10)
	case Byte:
		res = strconv.FormatInt(int64(v), 10)
	case KiB, MiB, GiB, TiB, KB, MB, GB, TB:
		res = fmt.Sprintf("%.1f", v*float64(Byte)/float64(s.sizeUint))
	case Auto:
		maxUnit, gap := TiB, SizeUnit(1024)
		if s.isSi {
			maxUnit, gap = TB, SizeUnit(1000)
		}
		for i := Byte; i <= maxUnit; i *= gap {
			if v < float64(gap) {
				switch i {
				case Byte:
					res = strconv.FormatInt(int64(v), 10)
				case Bit:
					res = strconv.FormatInt(int64(v*8), 10)
				default:
					res = strconv.FormatFloat(v, 'f', 1, 64)
				}
				if res == "0" || res == "0.0" {
					// make align
					return "       -", actualUnit
				}
				res += " " + Convert2SizeString(i)
				actualUnit = i
				return res, actualUnit
			}
			v /= float64(gap)
		}
		panic("too large")
	default:
		panic("invalid size uint" + strconv.Itoa(int(s.sizeUint)))
	}

	if res == "0" {
		res = "-"
	} else {
		res += " " + Convert2SizeString(s.sizeUint)
	}
	return res, actualUnit
}

const RecursiveSizeName = "recursive_size"

func (s *SizeEnabler) EnableSize(size SizeUnit, renderer *render.Renderer) ContentOption {
	s.sizeUint = size
	align.RegisterHeaderFooter(SizeName)
	return func(info *item.FileInfo) (string, string) {
		var v int64
		if s.recursive != nil {
			if r, ok := info.Cache[RecursiveSizeName]; ok {
				// convert []byte to int64
				v, _ = strconv.ParseInt(string(r), 10, 64)
			} else {
				v = util.RecursivelySizeOf(info, s.recursive.depth)
			}
		} else {
			v = info.Size()
		}
		if s.enableTotal {
			s.total.Add(v)
		}
		res, unit := s.Size2String(v)
		return renderer.Size(res, Convert2SizeString(unit)), SizeName
	}
}

type BlockSizeEnabler struct{}

func NewBlockSizeEnabler() *BlockSizeEnabler {
	return &BlockSizeEnabler{}
}

const BlockSizeName = "Blocks"

func (b *BlockSizeEnabler) Enable(renderer *render.Renderer) ContentOption {
	align.RegisterHeaderFooter(BlockSizeName)
	return func(info *item.FileInfo) (string, string) {
		res := ""
		bs := osbased.BlockSize(info)
		if bs == 0 {
			res = "-"
		} else {
			res = strconv.FormatInt(bs, 10)
		}
		return renderer.BlockSize(res), BlockSizeName
	}
}
