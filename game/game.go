package game

type Game interface {
	GetNextGen()
	GetCurGen() [][]bool
	
}


