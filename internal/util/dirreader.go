package util

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Equationzhao/g/internal/item"
)

// BatchFileInfo provides batch file information retrieval to reduce system calls
type BatchFileInfo struct {
	entries []fs.DirEntry
	parent  string
}

// NewBatchFileInfo creates a new batch file info reader
func NewBatchFileInfo(parent string) (*BatchFileInfo, error) {
	entries, err := os.ReadDir(parent)
	if err != nil {
		return nil, err
	}
	
	return &BatchFileInfo{
		entries: entries,
		parent:  parent,
	}, nil
}

// GetFileInfos retrieves file information in batch with pre-allocated memory and string optimizations
func (b *BatchFileInfo) GetFileInfos() ([]*item.FileInfo, []error) {
	// Pre-allocate slices to reduce memory reallocations
	infos := make([]*item.FileInfo, 0, len(b.entries))
	errors := make([]error, 0, 2) // Most cases won't have many errors
	
	// Pre-allocate string builder for path concatenation to avoid repeated allocations
	var pathBuilder strings.Builder
	pathBuilder.Grow(len(b.parent) + 64) // Estimate path length
	
	for _, entry := range b.entries {
		// Use DirEntry.Info() to get file information in one system call
		fileInfo, err := entry.Info()
		if err != nil {
			errors = append(errors, err)
			continue
		}
		
		// Optimize string concatenation to reduce memory allocations
		pathBuilder.Reset()
		pathBuilder.WriteString(b.parent)
		if !strings.HasSuffix(b.parent, string(filepath.Separator)) {
			pathBuilder.WriteByte(filepath.Separator)
		}
		pathBuilder.WriteString(entry.Name())
		fullPath := pathBuilder.String()
		
		// Use optimized constructor with pre-allocated memory reuse
		info, err := item.NewFileInfoWithOption(
			item.WithAbsPath(fullPath),
			item.WithFileInfo(fileInfo),
		)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		
		info.ParentPath = b.parent
		infos = append(infos, info)
	}
	
	return infos, errors
}

// GetCount returns the number of files in the directory for pre-allocation
func (b *BatchFileInfo) GetCount() int {
	return len(b.entries)
}

// HasSubdirectories checks if there are subdirectories for recursive optimization
func (b *BatchFileInfo) HasSubdirectories() bool {
	for _, entry := range b.entries {
		if entry.IsDir() {
			return true
		}
	}
	return false
}

// GetDirectoryCount returns the number of subdirectories for more accurate pre-allocation
func (b *BatchFileInfo) GetDirectoryCount() int {
	count := 0
	for _, entry := range b.entries {
		if entry.IsDir() {
			count++
		}
	}
	return count
}