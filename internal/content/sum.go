package content

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"

	"github.com/Equationzhao/g/internal/align"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/render"
)

type SumType string

const (
	SumTypeMd5    SumType = "MD5"
	SumTypeSha1   SumType = "SHA1"
	SumTypeSha224 SumType = "SHA224"
	SumTypeSha256 SumType = "SHA256"
	SumTypeSha384 SumType = "SHA384"
	SumTypeSha512 SumType = "SHA512"
	SumTypeCRC32  SumType = "CRC32"
)

type SumEnabler struct{}

func (s SumEnabler) EnableSum(renderer *render.Renderer, sumTypes ...SumType) []ContentOption {
	options := make([]ContentOption, 0, len(sumTypes))
	factory := func(sumType SumType) ContentOption {
		return func(info *item.FileInfo) (string, string) {
			if info.IsDir() {
				return "", string(sumType)
			}

			var content []byte
			if content_, ok := info.Cache["content"]; ok {
				content = content_
			} else {
				file, err := os.Open(info.FullPath)
				if err != nil {
					return "", string(sumType)
				}
				defer func() {
					_ = file.Close()
				}()
				content, err = io.ReadAll(file)
				if err != nil {
					return "", string(sumType)
				}
				info.Cache["content"] = content
			}

			var hashed hash.Hash
			switch sumType {
			case SumTypeMd5:
				hashed = md5.New()
			case SumTypeSha1:
				hashed = sha1.New()
			case SumTypeSha224:
				hashed = sha256.New224()
			case SumTypeSha256:
				hashed = sha256.New()
			case SumTypeSha384:
				hashed = sha512.New384()
			case SumTypeSha512:
				hashed = sha512.New()
			case SumTypeCRC32:
				hashed = crc32.NewIEEE()
			}
			if _, err := io.Copy(hashed, bytes.NewReader(content)); err != nil {
				return "", string(sumType)
			}
			return renderer.Checksum(fmt.Sprintf("%x", hashed.Sum(nil))), string(sumType)
		}
	}
	for _, sumType := range sumTypes {
		align.RegisterHeaderFooter(string(sumType))
		options = append(options, factory(sumType))
	}
	return options
}
