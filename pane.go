package main

import (
	"io"
	"os"
	"bufio"
	"github.com/chrisseto/pty"
	"os/exec"
)

type Pane struct {
	*exec.Cmd

  cx, cy int
  width, height uint16

  prog string
  args []string

	Pty    *os.File
	output io.Reader
	Output chan []byte
}

func CreatePane(width, height uint16, prog string, args ...string) *Pane {
	return &Pane{
    Cmd: exec.Command(prog, args...),
    cx: 0, cy: 0,
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
  if err = pty.Setsize(pterm, p.width, p.height); err != nil {
    panic(err)
  }
	p.Pty = pterm
	p.Output = make(chan []byte, 20)
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
			p.Output <- buf[0:nr]
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

	}
	close(p.Output)
}
