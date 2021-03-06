package pane

import (
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/chrisseto/pty"
	"github.com/chrisseto/sux/pansi"
	"github.com/nsf/termbox-go"
)

type Pane struct {
	cmd           *exec.Cmd
	index         int
	screen        Screen
	width, height int
	length        int
	pty           *os.File
	cursor        *Cursor
	//Color state
	fg termbox.Attribute
	bg termbox.Attribute

	//Launch parameters
	prog string
	args []string

	//TODO Remove this
	ShouldRedraw chan struct{}
}

func CreatePane(prog string, args []string, width, height int) *Pane {
	p := &Pane{
		cmd:          exec.Command(prog, args...),
		pty:          nil,
		prog:         prog,
		args:         args,
		ShouldRedraw: make(chan struct{}, 2), // TODO Fix this, there really shouldn't be a reason to buffer this
		screen:       NewScreen(width, height),
		height:       height,
		width:        width,
		cursor:       NewCursor(width, height-1),
	}

	return p
}

func (p *Pane) Start() error {
	var err error
	p.pty, err = pty.Start(p.cmd)
	if err != nil {
		panic(err)
	}

	if err = pty.Setsize(p.pty, uint16(p.height), uint16(p.width)); err != nil {
		panic(err)
	}

	go p.main()

	return nil
}

func (p *Pane) Stop() error {
	if err := p.cmd.Process.Kill(); err != nil {
		return err
	}

	// if err := p.pty.Close(); err != nil {
	// 	return err
	// }

	return nil
}

func (p *Pane) Kill() error {
	return nil
}

func (p *Pane) Prog() string {
	return p.prog
}

func (p *Pane) Args() []string {
	return p.args
}

func (p *Pane) Send(input []byte) (int, error) {
	return p.pty.Write(input)
}

func (p *Pane) Draw(xOffset, yOffset int) {
	for y, line := range p.screen.Cells() {
		for x, cell := range line {
			termbox.SetCell(x+xOffset, y+yOffset, cell.Ch, cell.Fg, cell.Bg)
		}
	}
	// Finally Position the cursor
	termbox.SetCursor(p.cursor.Get())
}

func (p *Pane) redraw() {
	select {
	case p.ShouldRedraw <- struct{}{}:
	default: //Failed to send, a redraw is already happening
		log.Printf("Redraw request unaccepted\n")
	}
}

func (p *Pane) main() {
	lexer := pansi.NewLexer()
	buf := make([]byte, 32*1024)

	for {
		nr, err := p.pty.Read(buf)

		if err != nil {
			if err == io.EOF {
				break // This pane's process has terminated
			}
			panic(err)
		}

		if nr < 1 {
			continue // Nothing to do, no output from the proccess
		}

		for _, char := range buf[:nr] {
			lexer.Feed(char)

			if res := lexer.Result(); res != nil {
				p.handleEscapeCode(res)
				lexer.Clear()
				continue
			}

			if lexer.State() != pansi.Ground {
				continue
			}

			p.handleByte(char)
		}

		p.redraw()
	}
}
