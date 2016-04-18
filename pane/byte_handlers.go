package pane

import (
	"github.com/nsf/termbox-go"
)

const (
	NEW_LINE        = 0xA
	TERMINAL_BELL   = 0x7
	BACKSPACE       = 0x8
	CARRIAGE_RETURN = 0xD
)

type ByteHandler func(p *Pane)

var BYTE_HANDLERS = map[byte]ByteHandler{
	NEW_LINE:        (*Pane).NewLine,
	BACKSPACE:       (*Pane).Backspace,
	TERMINAL_BELL:   (*Pane).TerminalBell,
	CARRIAGE_RETURN: (*Pane).CarriageReturn,
}

func (p *Pane) handleByte(b byte) {
	if handler, ok := BYTE_HANDLERS[b]; ok {
		handler(p)
	} else {
		p.defaultByteHandler(b)
	}
}

func (p *Pane) defaultByteHandler(b byte) {
	*p.Cell(p.cursor.Get()) = termbox.Cell{rune(b), p.fg, p.bg}
	p.cursor.Right(1)
}

func (p *Pane) NewLine() {
	p.buffer.Append(make([]termbox.Cell, p.width))
	p.cursor.Down(1)
}

func (p *Pane) CarriageReturn() {
	p.cursor.SetX(0)
}

func (p *Pane) Backspace() {
	p.cursor.Left(1)
	//Should this always be 8, 1?
	*p.Cell(p.cursor.Get()) = termbox.Cell{' ', p.fg, p.bg}
}

//This function intentionally left blank
func (p *Pane) TerminalBell() {
}
