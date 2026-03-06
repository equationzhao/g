/*
MIT License

Copyright (c) 2022 Skyascii

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

// code from github.com/savioxavier/termlink

package util

import (
	"fmt"
	"os"
	"slices"
)

// v struct represents a semver version (usually, with some exceptions)
// with major, minor, and patch segments
type v struct {
	major int
	minor int
	patch int
}

// parseVersion takes a string "version" number and returns
// a Version struct with the major, minor, and patch
// segments parsed from the string.
// If a version number is not provided
func parseVersion(version string) v {
	var major, minor, patch int
	_, _ = fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)
	return v{
		major: major,
		minor: minor,
		patch: patch,
	}
}

// hasEnvironmentVariables returns true if the environment variable "name"
// is present in the environment, false otherwise
func hasEnv(name string) bool {
	_, envExists := os.LookupEnv(name)

	return envExists
}

// checkAllEnvs returns true if any of the environment variables in the "vars"
// string slice are actually present in the environment, false otherwise
func checkAllEnvs(vars []string) bool {
	return slices.ContainsFunc(vars, hasEnv)
}

// getEnv returns the value of the environment variable, if it exists
func getEnv(name string) string {
	envValue, _ := os.LookupEnv(name)

	return envValue
}

// matchesEnv returns true if the environment variable "name" matches any
// of the given values in the "values" string slice, false otherwise
func matchesEnv(name string, values []string) bool {
	if hasEnv(name) {
		if slices.Contains(values, getEnv(name)) {
			return true
		}
	}
	return false
}

func SupportsHyperlinks() bool {
	// Allow hyperlinks to be forced, independent of any environment variables
	// Instead of checking whether it is equal to anything other than "0",
	// a set of allowed values are provided, as something like
	// FORCE_HYPERLINK="do-not-enable-it" wouldn't make sense if it returned true
	if matchesEnv("FORCE_HYPERLINK", []string{"1", "true", "always", "enabled"}) {
		return true
	}

	// VTE-based terminals (Gnome Terminal, Guake, ROXTerm, etc)
	// VTE_VERSION is rendered as four-digit version string
	// eg: 0.52.2 => 5202
	// parseVersion will parse it with a standalone major segment
	// with minor and patch segments set to 0
	// 0.50.0 (parsed as 5000) was supposed to support hyperlinks, but throws a segfault
	// so we check if the "major" version is greater than 5000 (5000 exclusive)
	if hasEnv("VTE_VERSION") {
		v := parseVersion(getEnv("VTE_VERSION"))
		return v.major > 5000
	}

	// Terminals which have a TERM_PROGRAM variable set
	// This is the most versatile environment variable as it also provides another
	// variable called TERM_PROGRAM_VERSION, which helps us to determine
	// the exact version of the program, and allow for stricter variable checks
	if hasEnv("TERM_PROGRAM") {
		v := parseVersion(getEnv("TERM_PROGRAM_VERSION"))

		switch term := getEnv("TERM_PROGRAM"); term {
		case "iTerm.app":
			if v.major == 3 {
				return v.minor >= 1
			}
			return v.major > 3
		case "WezTerm":
			// Even though WezTerm's version is something like 20200620-160318-e00b076c
			// parseVersion will still parse it with a standalone major segment (ie: 20200620)
			// with minor and patch segments set to 0
			return v.major >= 20200620
		case "vscode":
			return v.major > 1 || (v.major == 1 && v.minor >= 72)
		case "ghostty":
			// It is unclear when during the private beta that ghostty started supporting hyperlinks,
			// so we'll start from the public release.
			return v.major >= 1

			// Hyper Terminal used to be included in this list, and it even supports hyperlinks
			// but the hyperlinks are pseudo-hyperlinks and are actually not clickable
		}
	}

	// Terminals which have a TERM variable set
	if matchesEnv("TERM", []string{"xterm-kitty", "alacritty", "alacritty-direct", "xterm-ghostty"}) {
		return true
	}

	// Terminals which have a COLORTERM variable set
	if matchesEnv("COLORTERM", []string{"xfce4-terminal"}) {
		return true
	}

	// Terminals in JetBrains IDEs
	if matchesEnv("TERMINAL_EMULATOR", []string{"JetBrains-JediTerm"}) {
		return true
	}

	// Match standalone environment variables
	// ie, those which do not require any special handling
	// or version checking
	if checkAllEnvs([]string{
		"DOMTERM",
		"WT_SESSION",
		"KONSOLE_VERSION",
	}) {
		return true
	}

	return false
}
