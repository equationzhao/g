package request

import "fmt"

type When uint8

const (
	WhenAuto When = iota
	WhenAlways
	WhenNever
)

func (w When) String() string {
	switch w {
	case WhenAlways:
		return "always"
	case WhenNever:
		return "never"
	default:
		return "auto"
	}
}

type Format uint8

const (
	FormatGrid Format = iota
	FormatAcross
	FormatOneline
	FormatComma
	FormatTree
	FormatJSON
)

func (f Format) String() string {
	switch f {
	case FormatAcross:
		return "across"
	case FormatOneline:
		return "oneline"
	case FormatComma:
		return "comma"
	case FormatTree:
		return "tree"
	case FormatJSON:
		return "json"
	default:
		return "grid"
	}
}

type SortKey uint8

const (
	SortName SortKey = iota
	SortSize
	SortTime
	SortExt
	SortVersion
	SortNone
)

type TimeField uint8

const (
	TimeModified TimeField = iota
	TimeAccessed
	TimeChanged
	TimeBirth
)

type Visibility uint8

const (
	VisHidden Visibility = iota
	VisAlmostAll
	VisAll
)

type KindFilter uint8

const (
	KindAll KindFilter = iota
	KindDirsOnly
	KindFilesOnly
)

type DirOrder uint8

const (
	DirOrderNone DirOrder = iota
	DirOrderFirst
	DirOrderLast
)

type IconSet uint8

const (
	IconNerd IconSet = iota
	IconUnicode
)

type ConfigSrc uint8

const (
	ConfigSearch ConfigSrc = iota
	ConfigNone
	ConfigPath
)

type QuoteMode uint8

const (
	QuoteDefault QuoteMode = iota
	QuoteAlways
	QuoteLiteral
)

type TimeStyle struct {
	Named  string
	Custom string
}

// Request is the immutable-after-parse listing configuration.
type Request struct {
	Paths []string

	Format      Format
	FormatSet   bool
	Long        bool
	LongNoOwner bool
	LongNoGroup bool

	Visibility Visibility
	DirSelf    bool
	Recurse    bool
	Depth      int
	HasDepth   bool

	Ignore        []string
	KindFilter    KindFilter
	IgnoreBackups bool
	GitIgnore     bool

	Sort     SortKey
	Reverse  bool
	DirOrder DirOrder

	HumanReadable bool
	SI            bool
	Inode         bool
	Links         bool
	NumericIDs    bool
	Blocks        bool
	Header        bool

	TimeField TimeField
	TimeStyle TimeStyle

	Color       When
	Icons       When
	Hyperlink   When
	ThemePath   string
	IconSet     IconSet
	Classify    When
	ClassifySet bool
	Quote       QuoteMode
	Zero        bool
	Width       int
	WidthSet    bool

	Git         bool
	Dereference bool

	Config     ConfigSrc
	ConfigPath string
	Help       bool
	Version    bool
}

type Runtime struct {
	StdoutTTY  bool
	WidthIOCTL int
	COLUMNS    int
	NOCOLOR    bool
}

func Default() Request {
	return Request{
		TimeStyle: TimeStyle{Named: "default"},
	}
}

func (r Request) Validate() error {
	if r.WidthSet && r.Width <= 0 {
		return fmt.Errorf("--width must be a positive integer")
	}
	if r.HasDepth && r.Depth < 0 {
		return fmt.Errorf("--depth must be >= 0")
	}
	if r.Zero && r.Format != FormatOneline {
		return fmt.Errorf("-0 requires --format=oneline")
	}
	if r.Zero && (r.Color == WhenAlways || r.Icons == WhenAlways || r.Hyperlink == WhenAlways || r.Classify == WhenAlways || r.Quote == QuoteAlways) {
		return fmt.Errorf("-0 cannot be combined with always color/icons/hyperlink/classify or -Q")
	}
	if r.Format == FormatJSON && r.Hyperlink == WhenAlways {
		return fmt.Errorf("--format=json cannot be combined with --hyperlink=always")
	}
	if r.Config == ConfigPath && r.ConfigPath == "" {
		return fmt.Errorf("--config= requires a path")
	}
	return nil
}

func Resolve(r Request, rt Runtime) Request {
	if !r.FormatSet {
		if rt.StdoutTTY {
			r.Format = FormatGrid
		} else {
			r.Format = FormatOneline
		}
	}
	if r.Long && (r.Format == FormatGrid || r.Format == FormatAcross || r.Format == FormatComma) {
		r.Format = FormatOneline
	}
	if r.DirSelf {
		r.Recurse = false
		r.HasDepth = false
	}
	if !r.WidthSet {
		if rt.COLUMNS > 0 {
			r.Width = rt.COLUMNS
		} else if rt.WidthIOCTL > 0 {
			r.Width = rt.WidthIOCTL
		} else {
			r.Width = 80
		}
	}
	r.Color = resolveWhen(r.Color, rt.StdoutTTY, rt.NOCOLOR)
	if r.Color == WhenNever {
		r.Icons = WhenNever
	} else {
		r.Icons = resolveWhen(r.Icons, rt.StdoutTTY, false)
	}
	if r.Format == FormatJSON {
		r.Hyperlink = WhenNever
		r.Classify = WhenNever
	} else {
		r.Hyperlink = resolveWhen(r.Hyperlink, rt.StdoutTTY, false)
		if !r.ClassifySet {
			r.Classify = WhenNever
		} else if r.Classify == WhenAuto {
			if rt.StdoutTTY {
				r.Classify = WhenAlways
			} else {
				r.Classify = WhenNever
			}
		}
	}
	if r.Zero {
		r.Header = false
	}
	return r
}

func resolveWhen(w When, tty, nocolor bool) When {
	if w == WhenAlways || w == WhenNever {
		return w
	}
	if nocolor {
		return WhenNever
	}
	if tty {
		return WhenAlways
	}
	return WhenNever
}
