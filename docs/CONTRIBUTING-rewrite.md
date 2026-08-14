# Contributing to the g rewrite (`rewrite/v1`)

Deleting a released flag is a breaking change and a major version bump.
That is why adding a flag is hard. Read this before opening an issue or PR.

Production `master` (v0.31.x), Homebrew `g-ls`, scoop, and `script/install.sh` are **out of scope**. All rewrite work stays on `rewrite/v1`.

## Five gates (any one fails → reject)

1. **Is this listing?** If `find`, `du`, `stat`, or `xargs` already does it, no.
2. **New dimension, or a new value on an existing one?** Prefer a new `--sort=` / `--time=` / `--format=` value (budget 0).
3. **Can it not be a flag?** Order: theme JSON key → config-only key → new enum value → new flag. Only things that must change **per invocation** belong on argv.
4. **Does GNU / eza / lsd already have it?** Match GNU names and semantics first. A made-up flag has a higher burden of proof. Short letters: see `docs/flag-registry.md`. GNU-reserved letters are not free.
5. **Budget.** `parse.Budget == 40` is a test. One new dimension means one-in-one-out, or a new Key Decision that revises the cap. Do not bump `Budget` in the same PR that adds the flag.

Rejected names live in `docs/rejected-flags.md`. Search there before debating an old request again.

## After a flag is accepted

1. Add a `Spec` in `internal/parse/specs.go`.
2. Fill every empty cell this creates: dimension row/column in `internal/parse/testdata/interactions.tsv`, plus any flag-level row in `exceptions.tsv`.
3. Answer every question in `.github/PULL_REQUEST_TEMPLATE/new-flag.md`.
4. Run `go generate ./internal/parse` and `go test ./internal/parse`.
5. Write a Key Decision with the rejected alternative.

CI fails the PR if `Specs()` is not 40, if the matrix has a hole, if a rejected name is revived, if generated docs are stale, or if a spec is missing from the man OPTIONS block.

## Tree gates (not just flags)

`go test ./internal/rewritecheck` also fails the PR when a **rewrite** package (listed in `internal/rewritecheck/testdata/packages.tsv`):

- imports a banned module (urfave, gomonkey, go-pretty, …) or any **legacy** `internal/*` package
- imports a third-party module outside the six-dep allowlist
- declares package-level mutable state (except `Err*` and immutable lookup tables)
- has `init()`, writes `os.Args`, or calls `os.Chdir`
- uses feature build tags `fuzzy` / `mounts` / `lite`
- has tests that import `net` / `net/http`

A new `internal/foo` directory that is in neither `rewrite` nor `legacy` fails CI. Promote it to `rewrite` when you add it.

`go.mod` must stay `module github.com/Equationzhao/g`. `LICENSE` must stay MIT.

These do **not** yet apply to legacy packages; they still carry the 0.31 tree until PR17.

## Versioning

- New enum value or new flag: minor (after 1.0).
- Change of published semantics (defaults, sort direction, interaction class): major.
- The diff of `interactions.tsv` is the “behavior changes” section of the release notes.
