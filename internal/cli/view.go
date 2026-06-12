package cli

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	contents "github.com/Equationzhao/g/internal/content"
	"github.com/Equationzhao/g/internal/display"
	"github.com/Equationzhao/g/internal/filter"
	"github.com/gabriel-vasile/mimetype"
	"github.com/urfave/cli/v2"
)

var viewFlag = []cli.Flag{
	// VIEW
	&cli.BoolFlag{
		Name:               "header",
		Aliases:            []string{"title"},
		Usage:              "add a header row",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "footer",
		Usage:              "add a footer row",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "statistic",
		Usage:              "show statistic info",
		DisableDefaultText: true,
		Category:           "VIEW",
	},
		&cli.StringSliceFlag{
			Name:    "time-type",
			Usage:   "time type, mod(default), create, access, all, birth[macOS only]",
			EnvVars: []string{"TIME_TYPE"},
			Action: func(context *cli.Context, ss []string) error {
				_ = context.Set("time", "1")
				timeType = make([]string, 0, len(ss))
				for _, s := range ss {
					normalized, ok := normalizeTimeType(s)
					if !ok {
						ReturnCode = 2
						return errors.New("invalid time type")
					}
					if normalized == "" {
						timeType = []string{"mod", "create", "access"}
						continue
					}
					timeType = append(timeType, normalized)
				}
				return nil
			},
			Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "access",
		Aliases:            []string{"ac", "accessed"},
		Usage:              "accessed time",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				_ = context.Set("time", "1")
				timeType = append(timeType, "access")
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "modify",
		Aliases:            []string{"mod", "modified"},
		Usage:              "modified time",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				_ = context.Set("time", "1")
				timeType = append(timeType, "mod")
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "create",
		Aliases:            []string{"cr", "created"},
		Usage:              "created time",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				_ = context.Set("time", "1")
				timeType = append(timeType, "create")
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "birth",
		Usage:              "birth time[macOS only]",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if runtime.GOOS != "darwin" {
				return errors.New("birth is only supported in darwin")
			}
			if b {
				_ = context.Set("time", "1")
				timeType = append(timeType, "birth")
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name: "si",
		Usage: `use powers of 1000 not 1024 for size format
		eg: 1K = 1000 bytes`,
		EnvVars: []string{"SI"},
		Action: func(context *cli.Context, b bool) error {
			if b {
				sizeEnabler.SetSI()
			}
			return nil
		},
	},
	&cli.StringFlag{
		Name:    "size-unit",
		Aliases: []string{"su", "block-size"},
		Usage: `size unit:
			bit, b, k, m, g, t, auto`,
		Action: func(context *cli.Context, s string) error {
			_ = context.Set("size", "1")
			if strings.EqualFold(s, "auto") {
				return nil
			}
			si := context.Bool("si")
			sizeUint = contents.ConvertFromSizeString(s, si)
			if sizeUint == contents.Unknown {
				ReturnCode = 2
				return fmt.Errorf("invalid size unit: %s", s)
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.StringFlag{
		Name: "time-style",
		Usage: `time/date format with -l, 
	valid timestamp styles are default, iso, long-iso, full-iso, locale, 
	custom +FORMAT like date(1). 
	(default: +%b %d %H:%M ,like Jan 02 15:04)`,
		EnvVars: []string{"TIME_STYLE"},
		Action: func(context *cli.Context, s string) error {
			_ = context.Set("time", "1")
			/*
				The TIME_STYLE argument can be full-iso, long-iso, iso, locale, or  +FORMAT.
				FORMAT is interpreted like in date(1).
				Also, the TIME_STYLE environment variable sets the default style to use.
			*/
			if strings.HasPrefix(s, "+") {
				timeFormat = s
				return nil
			}

			switch s {
			case "full-iso":
				timeFormat = "2006-01-02 15:04:05.000000000 -0700"
			case "long-iso":
				timeFormat = "2006-01-02 15:04"
			case "locale":
				timeFormat = "Jan 02 15:04"
			case "iso":
				timeFormat = "01-02 15:04"
			case "default":
				timeFormat = "Jan 02 15:04"
			default:
				ReturnCode = 2
				return errors.New("invalid time-style")
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "full-time",
		Usage:              "like -all/l --time-style=full-iso",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				_ = context.Set("time", "1")
				timeFormat = "2006-01-02 15:04:05.000000000 -0700"
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "#",
		DisableDefaultText: true,
		Usage:              "print entry Number for each entry",
		Category:           "DISPLAY",
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(contentFunc, contents.NewIndexEnabler().Enable())
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:               "inode",
		Aliases:            []string{"i"},
		Usage:              "show inode[linux/darwin only]",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			i := contents.NewInodeEnabler()
			contentFunc = append(contentFunc, i.Enable(r))
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "uid",
		Usage:              "show uid instead of username [sid in windows]",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				ownerEnabler.EnableNumeric()
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "gid",
		Usage:              "show gid instead of groupname [sid in windows]",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				groupEnabler.EnableNumeric()
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "numeric",
		Aliases:            []string{"numeric-uid-gid"},
		Usage:              "list numeric user and group IDs instead of name [sid in windows]",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				ownerEnabler.EnableNumeric()
				groupEnabler.EnableNumeric()
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "octal-perm",
		Aliases:            []string{"octal-permission"},
		Usage:              "list each file's permission in octal format",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(contentFunc, contents.EnableFileOctalPermissions(r))
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "perm",
		Aliases:            []string{"permission"},
		Usage:              "show permission",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(contentFunc, contents.EnableFileMode(r))
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "size",
		Usage:              "show file/dir size",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(contentFunc, sizeEnabler.EnableSize(sizeUint, r))
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "recursive-size",
		Usage:              "show recursive size of dir, only work with --size",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				_ = context.Set("size", "1")
				n := context.Int("depth")
				sizeEnabler.SetRecursive(contents.NewSizeRecursive(n))
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "block",
		Aliases:            []string{"blocks"},
		Usage:              "show block size",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(contentFunc, blockEnabler.Enable(r))
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "lh",
		Aliases:            []string{"human-readable"},
		DisableDefaultText: true,
		Usage:              "show human readable size",
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(
					contentFunc, contents.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint, r),
					ownerEnabler.EnableOwner(r), groupEnabler.EnableGroup(r),
				)
				for _, s := range timeType {
					contentFunc = append(contentFunc, contents.EnableTime(timeFormat, s, r))
				}
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "H",
		Aliases:            []string{"link"},
		Usage:              "list each file's number of hard links",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				link := contents.NewLinkEnabler()
				contentFunc = append(contentFunc, link.Enable(r))
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "owner",
		Aliases:            []string{"author"},
		Usage:              "show owner",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(contentFunc, ownerEnabler.EnableOwner(r))
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "group",
		Usage:              "show group",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(contentFunc, groupEnabler.EnableGroup(r))
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "smart-group",
		Usage:              "only show group if it has a different name from owner",
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "time",
		Usage:              "show time",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				for _, s := range timeType {
					contentFunc = append(contentFunc, contents.EnableTime(timeFormat, s, r))
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "no-icon",
		Usage:              "disable icon(always override --icon)",
		Aliases:            []string{"noicon", "ni"},
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "icon",
		Usage:              "show icon",
		Aliases:            []string{"icons"},
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "fp",
		Usage:              "show full path",
		Aliases:            []string{"full-path", "fullpath"},
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.StringFlag{
		Name:        "relative-to",
		Usage:       "show relative path to the given path",
		DefaultText: "current directory",
		Category:    "VIEW",
	},
	&cli.BoolFlag{
		Name:               "total-size",
		Usage:              "show total size",
		DisableDefaultText: true,
		Category:           "VIEW",
		Action: func(context *cli.Context, b bool) error {
			if b {
				_ = context.Set("size", "1")
				sizeEnabler.SetEnableTotal()
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:               "mime",
		Usage:              "show mime file type",
		Aliases:            []string{"mime-type", "mimetype"},
		DisableDefaultText: true,
		Category:           "VIEW",
		Action: func(context *cli.Context, b bool) error {
			if b {
				exact := contents.NewMimeFileTypeEnabler()
				setLimit()
				contentFunc = append(contentFunc, exact.Enable(r))
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:               "mime-parent",
		Usage:              "show mime parent type",
		Aliases:            []string{"mime-parent-type", "mimetype-parent"},
		Category:           "VIEW",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				exact := contents.NewMimeFileTypeEnabler()
				exact.ParentOnly = true
				setLimit()
				contentFunc = append(contentFunc, exact.Enable(r))
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:               "charset",
		Usage:              "show charset of text file in mime type field",
		DisableDefaultText: true,
		Category:           "VIEW",
		Action: func(context *cli.Context, b bool) error {
			if b {
				charset := contents.NewCharsetEnabler()
				setLimit()
				contentFunc = append(contentFunc, charset.Enable(r))
			}
			return nil
		},
	},
	&cli.StringSliceFlag{
		Name: "checksum-algorithm",
		Usage: `show checksum of file with algorithm: 
	md5, sha1, sha224, sha256, sha384, sha512, crc32`,
		Aliases:     []string{"ca"},
		DefaultText: "sha1",
		Value:       cli.NewStringSlice("sha1"),
		Category:    "VIEW",
		Action: func(context *cli.Context, i []string) error {
			_ = context.Set("checksum", "1")
			return nil
		},
	},
	&cli.BoolFlag{
		Name:               "checksum",
		Usage:              `show checksum of file with algorithm, see --checksum-algorithm`,
		Aliases:            []string{"cs"},
		DisableDefaultText: true,
		Category:           "VIEW",
		Action: func(context *cli.Context, b bool) error {
			ss := context.StringSlice("checksum-algorithm")
			if ss == nil {
				ss = []string{"sha1"}
			}
			sums := make([]contents.SumType, 0, len(ss))
			for _, s := range ss {
				switch s {
				case "md5":
					sums = append(sums, contents.SumTypeMd5)
				case "sha1":
					sums = append(sums, contents.SumTypeSha1)
				case "sha224":
					sums = append(sums, contents.SumTypeSha224)
				case "sha256":
					sums = append(sums, contents.SumTypeSha256)
				case "sha384":
					sums = append(sums, contents.SumTypeSha384)
				case "sha512":
					sums = append(sums, contents.SumTypeSha512)
				case "crc32":
					sums = append(sums, contents.SumTypeCRC32)
				default:
					return fmt.Errorf("invalid checksum algorithm: %s", s)
				}
			}

			if b {
				contentFunc = append(contentFunc, contents.SumEnabler{}.EnableSum(r, sums...)...)
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:               "git",
		Usage:              "show git status [if git is installed]",
		Aliases:            []string{"git-status"},
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "git-detail",
		Usage:              "show git commit detail with hash, author, author date [if git is installed]",
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "git-repo-branch",
		Usage:              "list root of git-tree branch [if git is installed]",
		Aliases:            []string{"branch"},
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "git-repo-status",
		Usage:              "list root of git-tree status [if git is installed]",
		Aliases:            []string{"repo-status"},
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "Q",
		Aliases:            []string{"quote-name"},
		Usage:              "enclose entry names in double quotes(overridden by --literal)",
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "mounts",
		Usage:              "show mount details",
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "N",
		Aliases:            []string{"literal"},
		Usage:              "print entry names without quoting",
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "no-dereference",
		Usage:              "do not follow symbolic links",
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "dereference",
		Usage:              "dereference symbolic links",
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.StringFlag{
		Name:        "hyperlink",
		Usage:       "attach hyperlink to filenames [auto|always|never]",
		Category:    "VIEW",
		DefaultText: "auto",
		Action: func(context *cli.Context, s string) error {
			if strings.EqualFold(s, "auto") {
				_ = context.Set("hyperlink", "auto")
			} else if strings.EqualFold(s, "always") {
				_ = context.Set("hyperlink", "always")
			} else if strings.EqualFold(s, "never") {
				_ = context.Set("hyperlink", "never")
			} else {
				return fmt.Errorf("invalid hyperlink value: %s", s)
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:               "o",
		DisableDefaultText: true,
		Usage:              "like -all, but do not list group information",
		Action: func(context *cli.Context, b bool) error {
			if b {
				// remove filter.RemoveHidden
				newFF := make([]*filter.ItemFilterFunc, 0, len(itemFilterFunc))
				for _, typeFunc := range itemFilterFunc {
					if typeFunc != &filter.RemoveHidden {
						newFF = append(newFF, typeFunc)
					}
				}
				itemFilterFunc = newFF
				contentFunc = append(
					contentFunc, contents.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint, r),
					ownerEnabler.EnableOwner(r),
				)
				for _, s := range timeType {
					contentFunc = append(contentFunc, contents.EnableTime(timeFormat, s, r))
				}
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "g",
		DisableDefaultText: true,
		Usage:              "like -all, but do not list owner",
		Action: func(context *cli.Context, b bool) error {
			if b {
				// remove filter.RemoveHidden
				newFF := make([]*filter.ItemFilterFunc, 0, len(itemFilterFunc))
				for _, typeFunc := range itemFilterFunc {
					if typeFunc != &filter.RemoveHidden {
						newFF = append(newFF, typeFunc)
					}
				}
				itemFilterFunc = newFF
				contentFunc = append(
					contentFunc, contents.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint, r),
					groupEnabler.EnableGroup(r),
				)
				for _, s := range timeType {
					contentFunc = append(contentFunc, contents.EnableTime(timeFormat, s, r))
				}
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "G",
		DisableDefaultText: true,
		Aliases:            []string{"no-group"},
		Usage:              "in a long listing, don't print group names",
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "O",
		DisableDefaultText: true,
		Aliases:            []string{"no-owner"},
		Usage:              "in a long listing, don't print owner names",
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "l",
		Aliases:            []string{"long"},
		Usage:              "use a long listing format",
		Category:           "VIEW",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(
					contentFunc, contents.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint, r),
				)
				if !context.Bool("O") {
					contentFunc = append(contentFunc, ownerEnabler.EnableOwner(r))
				}
				if !context.Bool("G") {
					contentFunc = append(contentFunc, groupEnabler.EnableGroup(r))
				}
				for _, s := range timeType {
					contentFunc = append(contentFunc, contents.EnableTime(timeFormat, s, r))
				}
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:               "all",
		Aliases:            []string{"la"},
		Usage:              "show all info/use a long listing format",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				// remove filter.RemoveHidden
				newFF := make([]*filter.ItemFilterFunc, 0, len(itemFilterFunc))
				for _, typeFunc := range itemFilterFunc {
					if typeFunc != &filter.RemoveHidden {
						newFF = append(newFF, typeFunc)
					}
				}
				itemFilterFunc = newFF
				contentFunc = append(
					contentFunc, contents.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint, r),
				)
				if !context.Bool("O") {
					contentFunc = append(contentFunc, ownerEnabler.EnableOwner(r))
				}
				if !context.Bool("G") {
					contentFunc = append(contentFunc, groupEnabler.EnableGroup(r))
				}
				for _, s := range timeType {
					contentFunc = append(contentFunc, contents.EnableTime(timeFormat, s, r))
				}
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "no-total-size",
		Usage:              "disable total size(always override --total-size)",
		DisableDefaultText: true,
		Category:           "VIEW",
		Action: func(context *cli.Context, b bool) error {
			if b {
				sizeEnabler.DisableTotal()
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:               "rt",
		Aliases:            []string{"relative-time"},
		Usage:              "show relative time",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				for _, s := range timeType {
					rt := contents.NewRelativeTimeEnabler()
					rt.Mode = s
					contentFunc = append(contentFunc, rt.Enable(r))
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "stdin",
		Usage:              "read path from stdin, split by newline",
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "flags",
		Usage:              "list file flags[macOS only]",
		DisableDefaultText: true,
		Category:           "VIEW",
	},
}

func normalizeTimeType(s string) (string, bool) {
	switch strings.ToLower(s) {
	case "mod", "modified":
		return "mod", true
	case "create", "cr":
		return "create", true
	case "access", "ac":
		return "access", true
	case "birth":
		return "birth", true
	case "all":
		return "", true
	default:
		return "", false
	}
}

func setLimit() {
	mimetype.SetLimit(1024 * 1024)
}
