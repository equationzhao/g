package parse

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Equationzhao/g/internal/request"
)

type Error struct {
	Msg  string
	Code int
}

func (e *Error) Error() string { return e.Msg }

func Usage() string {
	var b strings.Builder
	b.WriteString("Usage: g [OPTION]... [PATH]...\n\n")
	for _, s := range Specs() {
		names := s.Name
		if s.Short != 0 && s.Name != "-"+string(s.Short) {
			names = "-" + string(s.Short) + ", " + s.Name
		}
		fmt.Fprintf(&b, "  %-28s %s\n", names, s.Help)
	}
	return b.String()
}

const Version = "1.0.0"

// Parse converts argv (without argv0) and env into a Request. It never reads or writes os.Args.
func Parse(args, env []string) (request.Request, error) {
	r := request.Default()
	byToken := tokenIndex()
	i := 0
	for i < len(args) {
		tok := args[i]
		if tok == "--" {
			r.Paths = append(r.Paths, args[i+1:]...)
			break
		}
		if tok == "-" || !strings.HasPrefix(tok, "-") {
			r.Paths = append(r.Paths, expandTilde(tok, env))
			i++
			continue
		}
		if strings.HasPrefix(tok, "--") {
			name, val, hasVal := splitLong(tok)
			spec, ok := byToken[name]
			if !ok {
				return r, &Error{Msg: "unknown option " + name, Code: 2}
			}
			if name == "--no-dereference" {
				r.Dereference = false
				i++
				continue
			}
			if spec.TakesValue && !hasVal && impliedValue(name) == "" {
				if i+1 >= len(args) {
					return r, &Error{Msg: name + " requires a value", Code: 2}
				}
				i++
				val = args[i]
				hasVal = true
			}
			if err := apply(spec, name, val, hasVal, &r); err != nil {
				return r, err
			}
			i++
			continue
		}
		// short cluster: -lah, -1, -I'*.o'
		body := tok[1:]
		for j := 0; j < len(body); {
			sh := rune(body[j])
			name := "-" + string(sh)
			spec, ok := byToken[name]
			if !ok {
				return r, &Error{Msg: "unknown option " + name, Code: 2}
			}
			if impliedValue(name) != "" {
				if err := apply(spec, name, "", false, &r); err != nil {
					return r, err
				}
				j++
				continue
			}
			if spec.TakesValue && !optionalWhen(spec.Name) {
				rest := body[j+1:]
				var val string
				if rest != "" {
					val = rest
				} else {
					if i+1 >= len(args) {
						return r, &Error{Msg: name + " requires a value", Code: 2}
					}
					i++
					val = args[i]
				}
				if err := apply(spec, name, val, true, &r); err != nil {
					return r, err
				}
				break
			}
			if err := apply(spec, name, "", false, &r); err != nil {
				return r, err
			}
			j++
		}
		i++
	}
	if len(r.Paths) == 0 {
		r.Paths = []string{"."}
	}
	if ts := envVal(env, "TIME_STYLE"); ts != "" && r.TimeStyle.Named == "default" && r.TimeStyle.Custom == "" {
		r.TimeStyle.Named = ts
	}
	if err := r.Validate(); err != nil {
		return r, &Error{Msg: err.Error(), Code: 2}
	}
	return r, nil
}

func tokenIndex() map[string]Spec {
	return AllNames()
}

func splitLong(tok string) (name, val string, hasVal bool) {
	if i := strings.IndexByte(tok, '='); i >= 0 {
		return tok[:i], tok[i+1:], true
	}
	return tok, "", false
}

func optionalWhen(specName string) bool {
	switch specName {
	case "--color", "--icons", "--hyperlink", "--classify":
		return true
	default:
		return false
	}
}

func impliedValue(token string) string {
	switch token {
	case "-C":
		return "grid"
	case "-x":
		return "across"
	case "-1":
		return "oneline"
	case "-m":
		return "comma"
	case "-T":
		return "tree"
	case "--json":
		return "json"
	case "-t":
		return "time"
	case "-S":
		return "size"
	case "-X":
		return "ext"
	case "-U":
		return "none"
	case "-v":
		return "version"
	case "--dir-first", "--group-directories-first":
		return "first"
	case "--no-config":
		return "none"
	default:
		return ""
	}
}

func apply(spec Spec, token, val string, hasVal bool, r *request.Request) error {
	if imp := impliedValue(token); imp != "" && !hasVal {
		val = imp
		hasVal = true
	}
	switch spec.Name {
	case "--help":
		r.Help = true
	case "--version":
		r.Version = true
	case "--config":
		if !hasVal || val == "" {
			return &Error{Msg: "--config requires PATH or none", Code: 2}
		}
		if val == "none" {
			r.Config = request.ConfigNone
			r.ConfigPath = ""
		} else {
			r.Config = request.ConfigPath
			r.ConfigPath = val
		}
	case "--format":
		f, err := parseFormat(val)
		if err != nil {
			return err
		}
		r.Format = f
		r.FormatSet = true
	case "--long":
		r.Long = true
		if token == "-o" {
			r.LongNoGroup = true
		}
	case "--all":
		r.Visibility = request.VisAll
	case "--almost-all":
		r.Visibility = request.VisAlmostAll
	case "--directory":
		r.DirSelf = true
	case "--recursive":
		r.Recurse = true
	case "--depth":
		n, err := strconv.Atoi(val)
		if err != nil {
			return &Error{Msg: "--depth must be an integer", Code: 2}
		}
		r.Depth = n
		r.HasDepth = true
	case "--ignore":
		if val == "" {
			return &Error{Msg: "-I requires a glob", Code: 2}
		}
		r.Ignore = append(r.Ignore, val)
	case "--only-dirs":
		r.KindFilter = request.KindDirsOnly
	case "--only-files":
		r.KindFilter = request.KindFilesOnly
	case "--ignore-backups":
		r.IgnoreBackups = true
	case "--git-ignore":
		r.GitIgnore = true
	case "--sort":
		k, err := parseSort(val)
		if err != nil {
			return err
		}
		r.Sort = k
	case "--reverse":
		r.Reverse = true
	case "--dir-order":
		o, err := parseDirOrder(val)
		if err != nil {
			return err
		}
		r.DirOrder = o
	case "--human-readable":
		r.HumanReadable = true
	case "--si":
		r.SI = true
	case "--inode":
		r.Inode = true
	case "--links":
		r.Links = true
	case "--numeric-uid-gid":
		r.NumericIDs = true
	case "--no-group":
		r.LongNoGroup = true
	case "--blocks":
		r.Blocks = true
	case "--header":
		r.Header = true
	case "--time":
		tf, err := parseTimeField(val)
		if err != nil {
			return err
		}
		r.TimeField = tf
	case "--time-style":
		if strings.HasPrefix(val, "+") {
			r.TimeStyle = request.TimeStyle{Custom: val[1:]}
		} else {
			r.TimeStyle = request.TimeStyle{Named: val}
		}
	case "--color":
		w, err := parseWhen(val, hasVal)
		if err != nil {
			return err
		}
		r.Color = w
	case "--icons":
		w, err := parseWhen(val, hasVal)
		if err != nil {
			return err
		}
		r.Icons = w
	case "--hyperlink":
		w, err := parseWhen(val, hasVal)
		if err != nil {
			return err
		}
		r.Hyperlink = w
	case "--theme":
		if !hasVal || val == "" {
			return &Error{Msg: "--theme requires a path", Code: 2}
		}
		r.ThemePath = val
	case "--classify":
		w, err := parseWhen(val, hasVal)
		if err != nil {
			return err
		}
		r.Classify = w
		r.ClassifySet = true
	case "--quote-name":
		r.Quote = request.QuoteAlways
	case "--literal":
		r.Quote = request.QuoteLiteral
	case "--zero":
		applyZero(r)
	case "--width":
		n, err := strconv.Atoi(val)
		if err != nil {
			return &Error{Msg: "--width must be an integer", Code: 2}
		}
		r.Width = n
		r.WidthSet = true
	case "--git":
		r.Git = true
	case "--dereference":
		r.Dereference = token != "--no-dereference"
	case "-g":
		r.Long = true
		r.LongNoOwner = true
	default:
		return &Error{Msg: "unhandled option " + spec.Name, Code: 2}
	}
	return nil
}

func applyZero(r *request.Request) {
	r.Zero = true
	r.Format = request.FormatOneline
	r.FormatSet = true
	r.Color = request.WhenNever
	r.Icons = request.WhenNever
	r.Hyperlink = request.WhenNever
	r.Classify = request.WhenNever
	r.ClassifySet = true
	r.Quote = request.QuoteLiteral
}

func parseFormat(v string) (request.Format, error) {
	switch v {
	case "grid":
		return request.FormatGrid, nil
	case "across":
		return request.FormatAcross, nil
	case "oneline":
		return request.FormatOneline, nil
	case "comma":
		return request.FormatComma, nil
	case "tree":
		return request.FormatTree, nil
	case "json":
		return request.FormatJSON, nil
	default:
		return 0, &Error{Msg: "unknown --format value " + v, Code: 2}
	}
}

func parseSort(v string) (request.SortKey, error) {
	switch v {
	case "name":
		return request.SortName, nil
	case "size":
		return request.SortSize, nil
	case "time":
		return request.SortTime, nil
	case "ext":
		return request.SortExt, nil
	case "version":
		return request.SortVersion, nil
	case "none":
		return request.SortNone, nil
	default:
		return 0, &Error{Msg: "unknown --sort value " + v, Code: 2}
	}
}

func parseDirOrder(v string) (request.DirOrder, error) {
	switch v {
	case "none":
		return request.DirOrderNone, nil
	case "first":
		return request.DirOrderFirst, nil
	case "last":
		return request.DirOrderLast, nil
	default:
		return 0, &Error{Msg: "unknown --dir-order value " + v, Code: 2}
	}
}

func parseTimeField(v string) (request.TimeField, error) {
	switch v {
	case "modified":
		return request.TimeModified, nil
	case "accessed":
		return request.TimeAccessed, nil
	case "changed", "created":
		return request.TimeChanged, nil
	case "birth":
		return request.TimeBirth, nil
	default:
		return 0, &Error{Msg: "unknown --time value " + v, Code: 2}
	}
}

func parseWhen(v string, hasVal bool) (request.When, error) {
	if !hasVal || v == "" {
		return request.WhenAlways, nil
	}
	switch v {
	case "always":
		return request.WhenAlways, nil
	case "auto", "automatic":
		return request.WhenAuto, nil
	case "never":
		return request.WhenNever, nil
	default:
		return 0, &Error{Msg: "unknown WHEN value " + v, Code: 2}
	}
}

func envVal(env []string, key string) string {
	prefix := key + "="
	for i := len(env) - 1; i >= 0; i-- {
		if strings.HasPrefix(env[i], prefix) {
			return env[i][len(prefix):]
		}
	}
	return ""
}

func expandTilde(p string, env []string) string {
	if p != "~" && !strings.HasPrefix(p, "~/") {
		return p
	}
	home := envVal(env, "HOME")
	if home == "" {
		h, err := os.UserHomeDir()
		if err != nil {
			return p
		}
		home = h
	}
	if p == "~" {
		return home
	}
	return home + p[1:]
}

func EnvNOCOLOR(env []string) bool {
	return envVal(env, "NO_COLOR") != ""
}

func EnvCOLUMNS(env []string) int {
	s := envVal(env, "COLUMNS")
	if s == "" {
		return 0
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return 0
	}
	return n
}
