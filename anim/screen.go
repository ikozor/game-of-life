package anim

import (
	"github.com/gdamore/tcell/v2"
)

type screen struct {
	width, height int
	screen        tcell.Screen
	defStyle      tcell.Style
}

func CreateScreen(width, height int) (*screen, error) {
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
	return &screen{screen: s, width: width, height: height, defStyle: defStyle}, nil
}

func (s *screen) UpdateWithMatrix(matrix [][]rune) {
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			s.screen.SetContent(x, y, matrix[y][x], nil, s.defStyle)
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
