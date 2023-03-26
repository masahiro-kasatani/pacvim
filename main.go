package main

import (
	"bytes"
	"embed"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	termbox "github.com/nsf/termbox-go"
)

const (
	pose int = iota
	continuing
	quit
	win
	lose

	maxNumOfGhosts = 4

	chGhost  = 'G'
	chTarget = 'o'
	chPoison = 'X'
	chWall1  = '#'
	chWall2  = '|'
	chWall3  = '-'

	sceneDir = "files/scene/"
)

var (
	gameState           = 0
	targetScore         = 0
	score               = 0
	level               = 1
	life                = 3
	inputNum            = 0
	isLowercaseGEntered = false
	gameSpeed           = time.Second

	// For command gg or G
	firstTargetY int
	lastTargetY  int

	//go:embed files
	static embed.FS
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	err := termbox.Init()
	if err != nil {
		return err
	}

	if err = termbox.Clear(termbox.ColorWhite, termbox.ColorBlack); err != nil {
		return err
	}

	defer termbox.Close()

	stageMaps, err := dirwalk("./files/stage")
	if err != nil {
		return err
	}
	maxLevel := len(stageMaps)

	// スタート画面表示
	if err := switchScene(sceneDir + "start.txt"); err != nil {
		return err
	}

game:
	for {
		if err := stage(stageMaps[level]); err != nil {
			return err
		}
		switch gameState {
		case win:
			if err := switchScene(sceneDir + "youwin.txt"); err != nil {
				return err
			}
			level++
			gameSpeed = time.Duration(1000-(level-1)*50) * time.Millisecond
			if level == maxLevel {
				err = switchScene(sceneDir + "congrats.txt")
				if err != nil {
					return err
				}
				break game
			}
		case lose:
			err = switchScene(sceneDir + "youlose.txt")
			if err != nil {
				return err
			}
			life--
			if life < 0 {
				break game
			}
		case quit:
			break game
		}
	}
	err = switchScene(sceneDir + "goodbye.txt")
	if err != nil {
		return err
	}

	return err
}

func stage(stageMap string) error {
	b, w, err := initScene(stageMap)
	if err != nil {
		return err
	}
	if err = w.show(b); err != nil {
		return err
	}

	// ゲーム情報の初期化
	gameState = pose
	score = 0
	targetScore = 0
	b.plotStageMap()

	// プレイヤー初期化
	p := new(player)
	p.initPosition(b)

	// ゲーム情報の表示
	b.plotScore()
	b.plotSubInfo()

	// ゴーストを作成
	gList, err := initPosition(b)
	if err != nil {
		return err
	}
	// ステージマップを表示
	if err = termbox.Flush(); err != nil {
		return err
	}

	// ゲーム開始待ち状態
	standBy()

	// プレイヤーゴルーチン開始
	ch1 := make(chan bool)
	go p.control(ch1, b, w)

	// ゴーストゴルーチン開始
	ch2 := make(chan bool)
	go control(ch2, p, gList)

	// プレイヤーとゴーストの同期を取る
	<-ch1
	<-ch2

	return err
}

func switchScene(fileName string) error {
	termbox.HideCursor()
	_, w, err := initScene(fileName)
	if err != nil {
		return err
	}
	if err = termbox.Clear(termbox.ColorWhite, termbox.ColorBlack); err != nil {
		return err
	}
	for y, l := range w.lines {
		for x, r := range l.text {
			termbox.SetCell(x, y, r, termbox.ColorYellow, termbox.ColorBlack)
		}
	}
	if err = termbox.Flush(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	return err
}

func initScene(fileName string) (*buffer, *window, error) {
	f, err := static.ReadFile(fileName)
	if err != nil {
		return nil, nil, err
	}

	b := new(buffer)
	b.save(bytes.NewReader(f))

	w := new(window)
	w.copy(b)

	return b, w, err
}

func standBy() {
	for {
		ev := termbox.PollEvent()
		if ev.Key == termbox.KeyEnter {
			gameState = continuing
			break
		}
		if ev.Ch == 'q' {
			gameState = quit
			break
		}
	}
}

func dirwalk(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, file := range files {
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths, err
}

func random(min, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(max-min+1) + min
}
