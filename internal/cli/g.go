package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Equationzhao/g/internal/align"
	"github.com/Equationzhao/g/internal/config"
	contents "github.com/Equationzhao/g/internal/content"
	"github.com/Equationzhao/g/internal/display"
	"github.com/Equationzhao/g/internal/filter"
	"github.com/Equationzhao/g/internal/global"
	"github.com/Equationzhao/g/internal/index"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/osbased"
	"github.com/Equationzhao/g/internal/render"
	"github.com/Equationzhao/g/internal/shell"
	"github.com/Equationzhao/g/internal/sorter"
	"github.com/Equationzhao/g/internal/theme"
	"github.com/Equationzhao/g/internal/util"
	"github.com/Equationzhao/pathbeautify"
	"github.com/hako/durafmt"
	"github.com/urfave/cli/v2"
)

var (
	itemFilterFunc = make([]*filter.ItemFilterFunc, 0)
	contentFunc    = make([]contents.ContentOption, 0)
	noOutputFunc   = make([]contents.NoOutputOption, 0)
	r              = render.NewRenderer(&theme.DefaultAll)
	p              = display.NewFitTerminal()
	timeFormat     = "Jan 02 15:04"
	// ReturnCode - Exit status:
	//  0  if OK,
	//  1  if minor problems (e.g., cannot access subdirectory),
	//  2  if serious trouble (e.g., cannot access command-line argument).
	ReturnCode       = 0
	contentFilter    = contents.NewContentFilter()
	sort             = sorter.NewSorter()
	timeType         = []string{"mod"}
	sizeUint         = contents.Auto
	sizeEnabler      = contents.NewSizeEnabler()
	blockEnabler     = contents.NewBlockSizeEnabler()
	ownerEnabler     = contents.NewOwnerEnabler()
	groupEnabler     = contents.NewGroupEnabler()
	gitEnabler       = contents.NewGitEnabler()
	gitRepoEnabler   = contents.NewGitRepoEnabler()
	gitCommitEnabler = contents.NewGitCommitEnabler()
	nameToDisplay    = contents.NewNameEnabler()
	flagsEnabler     = contents.NewFlagsEnabler()
	depthLimitMap    map[string]int
	hookOnce         = util.Once{}
	duplicateDetect  = contents.NewDuplicateDetect()
	hookPost         = make([]func(display.Printer, ...*item.FileInfo), 0)
	allPart          []string
)

var G *cli.App

func init() {
	itemFilterFunc = append(itemFilterFunc, &filter.RemoveHidden)
	G = &cli.App{
		Name:               "g",
		Usage:              "a powerful ls",
		UsageText:          "g [options] [path]",
		Version:            Version,
		SliceFlagSeparator: ",",
		HideHelpCommand:    true,
		Suggest:            true,
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			ReturnCode = 2
			str := err.Error()
			const prefix = "flag provided but not defined: "
			if strings.HasPrefix(str, prefix) {
				suggest := suggestFlag(cCtx.App.Flags, strings.TrimLeft(strings.TrimPrefix(str, prefix), "-"))
				if suggest != "" {
					str = fmt.Sprintf("%s, Did you mean %s?", str, suggest)
				}
			}
			_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(str))
			return nil
		},
		Flags:  make([]cli.Flag, 0, len(viewFlag)+len(filteringFlag)+len(sortingFlags)+len(displayFlag)+len(indexFlags)),
		Action: logic,
	}
	G.Flags = append(
		G.Flags,
		&cli.BoolFlag{
			Name:               "no-path-transform",
			Aliases:            []string{"np"},
			DisableDefaultText: true,
			Usage: `By default, .../a/b/c will be transformed to ../../a/b/c, and ~ will be replaced by homedir, 
	using this flag to disable this feature`,
		},
		&cli.BoolFlag{
			Name:               "duplicate",
			Aliases:            []string{"dup"},
			Usage:              "show duplicate files",
			DisableDefaultText: true,
			Action: func(context *cli.Context, b bool) error {
				if b {
					noOutputFunc = append(noOutputFunc, duplicateDetect.Enable())
					hookPost = append(
						hookPost, func(p display.Printer, item ...*item.FileInfo) {
							duplicateDetect.Fprint(p)
							duplicateDetect.Reset()
						},
					)
				}
				return nil
			},
		},
		&cli.StringFlag{
			Name:  "init",
			Usage: `show the init script for shell, support zsh, bash, fish, powershell, nushell`,
			Action: func(context *cli.Context, s string) error {
				init, err := shell.Init(s)
				if err != nil {
					return err
				}
				_, _ = fmt.Println(init)
				return Err4Exit{}
			},
			Category: "SHELL",
		},
		&cli.BoolFlag{
			Name:               "no-config",
			Usage:              "do not load config file",
			DisableDefaultText: true,
		},
		&cli.BoolFlag{
			Name:               "bug",
			Usage:              "report bug",
			DisableDefaultText: true,
			Action: func(context *cli.Context, b bool) error {
				_, _ = fmt.Println("please report bug to equationzhao@foxmail.com\nor file an issue at https://github.com/Equationzhao/g/issues")
				return Err4Exit{}
			},
		},
	)

	G.Flags = append(G.Flags, viewFlag...)
	G.Flags = append(G.Flags, displayFlag...)
	G.Flags = append(G.Flags, filteringFlag...)
	G.Flags = append(G.Flags, sortingFlags...)
	G.Flags = append(G.Flags, indexFlags...)

	initHelpTemp()

	initVersionHelpFlags()
}

func fuzzyUpdate(path string) error {
	err := index.Update(path)
	if err != nil {
		return err
	}
	return nil
}

// fuzzyPath find the fuzzy path in index
// if error occurs, return empty string and error
func fuzzyPath(path string) (newPath string, minorErr error) {
	fuzzed, err := index.FuzzySearch(path)
	if err == nil {
		return fuzzed, nil
	}
	return "", err
}

// Err4Exit used for exiting without error print
type Err4Exit struct{}

func (c Err4Exit) Error() string {
	panic("it's an error defined to exit app, should not call this")
}

func initHelpTemp() {
	configDir, err := config.GetUserConfigDir()
	if err != nil {
		configDir = filepath.Join("$UserConfigDir", "g")
	}
	cli.AppHelpTemplate = fmt.Sprintf(
		`USAGE:
  g [options] [files...]

VERSION:
	{{.Version}}

CONFIG:
	Configuration: %s
	See More at: g.equationzhao.space

%s
`, filepath.Join(configDir, "g.yaml"), optionsHelp,
	)
}

const optionsHelp = `GLOBAL OPTIONS
   --bug                      report bug
   --duplicate, --dup         show duplicate files
   --no-config                do not load config file
   --no-path-transform, --np  By default, .../a/b/c will be transformed to ../../a/b/c, and ~ will be replaced by homedir,
                              using this flag to disable this feature
   --si                       use powers of 1000 for size format(default: false)

META
   --help, -h, -?       show help
   --version, -v        print the version

DISPLAY
   -#                                 print entry Number for each entry
   --CSV, --csv                       output in csv format
   --TSV, --tsv                       output in tsv format
   --byline, -1, --oneline            print by line
   --classic                          enable classic mode(no colors or icons)
   --color WHEN/LEVEL                 set terminal colors [always|auto|never][basic|256|24bit](default: auto)
   --colorless, --no-color        	  without color
   --depth NUM                        limit recursive/tree depth, negative -> infinity(default: infinity)
   --format FORMAT                    across  -x,  commas  -m, horizontal -x, long -l, single-column -1,
                                      verbose -l, vertical -C, table -tb, markdown -md, csv -csv, tsv -tsv, json -j, tree -T(default: C)

   --file-type                        like --classify, except do not append '*'
   --md, --markdown                   output in markdown-table format
   --tb, --table                      output in table format
   --table-style STYLE                set table style [ascii(default)/unicode]
   --term-width COLS                  set screen width (default: auto)
   --theme path/to/theme              apply theme path/to/theme
   --tree-style STYLE                 set tree style [ascii/unicode(default)/rectangle]
   --zero, -0                         end each output line with NUL, not newline
   -C, --vertical                     list entries by columns(default)
   -F, --classify                     append indicator (one of */=@|) to entries
   -R, --recurse                      recurse into directories
   -T, --tree                         recursively list in tree
   -d, --directory,                   list directories themselves, not their contents
   -j, --json                         output in json format
   -m, --comma                        fill width with a comma separated list of entries
   -x, --col, --across, --horizontal  list entries by lines instead of by columns

FILTERING
   --after TIME                  show items which was modified/access/created after given time, see --before
   --before TIME                 show items which was modified/access/created before given time, the time field is determined by --time-type,
                                 the time will be parsed using format:
                                   MM-dd, MM-dd HH:mm, HH:mm, YYYY-MM-dd, YYYY-MM-dd HH:mm, and the format set by --time-style
   --ext value                   show file which has target ext, eg: --ext=go,java
   --git-ignore                  hide git ignored file/dir [if git is installed]
   --no-dir, --file              do not show directory
   --no-ext value                show file which doesn't have target ext
   --only-mime value             only show file with given mime type
   --show-only-hidden, --hidden  show only hidden files(overridden by --show-hidden/-a/-A)
   -A, --almost-all              do not list implied . and ..
   -B, --ignore-backups          do not list implied entries ending with ~
   -D, --dir, --only-dir         show directory only
   -I GLOBS, --ignore GLOBS      ignore Glob patterns
   -M GLOBS, --match GLOBS       match Glob patterns
   -a, --show-hidden             show hidden files
   -n NUM, --limit NUM           limit display to a max of n items(n <=0 means unlimited)(default: unlimited)

INDEX
   --disable-index, --no-update      disable updating index
   --fuzzy, -f                       fuzzy search
   --list-index,                     list index
   --rebuild-index                   rebuild index
   --remove-current-path             remove current path from index
   --remove-index value, --rm value  remove paths from index
   --remove-invalid-path, --rip      remove invalid paths from index

SHELL
   --init value  show the init script for shell, support zsh, bash, fish, powershell, nushell

SORTING
   --sort SORT_FIELD                       sort by field, default: ascending and case-insensitive,
   available fields:                       nature(default),none(nosort),
                                           name,.name(sorts by name without a leading dot),
                                           size,time,owner,group,extension,inode,width,mime.
                                           append '-descend' to sort descending
                                           field beginning with an Uppercase letter is case-sensitive

   --dir-first, --group-directories-first  list directories before other files
   --sort-by-mime                          sort by mimetype
   --sort-by-mime-descend                  sort by mimetype, descending
   --sort-by-mime-parent                   sort by mimetype parent
   --sort-by-mime-parent-descend           sort by mimetype parent
   --sort-reverse, --reverse, -r           reverse the order of the sort
   --versionsort, --sort-by-version        sort by version numbers, ascending(default: false)
   --width                                 sort by entry name width
   -S, --sort-by-size, --sizesort          sort by file size, largest first(descending)
   -t                                      sort by time, newest first
   -U, --nosort, --no-sort                 do not sort; list entries in directory order.
   -X, --sort-by-ext                       sort alphabetically by entry extension

VIEW
   --access, --ac, --accessed              accessed time
   --all                                   show all info/use a long listing format
   --birth                                 birth time[macOS only]
   --block, --blocks                       show block size
   --charset                               show charset of text file in mime type field
   --checksum, --cs                        show checksum of file with algorithm, see --checksum-algorithm
   --checksum-algorithm value, --ca value  show checksum of file with algorithm:
                                             md5, sha1, sha224, sha256, sha384, sha512, crc32(default: sha1)
   --create, --cr, --created               created time
   --dereference                           dereference symbolic links
   --extended, -@                          list each file's extended attributes and sizes in long listing
   --footer                                add a footer row
   --fp, --full-path, --fullpath           show full path
   --full-time                             like -all/l --time-style=full-iso
   --gid                                   show gid instead of groupname [sid in windows]
   --git, --git-status                     show git status [if git is installed]
   --git-repo-branch, --branch             list root of git-tree branch [if git is installed]
   --git-repo-status, --repo-status        list root of git-tree status [if git is installed]
   --group                                 show group
   --header, --title                       add a header row
   --hyperlink value                       attach hyperlink to filenames [auto|always|never](default: auto)
   --icon, --icons                         show icon
   --inode, -i                             show inode[linux/darwin only]
   --lh, --human-readable                  show human readable size
   --mime, --mime-type, --mimetype         show mime file type
   --mime-parent, --mime-parent-type       show mime parent type
   --modify, --mod, --modified             modified time
   --mounts                                show mount details
   --no-dereference                        do not follow symbolic links
   --no-icon, --noicon, --ni               disable icon(always override --icon)
   --no-total-size                         disable total size(always override --total-size)
   --numeric, --numeric-uid-gid            list numeric user and group IDs instead of name [sid in windows]
   --octal-perm, --octal-permission        list each file's permission in octal format
   --owner, --author                       show owner
   --perm, --permission                    show permission
   --recursive-size                        show recursive size of dir, only work with --size
   --relative-to value                     show relative path to the given path (default: current directory)
   --rt, --relative-time                   show relative time
   --size                                  show file/dir size
   --size-unit value, --block-size value   size unit: bit, b, k, m, g, t, auto
   --smart-group                           only show group if it has a different name from owner
   --statistic                             show statistic info
   --stdin                                 read path from stdin, split by newline
   --time                                  show time
   --time-style TIME_TYPE                  time/date format with -l,
                                           valid TIME_TYPE are :
                                             default, iso, long-iso, full-iso, locale, 
                                             and custom +FORMAT like date(1).
                                             (default: +%%d.%%b'%%y %%H:%%M ,like 02.Jan'06 15:04)
   --time-type value                       time type, mod(default), create, access, all, birth[macOS only]
   --total-size                            show total size
   --uid                                   show uid instead of username [sid in windows]
   -G, --no-group                          in a long listing, don't print group names
   -H, --link                              list each file's number of hard links
   -N, --literal                           print entry names without quoting
   -O, --no-owner                          in a long listing, don't print owner names
   -Q, --quote-name                        enclose entry names in double quotes(overridden by --literal)
   -g                                      like -all, but do not list owner
   -l, --long                              use a long listing format
   -o                                      like -all, but do not list group information`

func initVersionHelpFlags() {
	cli.VersionPrinter = func(cCtx *cli.Context) {
		_, _ = fmt.Fprintf(os.Stdout, `💡 g - a powerful ls
 | Version                %s
 | Go Version             %s
 | Compiler               %s
 | Platform               %s

 | Copyright (C) 2024 Equationzhao. MIT License
 | This is free software: you are free to change and redistribute it.
 | There is NO WARRANTY, to the extent permitted by law.
`, Version, runtime.Version(), runtime.Compiler, runtime.GOOS+"/"+runtime.GOARCH)
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:               "version",
		Aliases:            []string{"v"},
		Usage:              "print the version",
		DisableDefaultText: true,
		Category:           "\b\b\b   META",
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:               "help",
		Aliases:            []string{"h", "?"},
		Usage:              "show help",
		DisableDefaultText: true,
		Category:           "\b\b\b   META",
	}
}

func MakeErrorStr(msg string) string {
	return fmt.Sprintf("%s × %s %s", global.Error, msg, global.Reset)
}

func checkErr(err error, start string) {
	var toPrint string
	if pathErr := new(os.PathError); errors.As(err, &pathErr) {
		if start != "" {
			pathErr.Path = start
		}
		toPrint = fmt.Sprintf("%s: %s (os error %d)", pathErr.Err, pathErr.Path, pathErr.Err)
	} else {
		toPrint = err.Error()
	}
	_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(toPrint))
}

// suggestFlag returns the suggested flag name
// modified from cli.suggestFlag
func suggestFlag(flags []cli.Flag, provided string) string {
	distance := 0.0
	suggestion := ""
	for _, flag := range flags {
		flagNames := flag.Names()
		for _, name := range flagNames {
			newDistance := util.JaroWinkler(name, provided, global.BoostThreshold, global.PrefixSize)
			if newDistance > distance {
				distance = newDistance
				suggestion = name
			}
		}
	}

	// one dash for short flags
	// two dashes for long flags
	if len(suggestion) == 1 {
		suggestion = "-" + suggestion
	} else if len(suggestion) > 1 {
		suggestion = "--" + suggestion
	}

	return suggestion
}

var logic = func(context *cli.Context) error {
	var (
		minorErr   = false
		seriousErr = false
	)

	path := context.Args().Slice()

	if !context.Bool("no-icon") && (context.Bool("icon") || context.Bool("all")) {
		nameToDisplay.SetIcon()
	}
	if context.Bool("F") {
		nameToDisplay.SetClassify()
	}
	if context.Bool("file-type") {
		nameToDisplay.SetClassify()
		nameToDisplay.SetFileType()
	}
	if _, ok := p.(*display.JsonPrinter); ok {
		nameToDisplay.SetJson()
	}
	git := context.Bool("git")
	if git {
		contentFunc = append(contentFunc, gitEnabler.Enable(r))
	}

	gitBranch := context.Bool("git-repo-branch")
	if gitBranch {
		contentFunc = append(contentFunc, gitRepoEnabler.Enable(r))
	}
	gitRepoStatus := context.Bool("git-repo-status")
	if gitRepoStatus {
		contentFunc = append(contentFunc, gitRepoEnabler.EnableStatus(r))
	}

	if context.Bool("git-detail") {
		contentFunc = append(contentFunc, gitCommitEnabler.EnableHash(r), gitCommitEnabler.EnableAuthor(r), gitCommitEnabler.EnableAuthorDateWithTimeFormat(r, timeFormat))
	}

	if context.Bool("flags") {
		contentFunc = append(contentFunc, flagsEnabler.Enable())
	}

	if context.Bool("no-dereference") {
		nameToDisplay.SetNoDeference()
	}

	fuzzy := context.Bool("fuzzy")
	if fuzzy {
		defer func() {
			for i := 0; i < 10; i++ {
				err := index.Close()
				if err != nil {
					continue
				}
				return
			}
		}()
	}

	disableIndex := context.Bool("di")
	wgUpdateIndex := sync.WaitGroup{}

	nameToDisplay.SetQuoteString(`'`)
	// set quote to always
	if context.Bool("Q") {
		nameToDisplay.SetQuote()
	}

	// if no quote, set quote to never
	// this will override the quote set by -Q
	if context.Bool("N") {
		nameToDisplay.UnsetQuote()
	}

	if context.Bool("mounts") {
		nameToDisplay.SetMounts()
	}

	// no path transform
	transformEnabled := !context.Bool("np")
	if rp := context.String("relative-to"); rp != "" {
		if transformEnabled {
			rp = pathbeautify.Beautify(rp)
		}
		if temp, err := filepath.Abs(rp); err == nil {
			rp = temp
		}
		nameToDisplay.SetRelativeTo(rp)
	} else if context.Bool("fp") {
		nameToDisplay.SetFullPath()
	}
	contentFunc = append(contentFunc, nameToDisplay.Enable(r))
	itemFilter := filter.NewItemFilter(itemFilterFunc...)

	gitignore := context.Bool("git-ignore")
	removeGitIgnore := new(filter.ItemFilterFunc)
	if gitignore {
		itemFilter.AppendTo(removeGitIgnore)
	}
	// if no path, use the current path
	if len(path) == 0 {
		path = append(path, ".")
	}

	depth := context.Int("depth")

	// flag: if d is set, display directory them self
	flagd := context.Bool("d")
	// flag: if A is set
	flagA := context.Bool("A")
	flagR := context.Bool("R")
	if flagR {
		depthLimitMap = make(map[string]int)
	}
	header := context.Bool("header")
	footer := context.Bool("footer")
	if context.Bool("statistic") {
		nameToDisplay.SetStatistics(&contents.Statistics{})
	}

	hyperlink := context.String("hyperlink")
	switch hyperlink {
	case "never":
	case "always":
		nameToDisplay.SetHyperlink()
		display.IncludeHyperlink = true
	default:
		fallthrough
	case "auto":
		if util.SupportsHyperlinks() {
			switch p.(type) {
			case display.PrettyPrinter:
			case *display.JsonPrinter:
			default:
				nameToDisplay.SetHyperlink()
				display.IncludeHyperlink = true
			}
		}
	}

	flagSharp := context.Bool("#")
	tree := context.Bool("tree")
	if tree {
		if _, ok := p.(*display.TreePrinter); !ok {
			p = display.NewTreePrinter()
			if flagSharp {
				p.(*display.TreePrinter).NO = true
			}
		}
	}

	smartGroup := context.Bool("smart-group")
	if smartGroup {
		groupEnabler.EnableSmartMode()
	}

	if n := context.Uint("n"); n > 0 && !tree {
		contentFilter.LimitN = n
	}

	longestEachPart := make(map[string]int)
	startDir, _ := os.Getwd()
	dereference := context.Bool("dereference")

	if !context.Bool("colorless") && !context.Bool("classic") && context.String("theme") == "" {
		if config.Default.ThemeLocation != "" {
			err := theme.GetTheme(config.Default.ThemeLocation)
			if err != nil {
				return err
			}
		}
	}

	theme.ConvertThemeColor()

	if len(path) != 0 && context.Bool("stdin") {
		newPath, err := getStdin()
		if err != nil {
			return err
		}
		path = newPath
	}

	// set sort func
	if sort.Len() == 0 {
		sort.AddOption(sorter.Default)
	}
	contentFilter.SetSortFunc(sort.Build())
	contentFilter.SetOptions(contentFunc...)
	contentFilter.SetNoOutputOptions(noOutputFunc...)
	for i := 0; i < len(path); i++ {
		start := time.Now()

		if len(path) > 1 {
			fmt.Println(r.DirPrompt(path[i]), ":")
		}

		if transformEnabled {
			_, err := os.Stat(path[i])
			if err != nil {
				path[i] = pathbeautify.Transform(path[i])
			}
		}
		originPath := path[i]

		infos := make([]*item.FileInfo, 0, 20)

		isFile := false

		// get the abs path
		absPath, err := filepath.Abs(path[i])
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(fmt.Sprintf("failed to get abs path: %s", absPath)))
			continue
		}
		path[i] = absPath

		stat, err := os.Stat(path[i])
		if err != nil {
			base := filepath.Base(path[i])
			if base == "-" {
				path[i] = os.Getenv("OLDPWD")
				stat, err = os.Stat(path[i])
				if err != nil {
					checkErr(err, "")
					seriousErr = true
					continue
				}
			} else if fuzzy { // no match
				// start fuzzy search
				if newPath, err := fuzzyPath(base); err != nil {
					_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
					minorErr = true
				} else {
					path[i] = newPath
					stat, err = os.Stat(path[i])
					fmt.Printf("%s:\n", path[i])
					if err != nil {
						checkErr(err, "")
						seriousErr = true
						continue
					}
				}
			} else {
				// output error
				seriousErr = true
				checkErr(err, originPath)
				continue
			}
		}
		if stat.IsDir() {
			if flagd {
				// when -d is set, treat dir as file
				info, err := item.NewFileInfoWithOption(item.WithFileInfo(stat), item.WithPath(path[i]))
				if err != nil {
					checkErr(err, originPath)
					seriousErr = true
					continue
				}
				infos = append(infos, info)
				isFile = true
			}
		} else {
			info, err := item.NewFileInfoWithOption(item.WithFileInfo(stat), item.WithPath(path[i]))
			if err != nil {
				checkErr(err, originPath)
				seriousErr = true
				continue
			}
			infos = append(infos, info)
			isFile = true
		}

		if !disableIndex {
			wgUpdateIndex.Add(1)
			i := i
			go func() {
				defer wgUpdateIndex.Done()
				if err = fuzzyUpdate(path[i]); err != nil {
					minorErr = true
				}
			}()
		}
		if gitignore {
			*removeGitIgnore = filter.RemoveGitIgnore(path[i])
		}
		if isFile {
			// remove non-display items
			infos = itemFilter.Filter(infos...)

			if tree {
				infos[0].Level = 0
			}
			goto final
		}

		if tree { // visit the dir recursively
			info, err := item.NewFileInfo(path[i])
			if err != nil {
				seriousErr = true
				checkErr(err, originPath)
				continue
			}
			infos = append(
				infos, info,
			)
			infos[0].Level = 0
			if depth >= 1 || depth < 0 {
				infoSlice := util.NewSlice[*item.FileInfo](10)
				errSlice := util.NewSlice[error](10)
				Dive(path[i], 1, depth, infoSlice, errSlice, itemFilter)
				infos = append(infos, *infoSlice.GetRaw()...)
				for _, err := range *errSlice.GetRaw() {
					if err != nil {
						minorErr = true
						checkErr(err, "")
					}
				}
			}
		} else {
			// Use adaptive strategy for optimal performance based on directory characteristics
			strategy := util.NewOptimizationStrategy()
			outputFormat := getOutputFormat(p)                   // We'll need to implement this helper
			hasComplexFiltering := hasComplexFilters(itemFilter) // We'll need this helper too

			processor := strategy.SelectProcessor(path[i], outputFormat, hasComplexFiltering)
			fileInfos, processingErrors := processor.ProcessDirectory(path[i])

			// Handle any processing errors
			for _, err := range processingErrors {
				if err != nil {
					minorErr = true
					checkErr(err, "")
				}
			}

			if !flagA && !tree { // if -A(almost-all) is not set, add the "."/".." info
				err := os.Chdir(path[i])
				if err != nil {
					_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
					seriousErr = true
				} else {
					FileInfoCurrent, err := item.NewFileInfo(".")
					if err != nil {
						seriousErr = true
						checkErr(err, ".")
					} else {
						infos = append(infos, FileInfoCurrent)
					}

					FileInfoParent, err := item.NewFileInfo("..")
					if err != nil {
						minorErr = true
						checkErr(err, "..")
					} else {
						infos = append(infos, FileInfoParent)
					}
				}
			}

			// Add processed files to main list
			infos = append(infos, fileInfos...)

			// remove non-display items
			infos = itemFilter.Filter(infos...)
		}

		// dereference
		if dereference {
			for i := range infos {
				if util.IsSymLink(infos[i]) || osbased.IsMacOSAlias(infos[i].FullPath) {
					symlinks, err := util.Evallinks(infos[i].FullPath)
					if err != nil {
						continue
					}
					info, err := os.Stat(symlinks)
					if err != nil {
						continue
					}
					infos[i].FileInfo = info
					infos[i].FullPath = symlinks
				}
			}
		}

		// if -R is set, add sub dir, insert into path[i+1]
		if flagR && !tree {

			// set depth
			dep, ok := depthLimitMap[path[i]]
			if !ok {
				depthLimitMap[path[i]] = depth
				dep = depth
			}
			if dep >= 2 || dep <= -1 {
				var j int
				for _, info := range infos {
					if info.IsDir() {
						if info.Name() == "." || info.Name() == ".." {
							continue
						}
						abs := filepath.Join(path[i], info.Name())
						newPath, err := filepath.Rel(startDir, abs)
						if err == nil {
							// if the path is relative, use it
							path = slices.Insert(path, i+1+j, newPath)
						} else {
							path = slices.Insert(path, i+1+j, abs)
						}
						j++
						depthLimitMap[abs] = dep - 1
					}
				}
			}
		}

	final:
		if git {
			repo := path[i]
			if isFile {
				repo = filepath.Dir(path[i])
			}
			gitEnabler.Path = repo
			gitEnabler.InitCache(repo)
		}

		contentFilter.GetDisplayItems(&infos)

		if len(infos) == 0 {
			goto clean
		}
		{
			// add total && statistics
			addTotalAndStatistic(nameToDisplay, start)

			// if is table printer, set title
			setTitleForPrettyPrinter(path[i])
			if flagSharp {
				setNumber(infos, tree)
			}

			// do scan
			// get max length for each Meta[key].Value
			if len(allPart) == 0 {
				allPart = infos[0].KeysByOrder()
			}

			// make longestEachPart empty
			for s := range longestEachPart {
				longestEachPart[s] = 0
			}

			if _, ok := p.(*display.JsonPrinter); !ok {
				for _, it := range infos {
					for _, part := range allPart {
						content, ok := it.Get(part)
						if ok && part != contents.NameName {
							l := display.WidthNoHyperLinkLen(content.String())
							if l > longestEachPart[part] {
								longestEachPart[part] = l
							}
						}
					}
				}

				// expand the length of each part using the scan result
				for _, it := range infos {
					for _, part := range allPart {
						if part != contents.NameName {
							content, _ := it.Get(part)
							l := display.WidthNoHyperLinkLen(content.String())
							if l < longestEachPart[part] {
								expand := content.SetPrefix
								if align.IsLeft(part) {
									expand = content.SetSuffix
								}
								// expand
								expand(strings.Repeat(" ", longestEachPart[part]-l))
								it.Set(part, content)
							}
						}
					}
				}
			}

			_ = hookOnce.Do(
				func() error {
					if header || footer {
						headerFooter := display.HeaderMaker{
							AllPart:         allPart,
							LongestEachPart: longestEachPart,
						}
						if header {
							headerFooter.Header = true
							headerFooter.IsBefore = true
							// pre scan
							p.AddBeforePrint(headerFooter.Make)
						}
						if footer {
							headerFooter.Footer = true
							if !header {
								// pre scan
								headerFooter.Header = false
								headerFooter.IsBefore = true
								p.AddBeforePrint(headerFooter.Make)
							}
							headerFooter.IsBefore = false
							p.AddAfterPrint(headerFooter.Make)
						}
					}
					p.AddAfterPrint(hookPost...)
					return nil
				},
			)
			p.Print(infos...)
		}

	clean:
		if i != len(path)-1 {
			fmt.Print("\n\n")
			// switch back to start dir
			if err = os.Chdir(startDir); err != nil {
				seriousErr = true
			}
			sizeEnabler.Reset()
		}
	}
	wgUpdateIndex.Wait()

	if seriousErr {
		ReturnCode = 2
	} else if minorErr {
		ReturnCode = 1
	}

	return nil
}

func getStdin() ([]string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var args []string
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			args = append(args, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return args, nil
}

func setNumber(infos []*item.FileInfo, isTree bool) {
	l := len(strconv.Itoa(len(infos)))
	for i, info := range infos {
		// if there is #, add No
		var no *display.ItemContent
		if !isTree {
			no = &display.ItemContent{
				No:      -1,
				Content: display.StringContent(strconv.Itoa(i)),
			}
			no.SetSuffix(strings.Repeat(" ", l-len(strconv.Itoa(i))))
		} else {
			no = &display.ItemContent{
				No:      -1,
				Content: display.StringContent(""),
			}
			no.SetSuffix(strings.Repeat(" ", l))
		}
		info.Set("#", no)
	}
}

func setTitleForPrettyPrinter(path string) {
	prettyPrinter, isPrettyPrinter := p.(display.PrettyPrinter)
	if isPrettyPrinter {
		switch p.(type) {
		case *display.CSVPrinter, *display.TSVPrinter:
			break
		default:
			prettyPrinter.SetTitle(path)
		}
	}
}

func addTotalAndStatistic(nameToDisplay *contents.Name, start time.Time) {
	{

		// if -l/show-total-size is set, add total size
		jp, isJsonPrinter := p.(*display.JsonPrinter)

		if isJsonPrinter {
			jp.Extra = make([]any, 0, 2)
		}

		if total, ok := sizeEnabler.Total(); ok {
			s, unit := sizeEnabler.Size2String(total)
			s = r.Size(s, contents.Convert2SizeString(unit))

			if isJsonPrinter {
				jp.Extra = append(
					jp.Extra, struct {
						Total string `json:"total"`
					}{
						Total: s,
					},
				)
			} else {
				_, _ = display.RawPrint(fmt.Sprintf("  total %s\n", s))
			}
		}
		if s := nameToDisplay.Statistics(); s != nil {
			t := r.Time(durafmt.Parse(time.Since(start)).LimitToUnit("ms").String())
			if isJsonPrinter {
				jp.Extra = append(
					jp.Extra, struct {
						Time      string               `json:"underwent"`
						Statistic *contents.Statistics `json:"statistic"`
					}{
						Time:      t,
						Statistic: s,
					},
				)
			} else {
				_, _ = display.RawPrint(
					fmt.Sprintf(
						"  underwent %s", t,
					),
				)
				_, _ = display.RawPrint(fmt.Sprintf("\n  statistic: %s\n", s))
			}
			s.Reset()
		}
	}
}
