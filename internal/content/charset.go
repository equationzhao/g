package content

import (
	"os"
	"strings"

	constval "github.com/Equationzhao/g/internal/global"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/render"
	"github.com/Equationzhao/g/internal/util"
	"github.com/gabriel-vasile/mimetype"
	"github.com/saintfish/chardet"
)

type CharsetEnabler struct{}

func NewCharsetEnabler() *CharsetEnabler {
	return &CharsetEnabler{}
}

const (
	Charset = constval.NameOfCharset
)

func (c *CharsetEnabler) Enable(renderer *render.Renderer) ContentOption {
	det := chardet.NewTextDetector()
	return func(info *item.FileInfo) (string, string) {
		// check cache
		if c, ok := info.Cache[Charset]; ok {
			return renderer.Charset(string(c)), Charset
		}
		// only text file has charset
		if !isTextFile(info) {
			return renderer.Charset("-"), Charset
		}
		// detect file type
		mtype, err := mimetype.DetectFile(info.FullPath)
		if err != nil {
			return renderer.Charset("failed_to_read"), Charset
		}
		// detect charset
		charset := detectCharset(mtype, info, det)
		info.Cache[Charset] = []byte(charset)
		return renderer.Charset(charset), Charset
	}
}

func isTextFile(info *item.FileInfo) bool {
	return !info.IsDir() &&
		!util.IsSymLink(info) &&
		info.Mode()&os.ModeNamedPipe == 0 &&
		info.Mode()&os.ModeSocket == 0
}

func detectCharset(mtype *mimetype.MIME, info *item.FileInfo, det *chardet.Detector) string {
	if tn := mtype.String(); strings.Contains(tn, ";") {
		return strings.SplitN(strings.SplitN(tn, ";", 2)[1], "=", 2)[1]
	}
	if p := mtype.Parent(); p != nil && strings.Contains(p.String(), "text") {
		content, err := readFileContent(info.FullPath)
		if err != nil {
			return err.Error()
		}
		best, err := det.DetectBest(content)
		if err != nil {
			return "failed_to_detect"
		}
		return best.Charset
	}
	return "-"
}

func readFileContent(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	content := make([]byte, 1024*1024)
	_, err = file.Read(content)
	return content, err
}
