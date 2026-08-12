package filter

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/osbased"
)

type testFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (t testFileInfo) Name() string       { return t.name }
func (t testFileInfo) Size() int64        { return t.size }
func (t testFileInfo) Mode() os.FileMode  { return t.mode }
func (t testFileInfo) ModTime() time.Time { return t.modTime }
func (t testFileInfo) IsDir() bool        { return t.mode.IsDir() }
func (t testFileInfo) Sys() any           { return nil }

func TestWhichTimeFiledAliases(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected any
	}{
		{name: "mod", input: "mod", expected: osbased.ModTime},
		{name: "modified", input: "modified", expected: osbased.ModTime},
		{name: "create", input: "create", expected: osbased.CreateTime},
		{name: "cr", input: "cr", expected: osbased.CreateTime},
		{name: "access", input: "access", expected: osbased.AccessTime},
		{name: "ac", input: "ac", expected: osbased.AccessTime},
		{name: "birth", input: "birth", expected: osbased.BirthTime},
		{name: "empty", input: "", expected: osbased.ModTime},
		{name: "unknown", input: "nope", expected: osbased.ModTime},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := WhichTimeFiled(tc.input)
			if got == nil {
				t.Fatalf("WhichTimeFiled(%q) returned nil", tc.input)
			}
			if reflect.ValueOf(got).Pointer() != reflect.ValueOf(tc.expected).Pointer() {
				t.Fatalf("WhichTimeFiled(%q) did not normalize correctly", tc.input)
			}
		})
	}
}

func TestBeforeTimeAndAfterTimeWithAlias(t *testing.T) {
	fi := &item.FileInfo{
		FileInfo: testFileInfo{name: "sample", modTime: time.Unix(100, 0)},
	}

	before := BeforeTime(time.Unix(200, 0), WhichTimeFiled("modified"))
	after := AfterTime(time.Unix(50, 0), WhichTimeFiled("modified"))

	if !before(fi) {
		t.Fatal("BeforeTime with modified alias should match")
	}
	if !after(fi) {
		t.Fatal("AfterTime with modified alias should match")
	}

	requireNoPanic := func(name string, fn func(*item.FileInfo) bool) {
		t.Helper()
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("%s panicked: %v", name, r)
			}
		}()
		_ = fn(fi)
	}

	requireNoPanic("BeforeTime(modified)", BeforeTime(time.Unix(200, 0), WhichTimeFiled("modified")))
	requireNoPanic("AfterTime(modified)", AfterTime(time.Unix(50, 0), WhichTimeFiled("modified")))
}
