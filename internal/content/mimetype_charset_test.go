package content

import "testing"

func TestSplitMimeAndCharset(t *testing.T) {
	cases := []struct {
		in, mime, charset string
	}{
		{"text/plain", "text/plain", ""},
		{"text/plain; charset=utf-8", "text/plain", "utf-8"},
		{"text/plain;charset=utf-8", "text/plain", "utf-8"},
		{"text/plain; foo", "text/plain", ""},
		{"text/plain;", "text/plain", ""},
	}
	for _, tc := range cases {
		mime, charset := splitMimeAndCharset(tc.in)
		if mime != tc.mime || charset != tc.charset {
			t.Fatalf("splitMimeAndCharset(%q) = (%q, %q), want (%q, %q)", tc.in, mime, charset, tc.mime, tc.charset)
		}
	}
}
