package cli

import (
	"github.com/Equationzhao/g/internal/display"
	"github.com/Equationzhao/g/internal/filter"
)

// getOutputFormat determines the output format being used
func getOutputFormat(printer display.Printer) string {
	switch printer.(type) {
	case *display.JsonPrinter:
		return "json"
	case *display.TablePrinter:
		return "table"
	default:
		return "standard"
	}
}

// hasComplexFilters determines if the item filter has complex filtering logic
func hasComplexFilters(itemFilter *filter.ItemFilter) bool {
	// Check if there are multiple filter functions or complex patterns
	// For now, we'll use a simple heuristic
	// This can be enhanced based on actual filter complexity

	// If the filter has been customized from default, consider it complex
	// This is a simplified check - in practice you'd want to examine
	// the actual filter functions being used
	return false // Placeholder - should be implemented based on actual filter analysis
}
