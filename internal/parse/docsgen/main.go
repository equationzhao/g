// Command docsgen writes the generated rewrite docs from Specs().
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Equationzhao/g/internal/parse"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		fatal(err)
	}
	// go generate runs with the package directory as cwd (internal/parse).
	docs := filepath.Join(root, "..", "..", "docs")
	writes := map[string]string{
		filepath.Join(docs, "flag-registry.md"):  parse.RenderFlagRegistry(),
		filepath.Join(docs, "rejected-flags.md"): parse.RenderRejected(),
		filepath.Join(docs, "rewrite-man.md"):    parse.RenderManOptions(),
	}
	for path, body := range writes {
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			fatal(err)
		}
		fmt.Println("wrote", path)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
