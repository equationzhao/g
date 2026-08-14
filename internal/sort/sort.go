package sort

import (
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/Equationzhao/g/internal/entry"
	"github.com/Equationzhao/g/internal/request"
)

func Apply(ents []entry.Entry, req request.Request) {
	sort.SliceStable(ents, func(i, j int) bool {
		a, b := ents[i], ents[j]
		if req.DirOrder != request.DirOrderNone {
			ad, bd := a.IsDir(), b.IsDir()
			if ad != bd {
				if req.DirOrder == request.DirOrderFirst {
					return ad
				}
				return bd
			}
		}
		cmp := compare(a, b, req)
		if req.Reverse {
			return cmp > 0
		}
		return cmp < 0
	})
}

func compare(a, b entry.Entry, req request.Request) int {
	switch req.Sort {
	case request.SortSize:
		if a.Size != b.Size {
			if a.Size > b.Size {
				return -1
			}
			return 1
		}
	case request.SortTime:
		ta, tb := timeOf(a, req), timeOf(b, req)
		if !ta.Equal(tb) {
			if ta.After(tb) {
				return -1
			}
			return 1
		}
	case request.SortExt:
		ea, eb := extOf(a.Name), extOf(b.Name)
		if c := strings.Compare(strings.ToLower(ea), strings.ToLower(eb)); c != 0 {
			return c
		}
	case request.SortVersion:
		if c := Filevercmp(a.Name, b.Name); c != 0 {
			return c
		}
	case request.SortNone:
		if a.ReadOrder != b.ReadOrder {
			if a.ReadOrder < b.ReadOrder {
				return -1
			}
			return 1
		}
		return 0
	}
	return nameCmp(a.Name, b.Name)
}

func timeOf(e entry.Entry, req request.Request) time.Time {
	switch req.TimeField {
	case request.TimeAccessed:
		return e.AccTime
	case request.TimeChanged:
		return e.ChangeTime
	case request.TimeBirth:
		if e.HasBirth {
			return e.Birth
		}
		return time.Time{}
	default:
		return e.ModTime
	}
}

func extOf(name string) string {
	i := strings.LastIndex(name, ".")
	if i <= 0 || i == len(name)-1 {
		return ""
	}
	return name[i+1:]
}

func nameCmp(a, b string) int {
	return strings.Compare(strings.ToLower(a), strings.ToLower(b))
}

// Filevercmp implements GNU coreutils filevercmp enough for the spec table.
func Filevercmp(a, b string) int {
	if a == b {
		return 0
	}
	if a == "" {
		return -1
	}
	if b == "" {
		return 1
	}
	ia, ib := 0, 0
	for ia < len(a) && ib < len(b) {
		ca, cb := a[ia], b[ib]
		da, db := isDigit(ca), isDigit(cb)
		if da && db {
			// skip leading zeros; longer zero prefix is smaller when values equal
			za, zb := 0, 0
			for ia < len(a) && a[ia] == '0' {
				za++
				ia++
			}
			for ib < len(b) && b[ib] == '0' {
				zb++
				ib++
			}
			sa, sb := ia, ib
			for ia < len(a) && isDigit(a[ia]) {
				ia++
			}
			for ib < len(b) && isDigit(b[ib]) {
				ib++
			}
			na, nb := a[sa:ia], b[sb:ib]
			if len(na) != len(nb) {
				if len(na) < len(nb) {
					return -1
				}
				return 1
			}
			if c := strings.Compare(na, nb); c != 0 {
				return c
			}
			if za != zb {
				if za > zb {
					return -1
				}
				return 1
			}
			continue
		}
		if da != db {
			// digit vs non-digit: file2 < file10 already handled; a1 < a1a means after equal prefix, missing < extra
			if da {
				return -1
			}
			return 1
		}
		// both non-digit: compare as unsigned bytes for non-UTF8, else runes
		ra, sa := decode(a, ia)
		rb, sb := decode(b, ib)
		if ra != rb {
			if ra < rb {
				return -1
			}
			return 1
		}
		ia += sa
		ib += sb
	}
	switch {
	case ia == len(a) && ib == len(b):
		return 0
	case ia == len(a):
		return -1
	default:
		return 1
	}
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

func decode(s string, i int) (rune, int) {
	if !utf8.ValidString(s[i:]) {
		return rune(s[i]), 1
	}
	r, n := utf8.DecodeRuneInString(s[i:])
	if r == utf8.RuneError && n == 1 {
		return rune(s[i]), 1
	}
	_ = unicode.ReplacementChar
	return r, n
}
