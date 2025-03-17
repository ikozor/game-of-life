package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ikozor/game-of-life/anim"
	"github.com/ikozor/game-of-life/mtx-impl"
	"gopkg.in/yaml.v3"
)

func main() {

	dead := ' '
	alive := '*'

	input, err := getUserInput()
	if err != nil {
		panic(err)
	}

	game, err := mtximpl.CreateNewGame(input.File, dead, alive)
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
		switch input.Method {
		case "matrix":
			screen.UpdateWithMatrix(game.GetCurGen())
		}
		<-ch
		time.Sleep(10 * time.Millisecond)
	}
	screen.Finished()
}

func getUserInput() (*userInput, error) {
	if len(os.Args) < 2 {
		return &userInput{
			Method: "matrix",
			File:   "./input_files/glider_gun.txt",
		}, nil
	}

	if os.Args[1] == "-c" && len(os.Args) == 3 {

		file, err := os.ReadFile(os.Args[2])
		if err != nil {
			return nil, err
		}
		input := &userInput{}
		err = yaml.Unmarshal(file, &input)
		if err != nil {
			return nil, err
		}
		fmt.Println(input)
		return input, nil

	}
	return nil, nil
}

type userInput struct {
	Method string `yaml:"method"`
	File   string `yaml:"file"`
}
