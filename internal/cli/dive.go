package cli

import (
	"github.com/Equationzhao/g/internal/filter"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/util"
)

// Dive performs optimized recursive directory traversal with reduced system calls
// and improved memory allocation patterns. This implementation is designed to
// improve performance through batch file information retrieval and better memory pre-allocation.
func Dive(parent string, depth, limit int, infos *util.Slice[*item.FileInfo], errSlice *util.Slice[error], itemFilter *filter.ItemFilter) {
	if limit > 0 && depth > limit {
		return
	}

	type dirState struct {
		path  string
		depth int
	}

	// Pre-allocate stack space to reduce memory allocations during recursion
	stack := make([]dirState, 0, 32) // Estimate reasonable depth
	stack = append(stack, dirState{path: parent, depth: depth})

	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]

		if limit > 0 && current.depth > limit {
			continue
		}

		// Use batch file info reader instead of individual entry.Info() calls
		batchReader, err := util.NewBatchFileInfo(current.path)
		if err != nil {
			errSlice.AppendTo(err)
			continue
		}

		// Retrieve file information in batch to reduce system calls
		fileInfos, errors := batchReader.GetFileInfos()
		for _, err := range errors {
			errSlice.AppendTo(err)
		}

		// Use accurate directory count for pre-allocation to further reduce memory allocations
		dirCount := batchReader.GetDirectoryCount()
		subDirs := make([]dirState, 0, dirCount)

		for _, info := range fileInfos {
			// Apply filters
			if !itemFilter.Match(info) {
				continue
			}

			// Set level information
			info.Level = current.depth
			infos.AppendTo(info)

			// Collect subdirectory information
			if info.IsDir() && (limit <= 0 || current.depth+1 <= limit) {
				subDirs = append(subDirs, dirState{
					path:  info.FullPath,
					depth: current.depth + 1,
				})
			}
		}

		// Add subdirectories to stack in reverse order to maintain traversal order
		for i := len(subDirs) - 1; i >= 0; i-- {
			stack = append(stack, subDirs[i])
		}
	}
}
