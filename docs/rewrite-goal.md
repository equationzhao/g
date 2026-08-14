# Goal: finish the g rewrite on `rewrite/v1`

Host command (paste in the TUI):

```
/goal Finish the g rewrite on rewrite/v1 per docs/rewrite-architecture.md. Stay on rewrite/v1; never merge to master or change Homebrew/scoop/install.sh. Implement the remaining PR plan in order (PR1 parser through PR18 exploration build). Done only when all of these are true and independently reproducible: (1) `g` lists directories — default TTY grid, `-l` long, `-T` tree; (2) GNU shorts `-l -a -A -h -1 -R -t -S -r -d -F` match §4; (3) `len(parse.Specs())==40` and `just rewrite-gates` is green; (4) no package-level mutable state in rewrite packages, no os.Args writes, no os.Chdir, no gomonkey; (5) main calls app.Run, not a stub; (6) git status is an optional fail-open column; (7) third-party deps are only the six allowed modules. Do not add flags beyond the 40. Do not revive names in docs/rejected-flags.md.
```

## Already done

- Design: `docs/rewrite-architecture.md`
- Gates: `internal/parse`, `internal/rewritecheck`, `rewrite-gates` CI
- v0 tree deleted on this branch; stub `main` exits 2
- Pushed: `origin/rewrite/v1`

## Remaining (PR plan)

1. `Request` + GNU argv `Parse` / `Resolve` / `Validate`
2. Filesystem seam + `entry` + `sys.Meta`
3. `collect.Walk`
4. filter + sort
5. text printers + quoting
6. long columns
7. tree
8. Theme / color / icons
9. Git seam
10. JSON printer
11. XDG config merge
12. OSC-8 + real NUL
13. `app.Run`
14. CI matrix on this branch
15. man from Specs + 10-usage README
16. point `main` at `app.Run`
17. (already deleted v0)
18. CGO-off just recipes only — no brew/scoop cutover

## Hard stops

- Product: listing only. Not find/du/fd/zoxide/TUI.
- Flag budget 40. New flag = five gates + matrix cell.
- Work stays on `rewrite/v1`.
