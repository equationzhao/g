# Rejected flags

Cemetery of names that must not come back as `Specs()` entries.
Generated from `internal/parse/testdata/rejected.tsv`.
Do not edit by hand; run `go generate ./internal/parse`.

When someone files "please add `--duplicate`", paste a link to the matching row and close the issue.

| Name | Use instead | Why |
| --- | --- | --- |
| `--duplicate` | fdupes | not listing; hashes whole files |
| `--dup` | fdupes | alias of --duplicate |
| `--fuzzy` | zoxide/fd | persistent index |
| `-f` | zoxide/fd | alias of --fuzzy |
| `--rebuild-index` | (none) | persistent index |
| `--list-index` | (none) | persistent index |
| `--remove-index` | (none) | persistent index |
| `--remove-current-path` | (none) | persistent index |
| `--remove-invalid-path` | (none) | persistent index |
| `--disable-index` | (none) | persistent index |
| `--checksum` | sha256sum | not listing |
| `--checksum-algorithm` | sha256sum | not listing |
| `--charset` | file(1) | not listing |
| `--mime` | file --mime | not listing |
| `--mime-parent` | file --mime | not listing |
| `--only-mime` | file --mime | not listing |
| `--mounts` | findmnt | not listing |
| `--ext` | -I | positive glob is find/fd |
| `--no-ext` | -I | positive glob is find/fd |
| `-M` | -I | alias of --match |
| `--match` | -I | positive glob is find/fd |
| `--no-dir` | --only-files | merged |
| `--file` | --only-files | merged |
| `--before` | find | time filter |
| `--after` | find | time filter |
| `--show-only-hidden` | -a + grep | rare orthogonal |
| `--hidden` | -a + grep | rare orthogonal |
| `--stdin` | xargs | not listing |
| `--init` | completions/ | main must not print shell scripts |
| `--bug` | GitHub issue | not listing |
| `--party` | (none) | toy |
| `--disco` | (none) | toy |
| `--statistic` | (none) | not listing |
| `--footer` | (none) | not listing |
| `--#` | (none) | line numbers are not listing |
| `--total-size` | du | not listing |
| `--no-total-size` | du | not listing |
| `--recursive-size` | du | not listing |
| `--git-detail` | git status | repo dashboard |
| `--git-repo-branch` | git status | repo dashboard |
| `--git-repo-status` | git status | repo dashboard |
| `--tb` | --format=json | structured format is json only |
| `--table` | --format=json | structured format is json only |
| `--md` | --format=json | structured format is json only |
| `--markdown` | --format=json | structured format is json only |
| `--CSV` | --format=json | structured format is json only |
| `--TSV` | --format=json | structured format is json only |
| `--table-style` | (none) | go-pretty |
| `--classic` | --color=never --icons=never | combination |
| `--colorless` | --color=never | alias |
| `--no-color` | --color=never | alias |
| `--ft` | -F | one classify |
| `--file-type` | -F | one classify |
| `--octal-perm` | stat | rwx is enough |
| `--smart-group` | (none) | hides group silently |
| `--flags` | ls -O / chflags | platform trinket |
| `--birth` | --time=birth | merged |
| `--owner` | -l | one bool per column banned |
| `--group` | -l | one bool per column banned |
| `--perm` | -l | one bool per column banned |
| `--size` | -l | one bool per column banned |
| `-O` | -g | use GNU -g |
| `--no-owner` | -g | use GNU -g |
| `--la` | -la | write the combo |
| `--lh` | -l -h | must not imply long |
| `--no-path-transform` | (none) | no magic paths |
| `--np` | (none) | no magic paths |
| `--relative-to` | (none) | path transform is not ls |
| `--fp` | (none) | use json path |
| `--full-path` | (none) | use json path |
| `--extended` | (none) | xattr is not the main path |
| `-@` | (none) | xattr is not the main path |
| `--tree-style` | (none) | locale picks the line set |
| `--sort-by-mime` | (none) | mime cut |
| `--sort-by-mime-parent` | (none) | mime cut |
| `--sort-by-mime-desc` | (none) | mime cut |
| `--sort-by-mime-parent-desc` | (none) | mime cut |
| `--git-repos` | (none) | not a per-file column |
| `--color-scale` | theme size colors | not a flag |
| `--icon-theme` | icon_set config key | not a flag |
| `--blocks-as-sort` | --sort=size | old --width sort key |
