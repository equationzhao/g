package content

import (
	"os"
	"strings"

	"github.com/Equationzhao/g/internal/align"

	constval "github.com/Equationzhao/g/internal/global"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/render"
	"github.com/Equationzhao/g/internal/util"
	"github.com/gabriel-vasile/mimetype"
)

type MimeFileTypeEnabler struct {
	ParentOnly bool
}

func NewMimeFileTypeEnabler() *MimeFileTypeEnabler {
	return &MimeFileTypeEnabler{
		ParentOnly: false,
	}
}

const (
	MimeTypeName       = constval.NameOfMIME
	ParentMimeTypeName = "Parent-Mime-type"
)

func (e *MimeFileTypeEnabler) Enable(renderer *render.Renderer) ContentOption {
	align.RegisterHeaderFooter(MimeTypeName, ParentMimeTypeName)
	return func(info *item.FileInfo) (string, string) {
		res, returnName := func() (string, string) {
			tn := ""
			returnName := MimeTypeName
			if e.ParentOnly {
				returnName = ParentMimeTypeName
			}
			if c, ok := info.Cache[MimeTypeName]; ok {
				tn = string(c)
			} else {
				if info.IsDir() {
					tn = "directory"
					return tn, returnName
				}
				if util.IsSymLink(info) {
					tn = "symlink"
					return tn, returnName
				}
				if info.Mode()&os.ModeNamedPipe != 0 {
					tn = "named_pipe"
					return tn, returnName
				}
				if info.Mode()&os.ModeSocket != 0 {
					tn = "socket"
					return tn, returnName
				}
				if m, ok := info.Cache[MimeTypeName]; ok {
					info.Cache[Charset] = m
					return string(m), returnName
				}

				file, err := os.Open(info.FullPath)
				if err != nil {
					return "failed_to_read", returnName
				}
				// nolint
				defer file.Close()
				mtype, err := mimetype.DetectReader(file)
				if err != nil {
					return err.Error(), returnName
				}
				tn = mtype.String()
			}

			if e.ParentOnly {
				tn = strings.SplitN(tn, "/", 2)[0]
			}

			tn, charset := splitMimeAndCharset(tn)
			if charset != "" {
				info.Cache[Charset] = []byte(charset)
			}

			return tn, returnName
		}()
		return renderer.Mime(res), returnName
	}
}

// splitMimeAndCharset strips an optional "; charset=..." suffix.
// Malformed parameters are ignored instead of panicking.
func splitMimeAndCharset(tn string) (mimeType, charset string) {
	if !strings.Contains(tn, ";") {
		return tn, ""
	}
	parts := strings.SplitN(tn, ";", 2)
	mimeType = strings.TrimSpace(parts[0])
	if len(parts) < 2 {
		return mimeType, ""
	}
	param := strings.TrimSpace(parts[1])
	kv := strings.SplitN(param, "=", 2)
	if len(kv) != 2 {
		return mimeType, ""
	}
	if strings.EqualFold(strings.TrimSpace(kv[0]), "charset") {
		return mimeType, strings.TrimSpace(kv[1])
	}
	return mimeType, ""
}
