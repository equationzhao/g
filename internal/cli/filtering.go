package cli

import (
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/Equationzhao/g/internal/filter"
	strftime "github.com/itchyny/timefmt-go"
	"github.com/urfave/cli/v2"
)

var filteringFlag = []cli.Flag{
	&cli.UintFlag{
		Name:        "n",
		Aliases:     []string{"limitN", "limit", "topN", "top"},
		Usage:       "limit display to a max of n items (n <=0 means unlimited)",
		Value:       0,
		DefaultText: "unlimited",
		Category:    "FILTERING",
	},
	&cli.StringSliceFlag{
		Name:    "I",
		Aliases: []string{"ignore"},
		Usage:   "ignore Glob patterns",
		Action: func(context *cli.Context, s []string) error {
			if len(s) > 0 {
				f, err := filter.RemoveGlob(s...)
				if err != nil {
					return err
				}
				itemFilterFunc = append(itemFilterFunc, &f)
			}
			return nil
		},
		Category: "FILTERING",
	},
	&cli.StringSliceFlag{
		Name:    "M",
		Aliases: []string{"match"},
		Usage:   "match Glob patterns",
		Action: func(context *cli.Context, s []string) error {
			if len(s) > 0 {
				f, err := filter.GlobOnly(s...)
				if err != nil {
					return err
				}
				itemFilterFunc = append(itemFilterFunc, &f)
			}
			return nil
		},
		Category: "FILTERING",
	},
	&cli.BoolFlag{
		Name:               "show-only-hidden",
		Aliases:            []string{"hidden"},
		DisableDefaultText: true,
		Usage:              "show only hidden files(overridden by --show-hidden/-a/-A)",
		Action: func(context *cli.Context, b bool) error {
			if b {
				itemFilterFunc = slices.DeleteFunc(itemFilterFunc, func(e *filter.ItemFilterFunc) bool {
					return e == &filter.RemoveHidden
				})
				itemFilterFunc = append(itemFilterFunc, &filter.HiddenOnly)
			}
			return nil
		},
		Category: "FILTERING",
	},
	&cli.BoolFlag{
		Name:               "a",
		Aliases:            []string{"sh", "show-hidden"},
		DisableDefaultText: true,
		Usage:              "show hidden files",
		Action: func(context *cli.Context, b bool) error {
			if b {
				// remove filter.RemoveHidden
				itemFilterFunc = slices.DeleteFunc(itemFilterFunc, func(e *filter.ItemFilterFunc) bool {
					return e == &filter.RemoveHidden
				})
			}
			return nil
		},
		Category: "FILTERING",
	},
	&cli.StringSliceFlag{
		Name:  "ext",
		Usage: "show file which has target ext, eg: --ext=go,java",
		Action: func(context *cli.Context, s []string) error {
			if len(s) > 0 {
				f := filter.ExtOnly(s...)
				itemFilterFunc = append(itemFilterFunc, &f)
			}
			return nil
		},
		Category: "FILTERING",
	},
	&cli.StringSliceFlag{
		Name:    "no-ext",
		Aliases: []string{"noext"},
		Usage:   "show file which doesn't have target ext",
		Action: func(context *cli.Context, s []string) error {
			if len(s) > 0 {
				f := filter.RemoveByExt(s...)
				itemFilterFunc = append(itemFilterFunc, &f)
			}
			return nil
		},
		Category: "FILTERING",
	},
	&cli.BoolFlag{
		Name:               "no-dir",
		Aliases:            []string{"nodir", "file"},
		DisableDefaultText: true,
		Usage:              "do not show directory",
		Action: func(context *cli.Context, b bool) error {
			if b {
				itemFilterFunc = append(itemFilterFunc, &filter.RemoveDir)
			}
			return nil
		},
		Category: "FILTERING",
	},
	&cli.BoolFlag{
		Name:               "D",
		Aliases:            []string{"dir", "only-dir"},
		DisableDefaultText: true,
		Usage:              "show directory only",
		Action: func(context *cli.Context, b bool) error {
			if b {
				itemFilterFunc = append(itemFilterFunc, &filter.DirOnly)
			}
			return nil
		},
		Category: "FILTERING",
	},
	&cli.BoolFlag{
		Name:               "B",
		Aliases:            []string{"ignore-backups"},
		DisableDefaultText: true,
		Usage:              "do not list implied entries ending with ~",
		Action: func(context *cli.Context, b bool) error {
			if b {
				itemFilterFunc = append(itemFilterFunc, &filter.RemoveBackups)
			}
			return nil
		},
		Category: "FILTERING",
	},
	&cli.BoolFlag{
		Name:               "A",
		Aliases:            []string{"almost-all"},
		DisableDefaultText: true,
		Usage:              "do not list implied . and ..",
		Action: func(context *cli.Context, b bool) error {
			if b {
				// remove filter.RemoveHidden
				itemFilterFunc = slices.DeleteFunc(itemFilterFunc, func(e *filter.ItemFilterFunc) bool {
					return e == &filter.RemoveHidden
				})
			}
			return nil
		},
		Category: "FILTERING",
	},
	&cli.StringSliceFlag{
		Name:     "only-mime",
		Usage:    "only show file with given mime type",
		Category: "FILTERING",
		Action: func(context *cli.Context, i []string) error {
			if len(i) > 0 {
				setLimit()
				eft := filter.MimeTypeOnly(i...)
				itemFilterFunc = append(itemFilterFunc, &eft)
			}
			return nil
		},
	},
	&cli.StringFlag{
		Name: "before",
		Usage: `show items which was modified/access/created before given time, the time field is determined by --time-type,
	the time will be parsed using format:
		MM-dd, MM-dd HH:mm, HH:mm, YYYY-MM-dd, YYYY-MM-dd HH:mm, and the format set by --time-style`,
		Category: "FILTERING",
		Action: func(ctx *cli.Context, s string) error {
			possibleTimeFormat := []string{"01-02", "01-02 15:04", "15:04", "2006-01-02", "2006-01-02 15:04", timeFormat}
			for _, f := range possibleTimeFormat {
				if after, ok := strings.CutPrefix(f, "+"); ok {
					t, err := strftime.Parse(s, after)
					if err != nil {
						continue
					}
					f := filter.BeforeTime(t, filter.WhichTimeFiled(timeType[0]))
					itemFilterFunc = append(itemFilterFunc, &f)
					return nil
				}
				t, err := time.ParseInLocation(f, s, time.Local)
				if err != nil {
					continue
				}
				if strings.HasPrefix(f, "01-02") {
					t = t.AddDate(time.Now().Year(), 0, 0)
				} else if strings.HasPrefix(f, "15:04") {
					t = combineTodayWithClock(t, time.Now())
				}
				f := filter.BeforeTime(t, filter.WhichTimeFiled(timeType[0]))
				itemFilterFunc = append(itemFilterFunc, &f)
				return nil
			}
			return errors.New("invalid time format")
		},
	},
	&cli.StringFlag{
		Name:     "after",
		Usage:    "show items which was modified/access/created after given time, see --before",
		Category: "FILTERING",
		Action: func(ctx *cli.Context, s string) error {
			possibleTimeFormat := []string{"01-02", "01-02 15:04", "15:04", "2006-01-02", "2006-01-02 15:04", timeFormat}
			for _, f := range possibleTimeFormat {
				if after, ok := strings.CutPrefix(f, "+"); ok {
					t, err := strftime.Parse(s, after)
					if err != nil {
						continue
					}
					f := filter.AfterTime(t, filter.WhichTimeFiled(timeType[0]))
					itemFilterFunc = append(itemFilterFunc, &f)
					return nil
				}
				t, err := time.ParseInLocation(f, s, time.Local)
				if err != nil {
					continue
				} else {
					if strings.HasPrefix(f, "01-02") {
						t = t.AddDate(time.Now().Year(), 0, 0)
					} else if strings.HasPrefix(f, "15:04") {
						t = combineTodayWithClock(t, time.Now())
					}
					f := filter.AfterTime(t, filter.WhichTimeFiled(timeType[0]))
					itemFilterFunc = append(itemFilterFunc, &f)
					return nil
				}
			}
			return errors.New("invalid time format")
		},
	},
}

// combineTodayWithClock keeps the clock from parsed and the calendar day from now.
// Used when the user passes HH:mm without a date.
func combineTodayWithClock(parsed, now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), parsed.Hour(), parsed.Minute(), parsed.Second(), parsed.Nanosecond(), now.Location())
}
