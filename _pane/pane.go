package pane

import (
	"github.com/chrisseto/pty"
	"github.com/nsf/termbox-go"
	// "log"
	"os"
	"os/exec"
)

type Pane struct {
	*exec.Cmd

	mode          int
	cx, cy        int
	fg, bg        termbox.Attribute
	width, height int
	scrollOffset  int
	drawOffset    int

	Prog string
	Args []string

	Pty          *os.File
	screen       Screen
	ShouldRedraw chan struct{}
}

func CreatePane(width, height int, prog string, args ...string) *Pane {
	return &Pane{
		Cmd: exec.Command(prog, args...),
		cx:  0, cy: 0,
		fg: 8, bg: 1,
		scrollOffset: 0,
		drawOffset:   0,
		Prog:         prog, Args: args,
		width: width, height: height,
		Pty:          nil,
		screen:       NewScreen(width, height),
		ShouldRedraw: make(chan struct{}),
	}
}

func (p *Pane) Start() error {
	pterm, err := pty.Start(p.Cmd)
	if err != nil {
		panic(err)
	}
	if err = pty.Setsize(pterm, uint16(p.height), uint16(p.width)); err != nil {
		panic(err)
	}
	p.Pty = pterm
	go p.mainLoop()
	return nil
}

func (p *Pane) Close() error {
	//TODO end process nicely
	return p.Process.Kill()
}

func (p *Pane) Cells() [][]termbox.Cell {
	return p.screen.Cells()
}

func (p *Pane) Width() int {
	return p.width
}

func (p *Pane) Height() int {
	return p.height
}

func (p *Pane) redraw() {
	select {
	case p.ShouldRedraw <- struct{}{}:
	default: //Failed to send, a redraw is already happening
	}
}

func (p *Pane) Scroll(far int) {
	p.screen.SetScrollOffset(far)
	p.redraw()
}

func (p *Pane) NewLine() {
	p.cy++
	if p.cy > p.height-1 {
		p.cy = p.height - 1
		p.screen.AppendRows(1)
	}
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
	p.cx = bound(p.cx, 0, p.width-1)
	p.cy = bound(p.cy, 0, p.height-1)
	return p.cx, p.cy
}
