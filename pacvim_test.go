package main

import (
	"testing"

	termbox "github.com/nsf/termbox-go"
	"github.com/stretchr/testify/assert"
)

func Test_switchScene(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	if err := termbox.Init(); err != nil {
		t.Error(err)
	}
	if err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack); err != nil {
		t.Error(err)
	}
	defer termbox.Close()
	scenes := []string{
		sceneStart,
		sceneYouwin,
		sceneYoulose,
		sceneCongrats,
		sceneGoodbye,
	}
	for _, s := range scenes {
		if err := switchScene(s); err != nil {
			t.Error(err)
		}
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

func Test_validateStageMap(t *testing.T) {
	cases := map[string]struct {
		mapPath  string
		expected string
	}{
		"normal/map01": {
			mapPath:  "files/test/map01.txt",
			expected: "",
		},
		"normal/map02": {
			mapPath:  "files/test/map02.txt",
			expected: "",
		},
		"error/map03": {
			mapPath:  "files/test/map03.txt",
			expected: "files/test/map03.txt: Make the stage map 20 to 50 columns: Stage Map Validation Error",
		},
		"error/map04": {
			mapPath:  "files/test/map04.txt",
			expected: "files/test/map04.txt: Make the stage map 10 to 20 lines: Stage Map Validation Error",
		},
		"error/map05": {
			mapPath:  "files/test/map05.txt",
			expected: "files/test/map05.txt: Make the stage map 20 to 50 columns: Stage Map Validation Error",
		},
		"error/map06": {
			mapPath:  "files/test/map06.txt",
			expected: "files/test/map06.txt: Make the stage map 10 to 20 lines: Stage Map Validation Error",
		},
		"error/map07": {
			mapPath: "files/test/map07.txt",
			expected: "the following errors occurred:\n" +
				" -  files/test/map07.txt: Make the width of the stage map uniform (line 6,8,10)\n" +
				" -  files/test/map07.txt: Create a boundary for the stage map with '+' (line 1,8,10,15): Stage Map Validation Error",
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if result := validateStageMap(tt.mapPath); result != nil {
				assert.EqualErrorf(t, result, tt.expected, "Error should be: %v, got: %v", tt.expected, result)
			}
		})
	}

}
