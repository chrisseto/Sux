package pane

import (
	"github.com/nsf/termbox-go"
)

// TODO Maybe Embedding is a better idea?
type Screen struct {
	RingBuffer
	// buffer        RingBuffer
	width, height int
}

func NewScreen(width, height int) Screen {
	s := Screen{
		RingBuffer: NewRingBuffer(make([][]termbox.Cell, 0, height*2)),
		// buffer: NewRingBuffer(make([][]termbox.Cell, 0, height*2)),
		// buffer: make([][]termbox.Cell, 0, height*2),
		width:  width,
		height: height,
	}

	s.Append(make([]termbox.Cell, s.width))

	return s
}

func (s *Screen) offset() int {
	if s.Length() < s.height {
		return 0
	} else {
		return s.Length() - s.height
	}
}

func (s *Screen) Row(index int) []termbox.Cell {
	return s.Get(index + s.offset())
}

func (s *Screen) Cell(x, y int) *termbox.Cell {
	return &s.Get(y + s.offset())[x]
}

func (s *Screen) Cells() [][]termbox.Cell {
	ret := s.Tail(s.height)

	lines := s.height - len(ret) // Don't change lists while iterating over them
	for i := 0; i < lines; i++ {
		ret = append(ret, make([]termbox.Cell, s.width))
	}

	return ret
}

func (s *Screen) NewLine() {
	s.Append(make([]termbox.Cell, s.width))
}
