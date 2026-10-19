package demo

import (
	"slices"
	"strings"
	"testing"
)

func intPointer(value int) *int { return &value }

func TestParseRequestValidatesBrowserBoundary(t *testing.T) {
	tooMany := strings.Repeat("0,", MaxItems) + "0"
	tests := []struct {
		name string
		json string
		want string
	}{
		{"valid", `{"algorithm":"quick","values":[3,-1,2]}`, ""},
		{"unknown field", `{"algorithm":"quick","values":[3],"extra":true}`, "unknown field"},
		{"trailing object", `{"algorithm":"quick","values":[3]} {}`, "expected one JSON object"},
		{"unknown algorithm", `{"algorithm":"bogus","values":[3]}`, "supported algorithm"},
		{"empty", `{"algorithm":"quick","values":[]}`, "at least one"},
		{"too many", `{"algorithm":"quick","values":[` + tooMany + `]}`, "no more than"},
		{"value range", `{"algorithm":"quick","values":[1000001]}`, "values must be"},
		{"missing target", `{"algorithm":"binary-search","values":[1,2]}`, "search target"},
		{"sort target", `{"algorithm":"quick","values":[1,2],"target":1}`, "cannot include"},
		{"target range", `{"algorithm":"linear-search","values":[1,2],"target":-1000001}`, "target must be"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ParseRequest([]byte(test.json))
			if test.want == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("got error %v, want text %q", err, test.want)
			}
		})
	}
}

func TestExecuteSortsWithPublicAlgorithms(t *testing.T) {
	input := []int{8, -2, 5, 5, 0, 13, 1}
	want := slices.Clone(input)
	slices.Sort(want)
	for algorithm := range sortAlgorithms {
		t.Run(algorithm, func(t *testing.T) {
			result := Execute(Request{Algorithm: algorithm, Values: input})
			if !slices.Equal(result.Values, want) || !result.ContractVerified {
				t.Fatalf("got %v verified=%v, want %v", result.Values, result.ContractVerified, want)
			}
			if !slices.Equal(input, []int{8, -2, 5, 5, 0, 13, 1}) {
				t.Fatalf("mutated caller input: %v", input)
			}
		})
	}
}

func TestExecuteSearchesWithPublicAlgorithms(t *testing.T) {
	input := []int{9, 2, 7, 2, -1}
	for algorithm := range searchAlgorithms {
		t.Run(algorithm, func(t *testing.T) {
			result := Execute(Request{Algorithm: algorithm, Values: input, Target: intPointer(2)})
			if result.Search == nil || !result.Search.Found || !result.ContractVerified {
				t.Fatalf("unexpected result: %+v", result)
			}
			if algorithm == "linear-search" {
				if result.Search.Index != 1 || result.PreparedAscending {
					t.Fatalf("linear result: %+v", result)
				}
			} else if result.Search.Index != 1 || !result.PreparedAscending || !slices.IsSorted(result.Values) {
				t.Fatalf("sorted-search result: %+v", result)
			}
		})
	}
}
