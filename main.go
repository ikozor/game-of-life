package main

import (
	"os"
	"time"

	"github.com/ikozor/game-of-life/anim"
	"github.com/ikozor/game-of-life/mtx-impl"
)

func main() {

	dead := '_'
	alive := '*'

	if len(os.Args) == 3 {
	}

	game, err := mtximpl.CreateNewGame("./input_files/glider_gun.txt", dead, alive)
	if err != nil {
		panic(err)
	}
	screen, err := anim.CreateScreen(dead, alive)
	if err != nil {
		panic("Cannot Create Screen")
	}

	for screen.CaptureEscape() {
		screen.UpdateWithMatrix(game.GetCurGen())
		game.CalcNextGen()
		time.Sleep(100 * time.Millisecond)
	}
	screen.Finished()
}
