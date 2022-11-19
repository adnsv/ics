package ics

import (
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/exp/slices"
)

func TestRuneSet_AsciiSplit(t *testing.T) {
	tests := []struct {
		name   string
		r      RuneSet
		want_a AsciiSet
		want_r RuneSet
	}{
		{"empty", RuneSet{}, AsciiSet{}, RuneSet{}},
		{"bounded", RuneSet{0x00, 0x01}, AsciiSet{0x00, 0x01}, RuneSet{}},
		{"bounded", RuneSet{0x00, 0x7f}, AsciiSet{0x00, 0x7f}, RuneSet{}},
		{"bounded", RuneSet{0x00, 0x80}, AsciiSet{0x00}, RuneSet{}},
		{"bounded", RuneSet{0x00, 0x90}, AsciiSet{0x00}, RuneSet{0x80, 0x90}},
		{"bounded", RuneSet{0x7f, 0x90}, AsciiSet{0x7f}, RuneSet{0x80, 0x90}},
		{"bounded", RuneSet{0x80, 0x90}, AsciiSet{}, RuneSet{0x80, 0x90}},
		{"bounded", RuneSet{0x85, 0x90}, AsciiSet{}, RuneSet{0x85, 0x90}},
		{"bounded", RuneSet{0x00, 0x01, 0x60, 0x7f}, AsciiSet{0x00, 0x01, 0x60, 0x7f}, RuneSet{}},
		{"bounded", RuneSet{0x00, 0x01, 0x60, 0x80}, AsciiSet{0x00, 0x01, 0x60}, RuneSet{}},
		{"bounded", RuneSet{0x00, 0x01, 0x60, 0x81}, AsciiSet{0x00, 0x01, 0x60}, RuneSet{0x80, 0x81}},
		{"bounded", RuneSet{0x00, 0x01, 0x7f, 0x81}, AsciiSet{0x00, 0x01, 0x7f}, RuneSet{0x80, 0x81}},
		{"bounded", RuneSet{0x00, 0x01, 0x80, 0x81}, AsciiSet{0x00, 0x01}, RuneSet{0x80, 0x81}},
		{"bounded", RuneSet{0x00, 0x01, 0x90, 0x91}, AsciiSet{0x00, 0x01}, RuneSet{0x90, 0x91}},
		{"bounded", RuneSet{0x00, 0x7f, 0x90, 0x91}, AsciiSet{0x00, 0x7f}, RuneSet{0x90, 0x91}},
		{"bounded", RuneSet{0x00, 0x80, 0x90, 0x91}, AsciiSet{0x00}, RuneSet{0x90, 0x91}},
		{"bounded", RuneSet{0x00, 0x81, 0x90, 0x91}, AsciiSet{0x00}, RuneSet{0x80, 0x81, 0x90, 0x91}},
		{"bounded", RuneSet{0x7e, 0x7f, 0x91, 0x92}, AsciiSet{0x7e, 0x7f}, RuneSet{0x91, 0x92}},
		{"bounded", RuneSet{0x7e, 0x80, 0x91, 0x92}, AsciiSet{0x7e}, RuneSet{0x91, 0x92}},
		{"bounded", RuneSet{0x7f, 0x80, 0x91, 0x92}, AsciiSet{0x7f}, RuneSet{0x91, 0x92}},
		{"bounded", RuneSet{0x7f, 0x81, 0x91, 0x92}, AsciiSet{0x7f}, RuneSet{0x80, 0x81, 0x91, 0x92}},
		{"bounded", RuneSet{0x80, 0x81, 0x91, 0x92}, AsciiSet{}, RuneSet{0x80, 0x81, 0x91, 0x92}},
		{"bounded", RuneSet{0x81, 0x91, 0x92, 0x93}, AsciiSet{}, RuneSet{0x81, 0x91, 0x92, 0x93}},
		{"open-ended", RuneSet{0x00}, AsciiSet{0x00}, RuneSet{0x80}},
		{"open-ended", RuneSet{0x7f}, AsciiSet{0x7f}, RuneSet{0x80}},
		{"open-ended", RuneSet{0x80}, AsciiSet{}, RuneSet{0x80}},
		{"open-ended", RuneSet{0x100}, AsciiSet{}, RuneSet{0x100}},
		{"open-ended", RuneSet{0x00, 0x20, 0x7f}, AsciiSet{0x00, 0x20, 0x7f}, RuneSet{0x80}},
		{"open-ended", RuneSet{0x00, 0x20, 0x80}, AsciiSet{0x00, 0x20}, RuneSet{0x80}},
		{"open-ended", RuneSet{0x00, 0x20, 0x90}, AsciiSet{0x00, 0x20}, RuneSet{0x90}},
		{"open-ended", RuneSet{0x00, 0x7f, 0x90}, AsciiSet{0x00, 0x7f}, RuneSet{0x90}},
		{"open-ended", RuneSet{0x00, 0x80, 0x90}, AsciiSet{0x00}, RuneSet{0x90}},
		{"open-ended", RuneSet{0x00, 0x85, 0x90}, AsciiSet{0x00}, RuneSet{0x80, 0x85, 0x90}},
		{"open-ended", RuneSet{0x7f, 0x85, 0x90}, AsciiSet{0x7f}, RuneSet{0x80, 0x85, 0x90}},
		{"open-ended", RuneSet{0x80, 0x85, 0x90}, AsciiSet{}, RuneSet{0x80, 0x85, 0x90}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got_a, got_r := tt.r.AsciiSplit()

			if !slices.Equal(got_a, tt.want_a) || !slices.Equal(got_r, tt.want_r) {
				t.Errorf("RuneSet[%v].AsciiSplit() = [%v] [%v], want [%v] [%v]", tt.r, got_a, got_r, tt.want_a, tt.want_r)
			}
		})
	}
}

func TestRuneSet_String(t *testing.T) {
	tests := []struct {
		s    RuneSet
		want string
	}{
		{RuneSet{}, ""},
		{RuneSet{0x00}, `\x00-\U0010FFFF`},
		{RuneSet{0x00, ' ' + 1}, `\x00- `},
		{RuneSet{'a', 'z' + 1}, `a-z`},
		{RuneSet{'α', 'ω' + 1}, `\u03B1-\u03C9`},
		{RuneSet{0x00, 0x7f + 1}, `\x00-\x7F`},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.s.String()
			name := fmt.Sprintf("%q", string(tt.s))
			if got != tt.want {
				t.Errorf("[%s].String() = %v, want %v", name, got, tt.want)
			} else {
				t.Logf("[%s].String() = %v", name, got)
			}
		})
	}
}

func TestAsciiSet_CountElements(t *testing.T) {
	tests := []struct {
		s    AsciiSet
		want int
	}{
		{AsciiSet{}, 0},
		{AsciiSet{0}, 128},
		{AsciiSet{126}, 2},
		{AsciiSet{127}, 1},
		{AsciiSet{0, 1}, 1},
		{AsciiSet{126, 127}, 1},
		{AsciiSet{0, 10, 117, 127}, 20},
		{AsciiSet{0, 10, 118}, 20},
	}
	for _, tt := range tests {
		t.Run(tt.s.String(), func(t *testing.T) {
			if got := tt.s.CountElements(); got != tt.want {
				t.Errorf("AsciiSet.CountElements() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAsciiSet_Hull(t *testing.T) {
	tests := []struct {
		s    AsciiSet
		want AsciiSet
	}{
		{AsciiSet{}, AsciiSet{}},
		{AsciiSet{0}, AsciiSet{0}},
		{AsciiSet{0, 127}, AsciiSet{0, 127}},
		{AsciiSet{0, 30, 127}, AsciiSet{0}},
		{AsciiSet{30, 40, 127}, AsciiSet{30}},
		{AsciiSet{20, 30, 50, 127}, AsciiSet{20, 127}},
	}
	for _, tt := range tests {
		t.Run(tt.s.String(), func(t *testing.T) {
			if got := tt.s.Hull(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsciiSet.Hull() = %v, want %v", got, tt.want)
			}
		})
	}
}
