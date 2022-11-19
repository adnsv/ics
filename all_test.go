package ics

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type byteset Set[byte]

func (b byteset) String() string {
	w := bytes.Buffer{}
	Write(&w, b)
	return w.String()
}

func TestContains(t *testing.T) {
	tests := []struct {
		s    string
		v    byte
		want bool
	}{
		// empty
		{"", 'f', false},

		// single open
		{"c", 'b', false},
		{"c", 'c', true},
		{"c", 'd', true},
		{"c", '\xff', true},

		// open for any
		{"\x00", '\x00', true},
		{"\x00", 'a', true},
		{"\x00", 'z', true},
		{"\x00", '\xff', true},

		// single bounded
		{"bd", '\x00', false},
		{"bd", 'a', false},
		{"bd", 'b', true},
		{"bd", 'c', true},
		{"bd", 'd', false},
		{"bd", 'e', false},
		{"bd", '\xff', false},

		// bounded then open
		{"bdf", '\x00', false},
		{"bdf", 'a', false},
		{"bdf", 'b', true},
		{"bdf", 'c', true},
		{"bdf", 'd', false},
		{"bdf", 'e', false},
		{"bdf", 'f', true},
		{"bdf", 'g', true},
		{"bdf", '\xff', true},

		// double bounded
		{"bdfg", '\x00', false},
		{"bdfg", 'a', false},
		{"bdfg", 'b', true},
		{"bdfg", 'c', true},
		{"bdfg", 'd', false},
		{"bdfg", 'e', false},
		{"bdfg", 'f', true},
		{"bdfg", 'g', false},
		{"bdfg", '\xff', false},
	}
	for _, tt := range tests {
		set := byteset(tt.s)
		name := fmt.Sprintf("%s_contains_%d", set, tt.v)
		t.Run(name, func(t *testing.T) {
			if got := Contains(set, tt.v); got != tt.want {
				t.Errorf("set: '%s' contains %d = %v, want %v", set, tt.v, got, tt.want)
			}
		})
	}
}

func Fuzz_search(f *testing.F) {
	f.Add(64, 64)
	f.Fuzz(func(t *testing.T, n_keys, n_vals int) {
		const rand_n = 1024

		m := map[int]struct{}{}
		for i := 0; i < n_keys; i++ {
			k := rand.Intn(rand_n)
			m[k] = struct{}{}
		}
		vec := maps.Keys(m)
		slices.Sort(vec)

		for i := 0; i < n_vals; i++ {
			v := rand.Intn(rand_n)
			idx, ok := slices.BinarySearch(vec, v)
			idx_l, ok_l := linear_search(vec, v)
			idx_b, ok_b := binary_search(vec, v)

			if idx != idx_l || ok != ok_l {
				t.Errorf("linear_search: %v for %v -> (%v, %v), want (%v, %v)", vec, v, idx_l, ok_l, idx, ok)
			}
			if idx != idx_b || ok != ok_b {
				t.Errorf("binary_search: %v for %v -> (%v, %v), want (%v, %v)", vec, v, idx_b, ok_b, idx, ok)
			}
		}
	})
}

func Test_Insert(t *testing.T) {

	check := func(name string, cs byteset, want string) {
		got := cs.String()
		if got != want {
			t.Errorf("%s got %s, want %s", name, got, want)
		}
		//t.Logf("%s succeeded", name)
	}
	prepare := func(cs *byteset, l, h byte, want string) {
		var name string
		if l < h {
			name = fmt.Sprintf("Insert [%v,%v) into %v", l, h, *cs)
		} else {
			name = fmt.Sprintf("Insert [%v... into %v", l, *cs)
		}
		InsertInterval(cs, l, h)
		check(name, *cs, want)
	}
	insert := func(cs byteset, l, h byte, want string) {
		var name string
		if l < h {
			name = fmt.Sprintf("Insert [%v,%v) into %v", l, h, cs)
		} else {
			name = fmt.Sprintf("Insert [%v... into %v", l, cs)
		}
		tmp := slices.Clone(cs)
		InsertInterval(&tmp, l, h)
		check(name, tmp, want)
	}

	var cs byteset

	// ------------------------------
	// TARGET: empty
	// ------------------------------
	cs = byteset{}

	// insert open
	insert(cs, 6, 6, "[6...")

	// insert bounded
	insert(cs, 4, 7, "[4,7)")

	// ------------------------------
	// TARGET: a single open interval
	// ------------------------------
	cs = byteset{}
	prepare(&cs, 6, 6, "[6...")

	// inserting open
	insert(cs, 0, 0, "[0...")
	insert(cs, 4, 4, "[4...")
	insert(cs, 6, 6, "[6...")
	insert(cs, 7, 7, "[6...")

	// inserting bounded
	insert(cs, 0, 5, "[0,5)[6...")
	insert(cs, 0, 6, "[0...")
	insert(cs, 0, 7, "[0...")
	insert(cs, 5, 6, "[5...")
	insert(cs, 6, 7, "[6...")
	insert(cs, 7, 8, "[6...")
	insert(cs, 5, 8, "[5...")
	insert(cs, 5, 8, "[5...")

	// -------------------------------
	// TARGET: single bounded interval
	// -------------------------------
	cs = byteset{}
	prepare(&cs, 3, 8, "[3,8)")

	// inserting open
	insert(cs, 0, 0, "[0...")
	insert(cs, 3, 3, "[3...")
	insert(cs, 5, 5, "[3...")
	insert(cs, 8, 8, "[3...")
	insert(cs, 9, 9, "[3,8)[9...")
	insert(cs, 10, 10, "[3,8)[10...")

	// inserting bounded
	insert(cs, 0, 1, "[0,1)[3,8)")
	insert(cs, 0, 3, "[0,8)")
	insert(cs, 0, 5, "[0,8)")
	insert(cs, 0, 8, "[0,8)")
	insert(cs, 0, 9, "[0,9)")
	//
	insert(cs, 3, 5, "[3,8)")
	insert(cs, 3, 8, "[3,8)")
	insert(cs, 3, 9, "[3,9)")
	//
	insert(cs, 5, 6, "[3,8)")
	insert(cs, 5, 8, "[3,8)")
	insert(cs, 5, 9, "[3,9)")
	//
	insert(cs, 8, 9, "[3,9)")
	insert(cs, 9, 10, "[3,8)[9,10)")

	// -------------------------------------------------------
	// TARGET: a bounded interval followed by an open interval
	// -------------------------------------------------------
	cs = byteset{}
	prepare(&cs, 3, 5, "[3,5)")
	prepare(&cs, 8, 8, "[3,5)[8...")

	// inserting open
	insert(cs, 0, 0, "[0...")
	insert(cs, 2, 2, "[2...")
	insert(cs, 3, 3, "[3...")
	insert(cs, 4, 4, "[3...")
	insert(cs, 5, 5, "[3...")
	insert(cs, 6, 6, "[3,5)[6...")
	insert(cs, 7, 7, "[3,5)[7...")
	insert(cs, 8, 8, "[3,5)[8...")
	insert(cs, 9, 9, "[3,5)[8...")

	// inserting bounded
	insert(cs, 0, 2, "[0,2)[3,5)[8...")
	insert(cs, 0, 3, "[0,5)[8...")
	insert(cs, 0, 4, "[0,5)[8...")
	insert(cs, 0, 5, "[0,5)[8...")
	insert(cs, 0, 6, "[0,6)[8...")
	insert(cs, 0, 7, "[0,7)[8...")
	insert(cs, 0, 8, "[0...")
	insert(cs, 0, 9, "[0...")
	//
	insert(cs, 0, 5, "[0,5)[8...")
	insert(cs, 0, 3, "[0,5)[8...")
	insert(cs, 2, 4, "[2,5)[8...")
	insert(cs, 2, 5, "[2,5)[8...")
	insert(cs, 2, 6, "[2,6)[8...")
	insert(cs, 2, 7, "[2,7)[8...")
	insert(cs, 2, 8, "[2...")
	insert(cs, 3, 8, "[3...")
	insert(cs, 4, 8, "[3...")
	insert(cs, 5, 8, "[3...")
	insert(cs, 6, 8, "[3,5)[6...")
	insert(cs, 7, 8, "[3,5)[7...")
	insert(cs, 8, 10, "[3,5)[8...")

	// -----------------------------------
	// TARGET: a pair of bounded intervals
	// -----------------------------------
	cs = byteset{}
	prepare(&cs, 1, 3, "[1,3)")
	prepare(&cs, 6, 8, "[1,3)[6,8)")

	// insert open
	insert(cs, 0, 0, "[0...")
	insert(cs, 1, 1, "[1...")
	insert(cs, 2, 2, "[1...")
	insert(cs, 3, 3, "[1...")
	insert(cs, 4, 4, "[1,3)[4...")
	insert(cs, 5, 5, "[1,3)[5...")
	insert(cs, 6, 6, "[1,3)[6...")
	insert(cs, 7, 7, "[1,3)[6...")
	insert(cs, 8, 8, "[1,3)[6...")
	insert(cs, 9, 9, "[1,3)[6,8)[9...")
	insert(cs, 10, 10, "[1,3)[6,8)[10...")

	// insert bounded
	insert(cs, 0, 1, "[0,3)[6,8)")
	insert(cs, 0, 3, "[0,3)[6,8)")
	insert(cs, 0, 4, "[0,4)[6,8)")
	insert(cs, 0, 6, "[0,8)")
	insert(cs, 0, 7, "[0,8)")
	insert(cs, 0, 8, "[0,8)")
	insert(cs, 0, 9, "[0,9)")
	//
	insert(cs, 1, 2, "[1,3)[6,8)")
	insert(cs, 1, 3, "[1,3)[6,8)")
	insert(cs, 1, 4, "[1,4)[6,8)")
	insert(cs, 1, 6, "[1,8)")
	insert(cs, 1, 7, "[1,8)")
	insert(cs, 1, 8, "[1,8)")
	insert(cs, 1, 9, "[1,9)")
	//
	insert(cs, 2, 3, "[1,3)[6,8)")
	insert(cs, 2, 4, "[1,4)[6,8)")
	insert(cs, 2, 6, "[1,8)")
	insert(cs, 2, 7, "[1,8)")
	insert(cs, 2, 8, "[1,8)")
	insert(cs, 2, 9, "[1,9)")
	//
	insert(cs, 3, 4, "[1,4)[6,8)")
	insert(cs, 3, 5, "[1,5)[6,8)")
	insert(cs, 3, 6, "[1,8)")
	insert(cs, 3, 7, "[1,8)")
	insert(cs, 3, 8, "[1,8)")
	insert(cs, 3, 9, "[1,9)")
	//
	insert(cs, 4, 5, "[1,3)[4,5)[6,8)")
	insert(cs, 4, 6, "[1,3)[4,8)")
	insert(cs, 4, 7, "[1,3)[4,8)")
	insert(cs, 4, 8, "[1,3)[4,8)")
	insert(cs, 4, 9, "[1,3)[4,9)")
	//
	insert(cs, 5, 6, "[1,3)[5,8)")
	//
	insert(cs, 6, 7, "[1,3)[6,8)")
	insert(cs, 6, 8, "[1,3)[6,8)")
	insert(cs, 6, 9, "[1,3)[6,9)")
	//
	insert(cs, 7, 8, "[1,3)[6,8)")
	insert(cs, 7, 9, "[1,3)[6,9)")
	//
	insert(cs, 8, 9, "[1,3)[6,9)")
	//
	insert(cs, 9, 10, "[1,3)[6,8)[9,10)")
}

func Test_InsertValue(t *testing.T) {

	check := func(name string, cs byteset, want string) {
		got := cs.String()
		if got != want {
			t.Errorf("%s got %s, want %s", name, got, want)
		}
		//t.Logf("%s succeeded", name)
	}
	prepare := func(cs *byteset, l, h byte, want string) {
		var name string
		if l < h {
			name = fmt.Sprintf("Insert [%v,%v) into %v", l, h, *cs)
		} else {
			name = fmt.Sprintf("Insert [%v... into %v", l, *cs)
		}
		InsertInterval(cs, l, h)
		check(name, *cs, want)
	}
	insert := func(cs byteset, v byte, want string) {
		name := fmt.Sprintf("Insert %v into %v", v, cs)
		tmp := slices.Clone(cs)
		InsertInterval(&tmp, v, v+1)
		check(name, tmp, want)
	}

	var cs byteset

	// ------------------------------
	// TARGET: empty
	// ------------------------------
	cs = byteset{}

	insert(cs, 0, "[0,1)")
	insert(cs, 42, "[42,43)")
	insert(cs, 255, "[255...")

	// ------------------------------
	// TARGET: unbounded interval
	// ------------------------------
	cs = byteset{}
	prepare(&cs, 5, 5, "[5...")

	insert(cs, 0, "[0,1)[5...")
	insert(cs, 3, "[3,4)[5...")
	insert(cs, 4, "[4...")
	insert(cs, 5, "[5...")
	insert(cs, 6, "[5...")
	insert(cs, 255, "[5...")

	// ------------------------------
	// TARGET: bounded interval
	// ------------------------------
	cs = byteset{}
	prepare(&cs, 3, 6, "[3,6)")

	insert(cs, 0, "[0,1)[3,6)")
	insert(cs, 3, "[3,6)")
	insert(cs, 4, "[3,6)")
	insert(cs, 6, "[3,7)")
	insert(cs, 8, "[3,6)[8,9)")
	insert(cs, 9, "[3,6)[9,10)")
	insert(cs, 255, "[3,6)[255...")

	// -------------------------------------------------------
	// TARGET: a bounded interval followed by an open interval
	// -------------------------------------------------------
	cs = byteset{}
	prepare(&cs, 3, 6, "[3,6)")
	prepare(&cs, 8, 8, "[3,6)[8...")

	insert(cs, 0, "[0,1)[3,6)[8...")
	insert(cs, 2, "[2,6)[8...")
	insert(cs, 3, "[3,6)[8...")
	insert(cs, 4, "[3,6)[8...")
	insert(cs, 5, "[3,6)[8...")
	insert(cs, 6, "[3,7)[8...")
	insert(cs, 7, "[3,6)[7...")
	insert(cs, 8, "[3,6)[8...")
	insert(cs, 9, "[3,6)[8...")
	insert(cs, 255, "[3,6)[8...")

	// -------------------------------------------------------
	// TARGET: a bounded interval followed by an open interval
	// -------------------------------------------------------
	cs = byteset{}
	prepare(&cs, 3, 6, "[3,6)")
	prepare(&cs, 7, 7, "[3,6)[7...")
	insert(cs, 5, "[3,6)[7...")
	insert(cs, 6, "[3...")
	insert(cs, 7, "[3,6)[7...")
}

func Test_Join(t *testing.T) {
	var cs byteset

	check := func(name string, cs byteset, want string) {
		got := cs.String()
		if got != want {
			t.Errorf("%s got %s, want %s", name, got, want)
		}
		//t.Logf("%s succeeded", name)
	}
	prepare := func(l, h byte, want string) {
		var name string
		if l < h {
			name = fmt.Sprintf("Insert [%v,%v) into %v", l, h, cs)
		} else {
			name = fmt.Sprintf("Insert [%v... into %v", l, cs)
		}
		InsertInterval(&cs, l, h)
		check(name, cs, want)
	}
	join := func(l, h byte, want string) {
		var name string
		if l < h {
			name = fmt.Sprintf("Joining %v with [%v,%v)", cs, l, h)
		} else {
			name = fmt.Sprintf("Joining %v with [%v...", cs, l)
		}
		tmp := slices.Clone(cs)
		Join(&tmp, l, h)
		check(name, tmp, want)
	}

	// ------------------------------
	// TARGET: a single open interval
	// ------------------------------
	cs = byteset{}
	prepare(10, 10, "[10...")

	join(0, 0, "[10...")
	join(10, 10, "[10...")
	join(20, 20, "[10...")
	join(0, 8, "[10...")
	join(0, 10, "[10...")
	join(8, 20, "[10...")
	join(10, 20, "[10...")

	// ------------------------------
	// TARGET: a single bounded interval
	// ------------------------------
	cs = byteset{}
	prepare(10, 20, "[10,20)")

	join(0, 0, "[10,20)")
	join(5, 5, "[10,20)")
	join(10, 10, "[10,20)")
	join(15, 15, "[10,20)")
	join(20, 20, "[10,20)")
	join(25, 25, "[10,20)")
	join(0, 5, "[10,20)")
	join(0, 10, "[10,20)")
	join(0, 15, "[10,20)")
	join(0, 20, "[10,20)")
	join(0, 25, "[10,20)")
	join(5, 10, "[10,20)")
	join(5, 15, "[10,20)")
	join(5, 20, "[10,20)")
	join(5, 25, "[10,20)")
	join(10, 15, "[10,20)")
	join(10, 20, "[10,20)")
	join(10, 25, "[10,20)")

	// ------------------------------
	// TARGET: a bounded interval followed by an open interval
	// ------------------------------
	cs = byteset{}
	prepare(10, 20, "[10,20)")
	prepare(30, 30, "[10,20)[30...")

	join(0, 0, "[10...")
	join(0, 5, "[10,20)[30...")
	join(0, 10, "[10,20)[30...")
	join(0, 15, "[10,20)[30...")
	join(0, 20, "[10,20)[30...")
	join(0, 25, "[10,20)[30...")
	join(0, 30, "[10...")
	join(0, 35, "[10...")

	join(5, 5, "[10...")
	join(5, 10, "[10,20)[30...")
	join(5, 15, "[10,20)[30...")
	join(5, 20, "[10,20)[30...")
	join(5, 25, "[10,20)[30...")
	join(5, 30, "[10...")
	join(5, 35, "[10...")

	join(10, 10, "[10...")
	join(10, 15, "[10,20)[30...")
	join(10, 20, "[10,20)[30...")
	join(10, 25, "[10,20)[30...")
	join(10, 30, "[10...")
	join(10, 35, "[10...")

	join(15, 15, "[10...")
	join(15, 16, "[10,20)[30...")
	join(15, 20, "[10,20)[30...")
	join(15, 20, "[10,20)[30...")
	join(15, 25, "[10,20)[30...")
	join(15, 30, "[10...")
	join(15, 35, "[10...")

	join(20, 20, "[10...")
	join(20, 25, "[10,20)[30...")
	join(20, 30, "[10...")
	join(20, 35, "[10...")

	join(25, 25, "[10,20)[30...")
	join(25, 26, "[10,20)[30...")
	join(25, 30, "[10,20)[30...")
	join(25, 35, "[10,20)[30...")

	join(30, 30, "[10,20)[30...")
	join(30, 35, "[10,20)[30...")

	join(35, 35, "[10,20)[30...")
	join(35, 40, "[10,20)[30...")

	// -----------------------------------
	// TARGET: a pair of bounded intervals
	// -----------------------------------
	cs = byteset{}
	prepare(10, 20, "[10,20)")
	prepare(30, 40, "[10,20)[30,40)")

	join(0, 0, "[10,40)")
	join(0, 10, "[10,20)[30,40)")
	join(0, 15, "[10,20)[30,40)")
	join(0, 20, "[10,20)[30,40)")
	join(0, 25, "[10,20)[30,40)")
	join(0, 30, "[10,40)")
	join(0, 35, "[10,40)")
	join(0, 40, "[10,40)")
	join(0, 45, "[10,40)")

	join(5, 5, "[10,40)")
	join(5, 10, "[10,20)[30,40)")
	join(5, 15, "[10,20)[30,40)")
	join(5, 20, "[10,20)[30,40)")
	join(5, 25, "[10,20)[30,40)")
	join(5, 30, "[10,40)")
	join(5, 35, "[10,40)")
	join(5, 40, "[10,40)")
	join(5, 45, "[10,40)")

	join(10, 10, "[10,40)")
	join(10, 15, "[10,20)[30,40)")
	join(10, 20, "[10,20)[30,40)")
	join(10, 25, "[10,20)[30,40)")
	join(10, 30, "[10,40)")
	join(10, 35, "[10,40)")
	join(10, 40, "[10,40)")
	join(10, 45, "[10,40)")

	join(15, 15, "[10,40)")
	join(15, 16, "[10,20)[30,40)")
	join(15, 20, "[10,20)[30,40)")
	join(15, 25, "[10,20)[30,40)")
	join(15, 30, "[10,40)")
	join(15, 35, "[10,40)")
	join(15, 40, "[10,40)")
	join(15, 45, "[10,40)")

	join(20, 20, "[10,40)")
	join(20, 25, "[10,20)[30,40)")
	join(20, 30, "[10,40)")
	join(20, 35, "[10,40)")
	join(20, 40, "[10,40)")
	join(20, 45, "[10,40)")

	join(25, 25, "[10,20)[30,40)")
	join(25, 26, "[10,20)[30,40)")
	join(25, 30, "[10,20)[30,40)")
	join(25, 35, "[10,20)[30,40)")
	join(25, 40, "[10,20)[30,40)")
	join(25, 45, "[10,20)[30,40)")

	join(30, 30, "[10,20)[30,40)")
	join(30, 35, "[10,20)[30,40)")
	join(30, 40, "[10,20)[30,40)")
	join(30, 45, "[10,20)[30,40)")

	join(35, 35, "[10,20)[30,40)")
	join(35, 36, "[10,20)[30,40)")
	join(35, 40, "[10,20)[30,40)")
	join(35, 45, "[10,20)[30,40)")

	join(40, 40, "[10,20)[30,40)")
	join(40, 45, "[10,20)[30,40)")

	join(45, 45, "[10,20)[30,40)")
	join(45, 50, "[10,20)[30,40)")
}
