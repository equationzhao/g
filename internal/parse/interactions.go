package parse

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Exception is a flag-level override of the dimension matrix.
type Exception struct {
	A, B     string
	Relation Relation
	Note     string
}

func testdata(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "testdata", name)
}

// LoadMatrix reads the 12×12 dimension relation table.
func LoadMatrix() (map[[2]Dimension]Relation, error) {
	b, err := os.ReadFile(testdata("interactions.tsv"))
	if err != nil {
		return nil, err
	}
	return parseMatrix(b)
}

func parseMatrix(b []byte) (map[[2]Dimension]Relation, error) {
	sc := bufio.NewScanner(bytes.NewReader(b))
	var cols []Dimension
	out := make(map[[2]Dimension]Relation, int(dimCount)*int(dimCount))
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, "\t")
		if cols == nil {
			if fields[0] != "dim" {
				return nil, fmt.Errorf("interactions.tsv:%d: first cell must be dim", lineNo)
			}
			cols = make([]Dimension, 0, len(fields)-1)
			for _, name := range fields[1:] {
				d, ok := parseDimension(name)
				if !ok {
					return nil, fmt.Errorf("interactions.tsv:%d: unknown dim %q", lineNo, name)
				}
				cols = append(cols, d)
			}
			if len(cols) != int(dimCount) {
				return nil, fmt.Errorf("interactions.tsv:%d: want %d dim columns, got %d", lineNo, dimCount, len(cols))
			}
			continue
		}
		if len(fields) != len(cols)+1 {
			return nil, fmt.Errorf("interactions.tsv:%d: want %d fields, got %d", lineNo, len(cols)+1, len(fields))
		}
		row, ok := parseDimension(fields[0])
		if !ok {
			return nil, fmt.Errorf("interactions.tsv:%d: unknown dim %q", lineNo, fields[0])
		}
		for i, cell := range fields[1:] {
			rel, ok := parseRelation(cell)
			if !ok {
				return nil, fmt.Errorf("interactions.tsv:%d: unknown relation %q", lineNo, cell)
			}
			out[[2]Dimension{row, cols[i]}] = rel
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if len(out) != int(dimCount)*int(dimCount) {
		return nil, fmt.Errorf("interactions.tsv: got %d cells, want %d", len(out), int(dimCount)*int(dimCount))
	}
	return out, nil
}

// LoadExceptions reads flag-level overrides.
func LoadExceptions() ([]Exception, error) {
	b, err := os.ReadFile(testdata("exceptions.tsv"))
	if err != nil {
		return nil, err
	}
	return parseExceptions(b)
}

func parseExceptions(b []byte) ([]Exception, error) {
	sc := bufio.NewScanner(bytes.NewReader(b))
	var out []Exception
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.Split(line, "\t")
		if len(f) < 3 {
			return nil, fmt.Errorf("exceptions.tsv:%d: want name_a name_b relation [note]", lineNo)
		}
		rel, ok := parseRelation(f[2])
		if !ok {
			return nil, fmt.Errorf("exceptions.tsv:%d: unknown relation %q", lineNo, f[2])
		}
		note := ""
		if len(f) > 3 {
			note = f[3]
		}
		out = append(out, Exception{A: f[0], B: f[1], Relation: rel, Note: note})
	}
	return out, sc.Err()
}

// RejectedFlag is a name that must not return as a Spec.
type RejectedFlag struct {
	Name        string
	Replacement string
	Reason      string
}

// LoadRejected reads the cemetery table.
func LoadRejected() ([]RejectedFlag, error) {
	b, err := os.ReadFile(testdata("rejected.tsv"))
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(bytes.NewReader(b))
	var out []RejectedFlag
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.Split(line, "\t")
		if len(f) < 3 {
			return nil, fmt.Errorf("rejected.tsv: want name replacement reason, got %q", line)
		}
		out = append(out, RejectedFlag{Name: f[0], Replacement: f[1], Reason: strings.Join(f[2:], "\t")})
	}
	return out, sc.Err()
}

// RelationOf returns the interaction class for two primary flag names.
func RelationOf(a, b string, matrix map[[2]Dimension]Relation, ex []Exception) (Relation, error) {
	if a == b {
		return RelSameDim, nil
	}
	sa, ok := specByName(a)
	if !ok {
		return 0, fmt.Errorf("unknown spec %q", a)
	}
	sb, ok := specByName(b)
	if !ok {
		return 0, fmt.Errorf("unknown spec %q", b)
	}
	for _, e := range ex {
		if (e.A == a && e.B == b) || (e.A == b && e.B == a) {
			return e.Relation, nil
		}
	}
	if sa.Slot != "" && sa.Slot == sb.Slot {
		return RelSameDim, nil
	}
	rel, ok := matrix[[2]Dimension{sa.Dimension, sb.Dimension}]
	if !ok {
		return 0, fmt.Errorf("matrix missing %s×%s", sa.Dimension, sb.Dimension)
	}
	return rel, nil
}

func specByName(name string) (Spec, bool) {
	for _, s := range Specs() {
		if s.Name == name {
			return s, true
		}
	}
	return Spec{}, false
}
