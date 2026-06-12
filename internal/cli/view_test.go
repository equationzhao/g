package cli

import "testing"

func TestNormalizeTimeType(t *testing.T) {
	tests := map[string]string{
		"mod":      "mod",
		"modified": "mod",
		"create":   "create",
		"cr":       "create",
		"access":   "access",
		"ac":       "access",
		"birth":    "birth",
		"all":      "",
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			got, ok := normalizeTimeType(input)
			if !ok {
				t.Fatalf("normalizeTimeType(%q) reported false", input)
			}
			if got != want {
				t.Fatalf("normalizeTimeType(%q) = %q, want %q", input, got, want)
			}
		})
	}
}
