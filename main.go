package main

import (
	"math/rand"
	"os"

	"github.com/ikozor/game-of-life/anim"
	"golang.org/x/term"
)

func main() {
	// Get terminal dimensions
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width, height = 50, 25
	}
	dead := rune(' ')
	alive := '*'

	if len(os.Args) == 3 {
	}
	screen, err := anim.CreateScreen(width, height)
	if err != nil {
		return
	}
	matrix := make([][]rune, height)
	for y := 0; y < height; y++ {
		matrix[y] = make([]rune, width)
		for x := 0; x < width; x++ {
			matrix[y][x] = dead 
		}
	}

	for screen.CaptureEscape() {
		matrix[rand.Intn(height)][rand.Intn(width)] = alive
		screen.UpdateWithMatrix(matrix)
	}
	screen.Finished()
}
