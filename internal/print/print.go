package print

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/Equationzhao/g/internal/entry"
	"github.com/Equationzhao/g/internal/quote"
	"github.com/Equationzhao/g/internal/request"
	"github.com/itchyny/timefmt-go"
	"github.com/mattn/go-runewidth"
)

type Root struct {
	Path    string
	Abs     string
	Entries []entry.Entry
}

type Job struct {
	Roots   []Root
	Req     request.Request
	Width   int
	Now     time.Time
	ColorOn bool
}

func Print(w io.Writer, job Job) error {
	req := job.Req
	if req.Format == request.FormatJSON {
		return printJSON(w, job)
	}
	multi := len(job.Roots) > 1 || req.Recurse
	for i, root := range job.Roots {
		if multi && !req.Zero && req.Format != request.FormatTree {
			if i > 0 {
				if _, err := io.WriteString(w, "\n"); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintf(w, "%s:\n", root.Path); err != nil {
				return err
			}
		}
		if err := printRoot(w, root, job); err != nil {
			return err
		}
	}
	return nil
}

func printRoot(w io.Writer, root Root, job Job) error {
	req := job.Req
	switch {
	case req.Format == request.FormatTree:
		return printTree(w, root, job)
	case req.Long || req.Format == request.FormatOneline:
		return printLines(w, root.Entries, job, false)
	case req.Format == request.FormatComma:
		return printComma(w, root.Entries, job)
	case req.Format == request.FormatAcross:
		return printGrid(w, root.Entries, job, true)
	default:
		return printGrid(w, root.Entries, job, false)
	}
}

func displayName(e entry.Entry, req request.Request) string {
	name := quote.Name(e.Name, req.Quote, req.Zero)
	if req.Classify == request.WhenAlways {
		name += classify(e)
	}
	if e.Kind == entry.KindSymlink || e.Kind == entry.KindBrokenSymlink {
		if e.Target != "" && !req.Dereference && (req.Long || req.Format == request.FormatOneline || req.Format == request.FormatTree) {
			name += " -> " + e.Target
		}
	}
	if req.Git && e.Git != "" {
		name = e.Git + " " + name
	}
	if req.Inode && !(req.Long && true) {
		// inode prefix when not absorbed into long columns differently
	}
	var prefix []string
	if req.Inode {
		prefix = append(prefix, padInode(e.Inode))
	}
	if req.Links && !req.Long {
		prefix = append(prefix, fmt.Sprintf("%d", e.Nlink))
	}
	if req.Git && !req.Long {
		// already prepended
	}
	if len(prefix) > 0 && !req.Long {
		name = strings.Join(prefix, " ") + " " + name
	}
	return name
}

func classify(e entry.Entry) string {
	switch e.Kind {
	case entry.KindDir:
		return "/"
	case entry.KindExec:
		return "*"
	case entry.KindSymlink, entry.KindBrokenSymlink:
		return "@"
	case entry.KindPipe:
		return "|"
	case entry.KindSocket:
		return "="
	default:
		return ""
	}
}

func padInode(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func printLines(w io.Writer, ents []entry.Entry, job Job, _ bool) error {
	req := job.Req
	nl := "\n"
	if req.Zero {
		nl = "\x00"
	}
	if req.Long && req.Header && !req.Zero {
		if _, err := io.WriteString(w, headerLine(req)+"\n"); err != nil {
			return err
		}
	}
	for _, e := range ents {
		var line string
		if req.Long {
			line = longLine(e, req, job.Now)
		} else {
			line = displayName(e, req)
		}
		if _, err := io.WriteString(w, line+nl); err != nil {
			return err
		}
	}
	return nil
}

func headerLine(req request.Request) string {
	var p []string
	if req.Inode {
		p = append(p, "INODE")
	}
	if req.Blocks {
		p = append(p, "BLOCKS")
	}
	p = append(p, "MODE", "NLINK")
	if !req.LongNoOwner {
		p = append(p, "USER")
	}
	if !req.LongNoGroup {
		p = append(p, "GROUP")
	}
	p = append(p, "SIZE", "TIME")
	if req.Git {
		p = append(p, "GIT")
	}
	p = append(p, "NAME")
	return strings.Join(p, " ")
}

func longLine(e entry.Entry, req request.Request, now time.Time) string {
	var p []string
	if req.Inode {
		p = append(p, padInode(e.Inode))
	}
	if req.Blocks {
		p = append(p, fmt.Sprintf("%d", e.Blocks))
	}
	p = append(p, modeString(e), fmt.Sprintf("%d", e.Nlink))
	if !req.LongNoOwner {
		if req.NumericIDs || e.User == "" {
			p = append(p, dash(e.UID))
		} else {
			p = append(p, e.User)
		}
	}
	if !req.LongNoGroup {
		if req.NumericIDs || e.Group == "" {
			p = append(p, dash(e.GID))
		} else {
			p = append(p, e.Group)
		}
	}
	p = append(p, sizeString(e, req), timeString(e, req, now))
	if req.Git {
		g := e.Git
		if g == "" {
			g = "--"
		}
		p = append(p, g)
	}
	p = append(p, displayName(e, req))
	return strings.Join(p, " ")
}

func dash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func modeString(e entry.Entry) string {
	m := fs.FileMode(e.Mode)
	return m.String()
}

func sizeString(e entry.Entry, req request.Request) string {
	if req.SI {
		return HumanSize(e.Size, true)
	}
	if req.HumanReadable {
		return HumanSize(e.Size, false)
	}
	return fmt.Sprintf("%d", e.Size)
}

func timeString(e entry.Entry, req request.Request, now time.Time) string {
	t := e.ModTime
	switch req.TimeField {
	case request.TimeAccessed:
		t = e.AccTime
	case request.TimeChanged:
		t = e.ChangeTime
	case request.TimeBirth:
		if e.HasBirth {
			t = e.Birth
		} else {
			return "-"
		}
	}
	if t.IsZero() {
		return "-"
	}
	if now.IsZero() {
		now = time.Now()
	}
	style := req.TimeStyle.Named
	if req.TimeStyle.Custom != "" {
		return timefmt.Format(t, req.TimeStyle.Custom)
	}
	switch style {
	case "iso":
		return t.Format("01-02 15:04")
	case "long-iso":
		return t.Format("2006-01-02 15:04")
	case "full-iso":
		return t.Format("2006-01-02 15:04:05.000000000 -0700")
	case "relative":
		return relative(now.Sub(t))
	default:
		if t.Year() == now.Year() {
			return t.Format("Jan 02 15:04")
		}
		return t.Format("Jan 02  2006")
	}
}

func relative(d time.Duration) string {
	if d < 0 {
		d = -d
	}
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%d days ago", int(d.Hours()/24))
	}
}

func printComma(w io.Writer, ents []entry.Entry, job Job) error {
	names := make([]string, len(ents))
	for i, e := range ents {
		names[i] = displayName(e, job.Req)
	}
	_, err := io.WriteString(w, strings.Join(names, ", ")+"\n")
	return err
}

func printGrid(w io.Writer, ents []entry.Entry, job Job, across bool) error {
	if len(ents) == 0 {
		return nil
	}
	width := job.Width
	if width <= 0 {
		width = 80
	}
	names := make([]string, len(ents))
	maxw := 0
	for i, e := range ents {
		names[i] = displayName(e, job.Req)
		if n := runewidth.StringWidth(names[i]); n > maxw {
			maxw = n
		}
	}
	colw := maxw + 2
	if colw < 1 {
		colw = 1
	}
	cols := width / colw
	if cols < 1 {
		cols = 1
	}
	if cols > len(ents) {
		cols = len(ents)
	}
	rows := (len(ents) + cols - 1) / cols
	if across {
		for i, n := range names {
			pad := colw - runewidth.StringWidth(n)
			if pad < 1 {
				pad = 1
			}
			if _, err := io.WriteString(w, n); err != nil {
				return err
			}
			if (i+1)%cols == 0 || i == len(names)-1 {
				if _, err := io.WriteString(w, "\n"); err != nil {
					return err
				}
			} else {
				if _, err := io.WriteString(w, strings.Repeat(" ", pad)); err != nil {
					return err
				}
			}
		}
		return nil
	}
	// down then across (GNU -C)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			i := c*rows + r
			if i >= len(names) {
				continue
			}
			n := names[i]
			if _, err := io.WriteString(w, n); err != nil {
				return err
			}
			if c < cols-1 && c*rows+r+rows < len(names) {
				pad := colw - runewidth.StringWidth(n)
				if pad < 1 {
					pad = 1
				}
				if _, err := io.WriteString(w, strings.Repeat(" ", pad)); err != nil {
					return err
				}
			}
		}
		if _, err := io.WriteString(w, "\n"); err != nil {
			return err
		}
	}
	return nil
}

func printJSON(w io.Writer, job Job) error {
	_, err := io.WriteString(w, `{"roots":[`)
	if err != nil {
		return err
	}
	for i, root := range job.Roots {
		if i > 0 {
			if _, err := io.WriteString(w, ","); err != nil {
				return err
			}
		}
		fmt.Fprintf(w, `{"path":%q,"entries":[`, root.Abs)
		for j, e := range root.Entries {
			if j > 0 {
				io.WriteString(w, ",")
			}
			fmt.Fprintf(w, `{"name":%q,"path":%q,"type":%q,"size":%d,"mtime":%q}`,
				e.Name, e.Path, typeName(e), e.Size, e.ModTime.UTC().Format(time.RFC3339))
		}
		io.WriteString(w, "]}")
	}
	_, err = io.WriteString(w, "]}\n")
	return err
}

func typeName(e entry.Entry) string {
	switch e.Kind {
	case entry.KindDir:
		return "dir"
	case entry.KindSymlink:
		return "symlink"
	case entry.KindBrokenSymlink:
		return "symlink"
	default:
		return "file"
	}
}

func printTree(w io.Writer, root Root, job Job) error {
	req := job.Req
	byParent := map[string][]entry.Entry{}
	var tops []entry.Entry
	for _, e := range root.Entries {
		if e.Depth == 0 {
			tops = append(tops, e)
			continue
		}
		byParent[e.Parent] = append(byParent[e.Parent], e)
	}
	if len(tops) == 0 {
		// children-only walk: synthesize root
		tops = []entry.Entry{{Name: filepath.Base(root.Path), Path: root.Abs, Kind: entry.KindDir, Depth: 0}}
		byParent[root.Abs] = root.Entries
	}
	unicode := job.Req.Format == request.FormatTree // default unicode; ASCII if we want
	for i, t := range tops {
		_ = i
		if err := writeTree(w, t, byParent, "", true, req, job.Now, unicode); err != nil {
			return err
		}
	}
	return nil
}

func writeTree(w io.Writer, e entry.Entry, byParent map[string][]entry.Entry, prefix string, isLast bool, req request.Request, now time.Time, unicode bool) error {
	branch := ""
	if prefix != "" || e.Depth > 0 {
		if unicode {
			if isLast {
				branch = "└── "
			} else {
				branch = "├── "
			}
		} else {
			if isLast {
				branch = "`-- "
			} else {
				branch = "|-- "
			}
		}
	}
	name := displayName(e, req)
	line := prefix + branch + name
	if req.Long && e.Depth >= 0 && prefix+branch != "" {
		line = prefix + branch + longLine(e, req, now)
	} else if req.Long && prefix == "" && e.Depth == 0 {
		line = longLine(e, req, now)
	}
	if _, err := io.WriteString(w, line+"\n"); err != nil {
		return err
	}
	kids := byParent[e.Path]
	nextPref := prefix
	if e.Depth > 0 || prefix != "" || branch != "" {
		if unicode {
			if isLast {
				nextPref = prefix + "    "
			} else {
				nextPref = prefix + "│   "
			}
		} else {
			if isLast {
				nextPref = prefix + "    "
			} else {
				nextPref = prefix + "|   "
			}
		}
	}
	for i, k := range kids {
		if err := writeTree(w, k, byParent, nextPref, i == len(kids)-1, req, now, unicode); err != nil {
			return err
		}
	}
	return nil
}
