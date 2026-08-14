---
name: New flag (rewrite)
about: Add a primary flag or a new enum value on rewrite/v1
---

This template is mandatory for any PR that changes `internal/parse.Specs()`, `Budget`, `interactions.tsv`, or `exceptions.tsv`.

## Gates

- [ ] Gate 1: listing only
- [ ] Gate 2: new value vs new dimension (which?)
- [ ] Gate 3: not a theme/config-only key
- [ ] Gate 4: GNU / eza / lsd names checked; short letter is not reserved
- [ ] Gate 5: budget still 40, or a separate KD revises it
- [ ] Name is **not** in `docs/rejected-flags.md`

## Nine questions

1. Dimension / slot? Matrix row filled?
2. Behavior under each `--format` (`grid` `across` `oneline` `comma` `tree` `json`)? Does json gain a field? Legal in a `-0` record stream?
3. With `-l` and without `-l` (column / prefix / ignored)?
4. Config key, zero value, merge (replace vs append)? Explicit exemption (`meta` / `cli-only`)?
5. Needs `Resolve` (TTY)? Needs `Validate` (exit 2)?
6. linux / darwin / windows (print `-` when missing)?
7. Affects quoting, width, or goldens?
8. Affects exit codes?
9. Help text in `Spec.Help` (feeds man OPTIONS) + one EXAMPLES line?

## Generated artifacts

- [ ] `go generate ./internal/parse`
- [ ] `go test ./internal/parse` is green
- [ ] Key Decision written (accepted + rejected alternative)
