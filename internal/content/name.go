package content

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"unicode"

	"github.com/Equationzhao/g/internal/display"
	"github.com/Equationzhao/g/internal/osbased"

	"github.com/Equationzhao/g/internal/global"

	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/render"
	"github.com/Equationzhao/g/internal/theme"
	"github.com/Equationzhao/g/internal/util"
	"github.com/valyala/bytebufferpool"
)

type Name struct {
	icon, classify, fileType, fullPath, noDeference, hyperLink, mounts, json bool
	statistics                                                               *Statistics
	relativeTo                                                               string
	Quote                                                                    string
	QuoteStatus                                                              int8 // >1 always quote || =0 default || <0 never quote
	Party                                                                    bool
	partyIndex                                                               atomic.Uint64
}

func neverQuote(qs int8) bool {
	return qs < 0
}

func alwaysQuote(qs int8) bool {
	return qs > 0
}

type Statistics struct {
	file, dir, link atomic.Uint64
}

func (s *Statistics) MarshalJSON() ([]byte, error) {
	Export := struct {
		File, Dir, Link uint64
	}{
		File: s.file.Load(),
		Dir:  s.dir.Load(),
		Link: s.link.Load(),
	}
	return json.Marshal(Export)
}

func (s *Statistics) Reset() {
	s.file.Store(0)
	s.dir.Store(0)
	s.link.Store(0)
}

func (s *Statistics) String() string {
	return fmt.Sprintf("%d file(s), %d dir(s), %d link(s)", s.file.Load(), s.dir.Load(), s.link.Load())
}

func (n *Name) SetNoDeference() *Name {
	n.noDeference = true
	return n
}

func (n *Name) UnsetNoDeference() *Name {
	n.noDeference = false
	return n
}

func (n *Name) SetHyperlink() *Name {
	n.hyperLink = true
	return n
}

func (n *Name) UnsetHyperlink() *Name {
	n.hyperLink = false
	return n
}

func (n *Name) SetJson() *Name {
	n.json = true
	return n
}

func (n *Name) UnsetJson() *Name {
	n.json = false
	return n
}

func (n *Name) FullPath() bool {
	return n.fullPath
}

func (n *Name) SetFullPath() *Name {
	n.fullPath = true
	return n
}

func (n *Name) UnsetFullPath() *Name {
	n.fullPath = false
	return n
}

func (n *Name) RelativeTo() string {
	return n.relativeTo
}

func (n *Name) SetRelativeTo(relativeTo string) {
	n.relativeTo = relativeTo
}

func (n *Name) Statistics() *Statistics {
	return n.statistics
}

func (n *Name) SetStatistics(Statistics *Statistics) *Name {
	n.statistics = Statistics
	return n
}

func (n *Name) SetQuoteString(quote string) *Name {
	n.Quote = quote
	return n
}

func (n *Name) SetQuote() *Name {
	n.QuoteStatus = 1
	return n
}

func (n *Name) UnsetQuote() *Name {
	n.QuoteStatus = -1
	return n
}

func (n *Name) SetParty() *Name {
	n.Party = true
	return n
}

func (n *Name) UnsetParty() *Name {
	n.Party = false
	return n
}

func (n *Name) UnsetIcon() *Name {
	n.icon = false
	return n
}

func (n *Name) UnsetClassify() *Name {
	n.classify = false
	return n
}

func (n *Name) UnsetFileType() *Name {
	n.fileType = false
	return n
}

func (n *Name) SetIcon() *Name {
	n.icon = true
	return n
}

func (n *Name) SetClassify() *Name {
	n.classify = true
	return n
}

// SetFileType set file type
// should set classify first
// if classify is false, fileType will be ignored
func (n *Name) SetFileType() *Name {
	n.fileType = true
	return n
}

func (n *Name) SetMounts() *Name {
	n.mounts = true
	return n
}

func (n *Name) UnsetMounts() *Name {
	n.mounts = false
	return n
}

func NewNameEnabler() *Name {
	return &Name{}
}

const NameName = global.NameOfName

func makeLink(abs, name string) string {
	return util.MakeLink(abs, name)
}

func checkIfEmpty(info *item.FileInfo) bool {
	f, err := os.Open(info.FullPath)
	if err != nil {
		return true
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	return err == io.EOF
}

/*
Enable
color + icon + file://quote+filename/relative-name+quote + classify + color-end + dereference + mounts
color: filetype->filename->fileext->file
*/

func partyColor(index int) string {
	const goldenAngle = 137.507764 // 180° * (3 - sqrt(5))
	hue := math.Mod(float64(index)*goldenAngle, 360)
	r, g, b := theme.HslToRgb(hue, 1.0, 0.5)

	switch theme.ColorLevel {
	case theme.TrueColor:
		rgbStr, _ := theme.RGB(r, g, b)
		return rgbStr
	case theme.C256:
		return theme.RGBTo256(r, g, b)
	case theme.Ascii:
		return theme.RGBToBasic(r, g, b)
	default:
		return ""
	}
}

func (n *Name) Enable(renderer *render.Renderer) ContentOption {
	return func(info *item.FileInfo) (stringContent, funcName string) {
		name, color, icon, classify, mounts := info.Name(), "", "", "", ""
		dereference := bytebufferpool.Get()
		defer bytebufferpool.Put(dereference)
		mode := info.Mode()
		underline, bold, italics, faint, blink := false, false, false, false, false

		if info.IsDir() {
			if n.statistics != nil {
				n.statistics.dir.Add(1)
			}

			empty := checkIfEmpty(info)
			style := renderer.Dir(name, empty)

			if n.icon {
				icon = style.Icon
			}
			if n.classify {
				classify = "/"
			}
			color, underline, bold, italics, faint, blink = style.Color, style.Underline, style.Bold, style.Italics, style.Faint, style.Blink
		} else if util.IsSymLinkMode(mode) || osbased.IsMacOSAlias(info.FullPath) {
			if n.statistics != nil {
				n.statistics.link.Add(1)
			}
			style := renderer.Symlink()
			if n.icon {
				icon = style.Icon
			}
			color, underline, bold, italics, faint, blink = style.Color, style.Underline, style.Bold, style.Italics, style.Faint, style.Blink
			if n.classify {
				classify = "@"
			}
			// color + arrow + color-end + color + path + color-end
			if !n.noDeference {
				if n.json { // "dereference": "symlinks"
					symlinks, err := util.Evallinks(info.FullPath)
					if err != nil {
						info.Meta.Set("dereference_err", &display.ItemContent{Content: display.StringContent(err.Error())})
						symlinks = n.checkDereferenceErr(err)
					}
					info.Meta.Set("dereference", &display.ItemContent{Content: display.StringContent(symlinks)})
				} else { // => "symlinks"
					arrowStyle := renderer.SymlinkArrow()
					_, _ = dereference.WriteString(arrowStyle.Color)
					checkNameDisplayEffect(arrowStyle, dereference)
					_, _ = dereference.WriteString(arrowStyle.Icon)
					_, _ = dereference.WriteString(renderer.Colorend())
					symlinks, err := util.Evallinks(info.FullPath)
					var linkStyle theme.Style
					dereferenceMounts := ""
					if err != nil {
						symlinks = n.checkDereferenceErr(err)
						linkStyle = renderer.SymlinkBroken()
					} else {
						symlinks, linkStyle = n.getSymlink(info, symlinks, linkStyle, renderer)
						if n.mounts {
							dereferenceMounts = MountsOn(symlinks)
						}
					}
					_, _ = dereference.WriteString(linkStyle.Color)
					checkNameDisplayEffect(linkStyle, dereference)
					hasQuote := false
					// if the name contains space and QuoteStatus >=0, add quote
					if strings.ContainsFunc(symlinks, contains) {
						if !neverQuote(n.QuoteStatus) {
							hasQuote = true
							_, _ = dereference.WriteString(n.Quote)
						}
					} else {
						// no space, but QuoteStatus == 1
						if alwaysQuote(n.QuoteStatus) {
							hasQuote = true
							_, _ = dereference.WriteString(n.Quote)
						}
					}
					_, _ = dereference.WriteString(symlinks)
					if hasQuote {
						_, _ = dereference.WriteString(n.Quote)
					}
					_, _ = dereference.WriteString(renderer.Colorend())
					if dereferenceMounts != "" {
						if n.json {
							info.Meta.Set("dereference_mounts", &display.ItemContent{Content: display.StringContent(dereferenceMounts)})
						} else {
							_ = dereference.WriteByte(' ')
							_, _ = dereference.WriteString(renderer.Mounts(dereferenceMounts))
						}
					}
				}
			}
		} else {
			if n.statistics != nil {
				n.statistics.file.Add(1)
			}
			if mode&os.ModeNamedPipe != 0 {
				style := renderer.Pipe()
				if n.icon {
					icon = style.Icon
				}
				if n.classify {
					classify = "|"
				}
				color, underline, bold, italics, faint, blink = style.Color, style.Underline, style.Bold, style.Italics, style.Faint, style.Blink
			} else if mode&os.ModeSocket != 0 {
				style := renderer.Socket()
				if n.icon {
					icon = style.Icon
				}
				if n.classify {
					classify = "="
				}
				color, underline, bold, italics, faint, blink = style.Color, style.Underline, style.Bold, style.Italics, style.Faint, style.Blink
			} else if mode&os.ModeDevice != 0 {
				s := theme.Style{}
				if mode&os.ModeCharDevice != 0 {
					s = renderer.Char()
				} else {
					s = renderer.Device()
				}

				if n.icon {
					icon = s.Icon
				}
				color, underline, bold, italics, faint, blink = s.Color, s.Underline, s.Bold, s.Italics, s.Faint, s.Blink
			} else {
				if s, ok := renderer.ByName(name); ok {
					if n.icon {
						icon = s.Icon
					}
					color, underline, bold, italics, faint, blink = s.Color, s.Underline, s.Bold, s.Italics, s.Faint, s.Blink
				} else {
					s, ok = renderer.ByExt(name)
					if ok {
						if n.icon {
							icon = s.Icon
						}
						color, underline, bold, italics, faint, blink = s.Color, s.Underline, s.Bold, s.Italics, s.Faint, s.Blink
					} else {
						if strings.HasPrefix(name, ".") {
							s = renderer.HiddenFile()
						} else {
							s = renderer.File()
						}
						if n.icon {
							icon = s.Icon
						}
						color, underline, bold, italics, faint, blink = s.Color, s.Underline, s.Bold, s.Italics, s.Faint, s.Blink
					}
				}
			}
		}

		exe := util.IsExecutableMode(mode) && !util.IsSymLinkMode(mode) && !info.IsDir() && mode&os.ModeNamedPipe == 0 && mode&os.ModeSocket == 0
		if n.classify {
			if classify == "" && (!n.fileType) && exe {
				classify = "*"
			}
		}
		if exe {
			s := renderer.Executable()
			color, underline, bold, italics, faint, blink = s.Color, s.Underline, s.Bold, s.Italics, s.Faint, s.Blink
		}

		if n.mounts {
			mounts = MountsOn(info.FullPath)
		}

		if n.relativeTo != "" {
			relativePath, err := filepath.Rel(n.relativeTo, info.FullPath)
			if err == nil {
				name = relativePath
			}
		} else if n.fullPath {
			name = info.FullPath
		}

		name = util.Escape(name)

		if n.Party && theme.ColorLevel != theme.None {
			color = partyColor(int(n.partyIndex.Add(1) - 1))
		}

		b := bytebufferpool.Get()
		defer bytebufferpool.Put(b)
		if color != "" {
			_, _ = b.WriteString(color)
		}
		if icon != "" {
			if n.json {
				info.Meta.Set("icon", &display.ItemContent{Content: display.StringContent(icon)})
			} else {
				_, _ = b.WriteString(icon)
				_ = b.WriteByte(' ')
			}
		}
		checkNameDisplayEffect(theme.Style{
			Underline: underline,
			Bold:      bold,
			Faint:     faint,
			Italics:   italics,
			Blink:     blink,
		}, b)
		hasQuote := false
		// when the name contains space:
		// if json == true:
		// 		default:no quote, alwaysQuote==true:quote
		// if json == false:
		//		default:quote, neverQuote==true:no quote
		hasQuote = alwaysQuote(n.QuoteStatus) || (!n.json && strings.ContainsFunc(name, contains) && !neverQuote(n.QuoteStatus))
		if hasQuote {
			if n.json {
				_, _ = dereference.WriteString(n.Quote)
			} else {
				_, _ = b.WriteString(n.Quote)
			}
		}
		if n.hyperLink {
			if n.json {
				info.Meta.Set("hyperlink", &display.ItemContent{Content: display.StringContent("file://" + info.FullPath)})
				_, _ = b.WriteString(name)
			} else {
				_, _ = b.WriteString(makeLink("file://"+info.FullPath, name))
			}
		} else {
			_, _ = b.WriteString(name)
		}
		if hasQuote {
			_, _ = b.WriteString(n.Quote)
		}
		if classify == "@" {
			_, _ = b.WriteString(classify)
		}
		_, _ = b.WriteString(renderer.Colorend())
		if classify != "" && classify != "@" {
			_, _ = b.WriteString(classify)
		}
		if d := dereference.String(); d != "" {
			_, _ = b.WriteString(d)
		}
		if mounts != "" {
			if n.json {
				info.Meta.Set("mounts", &display.ItemContent{Content: display.StringContent(mounts)})
			} else {
				_ = b.WriteByte(' ')
				_, _ = b.WriteString(renderer.Mounts(mounts))
			}
		}
		return b.String(), NameName
	}
}

func checkNameDisplayEffect(style theme.Style, buffer *bytebufferpool.ByteBuffer) {
	if style.Underline {
		_, _ = buffer.WriteString(global.Underline)
	}
	if style.Bold {
		_, _ = buffer.WriteString(global.Bold)
	}
	if style.Italics {
		_, _ = buffer.WriteString(global.Italics)
	}
	if style.Faint {
		_, _ = buffer.WriteString(global.Faint)
	}
	if style.Blink {
		_, _ = buffer.WriteString(global.Blink)
	}
}

func (n *Name) getSymlink(info *item.FileInfo, symlinks string, linkStyle theme.Style, renderer *render.Renderer) (string, theme.Style) {
	stat, err := os.Stat(symlinks)
	if err == nil {
		if stat.IsDir() {
			empty := checkIfEmpty(info)
			linkStyle = renderer.Dir(stat.Name(), empty)
		} else if mode := stat.Mode(); util.IsSymLinkMode(mode) {
			linkStyle = renderer.Symlink()
		} else if mode&os.ModeNamedPipe != 0 {
			linkStyle = renderer.Pipe()
		} else if mode&os.ModeSocket != 0 {
			linkStyle = renderer.Socket()
		} else if mode&os.ModeDevice != 0 {
			if mode&os.ModeCharDevice != 0 {
				linkStyle = renderer.Char()
			} else {
				linkStyle = renderer.Device()
			}
		} else if s, ok := renderer.ByName(stat.Name()); ok {
			linkStyle = s
		} else if s, ok = renderer.ByExt(stat.Name()); ok {
			linkStyle = s
		} else {
			linkStyle = renderer.SymlinkDereference()
		}
	}

	if n.relativeTo != "" {
		symlinksRel, err := filepath.Rel(n.relativeTo, symlinks)
		if err == nil {
			symlinks = symlinksRel
		}
	}
	return symlinks, linkStyle
}

// checkDereferenceErr checks if the error is a *fs.PathError and if it is, it returns the path of the error. Otherwise, it returns the error message.
// err should not be nil
func (n *Name) checkDereferenceErr(err error) (symlinks string) {
	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		if n.relativeTo != "" {
			if symlinksRel, err := filepath.Rel(n.relativeTo, pathErr.Path); err == nil {
				pathErr.Path = symlinksRel
			}
		}
		symlinks = pathErr.Path
	} else {
		symlinks = err.Error()
	}
	return symlinks
}

func contains(r rune) bool {
	return unicode.IsSpace(r)
}
