package main

import (
	"io/fs"
	"os"
	"reflect"
	"testing"

	"github.com/Equationzhao/g/internal/cli"
	"github.com/Equationzhao/g/internal/config"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ucli "github.com/urfave/cli/v2"
)

func Test_catchPanic(t *testing.T) {
	gomonkey.ApplyFunc(os.Exit, func(int) {})
	tests := []struct {
		name string
		err  any
	}{
		{
			name: "empty",
			err:  nil,
		},
		{
			name: "ErrExist",
			err:  fs.ErrExist,
		},
		{
			name: "ErrNotExist",
			err:  fs.ErrNotExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			catchPanic(tt.err)
		})
	}
}

func Test_preprocessArgs(t *testing.T) {
	settingArgs := []string{"--args2", "--args3"}
	patch := gomonkey.NewPatches()
	defer patch.Reset()
	patch.ApplyFunc(config.Load, func() (*config.Config, error) {
		return &config.Config{
			Args: settingArgs,
		}, nil
	})

	os.Args = []string{"g", "-no-config"}
	preprocessArgs() // this will remove -no-config
	assert.Equal(t, 1, len(os.Args))

	os.Args = []string{"g", "--no-config"}
	preprocessArgs() // this will remove -no-config
	assert.Equal(t, 1, len(os.Args))

	os.Args = []string{"g", "--args1"}
	preprocessArgs() // this will add args from config
	assert.Equal(t, 4, len(os.Args))
	assert.Equal(t, os.Args, []string{"g", "--args2", "--args3", "--args1"})

	os.Args = []string{"g", "--args1", "-no-config"}
	preprocessArgs()
	assert.Equal(t, 2, len(os.Args))

	settingArgs[0] = "-no-config"
	os.Args = []string{"g", "--args1"}
	preprocessArgs()
	assert.Equal(t, 2, len(os.Args))
}

func TestPreprocessArgs_DelimiterProtectsNoConfigPaths(t *testing.T) {
	originalArgs := os.Args
	t.Cleanup(func() { os.Args = originalArgs })
	tests := []struct {
		name             string
		args             []string
		expectConfigLoad bool
		expectedArgs     []string
	}{
		{
			name:             "suffix paths do not disable config",
			args:             []string{"g", "--", "--no-config", "-no-config"},
			expectConfigLoad: true,
			expectedArgs:     []string{"g", "--", "--no-config", "-no-config"},
		},
		{
			name:             "prefix flag disables config without removing suffix paths",
			args:             []string{"g", "--no-config", "--", "--no-config", "-no-config"},
			expectConfigLoad: false,
			expectedArgs:     []string{"g", "--", "--no-config", "-no-config"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configLoaded := false
			patch := gomonkey.ApplyFunc(config.Load, func() (*config.Config, error) {
				configLoaded = true
				return &config.Config{}, nil
			})
			defer patch.Reset()
			os.Args = append([]string{}, tt.args...)

			preprocessArgs()

			assert.Equal(t, tt.expectConfigLoad, configLoaded)
			require.Equal(t, tt.expectedArgs, os.Args)
		})
	}
}

func TestPreprocessArgs_NilConfigKeepsDelimiterPaths(t *testing.T) {
	originalArgs := os.Args
	t.Cleanup(func() { os.Args = originalArgs })
	patch := gomonkey.ApplyFunc(config.Load, func() (*config.Config, error) {
		return nil, nil
	})
	t.Cleanup(patch.Reset)
	os.Args = []string{"g", "--", "--no-config"}

	assert.NotPanics(t, preprocessArgs)

	require.Equal(t, []string{"g", "--", "--no-config"}, os.Args)
}

func TestSeparateArgs(t *testing.T) {
	originalFlags := cli.G.Flags
	defer func() { cli.G.Flags = originalFlags }()
	cli.G.Flags = []ucli.Flag{
		&ucli.BoolFlag{Name: "all", Aliases: []string{"a"}},
		&ucli.StringFlag{Name: "sort", Aliases: []string{"s"}},
		&ucli.IntFlag{Name: "term-width"},
	}

	tests := []struct {
		name          string
		args          []string
		expectedFlags []string
		expectedPaths []string
	}{
		{
			name:          "Basic case",
			args:          []string{"--all", "dir1", "dir2"},
			expectedFlags: []string{"--all"},
			expectedPaths: []string{"dir1", "dir2"},
		},
		{
			name:          "Flag with value",
			args:          []string{"--sort", "name", "dir1"},
			expectedFlags: []string{"--sort", "name"},
			expectedPaths: []string{"dir1"},
		},
		{
			name:          "Flag with equals",
			args:          []string{"--sort=name", "dir1"},
			expectedFlags: []string{"--sort=name"},
			expectedPaths: []string{"dir1"},
		},
		{
			name:          "Mixed flags and paths",
			args:          []string{"--all", "dir1", "--sort", "name", "dir2"},
			expectedFlags: []string{"--all", "--sort", "name"},
			expectedPaths: []string{"dir1", "dir2"},
		},
		{
			name:          "With double dash",
			args:          []string{"--all", "dir1", "--", "--sort", "name"},
			expectedFlags: []string{"--all", "--"},
			expectedPaths: []string{"dir1", "--sort", "name"},
		},
		{
			name:          "Delimiter protects all following paths",
			args:          []string{"--", "a", "--bad"},
			expectedFlags: []string{"--"},
			expectedPaths: []string{"a", "--bad"},
		},
		{
			name:          "Short flags",
			args:          []string{"-a", "-s", "name", "dir1"},
			expectedFlags: []string{"-a", "-s", "name"},
			expectedPaths: []string{"dir1"},
		},
		{
			name:          "Complex case",
			args:          []string{"--all", "dir1", "--term-width", "100", "-s", "name", "--", "--fake-flag", "dir2"},
			expectedFlags: []string{"--all", "--term-width", "100", "-s", "name", "--"},
			expectedPaths: []string{"dir1", "--fake-flag", "dir2"},
		},
		{
			name:          "Complex case with double dash",
			args:          []string{"--all", "dir1", "--term-width", "100", "-s", "name", "--", "-a", "-a", "-l", "dir2", "--", "--fake-flag"},
			expectedFlags: []string{"--all", "--term-width", "100", "-s", "name", "--"},
			expectedPaths: []string{"dir1", "-a", "-a", "-l", "dir2", "--", "--fake-flag"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags, paths := separateArgs(tt.args)
			if !reflect.DeepEqual(flags, tt.expectedFlags) {
				t.Errorf("flags = %v, want %v", flags, tt.expectedFlags)
			}
			if !reflect.DeepEqual(paths, tt.expectedPaths) {
				t.Errorf("paths = %v, want %v", paths, tt.expectedPaths)
			}
		})
	}
}

func TestBuildFlagsWithArgsMap(t *testing.T) {
	originalFlags := cli.G.Flags
	defer func() { cli.G.Flags = originalFlags }()
	cli.G.Flags = []ucli.Flag{
		&ucli.BoolFlag{Name: "all", Aliases: []string{"a"}},
		&ucli.StringFlag{Name: "sort", Aliases: []string{"s"}},
		&ucli.IntFlag{Name: "term-width"},
	}

	expected := map[string]bool{
		"all":        false,
		"a":          false,
		"sort":       true,
		"s":          true,
		"term-width": true,
	}

	result := buildFlagsWithArgsMap()
	assert.Equal(t, result, expected, "buildFlagsWithArgsMap() = %v, want %v", result, expected)
}

func TestHandleLongFlag(t *testing.T) {
	tests := []struct {
		name                string
		arg                 string
		args                []string
		i                   int
		flagsWithArgs       map[string]bool
		expectedFlags       []string
		expectedExpectValue bool
		expectedI           int
	}{
		{
			name:                "Flag with equals",
			arg:                 "--sort=name",
			args:                []string{"--sort=name", "dir1"},
			i:                   0,
			flagsWithArgs:       map[string]bool{"sort": true},
			expectedFlags:       []string{"--sort=name"},
			expectedExpectValue: false,
			expectedI:           0,
		},
		{
			name:                "Flag without value",
			arg:                 "--all",
			args:                []string{"--all", "dir1"},
			i:                   0,
			flagsWithArgs:       map[string]bool{"all": false},
			expectedFlags:       []string{"--all"},
			expectedExpectValue: false,
			expectedI:           0,
		},
		{
			name:                "Flag expecting value",
			arg:                 "--sort",
			args:                []string{"--sort", "name", "dir1"},
			i:                   0,
			flagsWithArgs:       map[string]bool{"sort": true},
			expectedFlags:       []string{"--sort"},
			expectedExpectValue: true,
			expectedI:           0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := []string{}
			expectValue := false
			resultI := handleLongFlag(tt.arg, tt.args, tt.i, &flags, &expectValue, tt.flagsWithArgs)

			assert.Equal(t, flags, tt.expectedFlags, "flags = %v, want %v", flags, tt.expectedFlags)
			assert.Equal(t, resultI, tt.expectedI, "resultI = %v, expectedI %v", expectValue, tt.expectedExpectValue)
			assert.Equal(t, expectValue, tt.expectedExpectValue, "expectValue = %v, want %v", expectValue, tt.expectedExpectValue)
		})
	}
}

func TestHandleShortFlag(t *testing.T) {
	tests := []struct {
		name                string
		arg                 string
		args                []string
		i                   int
		flagsWithArgs       map[string]bool
		expectedFlags       []string
		expectedExpectValue bool
		expectedI           int
	}{
		{
			name:                "Short flag without value",
			arg:                 "-a",
			args:                []string{"-a", "dir1"},
			i:                   0,
			flagsWithArgs:       map[string]bool{"a": false},
			expectedFlags:       []string{"-a"},
			expectedExpectValue: false,
			expectedI:           0,
		},
		{
			name:                "Short flag expecting value",
			arg:                 "-s",
			args:                []string{"-s", "name", "dir1"},
			i:                   0,
			flagsWithArgs:       map[string]bool{"s": true},
			expectedFlags:       []string{"-s"},
			expectedExpectValue: true,
			expectedI:           0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := []string{}
			expectValue := false
			resultI := handleShortFlag(tt.arg, tt.args, tt.i, &flags, &expectValue, tt.flagsWithArgs)
			assert.Equal(t, flags, tt.expectedFlags, "flags = %v, want %v", flags, tt.expectedFlags)
			assert.Equal(t, resultI, tt.expectedI, "resultI = %v, expectedI %v", expectValue, tt.expectedExpectValue)
			assert.Equal(t, expectValue, tt.expectedExpectValue, "expectValue = %v, want %v", expectValue, tt.expectedExpectValue)
		})
	}
}

func Test_main(t *testing.T) {
	patch := gomonkey.ApplyFunc(os.Exit, func(int) {})
	defer patch.Reset()
	os.Args = []string{"g", "."}
	main()
}

func Test_main_TimeTypeAliasesWithFilters(t *testing.T) {
	patch := gomonkey.ApplyFunc(os.Exit, func(int) {})
	defer patch.Reset()

	tmp := t.TempDir()
	origWD, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	defer func() {
		_ = os.Chdir(origWD)
	}()

	tests := []struct {
		name string
		args []string
	}{
		{name: "after", args: []string{"g", "-no-config", "-A", "--time-type", "modified", "--after", "2000-01-01", "."}},
		{name: "before", args: []string{"g", "-no-config", "-A", "--time-type", "modified", "--before", "2999-01-01", "."}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			cli.ReturnCode = 0
			main()
		})
	}
}
