package parse

import (
	"testing"

	"github.com/Equationzhao/g/internal/request"
)

func TestGNUShortsParse(t *testing.T) {
	tests := []struct {
		args  string
		check func(*testing.T, request.Request)
	}{
		{"-l", func(t *testing.T, r request.Request) {
			if !r.Long {
				t.Fatal("expected long")
			}
		}},
		{"-a", func(t *testing.T, r request.Request) {
			if r.Visibility != request.VisAll {
				t.Fatalf("vis=%v", r.Visibility)
			}
		}},
		{"-A", func(t *testing.T, r request.Request) {
			if r.Visibility != request.VisAlmostAll {
				t.Fatalf("vis=%v", r.Visibility)
			}
		}},
		{"-h", func(t *testing.T, r request.Request) {
			if !r.HumanReadable {
				t.Fatal("expected human")
			}
		}},
		{"-1", func(t *testing.T, r request.Request) {
			if r.Format != request.FormatOneline || !r.FormatSet {
				t.Fatalf("format=%v set=%v", r.Format, r.FormatSet)
			}
		}},
		{"-R", func(t *testing.T, r request.Request) {
			if !r.Recurse {
				t.Fatal("expected recurse")
			}
		}},
		{"-t", func(t *testing.T, r request.Request) {
			if r.Sort != request.SortTime {
				t.Fatalf("sort=%v", r.Sort)
			}
		}},
		{"-S", func(t *testing.T, r request.Request) {
			if r.Sort != request.SortSize {
				t.Fatalf("sort=%v", r.Sort)
			}
		}},
		{"-r", func(t *testing.T, r request.Request) {
			if !r.Reverse {
				t.Fatal("expected reverse")
			}
		}},
		{"-d", func(t *testing.T, r request.Request) {
			if !r.DirSelf {
				t.Fatal("expected dirself")
			}
		}},
		{"-F", func(t *testing.T, r request.Request) {
			if r.Classify != request.WhenAlways || !r.ClassifySet {
				t.Fatalf("classify=%v", r.Classify)
			}
		}},
		{"-lah", func(t *testing.T, r request.Request) {
			if !r.Long || r.Visibility != request.VisAll || !r.HumanReadable {
				t.Fatalf("combined -lah: %+v", r)
			}
		}},
		{"-a -A", func(t *testing.T, r request.Request) {
			if r.Visibility != request.VisAlmostAll {
				t.Fatal("-A should win")
			}
		}},
	}
	for _, tc := range tests {
		t.Run(tc.args, func(t *testing.T) {
			args := splitArgs(tc.args)
			r, err := Parse(args, nil)
			if err != nil {
				t.Fatal(err)
			}
			tc.check(t, r)
		})
	}
}

func TestZeroLastWins(t *testing.T) {
	r, err := Parse([]string{"--color=always", "-0"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.Color != request.WhenNever || !r.Zero {
		t.Fatalf("got color=%v zero=%v", r.Color, r.Zero)
	}
	if _, err := Parse([]string{"-0", "--color=always"}, nil); err == nil {
		t.Fatal("expected exit 2")
	}
}

func TestHelpNotH(t *testing.T) {
	r, err := Parse([]string{"-h"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.Help || !r.HumanReadable {
		t.Fatal("-h must be human-readable, not help")
	}
}

func splitArgs(s string) []string {
	if s == "" {
		return nil
	}
	return splitWS(s)
}

func splitWS(s string) []string {
	var out []string
	cur := ""
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
			continue
		}
		cur += string(s[i])
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}
