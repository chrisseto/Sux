package pane

import "github.com/nsf/termbox-go"

type Screen struct {
	cells            [][]termbox.Cell
	width, height    int
	scrollbackOffset int
}

func NewScreen(width, height int) Screen {
	s := Screen{
		cells:            make([][]termbox.Cell, height, height*10),
		width:            width,
		height:           height,
		scrollbackOffset: 0,
	}

	for i := 0; i < height; i++ {
		s.cells[i] = make([]termbox.Cell, width, width*2)
	}
	return s
}

func (s *Screen) Cells() [][]termbox.Cell {
	top, bottom := s.rowOffset()
	return s.cells[top:bottom]
}

func (s *Screen) rowOffset() (top, bottom int) {
	bottom = len(s.cells) - s.scrollbackOffset
	top = bottom - s.height
	if top < 0 {
		top = 0
	}
	return
}

func (s *Screen) Row(index int) *[]termbox.Cell {
	t, _ := s.rowOffset()
	return &s.cells[t+index]
}

func (s *Screen) Cell(index, row int) *termbox.Cell {
	return &(*s.Row(row))[index]
}

func (s *Screen) SetScrollOffset(offset int) {
	bound := len(s.cells) - s.height
	if offset > bound {
		offset = bound
	}
	s.scrollbackOffset = bound
}

func (s *Screen) AppendRows(n int) {
	toAppend := make([][]termbox.Cell, n)

	for i := 0; i < n; i++ {
		toAppend[i] = make([]termbox.Cell, s.width, s.width*2)
	}

	s.cells = append(s.cells, toAppend...)
}
