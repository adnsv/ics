package ics

import (
	"bytes"
	"unicode"
	"unicode/utf8"
)

// RuneSet is a containment set for unicode codepoints [U+00..U+10FFFF].
type RuneSet Set[rune]

// AsciiSet is a containment set for 7-bit ASCII characters [\x00..\x7f].
type AsciiSet Set[byte]

// Inverted returns a containment set with inverted logic.
func (s RuneSet) Inverted() RuneSet {
	if len(s) == 0 {
		return RuneSet{0}
	} else if s[0] == 0 {
		return s[1:]
	} else {
		return append(RuneSet{0}, s...)
	}
}

// Insert adds r to the set.
func (s *RuneSet) Insert(r rune) {
	if r < 0 || r > utf8.MaxRune {
		panic("unsupported rune value")
	}
	if r == utf8.MaxRune {
		InsertInterval(s, r, r)
	} else {
		InsertInterval(s, r, r+1)
	}
}

// InsertRange inserts an inclusive [rmin,rmax] range of unicode codepoints into
// the set. Notice, that the inserted ranges are fully inclusive on both ends,
// unlike intervals which are always half open.
func (s *RuneSet) InsertRange(rmin, rmax rune) {
	if rmax < rmin {
		panic("invalid rune range")
	}
	if rmin < 0 || rmax > utf8.MaxRune {
		panic("unsupported rune value")
	}
	if rmax == utf8.MaxRune {
		// insert open-ended interval
		InsertInterval(s, rmin, rmin)
	} else {
		// insert fully-bound interval
		InsertInterval(s, rmin, rmax+1)
	}
}

// EnumerateRanges is a functional enumerator for all the continuous inclusive
// [rmin,rmax] ranges contained within the set.
func (s RuneSet) EnumerateRanges(f func(rmin, rmax rune)) {
	i, n := 0, len(s)
	for i+1 < n {
		f(s[i], s[i+1]-1)
		i += 2
	}
	if i < n {
		f(s[i], utf8.MaxRune)
	}
}

// MergeRuneSets combines multiple containment sets into one.
func MergeRuneSets(ss ...RuneSet) (m RuneSet) {
	for _, s := range ss {
		s.EnumerateRanges(func(rmin, rmax rune) {
			m.InsertRange(rmin, rmax)
		})
	}
	return
}

func print_rune(w *bytes.Buffer, r rune) {
	switch r {
	case '\b':
		w.WriteString("\\b")
	case '\f':
		w.WriteString("\\f")
	case '\n':
		w.WriteString("\\n")
	case '\r':
		w.WriteString("\\r")
	case '\t':
		w.WriteString("\\t")
	default:
		if r > unicode.MaxRune {
			w.WriteString("#INVALID")
		} else if r > 0xffff {
			b := [10]byte{'\\', 'U',
				'0',
				'0',
				hex[(r>>20)&0xf],
				hex[(r>>16)&0xf],
				hex[(r>>12)&0xf],
				hex[(r>>8)&0xf],
				hex[(r>>4)&0xf],
				hex[(r>>0)&0xf],
			}
			w.Write(b[:10])
		} else if r > 0x7f {
			b := [6]byte{'\\', 'u',
				hex[(r>>12)&0xf],
				hex[(r>>8)&0xf],
				hex[(r>>4)&0xf],
				hex[(r>>0)&0xf],
			}
			w.Write(b[:6])
		} else if r < 0x20 || r == 0x7f {
			b := [4]byte{'\\', 'x', hex[(r>>4)&0xf], hex[(r>>0)&0xf]}
			w.Write(b[:4])
		} else {
			w.WriteRune(r)
		}
	}
}

// String produces a human-readable string with ranges and elements.
func (s RuneSet) String() string {
	w := bytes.Buffer{}
	s.EnumerateRanges(func(rmin, rmax rune) {
		print_rune(&w, rmin)
		if rmax > rmin {
			if rmax > rmin+1 {
				w.WriteByte('-')
			}
			print_rune(&w, rmax)
		}
	})
	return w.String()
}

// AsciiSplit splits r into ascii and non-ascii matchers.
func (s RuneSet) AsciiSplit() (AsciiSet, RuneSet) {
	n := len(s)
	if n == 0 {
		return AsciiSet{}, RuneSet{}
	}
	if s[0] >= 0x80 {
		return AsciiSet{}, s
	}
	if s[n-1] < 0x80 {
		a := make(AsciiSet, n)
		for i := 0; i < n; i++ {
			a[i] = byte(s[i])
		}
		if n&1 == 0 {
			return a, RuneSet{}
		} else {
			return a, RuneSet{0x80}
		}
	}

	var be bool
	n, be = binary_search(s, 0x80)

	a := make(AsciiSet, n)
	for i := 0; i < n; i++ {
		a[i] = byte(s[i])
	}

	zb := n&1 == 0
	if be && !zb {
		return a, s[n+1:]
	} else if !be && !zb {
		return a, append(RuneSet{0x80}, s[n:]...)
	} else {
		return a, s[n:]
	}
}

// Inverted returns a containment set with inverted logic.
func (s AsciiSet) Inverted() AsciiSet {
	if len(s) == 0 {
		return AsciiSet{0}
	} else if s[0] == 0 {
		return s[1:]
	} else {
		return append(AsciiSet{0}, s...)
	}
}

// Insert adds c to the set.
func (s *AsciiSet) Insert(c byte) {
	if c > 0x7f {
		panic("invalid ascii value")
	}
	if c == 0x7f {
		InsertInterval(s, c, c)
	} else {
		InsertInterval(s, c, c+1)
	}
}

// InsertRange inserts an inclusive [cmin,cmax] range of ascii characters into
// the set. Notice, that the inserted ranges are fully inclusive on both ends,
// unlike intervals which are always half open.
func (s *AsciiSet) InsertRange(cmin, cmax byte) {
	if cmax < cmin || cmax > 0x7f {
		panic("invalid ascii range")
	}
	if cmax == 0x7f {
		// insert open-ended interval
		InsertInterval(s, cmin, cmin)
	} else {
		// insert fully-bound interval
		InsertInterval(s, cmin, cmax+1)
	}
}

// EnumerateRanges is a functional enumerator for all the continuous inclusive
// [cmin,cmax] ranges contained within the set.
func (s AsciiSet) EnumerateRanges(f func(cmin, cmax byte)) {
	i, n := 0, len(s)
	for i+1 < n {
		f(s[i], s[i+1]-1)
		i += 2
	}
	if i < n {
		f(s[i], 0x7f)
	}
}

// MergeRuneSets combines multiple containment sets into one.
func MergeAsciiSets(ss ...AsciiSet) (m AsciiSet) {
	for _, s := range ss {
		s.EnumerateRanges(func(rmin, rmax byte) {
			m.InsertRange(rmin, rmax)
		})
	}
	return
}

func print_ascii(w *bytes.Buffer, c byte) {
	switch c {
	case '\b':
		w.WriteString("\\b")
	case '\f':
		w.WriteString("\\f")
	case '\n':
		w.WriteString("\\n")
	case '\r':
		w.WriteString("\\r")
	case '\t':
		w.WriteString("\\t")
	default:
		if c >= 0x80 {
			w.WriteString("#INVALID")
		} else if c >= 0x20 && c < 0x7f {
			w.WriteByte(c)
		} else {
			b := [4]byte{'\\', 'x', hex[(c>>4)&0xf], hex[(c>>0)&0xf]}
			w.Write(b[:4])
		}
	}
}

// String produces a human-readable string with ranges and elements.
func (a AsciiSet) String() string {
	w := bytes.Buffer{}
	a.EnumerateRanges(func(cmin, cmax byte) {
		print_ascii(&w, cmin)
		if cmax > cmin {
			if cmax > cmin+1 {
				w.WriteByte('-')
			}
			print_ascii(&w, cmax)
		}
	})
	return w.String()
}

const hex = "0123456789ABCDEF"
