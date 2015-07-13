package main

import (
	"bufio"
	"github.com/chrisseto/pty"
	"io"
	"os"
	// "github.com/nsf/termbox-go"
	"os/exec"
)

type Pane struct {
	*exec.Cmd

	cx, cy        int
	width, height uint16

	prog string
	args []string

	Pty    *os.File
	output io.Reader
	// cells []termbox.Cell
	Output chan []byte
}

func CreatePane(width, height uint16, prog string, args ...string) *Pane {
	return &Pane{
		Cmd: exec.Command(prog, args...),
		cx:  0, cy: 0,
		prog: prog, args: args,
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
	p.output = bufio.NewReader(p.Pty)
	go p.outputPipe()
	return nil
}

func (p *Pane) Close() error {
	return p.Process.Kill()
}

func (p *Pane) outputPipe() {
	buf := make([]byte, 32*1024)
	for {
		nr, err := p.output.Read(buf)
		if nr > 0 {
			c := make([]byte, nr)
			copy(c, buf[:nr])
			p.Output <- c
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
