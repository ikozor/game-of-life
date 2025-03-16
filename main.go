package main

import (
	"os"
	"time"

	"github.com/ikozor/game-of-life/anim"
	"github.com/ikozor/game-of-life/mtx-impl"
)

func main() {

	dead := ' '
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

	screen.UpdateWithMatrix(game.GetCurGen())
	ch := make(chan bool)
	for screen.CaptureEscape() {
		go game.CalcNextGen(ch)
		screen.UpdateWithMatrix(game.GetCurGen())
		<- ch
		time.Sleep(10 * time.Millisecond)
	}
	screen.Finished()
}
