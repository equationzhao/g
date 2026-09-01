package display

import (
	"bytes"
	"testing"
	"time"

	"github.com/Equationzhao/g/internal/item"
)

// setEntryContent injects a single ordered content field so the Zero printer
// emits a known string per entry, independent of which default metadata fields
// NewFileInfoWithOption populates.
func setEntryContent(t *testing.T, info *item.FileInfo, text string) {
	t.Helper()
	info.Set("test-entry", &ItemContent{No: 0, Content: StringContent(text)})
}

// TestZeroPrinterTerminatesEachEntryWithNUL guards the fix for #333: the Zero
// printer must write a NUL byte after every entry so consumers can split the
// stream safely. Before the fix it wrote each entry's content but no NUL at
// all, concatenating entries with no delimiter.
func TestZeroPrinterTerminatesEachEntryWithNUL(t *testing.T) {
	now := time.Now()
	entries := []string{"a", "b", "c"}
	items := make([]*item.FileInfo, len(entries))
	for i, name := range entries {
		fi := newTreeFileInfo(t, "/root/"+name, "/root", 1, false, now)
		setEntryContent(t, fi, name)
		items[i] = fi
	}

	orig := Output
	Output = &bytes.Buffer{}
	t.Cleanup(func() { Output = orig })

	NewZero().Print(items...)

	got := Output.(*bytes.Buffer).Bytes()
	want := []byte("a\x00b\x00c\x00")
	if !bytes.Equal(got, want) {
		t.Errorf("Zero output = % x (%q), want % x (%q)", got, got, want, want)
	}
}

// TestZeroPrinterSingleEntryStillTerminated confirms a single entry is still
// NUL-terminated, so a consumer reading one entry does not block waiting for a
// delimiter that never comes.
func TestZeroPrinterSingleEntryStillTerminated(t *testing.T) {
	now := time.Now()
	fi := newTreeFileInfo(t, "/root/only", "/root", 1, false, now)
	setEntryContent(t, fi, "only")

	orig := Output
	Output = &bytes.Buffer{}
	t.Cleanup(func() { Output = orig })

	NewZero().Print(fi)

	got := Output.(*bytes.Buffer).Bytes()
	if want := []byte("only\x00"); !bytes.Equal(got, want) {
		t.Errorf("single-entry Zero output = % x (%q), want % x (%q)", got, got, want, want)
	}
}
