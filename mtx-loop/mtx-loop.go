package mtxloop

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ikozor/game-of-life/game"
)

type MatrixLoop struct {
	curGen [][]int8
}

func CreateNewGame(inputPath string, dead, alive rune) (*MatrixLoop, error) {
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

	UpdateLoopingWall(matrix)

	return &MatrixLoop{curGen: matrix}, nil
}

func UpdateLoopingWall(matrix [][]int8) {
	var (
		firstWallXLoc = 0
		firstWallYLoc = 0
		firstGameYLoc = 1
		firstGameXLoc = 1
		lastWallYLoc  = len(matrix) - 1
		lastWallXLoc  = len(matrix[0]) - 1
		lastGameYLoc  = len(matrix) - 2
		lastGameXLoc  = len(matrix[0]) - 2
	)

	// After all the data is scanned, begin to fill in the wall with data to make the loop work
	// Fill each value with the opposite position of the data (e.g. data in 1,1 (or 0,0 of game data) will go to lastX, lastY of the buffer wall
	for index, row := range matrix {
		switch index {
		case firstWallYLoc:
			// copy bottom row with corners if we are on the upper wall
			matrix[firstWallYLoc][firstWallXLoc] = matrix[lastGameYLoc][lastGameXLoc]
			matrix[firstWallYLoc][lastWallXLoc] = matrix[lastGameYLoc][firstGameXLoc]
			for i, e := range matrix[lastGameYLoc][firstGameXLoc:lastWallXLoc] {
				matrix[firstWallYLoc][i+1] = e
			}
		case lastWallYLoc:
			// copy bottom row with corners if we are on the upper wall
			matrix[lastWallYLoc][lastWallXLoc] = matrix[firstGameYLoc][firstGameXLoc]
			matrix[lastWallYLoc][firstWallXLoc] = matrix[firstGameYLoc][lastGameXLoc]
			for i, e := range matrix[firstGameYLoc][firstGameXLoc:lastWallXLoc] {
				matrix[lastWallYLoc][i+1] = e
			}
		default:
			// copy last game pos value into row,0 wall & vise versa
			matrix[index][firstWallXLoc] = row[lastGameXLoc]
			matrix[index][lastWallXLoc] = row[firstGameXLoc]
		}
	}

}

func (m *MatrixLoop) CalcNextGen(done chan bool) {
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
	UpdateLoopingWall(m.curGen)
	done <- true

}

func (m *MatrixLoop) PrintCurGen() {
	fmt.Println()
	for _, e := range m.curGen {
		fmt.Println(e)
	}
}

func (m *MatrixLoop) GetCurGen() game.GameData {
	return game.GameData{Matrix: m.curGen}
}
