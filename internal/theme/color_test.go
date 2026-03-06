package theme

import (
	"reflect"
	"testing"

	constval "github.com/Equationzhao/g/internal/global"
	"github.com/gookit/color"
)

func TestHexToRgb(t *testing.T) {
	type args struct {
		hex string
	}
	tests := []struct {
		name    string
		args    args
		wantRgb []uint8
	}{
		{
			name:    "black",
			args:    args{hex: "#000000"},
			wantRgb: []uint8{0, 0, 0},
		},
		{
			name:    "white",
			args:    args{hex: "#ffffff"},
			wantRgb: []uint8{255, 255, 255},
		},
		{
			name:    "red",
			args:    args{hex: "#ff0000"},
			wantRgb: []uint8{255, 0, 0},
		},
		{
			name:    "3 digits",
			args:    args{hex: "#f00"},
			wantRgb: []uint8{255, 0, 0},
		},
		{
			name:    "#0x",
			args:    args{hex: "#0xff0000"},
			wantRgb: []uint8{255, 0, 0},
		},
		{
			name:    "empty",
			args:    args{hex: ""},
			wantRgb: []uint8{0, 0, 0},
		},
		{
			name:    "invalid",
			args:    args{hex: "invalid"},
			wantRgb: []uint8{0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRgb := HexToRgb(tt.args.hex); !reflect.DeepEqual(gotRgb, tt.wantRgb) {
				t.Errorf("HexToRgb() = %v, want %v", gotRgb, tt.wantRgb)
			}
		})
	}
}

func TestBasicConverts(t *testing.T) {
	for c := range basicColor2str {
		_, _, _ = BasicToRGBInt(c)
		_ = BasicToRGB(c)
		_ = BasicTo256(c)
	}
}

func Test256Converts(t *testing.T) {
	for i := range 256 {
		_ = c256ToBasic(uint8(i))
		_ = c256ToRGB(uint8(i))
	}
}

func TestRGBConverts(t *testing.T) {
	for r := range 10 {
		for g := range 10 {
			for b := range 10 {
				_ = RGBToBasic(uint8(r), uint8(g), uint8(b))
				_ = RGBTo256(uint8(r), uint8(g), uint8(b))
			}
		}
	}
}

func TestConvertColorIfGreaterThanExpect(t *testing.T) {
	type args struct {
		to  color.Level
		src string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "none",
			args:    args{to: None, src: constval.Black},
			want:    "",
			wantErr: false,
		},
		{
			name:    "basic2basic",
			args:    args{to: Ascii, src: constval.Black},
			want:    constval.Black,
			wantErr: false,
		},
		{
			name:    "basicTo256",
			args:    args{to: C256, src: constval.Black},
			want:    constval.Black,
			wantErr: false,
		},
		{
			name:    "basicToRGB",
			args:    args{to: TrueColor, src: constval.Black},
			want:    constval.Black,
			wantErr: false,
		},
		{
			name:    "256toBasic",
			args:    args{to: Ascii, src: color256(100)},
			want:    c256ToBasic(100),
			wantErr: false,
		},
		{
			name: "256to256",
			args: args{to: C256, src: color256(100)},
			want: color256(100),
		},
		{
			name: "256toRGB",
			args: args{to: TrueColor, src: color256(100)},
			want: color256(100),
		},
		{
			name: "RGBtoBasic",
			args: args{to: Ascii, src: rgb(100, 200, 255)},
			want: RGBToBasic(100, 200, 255),
		},
		{
			name: "RGBto256",
			args: args{to: C256, src: rgb(100, 200, 255)},
			want: RGBTo256(100, 200, 255),
		},
		{
			name: "RGBtoRGB",
			args: args{to: TrueColor, src: rgb(100, 200, 255)},
			want: rgb(100, 200, 255),
		},
		{
			name: "unknown",
			args: args{to: 100, src: constval.Black},
			want: constval.Black,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertColorIfGreaterThanExpect(tt.args.to, tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertColorIfGreaterThanExpect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertColorIfGreaterThanExpect() got = %v, want %v", got, tt.want)
			}
		})
	}
}
