package sort

import "testing"

func TestFilevercmpSpecTable(t *testing.T) {
	pairs := [][2]string{
		{"", "a"},
		{"a", "a0"},
		{"a0", "a1"},
		{"a1", "a1a"},
		{"file2", "file10"},
		{"file00", "file0"},
		{"01", "1"},
		{"a.1", "a.2"},
		{"a.2", "a.10"},
		{"α2", "α10"},
	}
	for _, p := range pairs {
		if c := Filevercmp(p[0], p[1]); c >= 0 {
			t.Fatalf("Filevercmp(%q,%q)=%d want < 0", p[0], p[1], c)
		}
	}
}
