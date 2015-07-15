package main

import (
	"bufio"
	"github.com/chrisseto/pty"
	"github.com/nsf/termbox-go"
	"io"
	"os"
	"os/exec"
  "fmt"
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

	Prog string
	args []string

	Pty        *os.File
	output     io.Reader
	cells      [][]termbox.Cell
	Output     chan []byte
	CellOutput chan []Cell
}

func CreatePane(width, height uint16, prog string, args ...string) *Pane {
	return &Pane{
		Cmd: exec.Command(prog, args...),
		cx:  0, cy: 0,
		sx:  0, sy: 0,
		Prog: prog, args: args,
		width: width, height: height,
		Pty: nil, output: nil, Output: nil,
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
	p.Output = make(chan []byte, 32)
	p.CellOutput = make(chan []Cell, 32)
	p.output = bufio.NewReader(p.Pty)
	p.cells = make([][]termbox.Cell, 1, p.height)
	p.cells[0] = make([]termbox.Cell, 0, p.width)
	go p.outputPipe()
	return nil
}

func (p *Pane) Close() error {
	return p.Process.Kill()
}

func (p *Pane) Cells() [][]termbox.Cell {
	return p.cells
}

func (p *Pane) Width() uint16 {
	return p.width
}

func (p *Pane) Height() uint16 {
	return p.height
}

func (p *Pane) outputPipe() {
	buf := make([]byte, 32*1024)
  f, _ := os.Create("outpipe.log")
  defer f.Close()
	for {
		nr, err := p.output.Read(buf)
		if nr > 0 {
			b := make([]Cell, 0, nr)
			row := &p.cells[p.sy]

			for _, char := range buf[:nr] {
				switch char {
        case 0xA:
          p.sy++
					p.cells = append(p.cells, nil)
					*row = p.cells[p.sy]
				case 0xD:
					p.sx = 0
				case 0x8:
          p.sx--
					*row = (*row)[:p.sx]
          c := Cell{termbox.Cell{' ', 0x0, 0x0}, p.sx, p.sy}
          b = append(b, c)
				default:
          c := Cell{termbox.Cell{rune(char), 0x0, 0x0}, p.sx, p.sy}
					*row = append(*row, c.Cell)
          b = append(b, c)
					p.sx++
				}
			}
      fmt.Fprintf(f, "nr: %v\nbuf: %v\n%+v\n", nr, buf[:nr], b)

			p.CellOutput <- b
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

	}
	close(p.Output)
	p.Output = nil
}
