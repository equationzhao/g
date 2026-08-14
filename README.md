# g (`rewrite/v1`)

This branch is an **exploration rewrite**. It does not implement listing yet.

- Production 0.31.x stays on `master` (Homebrew `g-ls`, scoop, `go install` of tagged releases).
- Spec: [`docs/rewrite-architecture.md`](docs/rewrite-architecture.md)
- Gates: `just rewrite-gates` or `go test ./internal/parse ./internal/rewritecheck`
- How to add a flag: [`docs/CONTRIBUTING-rewrite.md`](docs/CONTRIBUTING-rewrite.md)

```bash
git checkout rewrite/v1
go test ./...
go build
./g   # exits 2: stub
```

Do not merge this branch to `master` until a later ship decision.
