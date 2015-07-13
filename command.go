package main

import (
	"bufio"
	"github.com/kr/pty"
	"io"
	"os"
	"os/exec"
)

type Cmd struct {
	*exec.Cmd
	Pty    *os.File
	output io.Reader
	Output chan []byte
}

func Command(prog string, args ...string) *Cmd {
	return &Cmd{exec.Command(prog, args...), nil, nil, nil}
}

func (c *Cmd) Start() error {
	pterm, err := pty.Start(c.Cmd)
	if err != nil {
		return err
	}
	c.Pty = pterm
	c.Output = make(chan []byte, 20)
	c.output = bufio.NewReader(c.Pty)
	go c.outputPipe()
	return nil
}

func (c *Cmd) Close() error {
	return c.Process.Kill()
}

func (c *Cmd) outputPipe() {
	buf := make([]byte, 32*1024)
	for {
		nr, err := c.output.Read(buf)
		if nr > 0 {
			c.Output <- buf[0:nr]
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

	}
	close(c.Output)
}
