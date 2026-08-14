package parse

import (
	"bufio"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

func TestBudgetIsForty(t *testing.T) {
	if Budget != 40 {
		t.Fatalf("Budget=%d: raising the cap is a Key Decision (docs/rewrite-architecture.md §17), not a drive-by", Budget)
	}
}

func TestSpecsCountEqualsBudget(t *testing.T) {
	if n := len(Specs()); n != Budget {
		t.Fatalf("len(Specs())=%d, Budget=%d", n, Budget)
	}
}

func TestSpecsUniqueNamesAndSlots(t *testing.T) {
	names := map[string]struct{}{}
	shorts := map[rune]string{}
	for _, s := range Specs() {
		if s.Name == "" {
			t.Fatal("empty Spec.Name")
		}
		if _, ok := names[s.Name]; ok {
			t.Errorf("duplicate Spec.Name %q", s.Name)
		}
		names[s.Name] = struct{}{}
		if s.Dimension >= dimCount {
			t.Errorf("%s: invalid dimension", s.Name)
		}
		if s.Slot == "" {
			t.Errorf("%s: empty slot", s.Name)
		}
		if s.ConfigExempt == ConfigRequired && s.ConfigKey == "" {
			t.Errorf("%s: missing config key (mark ConfigMeta/ConfigCLIOnly to exempt)", s.Name)
		}
		if s.ConfigExempt != ConfigRequired && s.ConfigKey != "" {
			t.Errorf("%s: exempt flag must not also have ConfigKey %q", s.Name, s.ConfigKey)
		}
		if s.Short != 0 {
			if other, ok := shorts[s.Short]; ok {
				t.Errorf("short -%c claimed by %s and %s", s.Short, other, s.Name)
			}
			shorts[s.Short] = s.Name
		}
		for _, a := range s.Aliases {
			if len(a) == 2 && a[0] == '-' && a[1] != '-' {
				r := rune(a[1])
				if other, ok := shorts[r]; ok && other != s.Name {
					t.Errorf("short -%c claimed by %s and %s", r, other, s.Name)
				}
				shorts[r] = s.Name
			}
		}
		for _, a := range s.Aliases {
			if _, ok := names[a]; ok {
				t.Errorf("alias %q of %s collides with a primary name", a, s.Name)
			}
		}
	}
}

func TestMatrixComplete(t *testing.T) {
	m, err := LoadMatrix()
	if err != nil {
		t.Fatal(err)
	}
	for i := Dimension(0); i < dimCount; i++ {
		for j := Dimension(0); j < dimCount; j++ {
			if _, ok := m[[2]Dimension{i, j}]; !ok {
				t.Errorf("empty matrix cell %s×%s", i, j)
			}
		}
	}
}

func TestEverySpecPairHasRelation(t *testing.T) {
	m, err := LoadMatrix()
	if err != nil {
		t.Fatal(err)
	}
	ex, err := LoadExceptions()
	if err != nil {
		t.Fatal(err)
	}
	specs := Specs()
	for i, a := range specs {
		for _, b := range specs[i:] {
			if _, err := RelationOf(a.Name, b.Name, m, ex); err != nil {
				t.Errorf("%s × %s: %v", a.Name, b.Name, err)
			}
		}
	}
}

func TestExceptionsReferToPrimaryNames(t *testing.T) {
	ex, err := LoadExceptions()
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range ex {
		if _, ok := specByName(e.A); !ok {
			t.Errorf("exceptions.tsv: %q is not Spec.Name", e.A)
		}
		if _, ok := specByName(e.B); !ok {
			t.Errorf("exceptions.tsv: %q is not Spec.Name", e.B)
		}
		if e.Note == "" {
			t.Errorf("exceptions.tsv: %s × %s missing note", e.A, e.B)
		}
	}
}

func TestRejectedNotInSpecs(t *testing.T) {
	rej, err := LoadRejected()
	if err != nil {
		t.Fatal(err)
	}
	if len(rej) < 40 {
		t.Fatalf("rejected.tsv looks too short: %d rows", len(rej))
	}
	ours := AllNames()
	for _, r := range rej {
		if _, ok := ours[r.Name]; ok {
			t.Errorf("rejected name %q is registered in Specs(); remove it from rejected.tsv or drop the spec", r.Name)
		}
	}
}

func TestGNUReservedShortsNotReused(t *testing.T) {
	reserved := GNUReservedShorts()
	for _, s := range Specs() {
		if s.Short == 0 {
			continue
		}
		if why, ok := reserved[s.Short]; ok {
			t.Errorf("%s uses -%c which is reserved (%s)", s.Name, s.Short, why)
		}
	}
}

func TestConfigKeysCovered(t *testing.T) {
	keys := map[string]struct{}{}
	for _, s := range Specs() {
		if s.ConfigKey != "" {
			keys[s.ConfigKey] = struct{}{}
		}
	}
	for _, k := range ConfigOnlyKeys() {
		if k.Key == "" {
			t.Fatal("empty ConfigOnlyKey")
		}
		if _, ok := keys[k.Key]; ok {
			t.Errorf("config-only %q also belongs to a flag", k.Key)
		}
		keys[k.Key] = struct{}{}
	}
	// Known YAML block in the architecture doc must list every key.
	doc, err := os.ReadFile(repoFile("docs/rewrite-architecture.md"))
	if err != nil {
		t.Fatal(err)
	}
	block := yamlExample(string(doc))
	if block == "" {
		t.Fatal("docs/rewrite-architecture.md: missing YAML example in §5.2")
	}
	for k := range keys {
		if !strings.Contains(block, k+":") {
			t.Errorf("config key %q missing from architecture YAML example", k)
		}
	}
}

func TestManOptionsCoverSpecs(t *testing.T) {
	doc, err := os.ReadFile(repoFile("docs/rewrite-man.md"))
	if err != nil {
		t.Fatal(err)
	}
	start := "<!-- BEGIN OPTIONS -->"
	end := "<!-- END OPTIONS -->"
	s := string(doc)
	i := strings.Index(s, start)
	j := strings.Index(s, end)
	if i < 0 || j < 0 || j <= i {
		t.Fatal("docs/rewrite-man.md: missing BEGIN/END OPTIONS markers")
	}
	block := s[i:j]
	for _, spec := range Specs() {
		if !strings.Contains(block, spec.Name) {
			t.Errorf("docs/rewrite-man.md OPTIONS missing %s", spec.Name)
		}
	}
}

func TestRegistryCoversShorts(t *testing.T) {
	doc, err := os.ReadFile(repoFile("docs/flag-registry.md"))
	if err != nil {
		t.Fatal(err)
	}
	body := string(doc)
	for _, s := range Specs() {
		if s.Short == 0 {
			continue
		}
		cell := "- " + string(s.Short)
		if !strings.Contains(body, cell) && !strings.Contains(body, "`-"+string(s.Short)+"`") {
			t.Errorf("docs/flag-registry.md missing used short -%c", s.Short)
		}
	}
	for r := range GNUReservedShorts() {
		if !strings.Contains(body, "`-"+string(r)+"`") {
			t.Errorf("docs/flag-registry.md missing reserved short -%c", r)
		}
	}
}

func TestGeneratedDocsFresh(t *testing.T) {
	root := repoRoot()
	checks := []struct {
		rel  string
		want string
	}{
		{"docs/flag-registry.md", RenderFlagRegistry()},
		{"docs/rejected-flags.md", RenderRejected()},
		{"docs/rewrite-man.md", RenderManOptions()},
	}
	for _, c := range checks {
		got, err := os.ReadFile(filepath.Join(root, c.rel))
		if err != nil {
			t.Errorf("%s: %v (run go generate ./internal/parse)", c.rel, err)
			continue
		}
		if string(got) != c.want {
			t.Errorf("%s is stale; run go generate ./internal/parse", c.rel)
		}
	}
}

func TestArchitectureMentionsGates(t *testing.T) {
	b, err := os.ReadFile(repoFile("docs/rewrite-architecture.md"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, needle := range []string{
		"## 17. Flag 准入与交互制度",
		"docs/rejected-flags.md",
		"docs/flag-registry.md",
		"internal/parse/testdata/interactions.tsv",
	} {
		if !strings.Contains(s, needle) {
			t.Errorf("architecture doc missing %q", needle)
		}
	}
}

func TestNoControlCharsInNames(t *testing.T) {
	for _, s := range Specs() {
		for _, r := range s.Name {
			if unicode.IsControl(r) {
				t.Errorf("Spec.Name %q contains control rune", s.Name)
			}
		}
		if !utf8.ValidString(s.Name) {
			t.Errorf("Spec.Name %q is not valid UTF-8", s.Name)
		}
	}
}

func repoRoot() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func repoFile(rel string) string {
	return filepath.Join(repoRoot(), rel)
}

func yamlExample(doc string) string {
	const mark = "```yaml"
	i := strings.Index(doc, mark)
	if i < 0 {
		return ""
	}
	rest := doc[i+len(mark):]
	j := strings.Index(rest, "```")
	if j < 0 {
		return ""
	}
	return rest[:j]
}

func TestTSVNoSpacesInsteadOfTabs(t *testing.T) {
	// Guard against editors turning the matrix into aligned spaces.
	f, err := os.Open(testdata("interactions.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	rows := 0
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "#") || strings.TrimSpace(line) == "" {
			continue
		}
		if !strings.Contains(line, "\t") {
			t.Fatalf("interactions.tsv line has no tab: %q", line)
		}
		rows++
	}
	if rows != int(dimCount)+1 {
		t.Fatalf("interactions.tsv: got %d data rows (header+dims), want %d", rows, int(dimCount)+1)
	}
}
