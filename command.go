package main

import (
	"os"
	// "io"
	// "bytes"
	"github.com/kr/pty"
	"os/exec"
	// "strings"
	"syscall"
)

type Cmd struct {
	*exec.Cmd
	Pty    *os.File
	Output chan []byte
}

func Command(prog string, args ...string) *Cmd {
	return &Cmd{exec.Command(prog, args...), nil, nil}
}

func (c *Cmd) Start() error {
	pterm, err := pty.Start(c.Cmd)
	if err != nil {
		return err
	}
	c.Pty = pterm
	go c.outputPipe()
	return nil
}

func (c *Cmd) Close() error {
	return c.Process.Kill()
}

func (c *Cmd) outputPipe() {
	fd := (int)(c.Pty.Fd())
	buf := make([]byte, 128)
	c.Output = make(chan []byte)
	syscall.SetNonblock(fd, true)
	defer syscall.SetNonblock(fd, false)

	for {
		n, err := syscall.Read(fd, buf)
		if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
			continue
		}
		c.Output <- buf[:n]
	}
	close(c.Output)
}
