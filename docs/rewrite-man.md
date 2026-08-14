# g — rewrite OPTIONS (generated)

This fragment is the OPTIONS source of truth for the rewrite.
Prose (DESCRIPTION, EXAMPLES) lives in `docs/rewrite-architecture.md`.
Do not edit the OPTIONS block by hand; run `go generate ./internal/parse`.

<!-- BEGIN OPTIONS -->

## Meta

`-?`, `--help`
: print usage to stdout and exit 0 (default: off)

`--version`
: print version to stdout and exit 0 (default: off)

`--config`, `--no-config`, `--config-file`
: config source: PATH or none (default: search)

## Format

`--format`, `-C`, `-x`, `-1`, `-m`, `-T`, `--json`
: grid|across|oneline|comma|tree|json (default: auto)

## Long

`-l`, `--long`, `-o`
: enable the long column set (default: false)

`-i`, `--inode`
: inode prefix or column (default: false)

`-H`, `--links`
: hard-link count prefix or column (default: false)

`-n`, `--numeric-uid-gid`
: print uid/gid or SID numbers (default: false)

`-G`, `--no-group`
: omit the group column (default: false)

`--blocks`
: allocated 512-byte blocks; long only (default: false)

`--header`
: print column names when long (default: false)

`-g`
: long without owner (default: false)

## Visibility

`-a`, `--all`
: show hidden including . and .. (default: hidden)

`-A`, `--almost-all`
: show hidden except . and .. (default: hidden)

## Walk

`-d`, `--directory`
: list directory arguments themselves (default: false)

`-R`, `--recursive`
: list subdirectories recursively (default: false)

`--depth`
: print nodes with Depth <= N (default: unlimited)

## Filter

`-I`, `--ignore`
: omit basenames matching GLOB; repeatable (default: )

`-D`, `--only-dirs`
: keep directories only (children, not argv roots) (default: all)

`--only-files`
: keep non-directories only (children, not argv roots) (default: all)

`-B`, `--ignore-backups`
: omit basenames ending in ~ (default: false)

`--git-ignore`
: omit gitignored children; fail-open (default: false)

## Sort

`--sort`, `-t`, `-S`, `-X`, `-U`, `-v`
: name|size|time|ext|version|none (default: name)

`-r`, `--reverse`
: reverse the primary key within groups (default: false)

`--dir-order`, `--dir-first`, `--group-directories-first`
: none|first|last (default: none)

## Size

`-h`, `--human-readable`
: human sizes, powers of 1024 (default: false)

`--si`
: human sizes, powers of 1000; wins over -h (default: false)

## Time

`--time`
: modified|accessed|changed|birth (default: modified)

`--time-style`
: default|iso|long-iso|full-iso|relative|+FORMAT (default: default)

## Present

`--color`
: always|auto|never (default: auto)

`--icons`
: always|auto|never (default: auto)

`--hyperlink`
: always|auto|never (default: auto)

`--theme`
: path to theme JSON (default: builtin)

`-F`, `--classify`
: always|auto|never; bare -F is always (default: never)

`-Q`, `--quote-name`
: always quote names (default: default)

`-N`, `--literal`
: never quote names (default: default)

`-0`, `--zero`
: NUL record separator (default: false)

`--width`
: screen width for grid/across/comma (default: auto)

## Git

`--git`
: two-character git status column (default: false)

## Deref

`-L`, `--dereference`, `--no-dereference`
: follow symlink/junction metadata (default: false)

<!-- END OPTIONS -->
