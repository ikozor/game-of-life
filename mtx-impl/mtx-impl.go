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
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	matrix := [][]int8{}
	scanner.Scan()

	// This is similar to creation a "wall" in the MatrixImpl
	// Create a temporary empty wall fill in with data after the first scan is complete
	width := len(scanner.Text()) + 2
	y := 0
	matrix = append(matrix, make([]int8, width))

	y++

	// first run because we already scanned first row to get width
	matrix = append(matrix, make([]int8, width))
	xLen := 2
	for x, e := range scanner.Text() {
		xLen++
		if e == alive {
			matrix[y][x+1] = 1
		} else {
			matrix[y][x+1] = 0
		}
	}
	if xLen != width {
		return nil, fmt.Errorf("Matrix rows is not uniform size, original width: %d, row %d width: %d",
			width,
			y, xLen)
	}
	y++

	// Continue scanning after first row
	for scanner.Scan() {
		xLen = 2
		matrix = append(matrix, make([]int8, width))
		for x, e := range scanner.Text() {
			xLen++
			if e == alive {
				matrix[y][x+1] = 1
			} else {
				matrix[y][x+1] = 0
			}
		}

		if xLen != width {
			return nil, fmt.Errorf("Matrix rows is not uniform size, original width: %d, row %d width: %d",
				width-2,
				y, xLen-2)
		}
		y++
	}

	// Final wall row
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
