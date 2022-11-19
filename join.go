package ics

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// Join checks if [l,h) overlaps with any of the intervals within s, then hulls
// overlapping intervals by merging them into one, effectively filling the gaps
// between them. This operation is useful when there is an interval [l,h) where
// results of the containment test are unimportant. Using Join with this
// interval then may potentially simplify and reduce the set in a way that still
// produces correct containment tests outside of [l, h).
func Join[S ~[]T, T constraints.Ordered](s *S, l, h T) {
	n := len(*s)
	if n == 0 || (*s)[n-1] <= l {
		return
	}

	if h <= l {
		// joining with open-ended interval [l...
		i, _ := linear_search(*s, l)
		if i&1 == 0 {
			i++
		}
		if n&1 == 1 {
			*s = (*s)[:i]
		} else {
			save_max := (*s)[n-1]
			*s = (*s)[:i+1]
			(*s)[i] = save_max
		}
		return
	} else if h <= (*s)[0] {
		return
	}

	li, _ := linear_search(*s, l)
	hi, h_be := linear_search((*s)[li:], h)
	hi += li

	if li&1 == 0 {
		li++
	}
	h_xe := (hi&1 == 1)
	h_bxe := h_be || h_xe
	h_b := h_be && !h_xe
	if h_b {
		hi++
	} else if !h_bxe {
		hi--
	}
	if hi > li {
		(*s) = slices.Delete(*s, li, hi)
	}
}
