// Package demo provides the input boundary and execution model for the browser demo.
package demo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/ZachariahAyers/algokit/v2/search"
	"github.com/ZachariahAyers/algokit/v2/sorting"
)

// MaxItems keeps quadratic algorithms responsive in a browser tab.
const MaxItems = 96

const maxMagnitude = 1_000_000

// Request describes one browser-demo execution.
type Request struct {
	Algorithm string `json:"algorithm"`
	Values    []int  `json:"values"`
	Target    *int   `json:"target,omitempty"`
}

// SearchResult reports the search contract for one target.
type SearchResult struct {
	Target int  `json:"target"`
	Index  int  `json:"index"`
	Found  bool `json:"found"`
}

// Result contains output produced by a real algokit call.
type Result struct {
	Algorithm          string        `json:"algorithm"`
	Values             []int         `json:"values"`
	ElapsedNanoseconds int64         `json:"elapsedNanoseconds"`
	ContractVerified   bool          `json:"contractVerified"`
	PreparedAscending  bool          `json:"preparedAscending"`
	Search             *SearchResult `json:"search,omitempty"`
}

var sortAlgorithms = map[string]func([]int){
	"bubble":    sorting.Bubble[int],
	"counting":  sorting.Counting[int],
	"heap":      sorting.Heap[int],
	"insertion": sorting.Insertion[int],
	"merge":     sorting.Merge[int],
	"quick":     sorting.Quick[int],
	"selection": sorting.Selection[int],
	"shell":     sorting.Shell[int],
}

var searchAlgorithms = map[string]func([]int, int) (int, bool){
	"binary-search":        search.Binary[int],
	"interpolation-search": search.Interpolation[int],
	"jump-search":          search.Jump[int],
	"linear-search":        search.Linear[int],
}

// ParseRequest validates untrusted JSON from the browser.
func ParseRequest(data []byte) (Request, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var request Request
	if err := decoder.Decode(&request); err != nil {
		return Request{}, fmt.Errorf("invalid request: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return Request{}, errors.New("invalid request: expected one JSON object")
	}
	if _, ok := sortAlgorithms[request.Algorithm]; !ok {
		if _, ok := searchAlgorithms[request.Algorithm]; !ok {
			return Request{}, errors.New("choose a supported algorithm")
		}
	}
	if len(request.Values) == 0 {
		return Request{}, errors.New("enter at least one integer")
	}
	if len(request.Values) > MaxItems {
		return Request{}, fmt.Errorf("enter no more than %d integers", MaxItems)
	}
	for _, value := range request.Values {
		if value < -maxMagnitude || value > maxMagnitude {
			return Request{}, fmt.Errorf("values must be between -%d and %d", maxMagnitude, maxMagnitude)
		}
	}
	_, isSearch := searchAlgorithms[request.Algorithm]
	if isSearch && request.Target == nil {
		return Request{}, errors.New("enter a search target")
	}
	if !isSearch && request.Target != nil {
		return Request{}, errors.New("sort requests cannot include a target")
	}
	if request.Target != nil && (*request.Target < -maxMagnitude || *request.Target > maxMagnitude) {
		return Request{}, fmt.Errorf("target must be between -%d and %d", maxMagnitude, maxMagnitude)
	}
	return request, nil
}

// Execute runs the selected public algokit function on a copy of the input.
func Execute(request Request) Result {
	values := slices.Clone(request.Values)
	if sortFn, ok := sortAlgorithms[request.Algorithm]; ok {
		started := time.Now()
		sortFn(values)
		elapsed := time.Since(started)
		return Result{
			Algorithm:          request.Algorithm,
			Values:             values,
			ElapsedNanoseconds: elapsed.Nanoseconds(),
			ContractVerified:   slices.IsSorted(values),
		}
	}

	searchFn := searchAlgorithms[request.Algorithm]
	prepared := request.Algorithm != "linear-search"
	if prepared {
		slices.Sort(values)
	}
	started := time.Now()
	index, found := searchFn(values, *request.Target)
	elapsed := time.Since(started)
	return Result{
		Algorithm:          request.Algorithm,
		Values:             values,
		ElapsedNanoseconds: elapsed.Nanoseconds(),
		ContractVerified:   searchContractHolds(values, *request.Target, index, found, prepared),
		PreparedAscending:  prepared,
		Search: &SearchResult{
			Target: *request.Target,
			Index:  index,
			Found:  found,
		},
	}
}

func searchContractHolds(values []int, target, index int, found, sorted bool) bool {
	if sorted {
		wantIndex, wantFound := slices.BinarySearch(values, target)
		return index == wantIndex && found == wantFound
	}
	if found {
		return index >= 0 && index < len(values) && values[index] == target &&
			slices.Index(values[:index], target) == -1
	}
	return index == len(values) && slices.Index(values, target) == -1
}
