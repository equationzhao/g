package parse

import (
	"fmt"
	"strings"
)

//go:generate go run ./docsgen

func RenderFlagRegistry() string {
	var b strings.Builder
	b.WriteString("# Flag short-letter registry\n\n")
	b.WriteString("Generated from `internal/parse.Specs()` and `GNUReservedShorts()`.\n")
	b.WriteString("Do not edit by hand; run `go generate ./internal/parse`.\n\n")
	b.WriteString("Letters GNU ls already assigned stay **reserved** even when we do not implement that feature.\n")
	b.WriteString("Inventing a new meaning for a reserved letter is forbidden (docs/rewrite-architecture.md §17.3).\n\n")

	used := map[rune]Spec{}
	for _, s := range Specs() {
		if s.Short != 0 {
			used[s.Short] = s
		}
		for _, a := range s.Aliases {
			if len(a) == 2 && a[0] == '-' && a[1] != '-' {
				used[rune(a[1])] = s
			}
		}
	}
	reserved := GNUReservedShorts()

	b.WriteString("## Used\n\n")
	b.WriteString("| Letter | Spec | Slot |\n| --- | --- | --- |\n")
	for _, r := range letters() {
		if s, ok := used[r]; ok {
			fmt.Fprintf(&b, "| `-%c` | `%s` | `%s` |\n", r, s.Name, s.Slot)
		}
	}

	b.WriteString("\n## Reserved (GNU meaning, not ours)\n\n")
	b.WriteString("| Letter | GNU meaning |\n| --- | --- |\n")
	for _, r := range letters() {
		if why, ok := reserved[r]; ok {
			fmt.Fprintf(&b, "| `-%c` | %s |\n", r, why)
		}
	}

	b.WriteString("\n## Free\n\n")
	var free []string
	for _, r := range letters() {
		if _, u := used[r]; u {
			continue
		}
		if _, rsv := reserved[r]; rsv {
			continue
		}
		if r == '?' {
			continue
		}
		free = append(free, fmt.Sprintf("`-%c`", r))
	}
	b.WriteString(strings.Join(free, ", "))
	b.WriteString("\n\nA free letter is not permission to add a flag. See §17 five gates.\n")
	return b.String()
}

func RenderRejected() string {
	rows, err := LoadRejected()
	if err != nil {
		return "# Rejected flags\n\nERROR: " + err.Error() + "\n"
	}
	var b strings.Builder
	b.WriteString("# Rejected flags\n\n")
	b.WriteString("Cemetery of names that must not come back as `Specs()` entries.\n")
	b.WriteString("Generated from `internal/parse/testdata/rejected.tsv`.\n")
	b.WriteString("Do not edit by hand; run `go generate ./internal/parse`.\n\n")
	b.WriteString("When someone files \"please add `--duplicate`\", paste a link to the matching row and close the issue.\n\n")
	b.WriteString("| Name | Use instead | Why |\n| --- | --- | --- |\n")
	for _, r := range rows {
		fmt.Fprintf(&b, "| `%s` | %s | %s |\n", r.Name, r.Replacement, r.Reason)
	}
	return b.String()
}

func RenderManOptions() string {
	var b strings.Builder
	b.WriteString("# g — rewrite OPTIONS (generated)\n\n")
	b.WriteString("This fragment is the OPTIONS source of truth for the rewrite.\n")
	b.WriteString("Prose (DESCRIPTION, EXAMPLES) lives in `docs/rewrite-architecture.md`.\n")
	b.WriteString("Do not edit the OPTIONS block by hand; run `go generate ./internal/parse`.\n\n")
	b.WriteString("<!-- BEGIN OPTIONS -->\n\n")
	var dims []Dimension
	seen := map[Dimension]bool{}
	for _, s := range Specs() {
		if !seen[s.Dimension] {
			dims = append(dims, s.Dimension)
			seen[s.Dimension] = true
		}
	}
	for _, d := range dims {
		fmt.Fprintf(&b, "## %s\n\n", d)
		for _, s := range Specs() {
			if s.Dimension != d {
				continue
			}
			names := []string{s.Name}
			if s.Short != 0 && s.Name != "-"+string(s.Short) {
				names = append([]string{"-" + string(s.Short)}, names...)
			}
			names = append(names, s.Aliases...)
			fmt.Fprintf(&b, "`%s`\n", strings.Join(uniq(names), "`, `"))
			fmt.Fprintf(&b, ": %s (default: %s)\n\n", s.Help, s.Default)
		}
	}
	b.WriteString("<!-- END OPTIONS -->\n")
	return b.String()
}

func letters() []rune {
	out := make([]rune, 0, 53)
	for r := 'a'; r <= 'z'; r++ {
		out = append(out, r)
	}
	for r := 'A'; r <= 'Z'; r++ {
		out = append(out, r)
	}
	for _, r := range []rune{'0', '1', '?', '#'} {
		out = append(out, r)
	}
	return out
}

func uniq(in []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
