package sorter

import "testing"

func TestByTimeSupportsBirth(t *testing.T) {
	for _, name := range []string{"birth"} {
		if ByTimeAscend(name) == nil {
			t.Fatalf("ByTimeAscend(%q) returned nil", name)
		}
		if ByTimeDescend(name) == nil {
			t.Fatalf("ByTimeDescend(%q) returned nil", name)
		}
	}
}
