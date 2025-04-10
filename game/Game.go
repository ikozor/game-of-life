package game

type Game interface {
	CalcNextGen(chan bool)
	PrintCurGen()
	GetCurGen() GameData
}

type GameData struct {
	Matrix [][]int8
}
