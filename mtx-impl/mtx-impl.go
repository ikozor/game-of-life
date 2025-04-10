package mtximpl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ikozor/game-of-life/game"
)

type MatrixImpl struct {
	curGen [][]int8
}

func CreateNewGame(inputPath string, dead, alive rune) (*MatrixImpl, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	matrix := [][]int8{}
	scanner.Scan()

	// to create a "wall"
	width := len(scanner.Text()) + 2
	y := 0
	matrix = append(matrix, make([]int8, width))

	y++

	// first run because we already scanned first row to get width
	matrix = append(matrix, make([]int8, width))
	for x, e := range scanner.Text() {
		if e == alive {
			matrix[y][x+1] = 1
		} else {
			matrix[y][x+1] = 0
		}
	}
	y++

	for scanner.Scan() {
		matrix = append(matrix, make([]int8, width))
		for x, e := range scanner.Text() {
			if e == alive {
				matrix[y][x+1] = 1
			} else {
				matrix[y][x+1] = 0
			}
		}
		y++
	}
	matrix = append(matrix, make([]int8, width))

	return &MatrixImpl{curGen: matrix}, nil
}

func (m *MatrixImpl) CalcNextGen(done chan bool) {
	nextGen := make([][]int8, len(m.curGen))
	for y := 1; y < len(m.curGen)-1; y++ {
		nextGen[y] = make([]int8, len(m.curGen[y]))
		for x := 1; x < len(m.curGen[y])-1; x++ {
			var aliveNearby int8

			aliveNearby += m.curGen[y-1][x-1]
			aliveNearby += m.curGen[y-1][x]
			aliveNearby += m.curGen[y-1][x+1]
			aliveNearby += m.curGen[y][x-1]
			aliveNearby += m.curGen[y][x+1]
			aliveNearby += m.curGen[y+1][x-1]
			aliveNearby += m.curGen[y+1][x]
			aliveNearby += m.curGen[y+1][x+1]

			if aliveNearby == 3 {
				nextGen[y][x] = 1
			} else if aliveNearby < 2 {
				nextGen[y][x] = 0
			} else if aliveNearby > 3 {
				nextGen[y][x] = 0
			} else if aliveNearby == 2 {
				nextGen[y][x] = m.curGen[y][x]
			}
		}
	}
	for i, e := range nextGen {
		copy(m.curGen[i], e)
	}
	done <- true

}

func (m *MatrixImpl) PrintCurGen() {
	fmt.Println()
	for _, e := range m.curGen[1 : len(m.curGen)-1] {
		fmt.Println(e[1 : len(e)-1])
	}
}

func (m *MatrixImpl) GetCurGen() game.GameData {
	return game.GameData{Matrix: m.curGen}
}
