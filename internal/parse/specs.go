package parse

// Spec is one budget-counted primary flag. Aliases do not count toward Budget.
type Spec struct {
	Name         string
	Short        rune // 0 if none
	Aliases      []string
	Dimension    Dimension
	Slot         string // same slot ⇒ last-wins replacement
	ConfigKey    string
	ConfigExempt ConfigExempt
	TakesValue   bool
	Default      string
	Help         string
}

// ConfigOnlyKey is a YAML field with no CLI flag (Gate 3: not worth argv).
type ConfigOnlyKey struct {
	Key  string
	Help string
}

// Specs returns the 40 primary flags. This is the single source of truth for
// names, dimensions, config keys, and man OPTIONS.
func Specs() []Spec {
	return []Spec{
		{Name: "--help", Short: '?', Aliases: []string{"-?"}, Dimension: DimMeta, Slot: "help", ConfigExempt: ConfigMeta, Default: "off", Help: "print usage to stdout and exit 0"},
		{Name: "--version", Dimension: DimMeta, Slot: "version", ConfigExempt: ConfigMeta, Default: "off", Help: "print version to stdout and exit 0"},
		{Name: "--config", Aliases: []string{"--no-config", "--config-file"}, Dimension: DimMeta, Slot: "config", ConfigExempt: ConfigMeta, TakesValue: true, Default: "search", Help: "config source: PATH or none"},

		{Name: "--format", Short: 0, Aliases: []string{"-C", "-x", "-1", "-m", "-T", "--json"}, Dimension: DimFormat, Slot: "format", ConfigKey: "format", TakesValue: true, Default: "auto", Help: "grid|across|oneline|comma|tree|json"},
		{Name: "--long", Short: 'l', Aliases: []string{"-l", "-o"}, Dimension: DimLong, Slot: "long", ConfigKey: "long", Default: "false", Help: "enable the long column set"},

		{Name: "--all", Short: 'a', Aliases: []string{"-a"}, Dimension: DimVisibility, Slot: "visibility", ConfigKey: "all", Default: "hidden", Help: "show hidden including . and .."},
		{Name: "--almost-all", Short: 'A', Aliases: []string{"-A"}, Dimension: DimVisibility, Slot: "visibility", ConfigKey: "almost_all", Default: "hidden", Help: "show hidden except . and .."},

		{Name: "--directory", Short: 'd', Aliases: []string{"-d"}, Dimension: DimWalk, Slot: "dirself", ConfigKey: "directory", Default: "false", Help: "list directory arguments themselves"},
		{Name: "--recursive", Short: 'R', Aliases: []string{"-R"}, Dimension: DimWalk, Slot: "recurse", ConfigKey: "recursive", Default: "false", Help: "list subdirectories recursively"},
		{Name: "--depth", Dimension: DimWalk, Slot: "depth", ConfigKey: "depth", TakesValue: true, Default: "unlimited", Help: "print nodes with Depth <= N"},

		{Name: "--ignore", Short: 'I', Aliases: []string{"-I"}, Dimension: DimFilter, Slot: "ignore", ConfigKey: "ignore", TakesValue: true, Default: "", Help: "omit basenames matching GLOB; repeatable"},
		{Name: "--only-dirs", Short: 'D', Aliases: []string{"-D"}, Dimension: DimFilter, Slot: "kind", ConfigKey: "only_dirs", Default: "all", Help: "keep directories only (children, not argv roots)"},
		{Name: "--only-files", Dimension: DimFilter, Slot: "kind", ConfigKey: "only_files", Default: "all", Help: "keep non-directories only (children, not argv roots)"},
		{Name: "--ignore-backups", Short: 'B', Aliases: []string{"-B"}, Dimension: DimFilter, Slot: "backups", ConfigKey: "ignore_backups", Default: "false", Help: "omit basenames ending in ~"},
		{Name: "--git-ignore", Dimension: DimFilter, Slot: "gitignore", ConfigKey: "git_ignore", Default: "false", Help: "omit gitignored children; fail-open"},

		{Name: "--sort", Short: 0, Aliases: []string{"-t", "-S", "-X", "-U", "-v"}, Dimension: DimSort, Slot: "sortkey", ConfigKey: "sort", TakesValue: true, Default: "name", Help: "name|size|time|ext|version|none"},
		{Name: "--reverse", Short: 'r', Aliases: []string{"-r"}, Dimension: DimSort, Slot: "reverse", ConfigKey: "reverse", Default: "false", Help: "reverse the primary key within groups"},
		{Name: "--dir-order", Aliases: []string{"--dir-first", "--group-directories-first"}, Dimension: DimSort, Slot: "dirorder", ConfigKey: "dir_order", TakesValue: true, Default: "none", Help: "none|first|last"},

		{Name: "--human-readable", Short: 'h', Aliases: []string{"-h"}, Dimension: DimSize, Slot: "human", ConfigKey: "human_readable", Default: "false", Help: "human sizes, powers of 1024"},
		{Name: "--si", Dimension: DimSize, Slot: "human", ConfigKey: "si", Default: "false", Help: "human sizes, powers of 1000; wins over -h"},

		{Name: "--inode", Short: 'i', Aliases: []string{"-i"}, Dimension: DimLong, Slot: "inode", ConfigKey: "inode", Default: "false", Help: "inode prefix or column"},
		{Name: "--links", Short: 'H', Aliases: []string{"-H"}, Dimension: DimLong, Slot: "nlink", ConfigKey: "links", Default: "false", Help: "hard-link count prefix or column"},
		{Name: "--numeric-uid-gid", Short: 'n', Aliases: []string{"-n"}, Dimension: DimLong, Slot: "numeric", ConfigKey: "numeric", Default: "false", Help: "print uid/gid or SID numbers"},
		{Name: "--no-group", Short: 'G', Aliases: []string{"-G"}, Dimension: DimLong, Slot: "nogroup", ConfigKey: "no_group", Default: "false", Help: "omit the group column"},
		{Name: "--blocks", Dimension: DimLong, Slot: "blocks", ConfigKey: "blocks", Default: "false", Help: "allocated 512-byte blocks; long only"},
		{Name: "--header", Dimension: DimLong, Slot: "header", ConfigKey: "header", Default: "false", Help: "print column names when long"},

		{Name: "--time", Dimension: DimTime, Slot: "timefield", ConfigKey: "time", TakesValue: true, Default: "modified", Help: "modified|accessed|changed|birth"},
		{Name: "--time-style", Dimension: DimTime, Slot: "timestyle", ConfigKey: "time_style", TakesValue: true, Default: "default", Help: "default|iso|long-iso|full-iso|relative|+FORMAT"},

		{Name: "--color", Dimension: DimPresent, Slot: "color", ConfigKey: "color", TakesValue: true, Default: "auto", Help: "always|auto|never"},
		{Name: "--icons", Dimension: DimPresent, Slot: "icons", ConfigKey: "icons", TakesValue: true, Default: "auto", Help: "always|auto|never"},
		{Name: "--hyperlink", Dimension: DimPresent, Slot: "hyperlink", ConfigKey: "hyperlink", TakesValue: true, Default: "auto", Help: "always|auto|never"},
		{Name: "--theme", Dimension: DimPresent, Slot: "theme", ConfigKey: "theme", TakesValue: true, Default: "builtin", Help: "path to theme JSON"},
		{Name: "--classify", Short: 'F', Aliases: []string{"-F"}, Dimension: DimPresent, Slot: "classify", ConfigKey: "classify", TakesValue: true, Default: "never", Help: "always|auto|never; bare -F is always"},
		{Name: "--quote-name", Short: 'Q', Aliases: []string{"-Q"}, Dimension: DimPresent, Slot: "quote", ConfigKey: "quote", Default: "default", Help: "always quote names"},
		{Name: "--literal", Short: 'N', Aliases: []string{"-N"}, Dimension: DimPresent, Slot: "quote", ConfigKey: "quote", Default: "default", Help: "never quote names"},
		{Name: "--zero", Short: '0', Aliases: []string{"-0"}, Dimension: DimPresent, Slot: "zero", ConfigKey: "zero", Default: "false", Help: "NUL record separator"},
		{Name: "--width", Dimension: DimPresent, Slot: "width", ConfigKey: "width", TakesValue: true, Default: "auto", Help: "screen width for grid/across/comma"},

		{Name: "--git", Dimension: DimGit, Slot: "git", ConfigKey: "git", Default: "false", Help: "two-character git status column"},
		{Name: "--dereference", Short: 'L', Aliases: []string{"-L", "--no-dereference"}, Dimension: DimDeref, Slot: "deref", ConfigKey: "dereference", Default: "false", Help: "follow symlink/junction metadata"},
		{Name: "-g", Short: 'g', Dimension: DimLong, Slot: "noowner", ConfigExempt: ConfigCLIOnly, Default: "false", Help: "long without owner"},
	}
}

// ConfigOnlyKeys are YAML fields with no CLI flag.
func ConfigOnlyKeys() []ConfigOnlyKey {
	return []ConfigOnlyKey{
		{Key: "icon_set", Help: "nerd|unicode; selects builtin glyph table"},
	}
}

// GNUReservedShorts are letters GNU ls already assigned. We must not invent
// a different meaning for them even if we do not implement that feature.
func GNUReservedShorts() map[rune]string {
	return map[rune]string{
		'b': "GNU --escape",
		'c': "GNU ctime / sort by ctime",
		'f': "GNU disable sort (and -l)",
		'k': "GNU block-size=1K",
		'p': "GNU classify directories with /",
		'q': "GNU hide control chars",
		's': "GNU size in blocks",
		'u': "GNU atime",
		'w': "GNU --width",
		'P': "GNU --no-dereference",
		'Z': "GNU --context",
	}
}

// AllNames returns primary names plus aliases, for "is this token ours?" checks.
func AllNames() map[string]Spec {
	out := make(map[string]Spec, 64)
	for _, s := range Specs() {
		out[s.Name] = s
		if s.Short != 0 {
			out["-"+string(s.Short)] = s
		}
		for _, a := range s.Aliases {
			out[a] = s
		}
	}
	return out
}
