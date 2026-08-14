package sorter

import (
	"cmp"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Equationzhao/g/internal/content"
	constval "github.com/Equationzhao/g/internal/global"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/osbased"
	"github.com/Equationzhao/g/internal/util"
	mt "github.com/gabriel-vasile/mimetype"
)

// ByNone
// Deprecated
//
//goland:noinspection GoUnusedParameter
func ByNone(a, b *item.FileInfo) int {
	return 0
}

func ByNameDescend(a, b *item.FileInfo) int {
	return cmp.Compare(strings.ToLower(b.Name()), strings.ToLower(a.Name()))
}

func ByNameAscend(a, b *item.FileInfo) int {
	return cmp.Compare(strings.ToLower(a.Name()), strings.ToLower(b.Name()))
}

func ByNameCaseSensitiveDescend(a, b *item.FileInfo) int {
	return cmp.Compare(b.Name(), a.Name())
}

func ByNameCaseSensitiveAscend(a, b *item.FileInfo) int {
	return cmp.Compare(a.Name(), b.Name())
}

func ByNameWithoutALeadingDotDescend(a, b *item.FileInfo) int {
	return cmp.Compare(
		strings.ToLower(strings.TrimPrefix(b.Name(), ".")), strings.ToLower(strings.TrimPrefix(a.Name(), ".")),
	)
}

func ByNameWithoutALeadingDotAscend(a, b *item.FileInfo) int {
	return cmp.Compare(
		strings.ToLower(strings.TrimPrefix(a.Name(), ".")), strings.ToLower(strings.TrimPrefix(b.Name(), ".")),
	)
}

func ByNameWithoutALeadingDotCaseSensitiveDescend(a, b *item.FileInfo) int {
	return cmp.Compare(strings.TrimPrefix(b.Name(), "."), strings.TrimPrefix(a.Name(), "."))
}

func ByNameWithoutALeadingDotCaseSensitiveAscend(a, b *item.FileInfo) int {
	return cmp.Compare(strings.TrimPrefix(a.Name(), "."), strings.TrimPrefix(b.Name(), "."))
}

func byInode(a, b *item.FileInfo) (int, int) {
	inodeA, _ := strconv.Atoi(osbased.Inode(a))
	inodeB, _ := strconv.Atoi(osbased.Inode(b))
	a.Cache["Inode"] = []byte(strconv.Itoa(inodeA))
	b.Cache["Inode"] = []byte(strconv.Itoa(inodeB))
	return inodeA, inodeB
}

func ByInodeDescend(a, b *item.FileInfo) int {
	inodeA, inodeB := byInode(a, b)
	return cmp.Compare(inodeB, inodeA)
}

func ByInodeAscend(a, b *item.FileInfo) int {
	inodeA, inodeB := byInode(a, b)
	return cmp.Compare(inodeA, inodeB)
}

func BySizeDescend(a, b *item.FileInfo) int {
	return cmp.Compare(b.Size(), a.Size())
}

func BySizeAscend(a, b *item.FileInfo) int {
	return cmp.Compare(a.Size(), b.Size())
}

func ByRecursiveSizeDescend(depth int) FileSortFunc {
	return func(a, b *item.FileInfo) int {
		return byRecursiveSize(a, b, depth, false)
	}
}

func ByRecursiveSizeAscend(depth int) FileSortFunc {
	return func(a, b *item.FileInfo) int {
		return byRecursiveSize(a, b, depth, true)
	}
}

func ByTimeAscend(timeType string) FileSortFunc {
	switch timeType {
	case "mod", "modified":
		return func(a, b *item.FileInfo) int {
			return osbased.ModTime(b).Compare(osbased.ModTime(a))
		}
	case "access", "ac":
		return func(a, b *item.FileInfo) int {
			return osbased.AccessTime(b).Compare(osbased.AccessTime(a))
		}
	case "create", "cr":
		return func(a, b *item.FileInfo) int {
			return osbased.CreateTime(b).Compare(osbased.CreateTime(a))
		}
	case "birth":
		return func(a, b *item.FileInfo) int {
			return osbased.BirthTime(b).Compare(osbased.BirthTime(a))
		}
	default:
		panic("invalid time type")
	}
}

func ByTimeDescend(timeType string) FileSortFunc {
	switch timeType {
	case "mod", "modified":
		return func(a, b *item.FileInfo) int {
			return osbased.ModTime(a).Compare(osbased.ModTime(b))
		}
	case "access", "ac":
		return func(a, b *item.FileInfo) int {
			return osbased.AccessTime(a).Compare(osbased.AccessTime(b))
		}
	case "create", "cr":
		return func(a, b *item.FileInfo) int {
			return osbased.CreateTime(a).Compare(osbased.CreateTime(b))
		}
	case "birth":
		return func(a, b *item.FileInfo) int {
			return osbased.BirthTime(a).Compare(osbased.BirthTime(b))
		}
	default:
		panic("invalid time type")
	}
}

func ByExtensionDescend(a, b *item.FileInfo) int {
	return cmp.Compare(filepath.Ext(b.Name()), filepath.Ext(a.Name()))
}

func ByExtensionAscend(a, b *item.FileInfo) int {
	return cmp.Compare(filepath.Ext(a.Name()), filepath.Ext(b.Name()))
}

func ByExtensionCaseSensitiveDescend(a, b *item.FileInfo) int {
	return cmp.Compare(strings.ToLower(filepath.Ext(b.Name())), strings.ToLower(filepath.Ext(a.Name())))
}

func ByExtensionCaseSensitiveAscend(a, b *item.FileInfo) int {
	return cmp.Compare(strings.ToLower(filepath.Ext(a.Name())), strings.ToLower(filepath.Ext(b.Name())))
}

func ByGroupDescend(a, b *item.FileInfo) int {
	return byGroupName(a, b, false)
}

func ByGroupAscend(a, b *item.FileInfo) int {
	return byGroupName(a, b, true)
}

func ByGroupCaseSensitiveDescend(a, b *item.FileInfo) int {
	return byGroupCaseSensitiveName(a, b, false)
}

func ByGroupCaseSensitiveAscend(a, b *item.FileInfo) int {
	return byGroupCaseSensitiveName(a, b, true)
}

func ByOwnerDescend(a, b *item.FileInfo) int {
	return byUserName(a, b, false)
}

func ByOwnerAscend(a, b *item.FileInfo) int {
	return byUserName(a, b, true)
}

func ByOwnerCaseSensitiveDescend(a, b *item.FileInfo) int {
	return byUserCaseSensitiveName(a, b, false)
}

func ByOwnerCaseSensitiveAscend(a, b *item.FileInfo) int {
	return byUserCaseSensitiveName(a, b, true)
}

func ByNameWidthDescend(a, b *item.FileInfo) int {
	return byNameWidth(a, b, false)
}

func ByNameWidthAscend(a, b *item.FileInfo) int {
	return byNameWidth(a, b, true)
}

func ByMimeTypeAscend(a, b *item.FileInfo) int {
	return byMimeType(a, b, true)
}

func ByMimeTypeDescend(a, b *item.FileInfo) int {
	return byMimeType(a, b, false)
}

func byMimeType(a, b *item.FileInfo, ascend bool) int {
	mimeAstr, mimeBstr := getMimeName(a, b)
	if ascend {
		return cmp.Compare(mimeAstr, mimeBstr)
	}
	return cmp.Compare(mimeBstr, mimeAstr)
}

const MimeTypeName = constval.NameOfMIME

func getMimeName(a, b *item.FileInfo) (string, string) {
	mimeAstr, mimeBstr := "", ""
	if c, ok := a.Cache[MimeTypeName]; ok {
		mimeAstr = string(c)
	} else {
		mimeA, err := mt.DetectFile(a.FullPath)
		if err != nil {
			if a.IsDir() {
				mimeAstr = "directory"
			} else if a.Mode()&os.ModeSymlink != 0 {
				mimeAstr = "symlink"
			} else if a.Mode()&os.ModeNamedPipe != 0 {
				mimeAstr = "named_pipe"
			} else if a.Mode()&os.ModeSocket != 0 {
				mimeAstr = "socket"
			}
		} else {
			mimeAstr = mimeA.String()
		}
		a.Cache[MimeTypeName] = []byte(mimeAstr)
	}

	if c, ok := b.Cache[MimeTypeName]; ok {
		mimeBstr = string(c)
	} else {
		mimeB, err := mt.DetectFile(b.FullPath)
		if err != nil {
			if b.IsDir() {
				mimeBstr = "directory"
			} else if b.Mode()&os.ModeSymlink != 0 {
				mimeBstr = "symlink"
			} else if b.Mode()&os.ModeNamedPipe != 0 {
				mimeBstr = "named_pipe"
			} else if b.Mode()&os.ModeSocket != 0 {
				mimeBstr = "socket"
			}
		} else {
			mimeBstr = mimeB.String()
		}
		b.Cache[MimeTypeName] = []byte(mimeBstr)
	}

	return mimeAstr, mimeBstr
}

func ByMimeTypeParentAscend(a, b *item.FileInfo) int {
	return byMimeTypeParent(a, b, true)
}

func ByMimeTypeParentDescend(a, b *item.FileInfo) int {
	return byMimeTypeParent(a, b, false)
}

func byMimeTypeParent(a, b *item.FileInfo, ascend bool) int {
	mimeAstr, mimeBstr := getMimeName(a, b)
	if ascend {
		return cmp.Compare(strings.SplitN(mimeAstr, "/", 2)[0], strings.SplitN(mimeBstr, "/", 2)[0])
	}
	return cmp.Compare(strings.SplitN(mimeBstr, "/", 2)[0], strings.SplitN(mimeAstr, "/", 2)[0])
}

const RecursiveSizeName = content.RecursiveSizeName

func byRecursiveSize(a, b *item.FileInfo, depth int, ascend bool) int {
	var sa []byte
	var sb []byte
	exist := false
	sai, sbi := int64(0), int64(0)
	if sa, exist = a.Cache[RecursiveSizeName]; !exist {
		sai = util.RecursivelySizeOf(a, depth)
		sa = []byte(strconv.FormatInt(sai, 10))
		a.Cache[RecursiveSizeName] = sa
	} else {
		sai, _ = strconv.ParseInt(string(sa), 10, 64)
	}
	if sb, exist = b.Cache[RecursiveSizeName]; !exist {
		sbi = util.RecursivelySizeOf(b, depth)
		sb = []byte(strconv.FormatInt(sbi, 10))
		b.Cache[RecursiveSizeName] = sb
	} else {
		sbi, _ = strconv.ParseInt(string(sb), 10, 64)
	}

	if ascend {
		return int(sai - sbi)
	}
	return int(sbi - sai)
}

func dirFirst(a, b *item.FileInfo) int {
	hdA := isHiddenDir(a)
	hdB := isHiddenDir(b)
	if hdA && !hdB { // a is hidden dir, b is not
		return -1
	}
	if !hdA && hdB { // a is not hidden dir, b is hidden dir
		return 1
	}
	// same hidden dir status
	dA := a.IsDir()
	dB := b.IsDir()
	if dA && !dB { // a is dir, b is not
		return -1
	}
	if !dA && dB { // a is not dir, b is dir
		return 1
	}
	return 0
}

func Default(a, b *item.FileInfo) int {
	return compareFileInfo(a, b)
}

type FileSortFunc = func(a, b *item.FileInfo) int

type Sorter struct {
	reverse  bool
	dirFirst bool
	option   []FileSortFunc
}

func (s *Sorter) Reset() {
	s.reverse = false
	s.dirFirst = false
	s.option = make([]FileSortFunc, 0, 10)
}

func (s *Sorter) DirFirst() {
	s.dirFirst = true
}

func (s *Sorter) UnsetDirFirst() {
	s.dirFirst = false
}

func (s *Sorter) Len() int {
	return len(s.option)
}

type Option = func(s *Sorter)

func WithSize(size int) Option {
	return func(s *Sorter) {
		s.option = make([]FileSortFunc, 0, size)
	}
}

func WithSortOption(option ...FileSortFunc) Option {
	return func(s *Sorter) {
		s.option = append(s.option, option...)
	}
}

func NewSorter(option ...Option) *Sorter {
	a := Sorter{
		reverse:  false,
		dirFirst: false,
		option:   make([]FileSortFunc, 0, 10),
	}

	for _, opt := range option {
		opt(&a)
	}
	return &a
}

func (s *Sorter) Reverse() {
	s.reverse = !s.reverse
}

func (s *Sorter) AddOption(option ...FileSortFunc) {
	s.option = append(s.option, option...)
}

func (s *Sorter) Build() FileSortFunc {
	return func(a, b *item.FileInfo) int {
		result := 0
		for _, sortFunc := range s.option {
			if s.dirFirst {
				result = dirFirst(a, b)
				if result != 0 {
					break
				}
			}
			result = sortFunc(a, b)
			if result != 0 {
				break
			}
		}

		if s.reverse {
			return -result
		}
		return result
	}
}

func isHidden(info *item.FileInfo) bool {
	return strings.HasPrefix(info.Name(), ".")
}

func isLink(info *item.FileInfo) bool {
	return info.Mode()&os.ModeSymlink != 0
}

func isHiddenDir(info *item.FileInfo) bool {
	return isHidden(info) && info.IsDir()
}

func compareFileInfo(a, b *item.FileInfo) int {
	hdA := isHiddenDir(a)
	hdB := isHiddenDir(b)
	// hidden dir comes first
	if hdA && !hdB {
		return -1
	}
	if !hdA && hdB {
		return 1
	}
	// same hidden dir status
	dA := a.IsDir()
	dB := b.IsDir()
	// dir comes first
	if dA && !dB {
		return -1
	}
	if !dA && dB {
		return 1
	}
	// same dir status
	lA := isLink(a)
	lB := isLink(b)
	switch {
	case lA && lB:
		// both are links, compare name
		// a<b
		return cmp.Compare(a.Name(), b.Name())
	case lA:
		// a is link, b is not link, b comes first
		return 1
	case lB:
		// a is not link, b is link, a comes first
		return -1
	default:
		// neither are links, compare name
		return cmp.Compare(a.Name(), b.Name())
	}
}

func byGroupName(a, b *item.FileInfo, Ascend bool) int {
	if Ascend {
		return cmp.Compare(strings.ToLower(osbased.Group(a)), strings.ToLower(osbased.Group(b)))
	}
	return cmp.Compare(strings.ToLower(osbased.Group(b)), strings.ToLower(osbased.Group(a)))
}

func byUserName(a, b *item.FileInfo, Ascend bool) int {
	if Ascend {
		return cmp.Compare(strings.ToLower(osbased.Owner(a)), strings.ToLower(osbased.Owner(b)))
	}
	return cmp.Compare(strings.ToLower(osbased.Owner(b)), strings.ToLower(osbased.Owner(a)))
}

func byGroupCaseSensitiveName(a, b *item.FileInfo, Ascend bool) int {
	if Ascend {
		return cmp.Compare(osbased.Group(a), osbased.Group(b))
	}
	return cmp.Compare(osbased.Group(b), osbased.Group(a))
}

func byUserCaseSensitiveName(a, b *item.FileInfo, Ascend bool) int {
	if Ascend {
		return cmp.Compare(osbased.Owner(a), osbased.Owner(b))
	}
	return cmp.Compare(osbased.Owner(b), osbased.Owner(a))
}

func byNameWidth(a, b *item.FileInfo, Ascend bool) int {
	if Ascend {
		return len(a.Name()) - len(b.Name())
	}
	return len(b.Name()) - len(a.Name())
}

func ByVersionAscend(a, b *item.FileInfo) int {
	return byVersion(a, b, true)
}

func ByVersionDescend(a, b *item.FileInfo) int {
	return byVersion(a, b, false)
}

// compare version number of two files
// like 1.10.0 > 1.9.0
func byVersion(a, b *item.FileInfo, ascend bool) int {
	av, bv := getVersion(a.Name()), getVersion(b.Name())
	if ascend {
		return av.Compare(bv)
	}
	return bv.Compare(av)
}

type version struct {
	major, minor, patch int
}

func (v *version) Compare(other *version) int {
	if v.major != other.major {
		return cmp.Compare(v.major, other.major)
	}
	if v.minor != other.minor {
		return cmp.Compare(v.minor, other.minor)
	}
	return cmp.Compare(v.patch, other.patch)
}

var (
	v3 = regexp.MustCompile(`\d+\.\d+\.\d+`)
	v2 = regexp.MustCompile(`\d+\.\d+`)
)

// extract version number
// possible formats:
// name_1.0.0.ext
// name-1.0.0.ext
func getVersion(name string) *version {
	v := &version{}
	if len(name) == 0 {
		return v
	}

	isV3 := true
	vs := v3.FindString(name)
	if len(vs) == 0 {
		isV3 = false
		vs = v2.FindString(name)
	}
	if len(vs) == 0 {
		return v
	}

	vsAfterSplit := strings.Split(vs, ".")
	if isV3 {
		v.major, _ = strconv.Atoi(vsAfterSplit[0])
		v.minor, _ = strconv.Atoi(vsAfterSplit[1])
		v.patch, _ = strconv.Atoi(vsAfterSplit[2])
	} else {
		v.major, _ = strconv.Atoi(vsAfterSplit[0])
		v.minor, _ = strconv.Atoi(vsAfterSplit[1])
	}

	return v
}
