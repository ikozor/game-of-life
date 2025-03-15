package anim

import (
	"github.com/gdamore/tcell/v2"
)

type screen struct {
	screen      tcell.Screen
	defStyle    tcell.Style
	Dead, Alive rune
}

func CreateScreen(dead, alive rune) (*screen, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := s.Init(); err != nil {
		return nil, err
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	s.Clear()
	return &screen{
		screen:   s,
		defStyle: defStyle,
		Dead:     dead,
		Alive:    alive,
	}, nil
}

func (s *screen) UpdateWithMatrix(matrix [][]int8) {
	for y := 1; y < len(matrix)-1; y++ {
		for x := 1; x < len(matrix[y])-1; x++ {
			if matrix[y][x] == 1 {
				s.screen.SetContent(x-1, y-1, s.Alive, nil, s.defStyle)
			} else {
				s.screen.SetContent(x-1, y-1, s.Dead, nil, s.defStyle)
			}
		}

	}
	s.screen.PostEvent(nil)
	s.screen.Show()
}

func (s *screen) CaptureEscape() bool {
	ev := s.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
			return false
		}
	}
	return true
}

func (s *screen) Finished() {
	s.screen.Fini()
}
