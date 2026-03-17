package man

import (
	"compress/gzip"
	"os"
	"path/filepath"

	"github.com/Equationzhao/g/internal/cli"
)

func GenMan() {
	// man
	man, _ := os.Create(filepath.Join("man", "g.1.gz"))
	s, _ := cli.G.ToMan()
	// compress to gzip
	manGz := gzip.NewWriter(man)
	defer func() {
		_ = manGz.Close()
	}()
	_, _ = manGz.Write([]byte(s))
	_ = manGz.Flush()
}
