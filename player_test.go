package main

import (
	"reflect"
	"testing"
)

func Test_squaresAround(t *testing.T) {
	expected := [][]int{
		{0, 0, 0}, {-1, 0, 1}, {1, 0, 1}, {0, 1, 1}, {0, -1, 1},
		{-1, -1, 2}, {0, 2, 2}, {2, 0, 2}, {-1, 1, 2}, {1, 1, 2}, {0, -2, 2}, {1, -1, 2}, {-2, 0, 2},
		{-2, 1, 3}, {1, -2, 3}, {-2, -1, 3}, {-1, 2, 3}, {1, 2, 3}, {2, -1, 3}, {-1, -2, 3}, {2, 1, 3},
		{-2, -2, 4}, {-2, 2, 4}, {2, -2, 4}, {2, 2, 4}}
	result := squaresAround(2)

	if !reflect.DeepEqual(expected, result) {
		t.Error("expected:", expected, "result:", result)
	}
}

func Test_isInputNum(t *testing.T) {
	type expectedValues struct {
		s string
		b bool
	}
	cases := map[string]struct {
		player   player
		arg      rune
		expected expectedValues
	}{
		"Argument cannot be converted to a number.": {
			player: player{
				inputNum: 0,
			},
			arg: 'k',
			expected: expectedValues{
				s: "k",
				b: false,
			},
		},
		"Argument can be converted to a number.": {
			player: player{
				inputNum: 0,
			},
			arg: '2',
			expected: expectedValues{
				s: "2",
				b: true,
			},
		},
		"0 is the number 0.": {
			player: player{
				inputNum: 1,
			},
			arg: '0',
			expected: expectedValues{
				s: "0",
				b: true,
			},
		},
		"0 is the special string 0": {
			player: player{
				inputNum: 0,
			},
			arg: '0',
			expected: expectedValues{
				s: "0",
				b: false,
			},
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s, b := tt.player.isInputNum(tt.arg)
			if s != tt.expected.s || b != tt.expected.b {
				t.Error("expected:", tt.expected.s, tt.expected.b, "result:", s, b)
			}
		})
	}
}