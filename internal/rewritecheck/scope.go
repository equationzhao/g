package rewritecheck

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Kind string

const (
	KindRewrite Kind = "rewrite"
	KindLegacy  Kind = "legacy"
)

func testdata(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "testdata", name)
}

func RepoRoot() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func LoadPackageKinds() (map[string]Kind, error) {
	b, err := os.ReadFile(testdata("packages.tsv"))
	if err != nil {
		return nil, err
	}
	out := map[string]Kind{}
	sc := bufio.NewScanner(bytes.NewReader(b))
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if line == "path\tkind" {
			continue
		}
		f := strings.Split(line, "\t")
		if len(f) != 2 {
			return nil, fmt.Errorf("packages.tsv:%d: want path\\tkind", lineNo)
		}
		switch Kind(f[1]) {
		case KindRewrite, KindLegacy:
			out[filepath.ToSlash(f[0])] = Kind(f[1])
		default:
			return nil, fmt.Errorf("packages.tsv:%d: kind must be rewrite or legacy", lineNo)
		}
	}
	return out, sc.Err()
}

func allowedThirdParty() []string {
	return []string{
		"gopkg.in/yaml.v3",
		"github.com/mattn/go-runewidth",
		"github.com/gobwas/glob",
		"github.com/itchyny/timefmt-go",
		"golang.org/x/sys",
		"golang.org/x/term",
	}
}

func bannedThirdParty() []string {
	return []string{
		"github.com/urfave/cli",
		"github.com/jedib0t/go-pretty",
		"github.com/gabriel-vasile/mimetype",
		"github.com/gookit/color",
		"github.com/syndtr/goleveldb",
		"github.com/sahilm/fuzzy",
		"github.com/agiledragon/gomonkey",
		"github.com/shirou/gopsutil",
		"github.com/saintfish/chardet",
		"github.com/Equationzhao/pathbeautify",
		"github.com/pkg/xattr",
		"github.com/alphadose/haxmap",
		"github.com/wk8/go-ordered-map",
		"github.com/valyala/bytebufferpool",
		"github.com/zeebo/xxh3",
		"github.com/acarl005/stripansi",
		"github.com/olekukonko/ts",
		"github.com/hako/durafmt",
		"github.com/stretchr/testify",
		"golang.org/x/exp",
	}
}

const ModulePath = "github.com/Equationzhao/g"
