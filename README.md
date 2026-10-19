# algokit

[![CI](https://github.com/ZachariahAyers/algokit/actions/workflows/ci.yml/badge.svg)](https://github.com/ZachariahAyers/algokit/actions/workflows/ci.yml)

A Go library of classic algorithms and data structures, for developers who want readable, tested implementations.

## Features

- Sorting, searching, dynamic programming, combinatorics, number theory and string algorithms.
- Lists, stack, queue, heap, BST and AVL trees, trie, hash table, set and weighted graph.
- Generic over your own types, with iterators for traversals.
- Time and space cost documented on every function.
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

How it is put together: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## License

GPL-3.0. See [LICENSE](LICENSE).

---

Built by [Zachariah Ayers](https://ayers.repair).
