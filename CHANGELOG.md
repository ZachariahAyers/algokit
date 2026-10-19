# Changelog

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and the project uses [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## 2.0.0 - 2026-10-19 (tag `v2.0.0`)

A rebuild. v1 held no working code; v2 is a tested Go library.

### Added

- Go module `github.com/ZachariahAyers/algokit/v2`, standard library only.
- Packages `sorting`, `search`, `numeric`, `array`, `matrix`,
  `combinatorics`, `dp`, `text`, `linkedlist`, `stack`, `queue`, `heap`,
  `tree`, `trie`, `hashtable`, `set` and `graph`, covering the v1 catalog.
  `docs/ARCHITECTURE.md` maps each v1 item to its v2 home.
- Tests for every algorithm, checked against the standard library, slice and
  map models, or brute force, plus overflow and edge cases.
- CI on every push and pull request: gofmt, `go vet`, golangci-lint and
  race-enabled tests on the two supported Go releases. Monthly Dependabot
  updates.
- `CLAUDE.md` and `.claude/settings.json` for Claude Code.

### Removed

- The fourteen per-language folders, the shared template and the generator
  script. They held comments only.
- Committed IDE settings (`.idea`) and a `.DS_Store` file.

## 1 - 2020-12-11 (tag `v1`)

- A planned catalog of algorithms and data structures, scaffolded as empty
  files in C, C#, C++, Go, Java, JavaScript, Kotlin, PHP, Python, Ruby, Rust,
  Scala, Swift and TypeScript.
- GPL-3.0 license.
