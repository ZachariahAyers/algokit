# algokit

[![CI](https://github.com/ZachariahAyers/algokit/actions/workflows/ci.yml/badge.svg)](https://github.com/ZachariahAyers/algokit/actions/workflows/ci.yml)

A Go library of classic algorithms and data structures, for developers who want readable, tested implementations.

## Features

- Sorting, searching, dynamic programming, combinatorics, number theory and string algorithms.
- Lists, stack, queue, heap, BST and AVL trees, trie, hash table, set and weighted graph.
- Generic over your own types, with iterators for traversals.
- Every algorithm documents its time and space cost.
- Overflow and bad input reported as errors.
- No dependencies outside the standard library.

## Tech stack

Go 1.26+, standard library, golangci-lint, GitHub Actions.

## Getting started

Requires Go 1.26 or newer.

```bash
go get github.com/ZachariahAyers/algokit/v2
go doc github.com/ZachariahAyers/algokit/v2/tree
```

## Development

```bash
go test -race ./...
go vet ./...
golangci-lint run
```

## Browser demo

The standalone browser lab calls the library's exported sorting and search
functions through a thin Go WebAssembly adapter. Build its static files with
the installed Go toolchain:

```bash
./scripts/build-demo.sh
```

The output is written to `demo/dist/`. Serve that directory from any static
host; no application server or JavaScript algorithm implementation is used.

How it is put together: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## License

GPL-3.0. See [LICENSE](LICENSE).

---

Built by [Zachariah Ayers](https://ayers.repair).
