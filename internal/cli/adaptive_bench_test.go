package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Equationzhao/g/internal/util"
)

// setupAdaptiveTestDir creates test directories of different sizes
func setupAdaptiveTestDir(tb testing.TB, dirName string, numFiles int) string {
	testDir := filepath.Join(os.TempDir(), dirName)
	os.RemoveAll(testDir)
	os.MkdirAll(testDir, 0o755)

	for i := range numFiles {
		fileName := filepath.Join(testDir, "file_"+string(rune('A'+(i/26)))+string(rune('a'+(i%26)))+".txt")
		f, _ := os.Create(fileName)
		f.WriteString("test content")
		f.Close()
	}

	return testDir
}

// BenchmarkAdaptiveStrategy_SmallDir tests adaptive strategy on small directories
func BenchmarkAdaptiveStrategy_SmallDir(b *testing.B) {
	testDir := setupAdaptiveTestDir(b, "adaptive_small", 25) // Below 50-file threshold
	defer os.RemoveAll(testDir)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		strategy := util.NewOptimizationStrategy()
		processor := strategy.SelectProcessor(testDir, "standard", false)
		fileInfos, _ := processor.ProcessDirectory(testDir)
		_ = fileInfos
	}
}

// BenchmarkAdaptiveStrategy_LargeDir tests adaptive strategy on large directories
func BenchmarkAdaptiveStrategy_LargeDir(b *testing.B) {
	testDir := setupAdaptiveTestDir(b, "adaptive_large", 200) // Above 50-file threshold
	defer os.RemoveAll(testDir)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		strategy := util.NewOptimizationStrategy()
		processor := strategy.SelectProcessor(testDir, "standard", false)
		fileInfos, _ := processor.ProcessDirectory(testDir)
		_ = fileInfos
	}
}

// BenchmarkAdaptiveStrategy_JSON tests adaptive strategy with JSON output
func BenchmarkAdaptiveStrategy_JSON(b *testing.B) {
	testDir := setupAdaptiveTestDir(b, "adaptive_json", 200) // Large dir but JSON format
	defer os.RemoveAll(testDir)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		strategy := util.NewOptimizationStrategy()
		processor := strategy.SelectProcessor(testDir, "json", false) // Should use traditional
		fileInfos, _ := processor.ProcessDirectory(testDir)
		_ = fileInfos
	}
}

// TestAdaptiveStrategy_DecisionLogic verifies the strategy selection logic
func TestAdaptiveStrategy_DecisionLogic(t *testing.T) {
	strategy := util.NewOptimizationStrategy()

	// Test small directory - should use traditional
	smallDir := setupAdaptiveTestDir(t, "test_small", 25)
	defer os.RemoveAll(smallDir)

	// Debug info
	estimatedSize := strategy.EstimateDirectorySize(smallDir)
	t.Logf("Small dir estimated size: %d, threshold: 50", estimatedSize)

	processor := strategy.SelectProcessor(smallDir, "standard", false)
	if _, ok := processor.(*util.TraditionalDirectoryProcessor); !ok {
		t.Errorf("Expected TraditionalDirectoryProcessor for small directory (size=%d)", estimatedSize)
	}

	// Test large directory - should use batch
	largeDir := setupAdaptiveTestDir(t, "test_large", 100)
	defer os.RemoveAll(largeDir)

	estimatedSize = strategy.EstimateDirectorySize(largeDir)
	t.Logf("Large dir estimated size: %d, threshold: 50", estimatedSize)

	processor = strategy.SelectProcessor(largeDir, "standard", false)
	if _, ok := processor.(*util.BatchDirectoryProcessor); !ok {
		t.Errorf("Expected BatchDirectoryProcessor for large directory (size=%d)", estimatedSize)
	}

	// Test JSON output - should use traditional regardless of size
	processor = strategy.SelectProcessor(largeDir, "json", false)
	if _, ok := processor.(*util.TraditionalDirectoryProcessor); !ok {
		t.Error("Expected TraditionalDirectoryProcessor for JSON output")
	}
}
