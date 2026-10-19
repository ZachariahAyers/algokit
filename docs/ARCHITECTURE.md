# Architecture

algokit is one Go module, `github.com/ZachariahAyers/algokit/v2`, with no
dependencies outside the standard library. Each package is one topic and can
be imported on its own.

## Packages

| Package | Contents |
| --- | --- |
| `sorting` | Bubble, insertion, selection, Shell, merge, quick, heap, counting and radix sort. All sort in place. |
| `search` | Linear, binary, jump and interpolation search. Sorted searches share the `slices.BinarySearch` contract. |
| `numeric` | Factorial, Fibonacci, Collatz steps, FizzBuzz, primality, sieve, power, integer square root, digit reversal, Roman numerals, English number names. |
| `array` | Flatten, sum, product, mean, median, duplicates, uniqueness, rotation, Fisher-Yates shuffle, move zeroes. |
| `matrix` | In-place 90 degree rotation of a square matrix. |
| `combinatorics` | Cartesian product, combinations, power set and distinct permutations as iterators; combination sum. |
| `dp` | Longest common, increasing and decreasing subsequences, shortest common supersequence, edit distance, 0/1 knapsack. |
| `text` | Palindromes (including Manacher's longest palindrome), anagrams, character frequency, reversal, rotation check, case conversion, Morse, Caesar cipher, polynomial hash. |
| `linkedlist` | Singly and doubly linked lists. |
| `stack`, `queue` | Slice-backed stack; ring-buffer queue. |
| `heap` | Binary heap with a less function; covers min-heap, max-heap and priority queue. |
| `tree` | Ordered maps on an unbalanced BST and an AVL tree, sharing lookups and traversals. |
| `trie` | Prefix tree over runes with sorted prefix queries. |
| `hashtable` | Separate-chaining hash map with a seeded hash and load-factor resizing. |
| `set` | Map-backed set. |
| `graph` | Weighted directed or undirected graph with BFS and DFS. |
| `internal/num` | Numeric type constraints shared by the packages above. |

`text` uses `combinatorics.Permutations` for anagrams; no other package imports
another, apart from `internal/num`.

## Design rules

- **Generics over interfaces.** Algorithms take `cmp.Ordered`, `comparable` or
  a constraint from `internal/num`, so callers keep their own types and pay no
  boxing.
- **Iterators for large or lazy output.** Enumerations and traversals return
  `iter.Seq` or `iter.Seq2`, so a caller can stop early and nothing exponential
  is built up front. Each yielded slice is a fresh allocation.
- **Errors instead of wrapped results.** Integer functions that can overflow
  return `numeric.ErrOverflow`; aggregates of an empty slice return
  `array.ErrEmpty`. Sentinel errors are compared with `errors.Is`.
- **Report, do not panic, on bad indexes.** Container methods return `false`
  for an empty container or an out-of-range index. The two deliberate panics
  are documented: `sorting.Counting` past `MaxCountingRange`, and a `heap.Heap`
  used without a constructor.
- **Deterministic randomness.** `array.Shuffle` takes a `*rand.Rand`; tests
  seed PCG so every run checks the same cases.
- **Complexity in every doc comment.** Each exported function states its time
  and, where it matters, space cost.
- **Zero values work.** Lists, stack, queue, trees, trie and set are usable as
  declared. `heap`, `hashtable` and `graph` need their constructor because
  they carry an ordering, a hash seed or a direction.

## Testing approach

Tests specify behaviour, not implementation. Where an independent oracle
exists the test uses it: sorts and searches are checked against `slices.Sort`
and `slices.BinarySearch`, lists, trees and the hash table against a slice or
map model under thousands of random operations, Manacher and the subsequence
algorithms against brute force on small random inputs. Structural invariants
(BST order, AVL heights and balance, hash table load factor) are checked
directly. Every overflow and domain boundary has a case on each side.

## Mapping from v1

v1 was a scaffold: fourteen language folders holding the same commented list
of planned algorithms, with no implementations. v2 builds that list in Go.
Items that Go's standard library or language already provides were not
re-implemented:

| v1 item | In v2 |
| --- | --- |
| `array/get-max-value`, `get-min-value` | `slices.Max`, `slices.Min` |
| `comparators/*` | Go operators, `cmp.Compare`, `reflect.DeepEqual` |
| `iterable/*` (map, filter, zip, range, generator, ...) | `iter`, `slices`, `maps`, range-over-func |
| `matrix/*` aggregates | `array.Flatten` then the `array` function |
| `number/is-even`, `is-odd`, `get-radian`, `get-bits`, `get-range`, `sum-up-to` | `%`, arithmetic, `math/bits`, `range n`, `n*(n+1)/2` |
| `number/get-steps-to` | `numeric.CollatzSteps` (read as steps of the Collatz map) |
| `number/remove-zeroes`, `shift-zeroes` | `slices.DeleteFunc`, `array.MoveZeroes` |
| `set/get-*-value` | `array` aggregates |
| `string/get-distance` | `dp.EditDistance([]rune(a), []rune(b))` |
| `string/get-duplicates`, `is-unique`, `rotate-string` | `array.Duplicates`, `array.IsUnique`, `array.Rotate` on `[]rune(s)` |
| `string/get-power-set` | `combinatorics.PowerSet([]rune(s))` |
| `string/format-string` | `fmt.Sprintf` |
| `struct/map`, `weak-map`, `weak-set`, `tuple` | built-in `map`, the `weak` package, structs |
| `struct/directed-graph` | `graph.New(true)` |
| `struct/min-heap`, `max-heap`, `priority-queue` | `heap.NewMin`, `heap.NewMax`, `heap.New` |
| `search/breadth-first-search`, `depth-first-search` | `graph.Graph.BFS`, `graph.Graph.DFS` |
| `sys/*` (CLI, files, processes, HTTP, timers, geometry) | out of scope: thin wrappers over `os`, `os/exec`, `net/http`, `time` |

Four v1 names had no definition to build from and were left out:
`number/get-sum-from`, `number/get-duplicates`, `number/is-unique` and
`set/get-min-knapsack`.
