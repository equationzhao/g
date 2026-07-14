//go:build !doc

package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/Equationzhao/g/internal/cli"
	"github.com/Equationzhao/g/internal/config"
	debugSetting "github.com/Equationzhao/g/internal/global/debug"
	"github.com/Equationzhao/g/internal/util"
	ucli "github.com/urfave/cli/v2"
)

func main() {
	// catch panic and print stack trace and version info
	defer func() {
		if !debugSetting.Enable {
			catchPanic(recover())
		}
	}()
	preprocessArgs()
	err := cli.G.Run(os.Args)
	if err != nil {
		if !errors.Is(err, cli.Err4Exit{}) {
			if cli.ReturnCode == 0 {
				cli.ReturnCode = 1
			}
			_, _ = fmt.Fprintln(os.Stderr, cli.MakeErrorStr(err.Error()))
		}
	}
}

func catchPanic(err any) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Version: v%s\n", cli.Version)
		_, _ = fmt.Fprintf(os.Stderr, "Please file an issue at %s with the following panic info\n\n", util.MakeLink("https://github.com/Equationzhao/g/issues/new/choose", "Github Repo"))
		_, _ = fmt.Fprintln(os.Stderr, cli.MakeErrorStr(fmt.Sprintf("error message:\n%v\n", err)))
		_, _ = fmt.Fprintln(os.Stderr, cli.MakeErrorStr(fmt.Sprintf("stack trace:\n%s", debug.Stack())))
		cli.ReturnCode = 1
	}
	os.Exit(cli.ReturnCode)
}

func preprocessArgs() {
	rearrangeArgs()
	argsBeforeDelimiter := os.Args
	if i := slices.Index(os.Args, "--"); i >= 0 {
		argsBeforeDelimiter = os.Args[:i]
	}
	// normal logic
	// load config if the args do not contains -no-config
	if !slices.ContainsFunc(argsBeforeDelimiter, hasNoConfig) {
		defaultArgs, err := config.Load()
		// if successfully load config and **the config.Args do not contain -no-config**
		if err == nil && defaultArgs != nil && !slices.ContainsFunc(defaultArgs.Args, hasNoConfig) {
			os.Args = slices.Insert(os.Args, 1, defaultArgs.Args...)
		} else if err != nil { // if failed to load config
			// if it's read error
			var errReadConfig config.ErrReadConfig
			if errors.As(err, &errReadConfig) {
				_, _ = fmt.Fprintln(os.Stderr, cli.MakeErrorStr(err.Error()))
			}
		}
	} else {
		// contains -no-config
		// remove it before the cli.G starts
		os.Args = append(slices.DeleteFunc(argsBeforeDelimiter, hasNoConfig), os.Args[len(argsBeforeDelimiter):]...)
	}
}

func rearrangeArgs() {
	if len(os.Args) <= 2 {
		return
	}
	flags, paths := separateArgs(os.Args[1:])
	newArgs := append([]string{os.Args[0]}, append(flags, paths...)...)
	os.Args = newArgs
}

func separateArgs(args []string) (flags, paths []string) {
	flagsWithArgs := buildFlagsWithArgsMap()
	expectValue, hasDoubleDash := false, false
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			hasDoubleDash = true
			paths = append(paths, args[i+1:]...)
			break
		}
		if strings.HasPrefix(arg, "--") {
			i = handleLongFlag(arg, args, i, &flags, &expectValue, flagsWithArgs)
		} else if strings.HasPrefix(arg, "-") {
			i = handleShortFlag(arg, args, i, &flags, &expectValue, flagsWithArgs)
		} else {
			if expectValue {
				flags = append(flags, arg)
				expectValue = false
			} else {
				paths = append(paths, arg)
			}
		}
	}
	if hasDoubleDash {
		flags = append(flags, "--")
	}
	return flags, paths
}

func buildFlagsWithArgsMap() map[string]bool {
	flagsWithArgs := make(map[string]bool)
	for _, flag := range cli.G.Flags {
		switch flag.(type) {
		case *ucli.BoolFlag:
			for _, s := range flag.Names() {
				flagsWithArgs[s] = false
			}
		default:
			for _, s := range flag.Names() {
				flagsWithArgs[s] = true
			}
		}
	}
	return flagsWithArgs
}

func handleLongFlag(arg string, args []string, i int, flags *[]string, expectValue *bool, flagsWithArgs map[string]bool) int {
	parts := strings.SplitN(arg, "=", 2)
	flagName := strings.TrimPrefix(parts[0], "--")
	*flags = append(*flags, arg)

	if len(parts) == 2 || !flagsWithArgs[flagName] {
		return i
	}

	if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
		*expectValue = true
	}
	return i
}

func handleShortFlag(arg string, args []string, i int, flags *[]string, expectValue *bool, flagsWithArgs map[string]bool) int {
	parts := strings.SplitN(arg, "=", 2)
	flagName := strings.TrimPrefix(parts[0], "-")
	*flags = append(*flags, arg)
	if len(parts) == 2 || !flagsWithArgs[flagName] {
		return i
	}
	if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
		*expectValue = true
	}
	return i
}

func hasNoConfig(s string) bool {
	if s == "-no-config" || s == "--no-config" {
		return true
	}
	return false
}
