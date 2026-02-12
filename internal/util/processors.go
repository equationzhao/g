package util

import (
	"os"
	"path/filepath"

	"github.com/Equationzhao/g/internal/item"
)

// TraditionalDirectoryProcessor implements traditional file-by-file processing
// Optimized for small directories and scenarios where batch overhead is not beneficial
type TraditionalDirectoryProcessor struct {
	dirPath string
}

// ProcessDirectory processes directory using traditional approach optimized for small datasets
func (tdp *TraditionalDirectoryProcessor) ProcessDirectory(dirPath string) ([]*item.FileInfo, []error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, []error{err}
	}

	// Pre-allocate with exact size for small directories
	infos := make([]*item.FileInfo, 0, len(entries))
	errors := make([]error, 0, 1) // Most small directories won't have errors

	for _, entry := range entries {
		fileInfo, err := entry.Info()
		if err != nil {
			errors = append(errors, err)
			continue
		}

		// Use filepath.Join for simple path concatenation (faster for small datasets)
		fullPath := filepath.Join(dirPath, entry.Name())

		info, err := item.NewFileInfoWithOption(
			item.WithAbsPath(fullPath),
			item.WithFileInfo(fileInfo),
		)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		info.ParentPath = dirPath
		infos = append(infos, info)
	}

	return infos, errors
}

// BatchDirectoryProcessor implements batch processing with optimizations
// Designed for larger directories and complex operations
type BatchDirectoryProcessor struct {
	dirPath string
}

// ProcessDirectory processes directory using batch approach for larger datasets
func (bdp *BatchDirectoryProcessor) ProcessDirectory(dirPath string) ([]*item.FileInfo, []error) {
	// Use existing BatchFileInfo for larger directories
	batchReader, err := NewBatchFileInfo(dirPath)
	if err != nil {
		return nil, []error{err}
	}

	return batchReader.GetFileInfos()
}
