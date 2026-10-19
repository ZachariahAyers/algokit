# algokit

Classic algorithms and data structures as one Go module with no dependencies
outside the standard library. Package map and design rules:
`docs/ARCHITECTURE.md`.

## Commands

- Test: `go test ./...` (CI adds `-race -shuffle=on`)
- One package: `go test ./sorting/`; one test: `go test ./tree/ -run TestAVL`
- Vet: `go vet ./...`
- Lint: `golangci-lint run` (v2.14.0, config in `.golangci.yml`)
- Format: `golangci-lint fmt` (gofumpt and goimports); CI also fails on `gofmt -l .`

Go 1.26 or newer; `go.mod` sets the minimum.

## Conventions

- One topic per package. Shared numeric constraints live in `internal/num`;
  never redefine them.
- Every exported function has a doc comment that states its behaviour, its
  edge cases and its time and space complexity.
- Generic type parameters, not `any` plus type switches. Return `iter.Seq` for
  enumerations and traversals.
- Report bad input with a sentinel error or a `false` result. Panic only where
  the doc comment says so.
- Do not reimplement what the standard library does (`slices.Max`,
  `maps.Keys`). Add the mapping to the table in `docs/ARCHITECTURE.md`.
- Tests compare against an independent oracle (stdlib, a map or slice model,
  brute force) with a fixed PCG seed, and cover empty, single, duplicate and
  overflow cases.
- Commits: Conventional Commits, lower case after the type, scope is the
  package (`feat(tree): ...`).

## Definition of done

1. `go vet ./...`, `go test -race ./...` and `golangci-lint run` all pass.
2. New behaviour has a test that fails without the change.
3. Doc comments state complexity; `docs/ARCHITECTURE.md` lists any new package.
4. `CHANGELOG.md` has an entry under Unreleased for user-visible changes.
