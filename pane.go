package main

import (
	"github.com/chrisseto/pty"
	"github.com/chrisseto/sux/pansi"
	"github.com/nsf/termbox-go"
	"io"
	"os"
	"os/exec"
)

type Cell struct {
	termbox.Cell
	x, y int
}

type Pane struct {
	*exec.Cmd

	cx, cy        int
	sx, sy        int
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
}

func (p *Pane) outputPipe() {
	fg, bg := 0x0, 0x0
	parser := pansi.NewParser()
	buf := make([]byte, 32*1024)

	for {
		nr, err := p.Pty.Read(buf)
		if nr > 0 {
			row := &p.cells[p.sy]

			for _, char := range buf[:nr] {
				parser.Feed(char)
				if parser.State() != nil {
					continue
				}
				if res := parser.Result(); res != nil {
					if len(res.Values) == 1 {
						fg, bg = 0x0, 0x0
					} else if res.Values[0] == 38 && res.Values[1] == 5 {
						fg = res.Values[2] + 1
					}
					parser.Clear()
					continue
				}

				switch char {
				case 0xA:
					p.sy++
					p.cells = append(p.cells, make([]termbox.Cell, p.width))
					row = &p.cells[p.sy]
				case 0xD:
					p.sx = 0
				case 0x8:
					p.sx--
					(*row)[p.sx] = termbox.Cell{' ', termbox.Attribute(fg), termbox.Attribute(bg)}
				default:
					p.sx++
					(*row)[p.sx] = termbox.Cell{rune(char), termbox.Attribute(fg), termbox.Attribute(bg)}
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
