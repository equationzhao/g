package theme

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/Equationzhao/g/internal/global"
)

var (
	dir   = global.BrightBlue
	pic   = global.BrightPurple
	video = color256(200)
	audio = color256(199)
	db    = color256(195)
	lang  = color256(158)
	text  = color256(153)
	doc   = color256(150)
	tar   = global.Red
	pkg   = global.BrightYellow
	bash  = global.BrightGreen
	lock  = color256(227)
	vim   = color256(41)
	key   = color256(214)
	conf  = global.White
	iso   = color256(88)
)

// permission: match file permission
// size: match file size unit
// user: match username
// group: match group name
// symlink: match symlink
// git: match git status
// name: match file name
// special: match file type: symlink, dir, executable
// ext: match file extension

type All struct {
	InfoTheme  Theme `json:"info,omitempty"`
	Permission Theme `json:"permission,omitempty"`
	Size       Theme `json:"size,omitempty"`
	User       Theme `json:"user,omitempty"`
	Group      Theme `json:"group,omitempty"`
	Symlink    Theme `json:"symlink,omitempty"`
	Git        Theme `json:"git,omitempty"`
	Name       Theme `json:"name,omitempty"`
	Special    Theme `json:"special,omitempty"`
	Ext        Theme `json:"ext,omitempty"`
}

var (
	genOnceAllField sync.Once
	allFieldTag     []string
	allField        []string
)

func genAllField() ([]string, []string) {
	genOnceAllField.Do(func() {
		// reflect on Style
		allTypeValue := reflect.TypeFor[All]()
		n := allTypeValue.NumField()
		allField = make([]string, 0, n)
		allFieldTag = make([]string, 0, n)
		for i := range n {
			f := allTypeValue.Field(i)
			if j := f.Tag.Get("json"); j != "" {
				jName, _, _ := strings.Cut(j, ",")
				if jName != "" {
					allFieldTag = append(allFieldTag, jName)
					allField = append(allField, f.Name)
				}
			}
		}
	})
	return allField, allFieldTag
}

func (a *All) UnmarshalJSON(bytes []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	all, allTags := genAllField()

	for key := range raw {
		index := 0
		for ; index < len(allTags); index++ {
			if allTags[index] == key {
				break
			}
		}
		if index == len(allTags) {
			return fmt.Errorf("unknown field: '%s'", key)
		}
		m := make(Theme)
		if err := json.Unmarshal(raw[key], &m); err != nil {
			return fmt.Errorf("failed at key '%s': %s", key, err.Error())
		}
		reflect.ValueOf(a).Elem().FieldByName(all[index]).Set(reflect.ValueOf(m))
	}
	return nil
}

func (a *All) Apply(f func(theme Theme)) {
	f(a.InfoTheme)
	f(a.Permission)
	f(a.Size)
	f(a.User)
	f(a.Group)
	f(a.Symlink)
	f(a.Git)
	f(a.Name)
	f(a.Special)
	f(a.Ext)
}

func (a *All) CheckLowerCase() {
	checkLowerCase := func(theme Theme) {
		for k := range theme {
			if strings.ToLower(k) != k {
				panic(fmt.Sprintf("theme key must be lowercase!: key=%s", k))
			}
		}
	}
	// key in a.Size can be UpperCase
	checkLowerCase(a.InfoTheme)
	checkLowerCase(a.Permission)
	checkLowerCase(a.User)
	checkLowerCase(a.Group)
	checkLowerCase(a.Symlink)
	checkLowerCase(a.Git)
	checkLowerCase(a.Name)
	checkLowerCase(a.Special)
	checkLowerCase(a.Ext)
}

var InfoTheme = Theme{
	"inode": Style{
		Color: global.Purple,
	},
	"time": Style{
		Color: global.Blue,
	},
	"reset": Style{
		Color: global.Reset,
	},
	"-": {
		Color: global.White,
	},
	"charset": {
		Color:   global.White,
		Italics: true,
	},
	"mime": {
		Color:   global.White,
		Italics: true,
	},
	"checksum": {
		Underline: true,
	},
}

var Ext = Theme{
	"deb": Style{
		Color: pkg,
		Icon:  "\uF306",
	},
	"aab": Style{
		Color: global.Green,
		Icon:  "\uF17B",
	},
	"aar": Style{
		Color: global.Green,
		Icon:  "\uF17B",
	},
	"apk": Style{
		Color: global.Green,
		Icon:  "\uF17B",
	},
	"dex": Style{
		Color: global.Green,
		Icon:  "\uF17B",
	},
	"odex": Style{
		Color: global.Green,
		Icon:  "\uF17B",
	},
	"aidl": Style{
		Color: lang,
		Icon:  "\uF17B",
	},
	"hidl": Style{
		Color: lang,
		Icon:  "\uF17B",
	},
	"app": Style{
		Color: global.Green,
		Icon:  "\ueb44",
	},
	"msi": Style{
		Color: pkg,
		Icon:  "\uF17A",
	},
	"dpkg": Style{
		Color: pkg,
		Icon:  "\ue77d",
	},
	"rpm": Style{
		Color: pkg,
		Icon:  "",
	},
	"exe": Style{
		Color: global.Green,
		Icon:  "\uF17A",
	},
	"rdp": Style{
		Color: global.White,
		Icon:  "\ueb39",
	},
	"link": Style{
		Color: global.Purple,
		Icon:  "\ueb15",
	},
	"links": Style{
		Color: global.Purple,
		Icon:  "\ueb15",
	},
	"lnk": Style{
		Color: global.Purple,
		Icon:  "\ueb15",
	},
	"go": Style{
		Color: lang,
		Icon:  "\uE626",
	},
	"rs": Style{
		Color: lang,
		Icon:  "\uE7A8",
	},
	"cargo": {
		Color: lang,
		Icon:  "\uE7A8",
	},
	"rustup": {
		Color: lang,
		Icon:  "\uE7A8",
	},
	"c": Style{
		Color: lang,
		Icon:  "\uE61E",
	},
	"cpp": Style{
		Color: lang,
		Icon:  "\uE61D",
	},
	"c++": Style{
		Color: lang,
		Icon:  "\uE61D",
	},
	"cc": Style{
		Color: lang,
		Icon:  "\uE61D",
	},
	"cxx": Style{
		Color: lang,
		Icon:  "\uE61D",
	},
	"h": Style{
		Color: lang,
		Icon:  "",
	},
	"hpp": Style{
		Color: lang,
		Icon:  "",
	},
	"h++": Style{
		Color: lang,
		Icon:  "",
	},
	"hh": Style{
		Color: lang,
		Icon:  "",
	},
	"hxx": Style{
		Color: lang,
		Icon:  "",
	},
	"cs": Style{
		Color: lang,
		Icon:  "\U000F031B",
	},
	"scala": Style{
		Color: lang,
		Icon:  "\uE737",
	},
	"swift": Style{
		Color: lang,
		Icon:  "\uE755",
	},
	"xcplayground": Style{
		Color: lang,
		Icon:  "\uE755",
	},
	"kt": Style{
		Color: lang,
		Icon:  "\uE634",
	},
	"m": Style{
		Color: lang,
		Icon:  "\uE61E",
	},
	"py": Style{
		Color: lang,
		Icon:  "\uE606",
	},
	"ipynb": Style{
		Color: lang,
		Icon:  "\uE606",
	},
	"pyc": Style{
		Color: lang,
		Icon:  "\uE606",
	},
	"whl": Style{
		Color: lang,
		Icon:  "\uE606",
	},
	"rb": Style{
		Color: lang,
		Icon:  "\uE21E",
	},
	"rubydoc": Style{
		Color: lang,
		Icon:  "\uE73B",
	},
	"astro": Style{
		Color: lang,
		Icon:  "\uf135",
	},
	"js": Style{
		Color: lang,
		Icon:  "\uE781",
	},
	"cjs": Style{
		Color: lang,
		Icon:  "\uE781",
	},
	"mjs": Style{
		Color: lang,
		Icon:  "\uE781",
	},
	"ts": Style{
		Color: lang,
		Icon:  "\uE628",
	},
	"http": Style{
		Color: lang,
		Icon:  "\ueb01",
	},
	"npm": Style{
		Color: lang,
		Icon:  "\ue71e",
	},
	"nix": Style{
		Color: global.Yellow,
		Icon:  "\uF313",
	},
	"asc": Style{
		Color: global.BrightGreen,
		Icon:  "\U000F099D",
	},
	"tf": {
		Color: lang,
		Icon:  "\U000F1062",
	},
	"ics": {
		Color: global.White,
		Icon:  "\uEAB0",
	},
	"env": {
		Color: global.Black,
		Icon:  "\uF462",
	},
	"jsx": Style{
		Color: lang,
		Icon:  "\ue7ba",
	},
	"htm": Style{
		Color: lang,
		Icon:  "\uF13B",
	},
	"html": Style{
		Color: lang,
		Icon:  "\uF13B",
	},
	"css": Style{
		Color: lang,
		Icon:  "\uE749",
	},
	"java": Style{
		Color: lang,
		Icon:  "\uE738",
	},
	"jar": Style{
		Color: lang,
		Icon:  "\uE738",
	},
	"json": Style{
		Color: conf,
		Icon:  "\uE60B",
	},
	"json5": Style{
		Color: conf,
		Icon:  "\uE60B",
	},
	"hson": Style{
		Color: conf,
		Icon:  "\uE60B",
	},
	"xml": Style{
		Color: conf,
		Icon:  "\ue796",
	},
	"cson": Style{
		Color: conf,
		Icon:  "\uE601",
	},
	"ttf": Style{
		Color: global.Cyan,
		Icon:  "\uf031",
	},
	"otf": Style{
		Color: global.Cyan,
		Icon:  "\uf031",
	},
	"woff": Style{
		Color: global.Cyan,
		Icon:  "\uf031",
	},
	"woff2": Style{
		Color: global.Cyan,
		Icon:  "\uf031",
	},
	"eot": Style{
		Color: global.Cyan,
		Icon:  "\uf031",
	},
	"properties": Style{
		Color: conf,
		Icon:  "\uE60B",
	},
	"asm": Style{
		Color: lang,
		Icon:  "\ue637",
	},
	"groovy": Style{
		Color: lang,
		Icon:  "\ue775",
	},
	"s": Style{
		Color: lang,
		Icon:  "\ue637",
	},
	"styl": Style{
		Color: lang,
		Icon:  "\uE600",
	},
	"iml": Style{
		Color: global.White,
		Icon:  "\uE7B5",
	},
	"gv": Style{
		Color: global.Cyan,
		Icon:  "\uE225",
	},
	"hs": Style{
		Color: lang,
		Icon:  "\ue777",
	},
	"d": Style{
		Color: lang,
		Icon:  "\ue7af",
	},
	"dart": Style{
		Color: lang,
		Icon:  "\uE798",
	},
	"erl": Style{
		Color: lang,
		Icon:  "\ue7b1",
	},
	"hrl": Style{
		Color: lang,
		Icon:  "\uE7B1",
	},
	"clj": Style{
		Color: lang,
		Icon:  "",
	},
	"bzl": Style{
		Color: lang,
		Icon:  "",
	},
	"avro": Style{
		Color: lang,
		Icon:  "\ue60b",
	},
	"svelte": Style{
		Color: lang,
		Icon:  "",
	},
	"mustache": Style{
		Color: lang,
		Icon:  "\ue60f",
	},
	"sass": Style{
		Color: lang,
		Icon:  "\ue603",
	},
	"coffee": Style{
		Color: lang,
		Icon:  "\uE751",
	},
	"iso": Style{
		Color: iso,
		Icon:  "\uE271",
	},
	"dmg": Style{
		Color: iso,
		Icon:  "\uE271",
	},
	"img": Style{
		Color: iso,
		Icon:  "\uE271",
	},
	"vhd": Style{
		Color: iso,
		Icon:  "\uE271",
	},
	"vhdx": Style{
		Color: iso,
		Icon:  "\uE271",
	},
	"lua": Style{
		Color: lang,
		Icon:  "\uE620",
	},
	"vue": Style{
		Color: lang,
		Icon:  "\uE6A0",
	},
	"php": Style{
		Color: lang,
		Icon:  "\uE73D",
	},
	"twig": Style{
		Color: lang,
		Icon:  "",
	},
	"fs": Style{
		Color: lang,
		Icon:  "",
	},
	"ex": Style{
		Color: lang,
		Icon:  "",
	},
	"cmake": Style{
		Color: lang,
		Icon:  "\uE20F",
	},
	"less": Style{
		Color: lang,
		Icon:  "\ue758",
	},
	"r": Style{
		Color: lang,
		Icon:  "\uF25D",
	},
	"f": Style{
		Color: lang,
		Icon:  "\uf121",
	},
	"atom": Style{
		Color: lang,
		Icon:  "\ue764",
	},
	"docker": Style{
		Color: global.Yellow,
		Icon:  "\ue7b0",
	},
	"md": Style{
		Color: global.BrightYellow,
		Icon:  "\uF48A",
	},
	"txt": Style{
		Color: text,
		Icon:  "\uF15C",
	},
	"zig": Style{
		Color: lang,
		Icon:  "\uF0E7",
	},
	"rss": Style{
		Color: lang,
		Icon:  "\uF09E",
	},
	"ko": Style{
		Color: lang,
		Icon:  "\uebc6",
	},
	"zip": Style{
		Color: tar,
		Icon:  "",
	},
	"zst": Style{
		Color: tar,
		Icon:  "",
	},
	"sitx": Style{
		Color: tar,
		Icon:  "",
	},
	"tar": Style{
		Color: tar,
		Icon:  "",
	},
	"gz": Style{
		Color: tar,
		Icon:  "",
	},
	"rar": Style{
		Color: tar,
		Icon:  "",
	},
	"7z": Style{
		Color: tar,
		Icon:  "",
	},
	"xz": Style{
		Color: tar,
		Icon:  "",
	},
	"jpg": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"jpeg": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"png": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"bmp": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"tif": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"tiff": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"gif": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"svg": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"webp": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"pcx": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"avif": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"psd": Style{
		Color: pic,
		Icon:  "\ue67f",
	},
	"mp4": Style{
		Color: video,
		Icon:  "",
	},
	"mkv": Style{
		Color: video,
		Icon:  "",
	},
	"avi": Style{
		Color: video,
		Icon:  "",
	},
	"m4v": Style{
		Color: video,
		Icon:  "",
	},
	"flv": Style{
		Color: video,
		Icon:  "",
	},
	"mov": Style{
		Color: video,
		Icon:  "",
	},
	"mpeg": Style{
		Color: video,
		Icon:  "",
	},
	"ai": Style{
		Color: global.Cyan,
		Icon:  "\uE7B4",
	},
	"mp3": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"m4a": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"mid": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"midi": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"flac": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"ape": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"alac": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"aac": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"wav": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"ogg": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"pdf": Style{
		Color: global.Cyan,
		Icon:  "\uF1C1",
	},
	"pl": Style{
		Color: lang,
		Icon:  "\uE67E",
	},
	"epub": Style{
		Color: global.Cyan,
		Icon:  "\uE28A",
	},
	"mobi": Style{
		Color: global.Cyan,
		Icon:  "\uF12D",
	},
	"azw": Style{
		Color: global.Cyan,
		Icon:  "\uF12D",
	},
	"azw3": Style{
		Color: global.Cyan,
		Icon:  "\uF12D",
	},
	"doc": Style{
		Color: doc,
		Icon:  "\uf1c2",
	},
	"docx": Style{
		Color: doc,
		Icon:  "\uf1c2",
	},
	"docm": Style{
		Color: doc,
		Icon:  "\uf1c2",
	},
	"pages": Style{
		Color: doc,
		Icon:  "\uf1c2",
	},
	"xls": Style{
		Color: doc,
		Icon:  "\uf1c3",
	},
	"xlsx": Style{
		Color: doc,
		Icon:  "\uf1c3",
	},
	"xlsm": Style{
		Color: doc,
		Icon:  "\uf1c3",
	},
	"numbers": Style{
		Color: doc,
		Icon:  "\uf1c3",
	},
	"csv": Style{
		Color: doc,
		Icon:  "\ue64a",
	},
	"ppt": Style{
		Color: doc,
		Icon:  "\uf1c4",
	},
	"pptm": Style{
		Color: doc,
		Icon:  "\uf1c4",
	},
	"pptx": Style{
		Color: doc,
		Icon:  "\uf1c4",
	},
	"dot": Style{
		Color: doc,
		Icon:  "\uf1c2",
	},
	"dotx": Style{
		Color: doc,
		Icon:  "\uf1c2",
	},
	"diff": Style{
		Color: global.Cyan,
		Icon:  "\uF440",
	},
	"ini": Style{
		Color: conf,
		Icon:  "\uE615",
	},
	"conf": Style{
		Color: conf,
		Icon:  "\uE615",
	},
	"config": Style{
		Color: conf,
		Icon:  "\uE615",
	},
	"cfg": Style{
		Color: conf,
		Icon:  "\uE615",
	},
	"yml": Style{
		Color: conf,
		Icon:  "\uF481",
	},
	"yaml": Style{
		Color: conf,
		Icon:  "\uF481",
	},
	"tex": Style{
		Color: lang,
		Icon:  "\uF034",
	},
	"typ": Style{
		Color: lang,
		Icon:  "∫",
	},
	"jl": Style{
		Color: lang,
		Icon:  "",
	},
	"toml": Style{
		Color: conf,
		Icon:  "\uE615",
	},
	"pkl": Style{
		Color: conf,
		Icon:  "\uE615",
	},
	"db": Style{
		Color: db,
		Icon:  "\uF1C0",
	},
	"accdb": Style{
		Color: db,
		Icon:  "\uF1C0",
	},
	"rdb": Style{
		Color: db,
		Icon:  "\uF1C0",
	},
	"sqlite": Style{
		Color: db,
		Icon:  "\uE7C4",
	},
	"sqlite3": Style{
		Color: db,
		Icon:  "\uE7C4",
	},
	"sql": Style{
		Color: db,
		Icon:  "\uE706",
	},
	"prql": Style{
		Color: db,
		Icon:  "\uE706",
	},
	"db3": Style{
		Color: db,
		Icon:  "\uF1C0",
	},
	"sh": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"bat": Style{
		Color: bash,
		Icon:  "\ue629",
	},
	"ps1": Style{
		Color: bash,
		Icon:  "\ue683",
	},
	"csh": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"fish": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"ksh": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"zsh": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"zsh_history": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"zshrc": Style{
		Color: bash,
		Icon:  "󱆃",
	},
	"zprofile": Style{
		Color: bash,
		Icon:  "󱆃",
	},
	"bash": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"bashrc": Style{
		Color: bash,
		Icon:  "󱆃",
	},
	"bash_history": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"bash_profile": Style{
		Color: bash,
		Icon:  "󱆃",
	},
	".profile": Style{
		Color: bash,
		Icon:  "󱆃",
	},
	"nu": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"vim": Style{
		Color: vim,
		Icon:  "\uE62B",
	},
	"viminfo": Style{
		Color: vim,
		Icon:  "\uE62B",
	},
	"vimrc": Style{
		Color: vim,
		Icon:  "\uE62B",
	},
	"gvimrc": Style{
		Color: vim,
		Icon:  "\uE62B",
	},
	"nvim": Style{
		Color: vim,
		Icon:  "\uE62B",
	},
	"nvimrc": Style{
		Color: vim,
		Icon:  "\uE62B",
	},
	"log": Style{
		Color: global.White,
		Icon:  "\uF18D",
	},
	"lock": Style{
		Color: lock,
		Icon:  "\uF023",
	},
	"code-workspace": Style{
		Color: lang,
		Icon:  "\uE70C",
	},
	"key": Style{
		Color: key,
		Icon:  "\uF084",
	},
	"pub": Style{
		Color: key,
		Icon:  "\uF084",
	},
	"id_rsa": Style{
		Color: key,
		Icon:  "\uF084",
	},
	"gpg": Style{
		Color: key,
		Icon:  "\uF084",
	},
	"cer": Style{
		Color: key,
		Icon:  "\uF084",
	},
	"crt": Style{
		Color: key,
		Icon:  "\uF084",
	},
	"pgp": Style{
		Color: key,
		Icon:  "\uF084",
	},
	"passwd": Style{
		Color: global.Cyan,
		Icon:  "\uF023",
	},
	"maintainers": Style{
		Color: global.Cyan,
		Icon:  "\uF0C0",
	},
	"so": Style{
		Color: global.Cyan,
		Icon:  "\uF121",
	},
	"dll": Style{
		Color: global.Cyan,
		Icon:  "\uF121",
	},
	"patch": Style{
		Color: global.White,
		Icon:  "",
	},
	"torrent": Style{
		Color: global.Green,
		Icon:  "\ueac2",
	},
	"service": Style{
		Color: global.White,
		Icon:  "",
	},
	"swp": Style{
		Color: global.White,
		Icon:  "\uebcb",
	},
	"graphql": Style{
		Color: global.White,
		Icon:  "\uE662",
	},
}

var Special = map[string]Style{
	"dir": {
		Color: dir,
		Icon:  "\uE5FF",
	},
	"empty-dir": {
		Color: dir,
		Icon:  "\uF115",
	},
	"dir-prompt": {
		Color: global.Yellow,
		Icon:  "► ",
	},
	"pipe": {
		Color: global.Cyan,
		Icon:  "\uF124",
	},
	"socket": {
		Color: global.Cyan,
		Icon:  "\uF1E6",
	},
	"device": {
		Color: global.Yellow,
		Icon:  "󰜫",
	},
	"char": {
		Color: global.Yellow,
		Icon:  "",
	},
	"exe": {
		Color: global.Green,
		Icon:  "\uF17A",
	},
	"link": {
		Color: global.Purple,
	},
	"file": {
		Color: global.White,
		Icon:  "\uF016",
	},
	"hidden-file": {
		Color: global.White,
		Icon:  "",
	},
	"mounts": {
		Color: global.BrightBlack,
	},
}

var Name = map[string]Style{
	"home": {
		Icon:  "\ue617",
		Color: global.White,
	},
	"desktop": {
		Icon:  "\uf108",
		Color: global.White,
	},
	"documents": {
		Icon:  "\uF02D",
		Color: global.White,
	},
	"doc": {
		Color: doc,
		Icon:  "\uf02d",
	},
	"links": {
		Icon:  "\uf0f6",
		Color: global.White,
	},
	"downloads": {
		Icon:  "\uf498",
		Color: global.White,
	},
	"trash": {
		Icon:  "\uf1f8",
		Color: global.BrightBlack,
	},
	".trash": {
		Icon:  "\uf1f8",
		Color: global.BrightBlack,
	},
	"searches": {
		Icon:  "\uf422",
		Color: global.White,
	},
	"microsoft": {
		Icon:  "\uF17A",
		Color: global.White,
	},
	"google": {
		Icon:  "\uf1a0",
		Color: global.White,
	},
	"onedrive": {
		Icon:  "\ue762",
		Color: global.Blue,
	},
	"onedrivetemp": {
		Icon:  "\ue762",
		Color: global.Blue,
	},
	"favorites": {
		Icon:  "\ue623",
		Color: global.Red,
	},
	"saved games": {
		Icon:  "\uf11b",
		Color: global.Red,
	},
	".wakatime": {
		Icon:  "\ue641",
		Color: global.White,
	},
	".azure": {
		Icon:  "\uebd8",
		Color: global.White,
	},
	"contacts": {
		Icon:  "\uf0c0",
		Color: global.White,
	},
	"users": {
		Icon:  "\uf0c0",
		Color: global.White,
	},
	"lib": {
		Icon:  "\uF121",
		Color: global.White,
	},
	"library": {
		Icon:  "\uF121",
		Color: global.White,
	},
	"bin": {
		Icon:  "\uE5FC",
		Color: global.Green,
	},
	"share": {
		Icon:  "\uf064",
		Color: global.White,
	},
	"license": {
		Icon:  "",
		Color: key,
	},
	"licence": {
		Icon:  "",
		Color: key,
	},
	"shell": {
		Icon:  "\uF489",
		Color: bash,
	},
	"config": {
		Icon:  "\uE615",
		Color: conf,
	},
	".ds_store": {
		Color: global.BrightBlack,
		Icon:  "\uf179",
	},
	"pkgbuild": {
		Color: pkg,
		Icon:  "\uf303",
	},
	".srcinfo": {
		Color: pkg,
		Icon:  "\uf303",
	},
	"applications": {
		Icon:  "\ueb44",
		Color: global.Green,
	},
	"android": {
		Icon:  "\uE70E",
		Color: global.White,
	},
	".idea": {
		Icon:  "\uE7B5",
		Color: global.White,
	},
	".github": {
		Icon:  "\uF408",
		Color: global.White,
	},
	".gitattributes": {
		Icon:  "\uE65D",
		Color: global.BrightBlack,
	},
	".gitmodules": {
		Icon:  "\uE65D",
		Color: global.BrightBlack,
	},
	".gitconfig": {
		Icon:  "\uE65D",
		Color: global.BrightBlack,
	},
	".vscode": {
		Icon:  "\uE70C",
		Color: global.White,
	},
	"include": {
		Icon:  "\ue5fc",
		Color: global.White,
	},
	".dotnet": {
		Icon:  "\ue72e",
		Color: global.White,
	},
	"src": {
		Icon:  "\ue796",
		Color: global.White,
	},
	"node_modules": {
		Icon:  "\ue5fa",
		Color: global.White,
	},
	"font": {
		Icon:  "\uf031",
		Color: global.White,
	},
	"fonts": {
		Icon:  "\uf031",
		Color: global.White,
	},
	".git": {
		Color: global.White,
		Icon:  "\uF1D3",
	},
	".gitignore": {
		Color: global.BrightBlack,
		Icon:  "\ue65d",
	},
	".gitignore_global": {
		Color: global.BrightBlack,
		Icon:  "\ue65d",
	},
	".gitlab-ci.yml": {
		Color: global.White,
		Icon:  "\uF296",
	},
	"cmakelists.txt": {
		Color: lang,
		Icon:  "\uE20F",
	},
	"makefile": {
		Color: lang,
		Icon:  "\uE20F",
	},
	"justfile": {
		Color: lang,
		Icon:  "\uE20F",
	},
	"history": {
		Icon:  "\uF1DA",
		Color: global.BrightBlack,
	},
	".history": {
		Icon:  "\uF1DA",
		Color: global.BrightBlack,
	},
	"recovery": {
		Icon:  "\uF1DA",
		Color: global.BrightBlack,
	},
	"apple": {
		Color: global.BrightBlack,
		Icon:  "\uF179",
	},
	"dockerfile": {
		Color: global.BrightYellow,
		Icon:  "\ue7b0",
	},
	"dockerignore": {
		Color: global.BrightYellow,
		Icon:  "\ue7b0",
	},
	"readme": {
		Color:     global.Yellow,
		Icon:      "\uF48A",
		Underline: true,
		Bold:      true,
	},
	"readme.md": {
		Color:     global.Yellow,
		Icon:      "\uF48A",
		Underline: true,
		Bold:      true,
	},
	"flake.nix": {
		Color:     global.Yellow,
		Icon:      "\uF313",
		Underline: true,
		Bold:      true,
	},
	"jenkinsfile": {
		Color: lang,
		Icon:  "\uE66E",
	},
	"brewfile": {
		Color: lang,
		Icon:  "\uF016",
	},
	"sdk": {
		Icon:  "\uF121",
		Color: global.White,
	},
	"pictures": {
		Icon:  "\uF1C5",
		Color: global.White,
	},
	"videos": {
		Icon:  "",
		Color: global.White,
	},
	"movies": {
		Icon:  "",
		Color: global.White,
	},
	"music": {
		Icon:  "\uF025",
		Color: global.White,
	},
	"audio": {
		Icon:  "\uF025",
		Color: global.White,
	},
	"cargo.lock": {
		Color:     lang,
		Icon:      "\uE7A8",
		Faint:     true,
		Underline: true,
	},
	"cargo.toml": {
		Color:     lang,
		Icon:      "\uE7A8",
		Underline: true,
	},
	"known_hosts": {
		Color: global.White,
		Icon:  "\uEB39",
	},
	"repo": {
		Color: global.White,
		Icon:  "\uea62",
	},
	".ssh": {
		Icon:  "\ueba9",
		Color: global.White,
	},
	"boot": {
		Icon:  "\uf287",
		Color: global.White,
	},
	"cache": {
		Icon:  "\uf49b",
		Color: global.White,
	},
	"passwd": {
		Color: global.Cyan,
		Icon:  "\uF023",
	},
	"vagrantfile": {
		Color: lang,
		Icon:  "⍱",
	},
	"package.json": {
		Color:     conf,
		Icon:      "\uE718",
		Underline: true,
	},
	"tsconfig.json": {
		Color:     conf,
		Icon:      "\uE628",
		Underline: true,
	},
	"go.mod": {
		Color:     lang,
		Icon:      "\uE626",
		Underline: true,
	},
	"go.sum": {
		Color:     lang,
		Icon:      "\uE626",
		Underline: true,
		Faint:     true,
	},
	".python_history": {
		Color: lang,
		Icon:  "\uE606",
	},
	".cfusertextencoding": {
		Color: global.BrightBlack,
		Icon:  "\uF179",
	},
	"maintainers": {
		Color: global.White,
		Icon:  "\uF0C0",
	},
	"__pycache__": {
		Color: lang,
		Icon:  "\uE606",
		Faint: true,
	},
	"requirements.txt": {
		Color: global.White,
		Icon:  "\uE606",
		Faint: true,
	},
	"robots.txt": {
		Color: global.White,
		Icon:  "\U000F06A9",
	},
	"docker-compose.yml": {
		Color: global.BrightYellow,
		Icon:  "\ue7b0",
	},
	"docker-compose.yaml": {
		Color: global.BrightYellow,
		Icon:  "\ue7b0",
	},
	"contributing": {
		Color: global.White,
		Icon:  "\uF0C0",
	},
	"contributing.md": {
		Color: global.White,
		Icon:  "\uF0C0",
	},
	"cron.minutely": {
		Color: global.White,
		Icon:  "\uE5FC",
	},
	"cron.d": {
		Color: global.White,
		Icon:  "\uE5FC",
	},
	"cron.daily": {
		Color: global.White,
		Icon:  "\uE5FC",
	},
	"cron.hourly": {
		Color: global.White,
		Icon:  "\uE5FC",
	},
	"cron.mouthly": {
		Color: global.White,
		Icon:  "\uE5FC",
	},
	"cron.weekly": {
		Color: global.White,
		Icon:  "\uE5FC",
	},
	"crontab": {
		Color: global.White,
		Icon:  "\uE5FC",
	},
	"code_of_conduct.md": {
		Color: global.White,
		Icon:  "",
	},
}

var Permission = map[string]Style{
	"block": {
		Color: global.Cyan,
		Bold:  true,
	},
	"directory": {
		Color: global.Blue,
		Bold:  true,
	},
	"link": {
		Color: global.Purple,
		Bold:  true,
	},
	"char": {
		Color: global.Yellow,
		Bold:  true,
	},
	"pipe": {
		Color: global.Yellow,
		Bold:  true,
	},
	"socket": {
		Color: global.Yellow,
		Bold:  true,
	},
	"setuid": {
		Color: global.Purple,
		Bold:  true,
	},
	"setgid": {
		Color: global.Purple,
		Bold:  true,
	},
	"read": {
		Color: global.Yellow,
		Bold:  true,
	},
	"write": {
		Color: global.Red,
		Bold:  true,
	},
	"exe": {
		Color: global.Green,
		Bold:  true,
	},
	"-": {
		Color: global.BrightBlack,
	},
	"octal": {
		Color: color256(208),
		Bold:  true,
	},
}

var Size = map[string]Style{
	"-": {
		Color: global.White,
	},
	"block": {
		Color: rgb(20, 255, 100),
	},
	"bit": {
		Color: rgb(20, 255, 100),
	},
	"B": {
		Color: rgb(20, 230, 100),
	},
	"KB": {
		Color: rgb(20, 207, 100),
	},
	"KiB": {
		Color: rgb(20, 207, 100),
	},
	"MB": {
		Color: rgb(20, 188, 100),
	},
	"MiB": {
		Color: rgb(20, 188, 100),
	},
	"GB": {
		Color: rgb(20, 170, 100),
	},
	"GiB": {
		Color: rgb(20, 170, 100),
	},
	"TB": {
		Color: rgb(20, 153, 100),
	},
	"TiB": {
		Color: rgb(20, 153, 100),
	},
}

var Git = map[string]Style{
	"git_unmodified": {
		Color: global.BrightBlack,
	},
	"git_modified": {
		Color: global.Yellow,
	},
	"git_renamed": {
		Color: global.Blue,
	},
	"git_copied": {
		Color: global.Purple,
	},
	"git_deleted": {
		Color: global.Red,
	},
	"git_added": {
		Color: global.Green,
	},
	"git_untracked": {
		Color: global.BrightBlack,
	},
	"git_ignored": {
		Color: global.BrightRed,
	},
	"git_type_changed": {
		Color: global.Yellow,
	},
	"git_updated_but_unmerged": {
		Color: global.BrightYellow,
	},
	"git-repo-skip": {
		Color: global.BrightBlack,
		Icon:  "-",
	},
	"git-repo-clean": {
		Color: global.Green,
		Icon:  "clean",
	},
	"git-repo-dirty": {
		Color: global.Yellow,
		Icon:  "dirty",
	},
	"git-branch-master": { // master and main
		Color: global.Green,
		Bold:  true,
	},
	"git-branch": {
		Color: global.Yellow,
	},
	"git-branch-none": {
		Color: global.BrightBlack,
		Icon:  "-",
	},
	"git-author": {
		Color: rgb(67, 202, 207),
	},
	"git-author-date": {
		Color: rgb(54, 102, 180),
	},
	"git-commit-hash": {
		Color: global.BrightBlack,
	},
}

var Owner = map[string]Style{
	"owner": {
		Color: global.Yellow,
		Bold:  true,
	},
	"root": {
		Color: global.Red,
		Bold:  true,
	},
}

var Group = map[string]Style{
	"group": {
		Color: global.Yellow,
		Bold:  true,
	},
	"root": {
		Color: global.Red,
		Bold:  true,
	},
}

var Symlink = map[string]Style{
	"symlink_path": {
		Color: global.Green,
	},
	"symlink_broken_path": {
		Color:     global.Red,
		Underline: true,
	},
	"symlink": {
		Color: global.Purple,
		Icon:  "\ueb15",
	},
	"symlink_arrow": {
		Color: global.BrightWhite,
		Icon:  " => ",
	},
	"link-num": {
		Color: global.Red,
	},
}

var (
	DefaultAll All
	_init      bool
)

func init() {
	if !_init {
		DefaultAll = All{
			InfoTheme:  InfoTheme,
			Permission: Permission,
			Size:       Size,
			User:       Owner,
			Group:      Group,
			Symlink:    Symlink,
			Git:        Git,
			Name:       Name,
			Special:    Special,
			Ext:        Ext,
		}
	}
}
