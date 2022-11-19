package ics

import (
	"golang.org/x/exp/constraints"
)

// Set is a compacted and flattened representation of a set of arithmetic
// intervals.
//
// All elements in the set are sorted and no duplicates are allowed. Each
// interval in the set is formed by pairing adjacent even/odd elements. Values
// with even indices are considered as inclusive lower boundaries, values with
// odd indices are considered as exclusive upper boundaries. If the total number
// of elements is odd, then the last element is treaded as open-ended interval:
//
//	elements: A B C D E F G
//	meaning:  [ ) [ ) [ ) [
type Set[T constraints.Ordered] []T

// Contains returns true if element e passes containment test within the
// interval set s.
func Contains[S ~[]T, T constraints.Ordered](s S, e T) bool {
	var i int
	var ok bool
	if len(s) < linear_search_threshold {
		i, ok = linear_search(s, e)
	} else {
		i, ok = binary_search(s, e)
	}
	return (i&1 == 0) == ok
}

// Hull returns a set that contains at most one interval that covers all
// intervals in s.
func Hull[S ~[]T, T constraints.Ordered](s S) S {
	n := len(s)
	if n <= 2 {
		return s
	} else if n&1 == 1 {
		return s[:1]
	} else {
		return S{s[0], s[n-1]}
	}
}

const linear_search_threshold = 64

func linear_search[S ~[]T, T constraints.Ordered](s S, e T) (int, bool) {
	i, n := 0, len(s)
	for i < n && s[i] < e {
		i++
	}
	return i, i < n && s[i] == e
}

func binary_search[S ~[]T, T constraints.Ordered](s S, e T) (int, bool) {
	i, j, n := 0, len(s), len(s)
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if s[h] < e {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	return i, i < n && s[i] == e
}
