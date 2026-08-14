package rewritecheck

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestEveryInternalPackageClassified(t *testing.T) {
	kinds, err := LoadPackageKinds()
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(RepoRoot(), "internal")
	ents, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, e := range ents {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		rel := filepath.ToSlash(filepath.Join("internal", e.Name()))
		seen[rel] = true
		if _, ok := kinds[rel]; !ok {
			t.Errorf("%s is not in testdata/packages.tsv (rewrite or legacy)", rel)
		}
	}
	for rel := range kinds {
		if !seen[rel] {
			t.Errorf("packages.tsv lists %s but the directory is gone; update the table", rel)
		}
	}
}

func TestModulePathAndLicense(t *testing.T) {
	b, err := os.ReadFile(filepath.Join(RepoRoot(), "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	first := strings.SplitN(string(b), "\n", 2)[0]
	if first != "module "+ModulePath {
		t.Fatalf("go.mod module line %q, want %q", first, "module "+ModulePath)
	}
	lic, err := os.ReadFile(filepath.Join(RepoRoot(), "LICENSE"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(lic), "MIT License") {
		t.Fatal("LICENSE is not MIT")
	}
}

func TestRewritePackagesAreClean(t *testing.T) {
	kinds, err := LoadPackageKinds()
	if err != nil {
		t.Fatal(err)
	}
	for rel, kind := range kinds {
		if kind != KindRewrite {
			continue
		}
		dir := filepath.Join(RepoRoot(), filepath.FromSlash(rel))
		checkRewriteDir(t, dir)
	}
}

func TestJustfileAndCIHaveNoMonkeyGCFlags(t *testing.T) {
	for _, rel := range []string{"justfile", ".github/workflows/rewrite-gates.yml"} {
		b, err := os.ReadFile(filepath.Join(RepoRoot(), rel))
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		// Legacy justfile still has the gomonkey crutch; rewrite recipes must not.
		if rel == "justfile" {
			i := strings.Index(s, "rewrite-gates:")
			if i < 0 {
				t.Fatal("justfile missing rewrite-gates recipe")
			}
			s = s[i:]
			if j := strings.Index(s[1:], "\n\n"); j >= 0 {
				s = s[:j+1]
			}
		}
		if strings.Contains(s, "-gcflags=all=-l") || strings.Contains(s, "-gcflags all=-l") {
			t.Errorf("%s enables -gcflags=all=-l (gomonkey crutch)", rel)
		}
	}
}

func TestArchitectureStillListsSixDeps(t *testing.T) {
	b, err := os.ReadFile(filepath.Join(RepoRoot(), "docs/rewrite-architecture.md"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, dep := range allowedThirdParty() {
		if !strings.Contains(s, dep) {
			t.Errorf("architecture doc dropped allowed dep %s", dep)
		}
	}
}

func checkRewriteDir(t *testing.T, dir string) {
	t.Helper()
	fset := token.NewFileSet()
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "testdata" || d.Name() == "docsgen" {
				return fs.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		f, err := parser.ParseFile(fset, path, src, parser.ParseComments)
		if err != nil {
			t.Errorf("%s: parse: %v", rel(path), err)
			return nil
		}
		checkBuildTags(t, path, string(src))
		checkImports(t, path, f, strings.HasSuffix(path, "_test.go"))
		checkNoPackageState(t, path, f)
		checkNoInit(t, path, f)
		checkNoOSMutation(t, path, f)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func checkBuildTags(t *testing.T, path, src string) {
	for _, line := range strings.Split(src, "\n") {
		trim := strings.TrimSpace(line)
		if trim == "package "+packageName(path) {
			break
		}
		if !strings.HasPrefix(trim, "//go:build") && !strings.HasPrefix(trim, "// +build") {
			continue
		}
		for _, tag := range []string{"fuzzy", "mounts", "lite"} {
			if hasBuildToken(trim, tag) {
				t.Errorf("%s: rewrite packages cannot use feature build tag %q", rel(path), tag)
			}
		}
	}
}

func packageName(path string) string {
	return filepath.Base(filepath.Dir(path))
}

func hasBuildToken(line, tag string) bool {
	for _, f := range strings.FieldsFunc(line, func(r rune) bool {
		return r == ' ' || r == ',' || r == '(' || r == ')' || r == '!'
	}) {
		if f == tag {
			return true
		}
	}
	return false
}

func checkImports(t *testing.T, path string, f *ast.File, isTest bool) {
	for _, imp := range f.Imports {
		p, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			continue
		}
		if isBanned(p) {
			t.Errorf("%s: banned import %s", rel(path), p)
		}
		if isLegacyInternal(p) {
			t.Errorf("%s: rewrite package imports legacy %s", rel(path), p)
		}
		if isThirdParty(p) && !isAllowedThirdParty(p) {
			t.Errorf("%s: third-party import %s is not in the six-module allowlist", rel(path), p)
		}
		if isTest && isNetworkImport(p) {
			t.Errorf("%s: tests must not import %s", rel(path), p)
		}
	}
}

func checkNoPackageState(t *testing.T, path string, f *ast.File) {
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs := spec.(*ast.ValueSpec)
			for _, name := range vs.Names {
				if name.Name == "_" {
					continue
				}
				if strings.HasPrefix(name.Name, "Err") {
					continue
				}
				if isLookupTable(vs) {
					continue
				}
				t.Errorf("%s:%s: package-level var %s (rewrite packages have no mutable package state; Err* and immutable lookup tables only)", rel(path), name.Name, name.Name)
			}
		}
	}
}

func checkNoInit(t *testing.T, path string, f *ast.File) {
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || fn.Name.Name != "init" {
			continue
		}
		if fn.Body != nil && len(fn.Body.List) == 0 {
			continue
		}
		t.Errorf("%s: init() is forbidden in rewrite packages (embed only, and not here)", rel(path))
	}
}

func checkNoOSMutation(t *testing.T, path string, f *ast.File) {
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			if sel, ok := x.Fun.(*ast.SelectorExpr); ok {
				if id, ok := sel.X.(*ast.Ident); ok && id.Name == "os" {
					switch sel.Sel.Name {
					case "Chdir", "Chmod", "Remove", "RemoveAll", "Rename":
						// tests may write temp dirs; only Chdir is banned everywhere
					}
					if sel.Sel.Name == "Chdir" {
						t.Errorf("%s: os.Chdir is forbidden", rel(path))
					}
				}
			}
		case *ast.AssignStmt:
			for _, lhs := range x.Lhs {
				if isOSArgs(lhs) {
					t.Errorf("%s: assignment to os.Args is forbidden", rel(path))
				}
			}
		}
		return true
	})
}

func isLookupTable(vs *ast.ValueSpec) bool {
	if vs.Type != nil {
		return isArrayOfBasic(vs.Type)
	}
	if len(vs.Values) != 1 {
		return false
	}
	cl, ok := vs.Values[0].(*ast.CompositeLit)
	if !ok {
		return false
	}
	return isArrayOfBasic(cl.Type)
}

func isArrayOfBasic(e ast.Expr) bool {
	switch t := e.(type) {
	case *ast.ArrayType:
		if t.Len == nil {
			return false // slice
		}
		id, ok := t.Elt.(*ast.Ident)
		return ok && isBasicIdent(id.Name)
	default:
		return false
	}
}

func isBasicIdent(name string) bool {
	switch name {
	case "bool", "string", "byte", "rune",
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return true
	default:
		return false
	}
}

func isOSArgs(e ast.Expr) bool {
	sel, ok := e.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	id, ok := sel.X.(*ast.Ident)
	return ok && id.Name == "os" && sel.Sel.Name == "Args"
}

func isBanned(p string) bool {
	for _, b := range bannedThirdParty() {
		if p == b || strings.HasPrefix(p, b+"/") {
			return true
		}
	}
	return false
}

func isAllowedThirdParty(p string) bool {
	for _, a := range allowedThirdParty() {
		if p == a || strings.HasPrefix(p, a+"/") {
			return true
		}
	}
	return false
}

func isThirdParty(p string) bool {
	if !strings.Contains(p, ".") {
		return false // stdlib
	}
	if strings.HasPrefix(p, ModulePath+"/") || p == ModulePath {
		return false
	}
	return true
}

func isLegacyInternal(p string) bool {
	const prefix = ModulePath + "/internal/"
	if !strings.HasPrefix(p, prefix) {
		return false
	}
	rest := strings.TrimPrefix(p, prefix)
	top := rest
	if i := strings.IndexByte(rest, '/'); i >= 0 {
		top = rest[:i]
	}
	kinds, err := LoadPackageKinds()
	if err != nil {
		return false
	}
	k := kinds["internal/"+top]
	return k == KindLegacy
}

func isNetworkImport(p string) bool {
	switch p {
	case "net", "net/http", "net/smtp", "net/rpc", "net/http/httptest":
		return true
	default:
		return false
	}
}

func rel(path string) string {
	r, err := filepath.Rel(RepoRoot(), path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(r)
}
