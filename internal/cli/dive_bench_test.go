package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Equationzhao/g/internal/filter"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/util"
)

// setupBenchmarkDir creates test directory structure for benchmarking
func setupBenchmarkDir(b *testing.B, dirName string, numDirs, numFiles int) string {
	testDir := filepath.Join(os.TempDir(), dirName)
	os.RemoveAll(testDir)
	os.MkdirAll(testDir, 0o755)

	for i := 0; i < numDirs; i++ {
		subDir := filepath.Join(testDir, "subdir"+string(rune('0'+i%10)))
		os.MkdirAll(subDir, 0o755)

		for j := 0; j < numFiles; j++ {
			fileName := filepath.Join(subDir, "file"+string(rune('0'+j%10))+".txt")
			f, _ := os.Create(fileName)
			f.Close()
		}
	}

	return testDir
}

// BenchmarkDive_Original tests the performance of the original dive function
func BenchmarkDive_Original(b *testing.B) {
	testDir := setupBenchmarkDir(b, "bench_original", 10, 50)
	defer os.RemoveAll(testDir)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		infos := util.NewSlice[*item.FileInfo](100)
		errSlice := util.NewSlice[error](10)
		itemFilter := filter.NewItemFilter()

		// Use the original dive function for comparison
		benchmarkDiveOriginal(testDir, 1, -1, infos, errSlice, itemFilter)
	}
}

// BenchmarkDive_Optimized tests the performance of the optimized dive function
func BenchmarkDive_Optimized(b *testing.B) {
	testDir := setupBenchmarkDir(b, "bench_optimized", 10, 50)
	defer os.RemoveAll(testDir)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		infos := util.NewSlice[*item.FileInfo](100)
		errSlice := util.NewSlice[error](10)
		itemFilter := filter.NewItemFilter()

		// Use the optimized dive function
		Dive(testDir, 1, -1, infos, errSlice, itemFilter)
	}
}

// BenchmarkFileInfoBatch compares traditional vs batched file info retrieval
func BenchmarkFileInfoBatch(b *testing.B) {
	testDir := setupBenchmarkDir(b, "bench_batch", 1, 100)
	defer os.RemoveAll(testDir)

	b.Run("Traditional", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkFileInfoTraditional(testDir)
		}
	})

	b.Run("Batched", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkFileInfoBatched(testDir)
		}
	})
}

// benchmarkFileInfoTraditional simulates traditional file info retrieval (individual calls)
func benchmarkFileInfoTraditional(dir string) {
	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		entry.Info()                     // Individual system calls
		filepath.Join(dir, entry.Name()) // Individual path concatenation
	}
}

// benchmarkFileInfoBatched uses the new batch file info retrieval method
func benchmarkFileInfoBatched(dir string) {
	batchReader, _ := util.NewBatchFileInfo(dir)
	batchReader.GetFileInfos() // Batch processing
}

// benchmarkDiveOriginal implements the original dive logic for benchmark comparison
// This function replicates the original implementation before optimization
func benchmarkDiveOriginal(parent string, depth, limit int, infos *util.Slice[*item.FileInfo], errSlice *util.Slice[error], itemFilter *filter.ItemFilter) {
	if limit > 0 && depth > limit {
		return
	}

	type dirState struct {
		path  string
		depth int
	}

	stack := []dirState{{path: parent, depth: depth}}
	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]
		if limit > 0 && current.depth > limit {
			continue
		}
		dir, err := os.ReadDir(current.path)
		if err != nil {
			errSlice.AppendTo(err)
			continue
		}
		subDirs := make([]dirState, 0, len(dir)) // Original implementation uses total count instead of directory count
		for _, entry := range dir {
			f, err := entry.Info() // Individual system calls
			if err != nil {
				errSlice.AppendTo(err)
				continue
			}
			nowAbs := filepath.Join(current.path, f.Name()) // Individual path concatenation
			info, _ := item.NewFileInfoWithOption(item.WithAbsPath(nowAbs), item.WithFileInfo(f))

			if !itemFilter.Match(info) {
				continue
			}

			info.ParentPath = current.path
			info.Level = current.depth
			infos.AppendTo(info)
			if f.IsDir() && (limit <= 0 || current.depth+1 <= limit) {
				subDirs = append(subDirs, dirState{path: info.FullPath, depth: current.depth + 1})
			}
		}
		for i := len(subDirs) - 1; i >= 0; i-- {
			stack = append(stack, subDirs[i])
		}
	}
}
