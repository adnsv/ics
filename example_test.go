package ics

import (
	"fmt"
)

func ExampleInsert() {
	any_of := func(chars string) (a AsciiSet) {
		for i := range chars {
			a.Insert(chars[i])
		}
		return
	}
	range_ := func(first, last byte) (a AsciiSet) {
		a.InsertRange(first, last)
		return
	}

	show := func(name string, a AsciiSet) {
		a_regex := "[" + a.String() + "]"
		a_set := string(a)
		not_a_set := string(a.Inverted())

		r := RuneSet{}
		a.EnumerateRanges(func(cmin, cmax byte) {
			r.InsertRange(rune(cmin), rune(cmax))
		})

		r_set := string(r)
		not_r_set := string(r.Inverted())

		fmt.Printf("%-8s %-16s %-16q %-18q %-18q %q\n",
			name, a_regex,
			a_set, not_a_set, r_set, not_r_set)
	}

	digits := range_('0', '9')
	lower := range_('a', 'z')
	upper := range_('A', 'Z')
	alphabetic := MergeAsciiSets(lower, upper)
	alphanumeric := MergeAsciiSets(lower, upper, digits)
	word := MergeAsciiSets(alphanumeric, any_of("_"))

	fmt.Printf("%-8s %-16s %-16s %-18s %-18s %s\n", "name", "regex", "AsciiSet", "^AsciiSet", "RuneSet", "^RuneSet")
	fmt.Printf("%-8s %-16s %-16s %-18s %-18s %s\n", "-------", "---------------", "---------------", "-----------------", "-----------------", "-----------------")
	show("alnum", alphanumeric)
	show("alpha", alphabetic)
	show("ascii", range_(0x00, 0x7f))
	show("blank", any_of("\t "))
	show("cntrl", MergeAsciiSets(range_('\x00', '\x1f'), any_of("\x7f")))
	show("digit", digits)
	show("graph", range_(0x21, 0x7e))
	show("lower", lower)
	show("print", range_(' ', '~'))
	show("punct", MergeAsciiSets(range_(0x21, 0x2f), range_(0x3a, 0x40), range_(0x5b, 0x60), range_(0x7b, 0x7e)))
	show("space", any_of("\t\n\v\f\r "))
	show("upper", upper)
	show("word", word)
	show("xdigit", MergeAsciiSets(range_('0', '9'), range_('A', 'F'), range_('a', 'f')))
	show("perl_s", any_of("\t\n\f\r ")) // does not have `\v`
	show("perl_w", MergeAsciiSets(alphanumeric, any_of("_")))

	// Output:
	// name     regex            AsciiSet         ^AsciiSet          RuneSet            ^RuneSet
	// -------  ---------------  ---------------  -----------------  -----------------  -----------------
	// alnum    [0-9A-Za-z]      "0:A[a{"         "\x000:A[a{"       "0:A[a{"           "\x000:A[a{"
	// alpha    [A-Za-z]         "A[a{"           "\x00A[a{"         "A[a{"             "\x00A[a{"
	// ascii    [\x00-\x7F]      "\x00"           ""                 "\x00\u0080"       "\u0080"
	// blank    [\t ]            "\t\n !"         "\x00\t\n !"       "\t\n !"           "\x00\t\n !"
	// cntrl    [\x00-\x1F\x7F]  "\x00 \x7f"      " \x7f"            "\x00 \x7f\u0080"  " \x7f\u0080"
	// digit    [0-9]            "0:"             "\x000:"           "0:"               "\x000:"
	// graph    [!-~]            "!\x7f"          "\x00!\x7f"        "!\x7f"            "\x00!\x7f"
	// lower    [a-z]            "a{"             "\x00a{"           "a{"               "\x00a{"
	// print    [ -~]            " \x7f"          "\x00 \x7f"        " \x7f"            "\x00 \x7f"
	// punct    [!-/:-@[-`{-~]   "!0:A[a{\x7f"    "\x00!0:A[a{\x7f"  "!0:A[a{\x7f"      "\x00!0:A[a{\x7f"
	// space    [\t-\r ]         "\t\x0e !"       "\x00\t\x0e !"     "\t\x0e !"         "\x00\t\x0e !"
	// upper    [A-Z]            "A["             "\x00A["           "A["               "\x00A["
	// word     [0-9A-Z_a-z]     "0:A[_`a{"       "\x000:A[_`a{"     "0:A[_`a{"         "\x000:A[_`a{"
	// xdigit   [0-9A-Fa-f]      "0:AGag"         "\x000:AGag"       "0:AGag"           "\x000:AGag"
	// perl_s   [\t\n\f\r ]      "\t\v\f\x0e !"   "\x00\t\v\f\x0e !" "\t\v\f\x0e !"     "\x00\t\v\f\x0e !"
	// perl_w   [0-9A-Z_a-z]     "0:A[_`a{"       "\x000:A[_`a{"     "0:A[_`a{"         "\x000:A[_`a{"

}
