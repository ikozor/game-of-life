package mtxloop

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

	y := 0
	for scanner.Scan() {
		matrix = append(matrix, make([]int8, len(scanner.Text())))
		for x, e := range scanner.Text() {
			if e == alive {
				matrix[y][x] = 1
			} else {
				matrix[y][x] = 0
			}
		}
		y++
	}
	if y < 2 && len(matrix[0]) < 2 {
		return nil, fmt.Errorf("game too small")
	}

	return &MatrixImpl{curGen: matrix}, nil
}

func (m *MatrixImpl) CalcNextGen(done chan bool) {
	nextGen := make([][]int8, len(m.curGen))
	for y := 0; y < len(m.curGen); y++ {
		nextGen[y] = make([]int8, len(m.curGen[y]))
		for x := 0; x < len(m.curGen[y]); x++ {
			var aliveNearby int8
			switch y {
			case 0:
				aliveNearby += m.curGen[y+1][x]
				switch x {
				case 0:
					aliveNearby += m.curGen[y+1][len(m.curGen[y])-1]
					aliveNearby += m.curGen[y+1][x+1]
				case len(m.curGen[y]) - 1:
					aliveNearby += m.curGen[y+1][0]
					aliveNearby += m.curGen[y+1][x-1]
				default:
					aliveNearby += m.curGen[y+1][x+1]
					aliveNearby += m.curGen[y+1][x-1]
				}
			case len(m.curGen) - 1:
				aliveNearby += m.curGen[y-1][x]
				switch x {
				case 0:
					aliveNearby += m.curGen[y-1][len(m.curGen[y])-1]
					aliveNearby += m.curGen[y-1][x+1]
				case len(m.curGen[y]) - 1:
					aliveNearby += m.curGen[y-1][x-1]
					aliveNearby += m.curGen[y-1][0]
				default:
					aliveNearby += m.curGen[y-1][x+1]
					aliveNearby += m.curGen[y-1][x-1]
				}
			default:
				aliveNearby += m.curGen[y-1][x]
				aliveNearby += m.curGen[y+1][x]
				switch x {
				case 0:
					aliveNearby += m.curGen[y-1][x+1]
					aliveNearby += m.curGen[y+1][x+1]
				case len(m.curGen[y]) - 1:
					aliveNearby += m.curGen[y-1][x-1]
					aliveNearby += m.curGen[y+1][x-1]
				default:
					aliveNearby += m.curGen[y-1][x+1]
					aliveNearby += m.curGen[y+1][x+1]
					aliveNearby += m.curGen[y-1][x-1]
					aliveNearby += m.curGen[y+1][x-1]
				}
			}

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
	for _, e := range m.curGen {
		fmt.Println(e)
	}
}

func (m *MatrixImpl) GetCurGen() game.GameData {
	return game.GameData{Matrix: m.curGen}
}
