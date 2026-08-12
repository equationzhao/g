package content

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"reflect"
	"testing"

	"github.com/Equationzhao/g/internal/item"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestName_checkDereferenceErr(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		wantSymlinks string
		panic        bool
	}{
		{
			name:         "not path PathError",
			err:          errors.New("not path PathError"),
			wantSymlinks: "not path PathError",
		},
		{
			name:         "path PathError",
			err:          &fs.PathError{Op: "lstat", Path: "nowhere", Err: errors.New("no such file or directory")},
			wantSymlinks: "nowhere",
		},
		{
			name:  "nil",
			err:   nil,
			panic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.panic {
						t.Errorf("checkDereferenceErr() panic = %v", r)
					}
				}
			}()
			var n Name
			if gotSymlinks := n.checkDereferenceErr(tt.err); gotSymlinks != tt.wantSymlinks {
				t.Errorf("checkDereferenceErr() = %v, want %v", gotSymlinks, tt.wantSymlinks)
			}
		})
	}
}

// TODO gomonkey seems not working ?
// func TestMountsOn(t *testing.T) {
// 	patch := gomonkey.NewPatches()
// 	defer patch.Reset()
// 	patch.ApplyFunc(disk.Partitions, func(all bool) ([]disk.PartitionStat, error) {
// 		return []disk.PartitionStat{
// 			{
// 				Device:     "/dev/sda1",
// 				Mountpoint: "/",
// 				Fstype:     "apfs",
// 				Opts:       []string{"rw", "relatime"},
// 			},
// 			{
// 				Device:     "/devfs",
// 				Mountpoint: "/dev",
// 				Fstype:     "apfs",
// 				Opts:       []string{"rw", "relatime"},
// 			},
// 		}, nil
// 	})
// 	tests := []struct {
// 		name string
// 		path string
// 		want string
// 	}{
// 		{
// 			name: "root",
// 			path: "/",
// 			want: "[/dev/sda1 (apfs)]",
// 		},
// 		{
// 			name: "dev",
// 			path: "/dev",
// 			want: "[/devfs (apfs)]",
// 		},
// 		{
// 			"not found",
// 			"/notfound",
// 			"",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := MountsOn(tt.path); got != tt.want {
// 				t.Errorf("MountsOn() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestStatistics_MarshalJSON(t *testing.T) {
	type fields struct {
		file uint64
		dir  uint64
		link uint64
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr error
	}{
		{
			name: "ok",
			fields: fields{
				file: 100,
				dir:  111,
				link: 123,
			},
			want:    []byte(`{"File":100,"Dir":111,"Link":123}`),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Statistics{}
			s.file.Add(tt.fields.file)
			s.dir.Add(tt.fields.dir)
			s.link.Add(tt.fields.link)
			got, err := s.MarshalJSON()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func Test_checkIfEmpty(t *testing.T) {
	patch := gomonkey.ApplyFuncReturn(os.Open, nil, io.EOF)
	defer patch.Reset()

	i := item.FileInfo{
		FullPath: "test",
	}
	assert.True(t, checkIfEmpty(&i))

	patch.ApplyFuncReturn(os.Open, nil, io.ErrUnexpectedEOF)
	assert.True(t, checkIfEmpty(&i))
}
