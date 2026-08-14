# Flag short-letter registry

Generated from `internal/parse.Specs()` and `GNUReservedShorts()`.
Do not edit by hand; run `go generate ./internal/parse`.

Letters GNU ls already assigned stay **reserved** even when we do not implement that feature.
Inventing a new meaning for a reserved letter is forbidden (docs/rewrite-architecture.md §17.3).

## Used

| Letter | Spec | Slot |
| --- | --- | --- |
| `-a` | `--all` | `visibility` |
| `-d` | `--directory` | `dirself` |
| `-g` | `-g` | `noowner` |
| `-h` | `--human-readable` | `human` |
| `-i` | `--inode` | `inode` |
| `-l` | `--long` | `long` |
| `-m` | `--format` | `format` |
| `-n` | `--numeric-uid-gid` | `numeric` |
| `-o` | `--long` | `long` |
| `-r` | `--reverse` | `reverse` |
| `-t` | `--sort` | `sortkey` |
| `-v` | `--sort` | `sortkey` |
| `-x` | `--format` | `format` |
| `-A` | `--almost-all` | `visibility` |
| `-B` | `--ignore-backups` | `backups` |
| `-C` | `--format` | `format` |
| `-D` | `--only-dirs` | `kind` |
| `-F` | `--classify` | `classify` |
| `-G` | `--no-group` | `nogroup` |
| `-H` | `--links` | `nlink` |
| `-I` | `--ignore` | `ignore` |
| `-L` | `--dereference` | `deref` |
| `-N` | `--literal` | `quote` |
| `-Q` | `--quote-name` | `quote` |
| `-R` | `--recursive` | `recurse` |
| `-S` | `--sort` | `sortkey` |
| `-T` | `--format` | `format` |
| `-U` | `--sort` | `sortkey` |
| `-X` | `--sort` | `sortkey` |
| `-0` | `--zero` | `zero` |
| `-1` | `--format` | `format` |
| `-?` | `--help` | `help` |

## Reserved (GNU meaning, not ours)

| Letter | GNU meaning |
| --- | --- |
| `-b` | GNU --escape |
| `-c` | GNU ctime / sort by ctime |
| `-f` | GNU disable sort (and -l) |
| `-k` | GNU block-size=1K |
| `-p` | GNU classify directories with / |
| `-q` | GNU hide control chars |
| `-s` | GNU size in blocks |
| `-u` | GNU atime |
| `-w` | GNU --width |
| `-P` | GNU --no-dereference |
| `-Z` | GNU --context |

## Free

`-e`, `-j`, `-y`, `-z`, `-E`, `-J`, `-K`, `-M`, `-O`, `-V`, `-W`, `-Y`, `-#`

A free letter is not permission to add a flag. See §17 five gates.
