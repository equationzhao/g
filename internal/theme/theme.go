package theme

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/Equationzhao/g/internal/global"
)

type Style struct {
	// Color of the text.
	Color string `json:"color,omitempty"`
	// unicode Icon
	Icon      string `json:"icon,omitempty"`
	Underline bool   `json:"underline,omitempty"`
	Bold      bool   `json:"bold,omitempty"`
	Faint     bool   `json:"faint,omitempty"`
	Italics   bool   `json:"italics,omitempty"`
	Blink     bool   `json:"blink,omitempty"`
}

var (
	genOnceStyleField sync.Once
	styleField        []string
)

func genStyleField() []string {
	genOnceStyleField.Do(func() {
		// reflect on Style
		styleTypeValue := reflect.TypeFor[Style]()
		n := styleTypeValue.NumField()
		styleField = make([]string, 0, n)
		for i := range n {
			if j := styleTypeValue.Field(i).Tag.Get("json"); j != "" {
				jName, _, _ := strings.Cut(j, ",")
				if jName != "" {
					styleField = append(styleField, jName)
				}
			}
		}
	})
	return styleField
}

func (s *Style) UnmarshalJSON(bytes []byte) error {
	type styleAlias Style // avoid recursion during JSON unmarshalling.
	var alias styleAlias
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	// check if there is any unknown field
	fields := genStyleField()
	for key := range raw {
		if !slices.Contains(fields, key) {
			return fmt.Errorf("unknown field: '%s'", key)
		}
	}

	if err := json.Unmarshal(bytes, &alias); err != nil {
		return err
	}

	*s = Style(alias)
	return nil
}

func (s *Style) FromReadable() error {
	c, err := str2color(s.Color)
	if err != nil {
		return err
	}
	s.Color = c
	return nil
}

type Theme map[string]Style

func (t Theme) UnmarshalJSON(bytes []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	for key := range raw {
		s := Style{}
		if err := json.Unmarshal(raw[key], &s); err != nil {
			return fmt.Errorf("failed at key '%s': %s", key, err.Error())
		}
		t[key] = s
	}
	return nil
}

var basicColor2str = map[string]string{
	global.Black:        "black",
	global.Red:          "red",
	global.Green:        "green",
	global.Yellow:       "yellow",
	global.Blue:         "blue",
	global.Purple:       "purple",
	global.Cyan:         "cyan",
	global.White:        "white",
	global.BrightRed:    "bright-red",
	global.BrightGreen:  "bright-green",
	global.BrightYellow: "bright-yellow",
	global.BrightBlue:   "bright-blue",
	global.BrightPurple: "bright-purple",
	global.BrightCyan:   "bright-cyan",
	global.BrightWhite:  "bright-white",
	global.BrightBlack:  "bright-black",
	global.Reset:        "reset",
	global.Underline:    "underline",
}

func color2str(color string) string {
	// basic colors
	if str, ok := basicColor2str[color]; ok {
		return str
	}

	// detect format:
	strReader := strings.NewReader(color)

	// 1.8bit/256 color
	var c uint8
	_, err := fmt.Fscanf(strReader, Color256Format, &c)
	if err == nil {
		return fmt.Sprintf("[%d]@256", c)
	}
	// 2.rgb
	var (
		r uint8 = 0
		g uint8 = 0
		b uint8 = 0
	)
	strReader = strings.NewReader(color)
	_, err = fmt.Fscanf(strReader, RGBFormat, &r, &g, &b)
	if err == nil {
		return fmt.Sprintf("[%d,%d,%d]@rgb", r, g, b)
	}

	color = strings.ReplaceAll(color, " ", "")
	if strings.HasPrefix(color, global.Underline) {
		return color2str(global.Underline) + " + " + color2str(color[len(global.Underline):])
	}
	return ""
}

// str2color convert string to color
// support: red, green, yellow, blue, purple, cyan, white, black, and their bright version
// Underline
// any color with underline, should be in the format of "Underline + [color]"
// [value]@256
// [values]@rgb
// [values]@hex (will be turned to rgb)
func str2color(str string) (string, error) {
	switch str {
	case "":
		return "", nil
	case "black", "Black":
		return global.Black, nil
	case "red", "Red":
		return global.Red, nil
	case "green", "Green":
		return global.Green, nil
	case "yellow", "Yellow":
		return global.Yellow, nil
	case "blue", "Blue":
		return global.Blue, nil
	case "purple", "Purple":
		return global.Purple, nil
	case "cyan", "Cyan":
		return global.Cyan, nil
	case "white", "White":
		return global.White, nil
	case "bright-red", "BrightRed":
		return global.BrightRed, nil
	case "bright-green", "BrightGreen":
		return global.BrightGreen, nil
	case "bright-yellow", "BrightYellow":
		return global.BrightYellow, nil
	case "bright-blue", "BrightBlue":
		return global.BrightBlue, nil
	case "bright-purple", "BrightPurple":
		return global.BrightPurple, nil
	case "bright-cyan", "BrightCyan":
		return global.BrightCyan, nil
	case "bright-white", "BrightWhite":
		return global.BrightWhite, nil
	case "bright-black", "BrightBlack":
		return global.BrightBlack, nil
	case "reset", "Reset":
		return global.Reset, nil
	default:
		// remove spaces
		str = strings.ReplaceAll(str, " ", "")

		// 256 color
		if strings.HasSuffix(str, "@256") {
			code, err := strconv.Atoi(strings.Trim(str[:len(str)-4], "[]"))
			if err != nil {
				return "", err
			}
			colorStr, err := Color256(code)
			if err != nil {
				return "", err
			}
			return colorStr, nil
		}

		// rgb color
		if strings.HasSuffix(str, "@rgb") {
			code := strings.Trim(str[:len(str)-4], "[]")
			rgb := strings.Split(code, ",")
			if len(rgb) != 3 {
				return "", errors.New("too many or too few rgb values")
			}
			r, err1 := strconv.ParseUint(rgb[0], 10, 8)
			g, err2 := strconv.ParseUint(rgb[1], 10, 8)
			b, err3 := strconv.ParseUint(rgb[2], 10, 8)
			if err1 != nil || err2 != nil || err3 != nil {
				return "", errors.New("rgb values must be numbers")
			}
			colorStr, err := RGB(uint8(r), uint8(g), uint8(b))
			if err != nil {
				return "", err
			}
			return colorStr, nil
		}

		// hex
		if strings.HasSuffix(str, "@hex") {
			code := strings.Trim(str[:len(str)-4], "[]")
			if !IsValidHexColor(code) {
				return "", errors.New("invalid hex color")
			}
			rgb := HexToRgb(code)
			colorStr, err := RGB(rgb[0], rgb[1], rgb[2])
			if err != nil {
				return "", errors.New("rgb values must be numbers")
			}
			return colorStr, nil
		}

		return global.Reset, nil
	}
}

func IsValidHexColor(color string) bool {
	match, _ := regexp.MatchString("^(#0x|#|0x)?([0-9a-fA-F]{3}){1,2}$", color)
	return match
}

/*
......
// if using 256 color, you can use color code like this:
[0-255]@256
// if using rgb color, you can use color code like this:
[0-255,0-255,0-255]@rgb
// if using hex color, you can use color code like this:
[hex]@hex
*/

type ErrBadColor struct {
	name string
	error
}

func (e ErrBadColor) Error() string {
	return fmt.Sprintf("bad color for %s: %s", e.name, e.error.Error())
}

type ErrOpenTheme struct {
	error
}

func (e ErrOpenTheme) Error() string {
	return fmt.Sprintf("load theme error: %s", e.error.Error())
}

func GetTheme(path string) error {
	themeJson, err := os.ReadFile(path)
	if err != nil {
		return ErrOpenTheme{err}
	}
	a, err, fatal := getTheme(themeJson)
	if fatal != nil {
		return ErrOpenTheme{fatal}
	}
	DefaultAll = a
	return err
}

func getTheme(themeJson []byte) (theme All, errSum, fatal error) {
	err := json.Unmarshal(themeJson, &theme)
	if err != nil {
		return All{}, nil, err
	}
	// check
	errSum = nil
	check := func(m Theme) {
		for key := range m {
			style := m[key]
			err := style.FromReadable()
			if err != nil {
				errSum = errors.Join(errSum, ErrBadColor{name: key, error: err})
				continue
			}
			m[key] = style
		}
	}
	theme.Apply(check)
	return theme, errSum, nil
}

func ConvertThemeColor() {
	convert := func(m Theme) {
		for key := range m {
			if key == "reset" {
				continue
			}
			color, err := ConvertColorIfGreaterThanExpect(ColorLevel, m[key].Color)
			if err != nil {
				continue
			}
			m[key] = Style{
				Icon:      m[key].Icon,
				Color:     color,
				Underline: m[key].Underline,
				Bold:      m[key].Bold,
				Faint:     m[key].Faint,
				Italics:   m[key].Italics,
				Blink:     m[key].Blink,
			}
		}
	}
	DefaultAll.Apply(convert)
}
