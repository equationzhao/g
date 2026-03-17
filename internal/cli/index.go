package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Equationzhao/g/internal/index"
	"github.com/Equationzhao/pathbeautify"
	"github.com/urfave/cli/v2"
)

var indexFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:               "disable-index",
		Aliases:            []string{"di", "no-update"},
		Usage:              "disable updating index",
		Category:           "INDEX",
		DisableDefaultText: true,
	},
	&cli.BoolFlag{
		Name:               "rebuild-index",
		Aliases:            []string{"ri", "remove-all"},
		Usage:              "rebuild index",
		DisableDefaultText: true,
		Category:           "INDEX",
		Action: func(context *cli.Context, b bool) error {
			if b {
				err := index.RebuildIndex()
				if err != nil {
					return err
				}
			}
			return Err4Exit{}
		},
	},
	&cli.BoolFlag{
		Name:               "fuzzy",
		Aliases:            []string{"fz", "f"},
		Usage:              "fuzzy search",
		DisableDefaultText: true,
		Category:           "INDEX",
	},
	&cli.StringSliceFlag{
		Name:     "remove-index",
		Aliases:  []string{"rm"},
		Usage:    "remove paths from index",
		Category: "INDEX",
		Action: func(context *cli.Context, i []string) error {
			var errSum error = nil
			beautification := !context.Bool("np") // --no-path-transform

			for _, s := range i {
				if beautification {
					s = pathbeautify.Transform(s)
				}

				// get absolute path
				r, err := filepath.Abs(s)
				if err != nil {
					errSum = errors.Join(errSum, fmt.Errorf("remove-path: %w", err))
					continue
				}

				err = index.Delete(r)
				if err != nil {
					errSum = errors.Join(errSum, fmt.Errorf("remove-path: %w", err))
				}
			}
			if errSum != nil {
				return errSum
			}
			return Err4Exit{}
		},
	},
	&cli.BoolFlag{
		Name:               "list-index",
		Aliases:            []string{"li"},
		Usage:              "list index",
		DisableDefaultText: true,
		Category:           "INDEX",
		Action: func(context *cli.Context, b bool) error {
			if b {
				keys, err := index.All()
				if err != nil {
					return err
				}
				for i := range keys {
					fmt.Println(keys[i])
				}
			}
			return Err4Exit{}
		},
	},
	&cli.BoolFlag{
		Name:               "remove-current-path",
		Aliases:            []string{"rcp", "rc", "rmc"},
		Usage:              "remove current path from index",
		Category:           "INDEX",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				r, err := os.Getwd()
				if err != nil {
					return err
				}
				err = index.Delete(r)
				if err != nil {
					return err
				}
			}
			return Err4Exit{}
		},
	},
	&cli.BoolFlag{
		Name:               "remove-invalid-path",
		Aliases:            []string{"rip"},
		Usage:              "remove invalid paths from index",
		Category:           "INDEX",
		DisableDefaultText: true,
		Action: func(ctx *cli.Context, b bool) error {
			if b {
				paths, err := index.All()
				if err != nil {
					return err
				}
				invalid := make([]string, 0, len(paths))
				for _, path := range paths {
					_, err := os.Stat(path)
					if err != nil {
						invalid = append(invalid, path)
					}
				}
				err = index.DeleteThose(invalid...)
				if err != nil {
					return err
				}
			}
			return Err4Exit{}
		},
	},
}
