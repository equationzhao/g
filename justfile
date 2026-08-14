# rewrite/v1 — exploration branch. Production installers stay on master.

rewrite-gates:
    go test -count=1 ./internal/parse ./internal/rewritecheck
    go generate ./internal/parse
    git diff --exit-code -- docs/flag-registry.md docs/rejected-flags.md docs/rewrite-man.md

test:
    CGO_ENABLED=0 go test -count=1 ./...

build:
    CGO_ENABLED=0 go build -o g .
