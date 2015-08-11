package pane

import (
	"github.com/chrisseto/sux/pansi"
	"github.com/nsf/termbox-go"
	"io"
)

const (
	NEW_LINE        = 0xA
	TERMINAL_BELL   = 0x7
	BACKSPACE       = 0x8
	CARRIAGE_RETURN = 0xD
)

func (p *Pane) HandleKey(key byte) {
	if p.cx >= int(p.width) {
		p.row = p.NewLine()
	}
	// if p.mode & MODE_WRAP {
	// 	p.row = p.newLine()
	// } else {
	// p.row = append(p.row, make([]termbox.Cell, p.width))
	// }
	// }
	(*p.row)[p.cx] = termbox.Cell{rune(key), p.fg, p.bg}
	p.cx++
}

func (p *Pane) BackSpace() {
	if p.cx != 0 {
		p.cx--
	}
	(*p.row)[p.cx] = termbox.Cell{' ', p.fg, p.bg} //Should this always be 8, 1?
}

func (p *Pane) Clear() {
	p.drawOffset = len(p.cells) - 1
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault) //Should not be here
	p.redraw()
}

func (p *Pane) mainLoop() {
	lexer := pansi.NewLexer()
	buf := make([]byte, 32*1024)

	for {
		nr, err := p.Pty.Read(buf)
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
					p.row = p.NewLine()
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
