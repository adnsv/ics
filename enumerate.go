package ics

import (
	"fmt"
	"io"
)

// Enumerate is a generic functional enumerator for a set. The callback is
// called with half-open boundaries for each interval within the set. If s has
// an odd number of elements, meaning the the last element is an open-ended
// interval, the callback is called with l = h.
func Enumerate[S ~[]T, T any](s S, f func(l, h T)) {
	i, n := 0, len(s)
	for i+1 < n {
		f(s[i], s[i+1])
		i += 2
	}
	if i < n {
		f(s[i], s[i])
	}
}

// Write writes interval to w in a human-readable form.
func Write[S ~[]T, T any](w io.Writer, s S) {
	i, n := 0, len(s)
	for i+1 < n {
		fmt.Fprintf(w, "[%v,%v)", s[i], s[i+1])
		i += 2
	}
	if i < n {
		fmt.Fprintf(w, "[%v...", s[i])
	}
}
