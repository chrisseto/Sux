package main

import (
	"github.com/chrisseto/pty"
	"github.com/chrisseto/sux/pansi"
	"github.com/nsf/termbox-go"
	"io"
	"os"
	"os/exec"
)

type Pane struct {
	*exec.Cmd

	cx, cy        int
	sx, sy        int
	fg, bg        termbox.Attribute
	width, height uint16
	scrollOffset  int

	Prog string
	Args []string

	Pty    *os.File
	output io.Reader
	cells  [][]termbox.Cell
}

func CreatePane(width, height uint16, prog string, args ...string) *Pane {
	return &Pane{
		Cmd: exec.Command(prog, args...),
		cx:  0, cy: 0,
		sx: 0, sy: 0,
		fg: 0, bg: 0,
		scrollOffset: 0,
		Prog:         prog, Args: args,
		width: width, height: height,
		Pty: nil,
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
	go p.outputPipe()
	return nil
}

func (p *Pane) Close() error {
	return p.Process.Kill()
}

func (p *Pane) Cells() [][]termbox.Cell {
	if offset := len(p.cells) + p.scrollOffset - int(p.height); offset > 0 {
		return p.cells[offset:]
	}
	return p.cells
}

func (p *Pane) Width() uint16 {
	return p.width
}

func (p *Pane) Height() uint16 {
	return p.height
}

func (p *Pane) Scroll(far int) {
	p.scrollOffset += far
	select {
	case Redraw <- struct{}{}:
	default: //Failed to send, a redraw is already happening
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

func (p *Pane) Cursor() (int, int) {
	return p.cx, p.cy
}

func (p *Pane) outputPipe() {
	lexer := pansi.NewLexer()
	buf := make([]byte, 32*1024)
	// f, _ := os.Create("output.log")
	for {
		nr, err := p.Pty.Read(buf)
		if nr > 0 {
			row := &p.cells[p.sy]
			// f.Write(buf[:nr])

			for _, char := range buf[:nr] {
				lexer.Feed(char)
				if lexer.State() != pansi.Ground {
					continue
				}
				if res := lexer.Result(); res != nil {
					p.handleEscapeCode(res)
					lexer.Clear()
					continue
				}

				switch char {
				case 0x7: //Terminal Bell. Skip for the moment
				case 0xA:
					p.sy++
					p.cy++
					p.cells = append(p.cells, make([]termbox.Cell, p.width))
					row = &p.cells[p.sy]
				case 0xD:
					p.sx = 0
					p.cx = 0
				case 0x8:
					p.sx--
					p.cx--
					(*row)[p.sx] = termbox.Cell{' ', p.fg, p.bg}
				default:
					(*row)[p.sx] = termbox.Cell{rune(char), p.fg, p.bg}
					p.sx++
					p.cx++
				}
			}

			Redraw <- struct{}{}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

	}
}

func (p *Pane) handleEscapeCode(c *pansi.AnsiEscapeCode) {
	switch c.Type {
	case pansi.SetGraphicMode:
		p.SetGraphicMode(c.Values)
	case pansi.CursorPosition:
		p.cx, p.cy = c.Values[1], c.Values[2]
	case pansi.CursorUp:
		p.cy--
	case pansi.CursorDown:
		p.cy++
	case pansi.CursorBackward:
		p.cx--
	case pansi.CursorForward:
		p.cx++
	case pansi.EraseLine:
		row := &p.cells[p.sy]
		for i := p.cx; i < len(*row); i++ {
			(*row)[i] = termbox.Cell{' ', p.fg, p.bg}
		}
	}
}

func (p *Pane) SetGraphicMode(vals []int) {
	for i := 0; i < len(vals); i++ {
		switch vals[i] {
		case 0:
			p.fg, p.bg = 0, 0
		case 1:
			p.fg |= termbox.AttrBold
		case 38:
			i++
			switch vals[i] {
			case 5:
				i++
				p.fg = termbox.Attribute(vals[i] + 1)
			case 2:
				i += 3 //TODO
			}
		case 48:
			i++
			switch vals[i] {
			case 5:
				i++
				p.bg = termbox.Attribute(vals[i] + 1)
			case 2:
				i += 3 //TODO
			}
		}
	}
}
