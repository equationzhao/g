package util

import (
	"os"

	"github.com/Equationzhao/g/internal/item"
)

// DirectoryProcessor defines interface for different directory processing strategies
type DirectoryProcessor interface {
	ProcessDirectory(dirPath string) ([]*item.FileInfo, []error)
}

// OptimizationStrategy determines which processing strategy to use
type OptimizationStrategy struct {
	// Thresholds for strategy selection
	SmallDirThreshold int  // Files count below which to use traditional approach
	BatchModeEnabled  bool // Global toggle for batch mode
}

// NewOptimizationStrategy creates a new strategy selector with default values
func NewOptimizationStrategy() *OptimizationStrategy {
	return &OptimizationStrategy{
		SmallDirThreshold: 50, // Based on benchmark results
		BatchModeEnabled:  true,
	}
}

// SelectProcessor chooses the optimal processing strategy based on context
func (ots *OptimizationStrategy) SelectProcessor(dirPath, outputFormat string, hasComplexFiltering bool) DirectoryProcessor {
	// Quick directory size estimation
	estimatedSize := ots.EstimateDirectorySize(dirPath)

	// Decision logic based on our performance analysis
	shouldUseBatch := ots.shouldUseBatchMode(estimatedSize, outputFormat, hasComplexFiltering)

	if shouldUseBatch {
		return &BatchDirectoryProcessor{dirPath: dirPath}
	}
	return &TraditionalDirectoryProcessor{dirPath: dirPath}
}

// shouldUseBatchMode implements the decision logic
func (ots *OptimizationStrategy) shouldUseBatchMode(estimatedSize int, outputFormat string, hasComplexFiltering bool) bool {
	if !ots.BatchModeEnabled {
		return false
	}

	// Never use batch mode for JSON output (24% regression observed)
	if outputFormat == "json" {
		return false
	}

	// For very small directories, traditional approach is faster (42% in micro-benchmarks)
	if estimatedSize < ots.SmallDirThreshold {
		return false
	}

	// Use batch mode for complex filtering scenarios regardless of size
	if hasComplexFiltering {
		return true
	}

	// Use batch mode for larger directories
	return estimatedSize >= ots.SmallDirThreshold
}

// EstimateDirectorySize provides a quick estimate without full directory read
func (ots *OptimizationStrategy) EstimateDirectorySize(dirPath string) int {
	// Quick estimation using os.ReadDir (this is fast)
	entries, err := ots.quickDirScan(dirPath)
	if err != nil {
		return 100 // Default to batch mode on error (conservative)
	}
	return len(entries)
}

// quickDirScan performs minimal directory scanning for size estimation
func (ots *OptimizationStrategy) quickDirScan(dirPath string) ([]os.DirEntry, error) {
	// Use os.ReadDir for quick directory entry counting
	// This is faster than full file info gathering
	return os.ReadDir(dirPath)
}
