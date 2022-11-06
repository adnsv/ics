package ics

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// InsertInterval merges an interval into s.
//
//   - if l < h, a bounded interval [l,h) is merged in
//   - if l >= h, a half-open interval [l,... is merged instead
func InsertInterval[S ~[]T, T constraints.Ordered](s *S, l, h T) {
	n := len(*s)

	if h <= l {
		// inserting open-ended interval [l...
		if n == 0 {
			*s = append(*s, l)
			return
		} else if (*s)[n-1] < l {
			if n&1 == 0 {
				*s = append(*s, l)
			}
			return
		}

		i, matched := linear_search(*s, l)
		if i&1 == 0 {
			*s = (*s)[:i+1]
			if !matched {
				(*s)[i] = l
			}
		} else {
			*s = (*s)[:i]
		}
		return
	}

	if n == 0 {
		*s = append(*s, l, h)
	} else if (*s)[n-1] < l {
		if n&1 == 0 {
			*s = append(*s, l, h)
		}
		return
	} else if h < (*s)[0] {
		*s = slices.Insert(*s, 0, l, h)
		return
	}

	//        [       )       [       )
	//  ..z.. b ..x.. e ..z.. b ..x.. e

	li, l_be := linear_search(*s, l)
	l_xe := (li&1 == 1)
	l_bxe := l_be || l_xe
	l_b := l_bxe && !l_xe

	hi, h_be := linear_search((*s)[li:], h)
	hi += li
	h_xe := (hi&1 == 1)
	h_bxe := h_be || h_xe
	h_b := h_be && !h_xe

	if l_bxe || h_bxe {
		if l_b {
			li++
		}
		if h_b {
			hi++
		}
		if l_bxe != h_bxe {
			hi--
		}

		(*s) = slices.Delete(*s, li, hi)
		if !h_bxe {
			(*s)[li] = h
		}
		if !l_bxe {
			(*s)[li] = l
		}
	} else if li == hi {
		*s = slices.Insert(*s, hi, l, h)
	} else {
		if li+2 < hi {
			(*s) = slices.Delete(*s, li+2, hi)
		}
		(*s)[li] = l
		(*s)[li+1] = h
	}
}
