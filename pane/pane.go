package pane

import (
	"github.com/chrisseto/pty"
	"github.com/nsf/termbox-go"
	"os"
	"os/exec"
)

type Pane struct {
	*exec.Cmd

	mode          int
	cx, cy        int
	row           *[]termbox.Cell
	fg, bg        termbox.Attribute
	width, height uint16
	scrollOffset  int
	drawOffset    int

	Prog string
	Args []string

	Pty          *os.File
	cells        [][]termbox.Cell
	ShouldRedraw chan struct{}
}

func CreatePane(width, height uint16, prog string, args ...string) *Pane {
	return &Pane{
		Cmd: exec.Command(prog, args...),
		cx:  0, cy: 0,
		fg: 8, bg: 1,
		scrollOffset: 0,
		drawOffset:   0,
		Prog:         prog, Args: args,
		width: width, height: height,
		Pty:          nil,
		row:          nil,
		ShouldRedraw: make(chan struct{}),
	}
}

func (p *Pane) Start() error {
	pterm, err := pty.Start(p.Cmd)
	if err != nil {
		panic(err)
	}
	if err = pty.Setsize(pterm, p.height, p.width); err != nil {
		panic(err)
	}
	p.Pty = pterm
	p.cells = make([][]termbox.Cell, 1, p.height)
	p.cells[0] = make([]termbox.Cell, p.width)
	p.row = p.bottomLine()
	go p.mainLoop()
	return nil
}

func (p *Pane) Close() error {
	//TODO end process nicely
	return p.Process.Kill()
}

func (p *Pane) Cells() [][]termbox.Cell {
	return p.cells[p.drawOffset:bound(p.drawOffset+int(p.height), p.drawOffset, len(p.cells))]
}

func (p *Pane) Width() uint16 {
	return p.width
}

func (p *Pane) Height() uint16 {
	return p.height
}

func (p *Pane) redraw() {
	select {
	case p.ShouldRedraw <- struct{}{}:
	default: //Failed to send, a redraw is already happening
	}
}

func (p *Pane) Scroll(far int) {
	p.scrollOffset = bound(p.scrollOffset+far, -len(p.cells), 0)
	p.redraw()
}

func (p *Pane) bottomLine() *[]termbox.Cell {
	return &p.cells[len(p.cells)-1]
}

func (p *Pane) NewLine() *[]termbox.Cell {
	p.cy++
	p.cells = append(p.cells, make([]termbox.Cell, p.width))
	if len(p.cells)-p.drawOffset > int(p.height) {
		p.drawOffset++
	}
	return p.bottomLine()
}

func (p *Pane) Redraw() {
	for y, line := range p.Cells() {
		for x, cell := range line {
			termbox.SetCell(x, y, cell.Ch, cell.Fg, cell.Bg)
		}
	}
	termbox.SetCursor(p.Cursor())
}

func bound(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func (p *Pane) Cursor() (int, int) {
	p.cx = bound(p.cx, 0, int(p.width)-1)
	p.cy = bound(p.cy, 0, int(p.height)-1)
	return p.cx, p.cy
}
