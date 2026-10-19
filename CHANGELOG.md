# Changelog

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and the project uses [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## 2.0.0 - 2026-10-19 (tag `v2.0.0`)

A tested Go library of classic algorithms and data structures.

### Added

- Go module `github.com/ZachariahAyers/algokit/v2`, standard library only.
- Packages `sorting`, `search`, `numeric`, `array`, `matrix`,
  `combinatorics`, `dp`, `text`, `linkedlist`, `stack`, `queue`, `heap`,
  `tree`, `trie`, `hashtable`, `set` and `graph`, covering the whole catalog.
  `docs/ARCHITECTURE.md` maps each catalog item to its package.
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
