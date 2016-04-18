package pane

import (
	"github.com/chrisseto/sux/pansi"
	"github.com/nsf/termbox-go"
	"io"
	"log"
	"os"
)

const (
	NEW_LINE        = 0xA
	TERMINAL_BELL   = 0x7
	BACKSPACE       = 0x8
	CARRIAGE_RETURN = 0xD
)

func (p *Pane) HandleKey(key byte) {
	if p.cx >= p.width {
		p.NewLine()
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic! %+v. Cursor at (%d, %d)", r, p.cx, p.cy)
			panic(r)
		}
	}()
	// if p.mode & MODE_WRAP {
	// 	p.row = p.newLine()
	// } else {
	// p.row = append(p.row, make([]termbox.Cell, p.width))
	// }
	// }
	*p.screen.Cell(p.cx, p.cy) = termbox.Cell{rune(key), p.fg, p.bg}
	p.cx++
}

func (p *Pane) BackSpace() {
	if p.cx != 0 {
		p.cx--
	}
	*p.screen.Cell(p.cx, p.cy) = termbox.Cell{' ', p.fg, p.bg} //Should this always be 8, 1?
}

func (p *Pane) Clear() {
	p.screen.AppendRows(p.height)
	p.redraw()
}

func (p *Pane) mainLoop() {
	lexer := pansi.NewLexer()
	buf := make([]byte, 32*1024)
	f, _ := os.Create("pane.raw")
	logfile, _ := os.Create("pane.log")
	log.SetOutput(logfile)
	log.Printf("Pane width, height: %d, %d\n", p.width, p.height)
	for {
		nr, err := p.Pty.Read(buf)
		f.Write(buf[:nr])
		if nr > 0 {
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

				switch char {
				case TERMINAL_BELL: //Skip for the moment
				case NEW_LINE:
					p.NewLine()
				case CARRIAGE_RETURN:
					p.cx = 0
				case BACKSPACE:
					p.BackSpace()
				default:
					p.HandleKey(char)
				}
			}

			p.redraw()
		}
		if err != nil {
			if err == io.EOF {
				break // This pane's process has terminated
			}
			panic(err)
		}
	}
}
