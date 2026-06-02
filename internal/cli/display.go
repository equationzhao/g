package cli

import (
	"fmt"

	"github.com/Equationzhao/g/internal/content"
	"github.com/Equationzhao/g/internal/display"
	"github.com/Equationzhao/g/internal/theme"
	"github.com/urfave/cli/v2"
)

var displayFlag = []cli.Flag{
	// DISPLAY
	&cli.StringFlag{
		Name:     "tree-style",
		Usage:    "set tree style [ascii/unicode(default)/rectangle]",
		Category: "DISPLAY",
		Action: func(context *cli.Context, s string) error {
			_ = context.Set("tree", "1")
			switch s {
			case "ascii", "ASCII", "Ascii":
				display.DefaultTreeStyle = display.TreeASCII
			case "unicode", "Unicode", "UNICODE":
			// no action needed
			case "rectangle", "Rectangle", "RECTANGLE":
				display.DefaultTreeStyle = display.TreeRectangle
			default:
				return fmt.Errorf("invalid tree style: %s", s)
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:               "T",
		Aliases:            []string{"tree"},
		Usage:              "recursively list in tree",
		DisableDefaultText: true,
		Category:           "DISPLAY",
	},
	&cli.StringFlag{
		Name:        "color",
		DefaultText: "auto",
		Usage:       "when to use terminal colors [always|auto|never][basic|256|24bit]",
		Action: func(context *cli.Context, s string) error {
			switch s {
			case "always", "force":
				if theme.ColorLevel == theme.None {
					theme.ColorLevel = theme.Ascii
				}
			case "auto", "tty":
			// skip
			case "never", "none", "off":
				_ = context.Set("no-color", "true")
			case "16", "basic":
				theme.ColorLevel = theme.Ascii
			case "256", "8bit":
				theme.ColorLevel = theme.C256
			case "24bit", "truecolor", "true-color", "24-bit", "16m":
				theme.ColorLevel = theme.TrueColor
			default:
				return fmt.Errorf("unknown color option:%s", s)
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.IntFlag{
		Name:        "depth",
		Aliases:     []string{"level"},
		Usage:       "limit recursive/tree depth, negative -> infinity",
		DefaultText: "infinity",
		Value:       -1,
		Category:    "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "R",
		Aliases:            []string{"recurse"},
		Usage:              "recurse into directories",
		DisableDefaultText: true,
		Category:           "DISPLAY",
		Action: func(context *cli.Context, b bool) error {
			if b {
				if context.Args().Len() > 1 {
					return fmt.Errorf("'--recurse' should not be used with more than one directory")
				}
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:               "byline",
		Aliases:            []string{"1", "oneline", "single-column"},
		Usage:              "print by line",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "zero",
		Aliases:            []string{"0"},
		Usage:              "end each output line with NUL, not newline",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.Zero); !ok {
					p = display.NewZero()
				}
				_ = context.Set("header", "0")
				_ = context.Set("footer", "0")
				_ = context.Set("statistic", "0")
				_ = context.Set("total-size", "0")
				sizeEnabler.DisableTotal()
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "m",
		Aliases:            []string{"comma"},
		Usage:              "fill width with a comma separated list of entries",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.CommaPrint); !ok {
					p = display.NewCommaPrint()
				}
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "x",
		Aliases:            []string{"col", "across", "horizontal"},
		Usage:              "list entries by lines instead of by columns",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.Across); !ok {
					p = display.NewAcross()
				}
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "C",
		Aliases:            []string{"vertical"},
		Usage:              "list entries by columns (default)",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.FitTerminal); !ok {
					p = display.NewFitTerminal()
				}
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "j",
		Aliases:            []string{"json"},
		Usage:              "output in json format",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.JsonPrinter); !ok {
					p = display.NewJsonPrinter()
				}
			}

			_ = context.Set("header", "0")
			theme.SetClassic()
			theme.ColorLevel = theme.None
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.StringFlag{
		Name:    "tb-style",
		Aliases: []string{"table-style"},
		Usage:   "set table style [ascii(default)/unicode]",
		Action: func(context *cli.Context, s string) error {
			_ = context.Set("table", "1")
			switch s {
			case "ascii", "ASCII", "Ascii":
				// no action needed
			case "unicode", "Unicode", "UNICODE":
				display.DefaultTBStyle = display.UNICODEStyle
			default:
				return fmt.Errorf("invalid table style: %s", s)
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "tb",
		Aliases:            []string{"table"},
		Usage:              "output in table format",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.TablePrinter); !ok {
					p = display.NewTablePrinter(display.DefaultTB)
				}
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "md",
		Aliases:            []string{"markdown", "Markdown"},
		Usage:              "output in markdown-table format",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.MDPrinter); !ok {
					p = display.NewMDPrinter()
					_ = context.Set("no-color", "1")
					err := context.Set("header", "1")
					if err != nil {
						return err
					}
				}
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "CSV",
		Aliases:            []string{"csv"},
		Usage:              "output in csv format",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.CSVPrinter); !ok {
					p = display.NewCSVPrinter()
					_ = context.Set("no-color", "1")
					_ = context.Set("no-icon", "1")
				}
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "TSV",
		Aliases:            []string{"tsv"},
		Usage:              "output in tsv format",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.TSVPrinter); !ok {
					p = display.NewTSVPrinter()
					_ = context.Set("no-color", "1")
					_ = context.Set("no-icon", "1")
				}
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.StringFlag{
		Name:        "format",
		DefaultText: "C",
		Usage: `across  -x,  commas  -m, horizontal -x, long -l, single-column -1,
	verbose -l, vertical -C, table -tb, markdown -md, csv -csv, tsv -tsv, json -j, tree -T`,
		Action: func(context *cli.Context, s string) error {
			switch s {
			case "across", "x", "horizontal":
				if _, ok := p.(*display.Across); !ok {
					p = display.NewAcross()
				}
			case "commas", "m":
				if _, ok := p.(*display.CommaPrint); !ok {
					p = display.NewCommaPrint()
				}
			case "long", "l", "verbose":
				contentFunc = append(
					contentFunc, content.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint, r),
				)
				if !context.Bool("O") {
					contentFunc = append(contentFunc, ownerEnabler.EnableOwner(r))
				}
				if !context.Bool("G") {
					contentFunc = append(contentFunc, groupEnabler.EnableGroup(r))
				}
				for _, s := range timeType {
					contentFunc = append(contentFunc, content.EnableTime(timeFormat, s, r))
				}
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			case "single-column", "1":
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			case "vertical", "C":
				if _, ok := p.(*display.FitTerminal); !ok {
					p = display.NewFitTerminal()
				}
			case "table", "tb":
				if _, ok := p.(*display.TablePrinter); !ok {
					p = display.NewTablePrinter(display.DefaultTB)
				}
			case "Markdown", "md", "MD", "markdown":
				if _, ok := p.(*display.MDPrinter); !ok {
					p = display.NewMDPrinter()
					_ = context.Set("no-color", "1")
					err := context.Set("header", "1")
					if err != nil {
						return err
					}
				}
			case "CSV", "csv":
				if _, ok := p.(*display.CSVPrinter); !ok {
					p = display.NewCSVPrinter()
					_ = context.Set("no-color", "1")
					_ = context.Set("no-icon", "1")
				}
			case "TSV", "tsv":
				if _, ok := p.(*display.TSVPrinter); !ok {
					p = display.NewTSVPrinter()
					_ = context.Set("no-color", "1")
					_ = context.Set("no-icon", "1")
				}
			case "json", "j":
				if _, ok := p.(*display.JsonPrinter); !ok {
					p = display.NewJsonPrinter()
					_ = context.Set("classic", "1")
					theme.SetClassic()
					theme.ColorLevel = theme.None
				}
			case "tree", "T":
				if _, ok := p.(*display.TreePrinter); !ok {
					p = display.NewTreePrinter()
				}
			default:
				return fmt.Errorf("unknown format option:%s", s)
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.StringFlag{
		Name:  "theme",
		Usage: "apply theme `path/to/theme`",
		Action: func(context *cli.Context, s string) error {
			err := theme.GetTheme(s)
			if err != nil {
				return err
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "colorless",
		Aliases:            []string{"no-color", "nocolor"},
		Usage:              "without color",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				theme.ColorLevel = theme.None
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "party",
		Aliases:            []string{"disco"},
		Usage:              "rainbow party mode — each file gets a unique color",
		DisableDefaultText: true,
		Category:           "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "classic",
		Usage:              "enable classic mode (no colors or icons)",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				theme.SetClassic()
				theme.ColorLevel = theme.None
				err := context.Set("no-icon", "1")
				if err != nil {
					return err
				}
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "d",
		Aliases:            []string{"directory", "list-dirs"},
		DisableDefaultText: true,
		Usage:              "list directories themselves, not their contents",
		Category:           "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "F",
		Aliases:            []string{"classify"},
		DisableDefaultText: true,
		Usage:              "append indicator (one of */=@|) to entries",
		Category:           "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "ft",
		Aliases:            []string{"file-type"},
		DisableDefaultText: true,
		Usage:              "like classify, except do not append '*'",
		Category:           "DISPLAY",
	},
	&cli.UintFlag{
		Name:        "term-width",
		DefaultText: "auto",
		Usage:       "set screen width",
		Category:    "DISPLAY",
		Action: func(context *cli.Context, u uint) error {
			display.CustomTermSize = u
			return nil
		},
	},
}
