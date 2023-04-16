package main

import (
	"reflect"
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

func Test_validateMimeType(t *testing.T) {
	cases := map[string]struct {
		mapPath  string
		expected string
	}{
		"normal/map01": {
			mapPath:  "files/test/validateMimeType/map01.txt",
			expected: "",
		},
		"error/map02": {
			mapPath:  "files/test/validateMimeType/map02.txt",
			expected: "MIME Type Validation Error: files/test/validateMimeType/map02.txt; Invalid mime type: application/octet-stream;",
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if result := validateMimeType(tt.mapPath); result != nil {
				assert.EqualErrorf(t, result, tt.expected, "Error should be: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func Test_validateFileSize(t *testing.T) {
	cases := map[string]struct {
		mapPath  string
		expected string
	}{
		"normal/map01": {
			mapPath:  "files/test/validateFileSize/map01.txt",
			expected: "",
		},
		"error/map02": {
			mapPath:  "files/test/validateFileSize/map02.txt",
			expected: "File Size Validation Error: files/test/validateFileSize/map02.txt; File size exceeded: 1049 (Max file size is 1024);",
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if result := validateFileSize(tt.mapPath); result != nil {
				assert.EqualErrorf(t, result, tt.expected, "Error should be: %v, got: %v", tt.expected, result)
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
			mapPath:  "files/test/validateStageMap/map01.txt",
			expected: "",
		},
		"normal/map02": {
			mapPath:  "files/test/validateStageMap/map02.txt",
			expected: "",
		},
		"error/map03": {
			mapPath:  "files/test/validateStageMap/map03.txt",
			expected: "Stage Map Validation Error: files/test/validateStageMap/map03.txt; Make the stage map 20 to 50 columns;",
		},
		"error/map04": {
			mapPath:  "files/test/validateStageMap/map04.txt",
			expected: "Stage Map Validation Error: files/test/validateStageMap/map04.txt; Make the stage map 10 to 20 lines;",
		},
		"error/map05": {
			mapPath:  "files/test/validateStageMap/map05.txt",
			expected: "files/test/validateStageMap/map05.txt: Make the stage map 20 to 50 columns: Stage Map Validation Error",
		},
		"error/map06": {
			mapPath:  "files/test/validateStageMap/map06.txt",
			expected: "Stage Map Validation Error: files/test/validateStageMap/map06.txt; Make the stage map 10 to 20 lines;",
		},
		"error/map07": {
			mapPath: "files/test/validateStageMap/map07.txt",
			expected: "Stage Map Validation Error: the following errors occurred:\n" +
				" -  files/test/validateStageMap/map07.txt; Make the width of the stage map uniform (line 6,8,10);\n" +
				" -  files/test/validateStageMap/map07.txt; Create a boundary for the stage map with '+' (line 1,8,10,15);",
		},
		"error/map08": {
			mapPath: "files/test/validateStageMap/map08.txt",
			expected: "Stage Map Validation Error: the following errors occurred:\n" +
				" -  files/test/validateStageMap/map08.txt; Please set one P;\n" +
				" -  files/test/validateStageMap/map08.txt; Please set one or more enemies;\n" +
				" -  files/test/validateStageMap/map08.txt; Please set one or more targets;",
		},
		"error/map09": {
			mapPath:  "files/test/validateStageMap/map09.txt",
			expected: "Stage Map Validation Error: files/test/validateStageMap/map09.txt; Please set only one P;",
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

func Test_getDigit(t *testing.T) {
	cases := map[string]struct {
		linenum  int
		expected int
	}{
		"1digit": {
			linenum:  1,
			expected: 1,
		},
		"2digit": {
			linenum:  10,
			expected: 2,
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result := getDigit(tt.linenum)
			if tt.expected != result {
				t.Error("expected:", tt.expected, "result:", result)
			}
		})
	}
}

func Test_makeLineNum(t *testing.T) {
	cases := map[string]struct {
		num      int
		maxDigit int
		expected []rune
	}{
		"1digit": {
			num:      1,
			maxDigit: 2,
			expected: []rune{' ', '1', ' '},
		},
		"2digit": {
			num:      10,
			maxDigit: 2,
			expected: []rune{'1', '0', ' '},
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result := makeLineNum(tt.num, tt.maxDigit, tt.maxDigit+1)
			if !reflect.DeepEqual(tt.expected, result) {
				t.Error("expected:", tt.expected, "result:", result)
			}
		})
	}
}
