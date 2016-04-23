package pane

import (
	"github.com/chrisseto/sux/pansi"
	"github.com/nsf/termbox-go"
	"log"
)

type EscapeCodeHandler func(p *Pane, c *pansi.AnsiEscapeCode)

var ESCAPE_HANDLERS = map[pansi.AnsiEscapeType]EscapeCodeHandler{
	pansi.EraseDisplay:   (*Pane).Clear,
	pansi.CursorPosition: (*Pane).CursorPosition,
}

func (p *Pane) handleEscapeCode(c *pansi.AnsiEscapeCode) {
	if handler, ok := ESCAPE_HANDLERS[c.Type]; ok {
		handler(p, c)
	} else {
		p.defaultEscapeCodeHandler(c)
	}
}

func (p *Pane) defaultEscapeCodeHandler(c *pansi.AnsiEscapeCode) {
	log.Printf("Go unhandled escape code %+v\n", c)
}

func (p *Pane) SetGraphicMode(c *pansi.AnsiEscapeCode) {
}

func (p *Pane) Clear(c *pansi.AnsiEscapeCode) {
	p.buffer.Clear()
	p.buffer.Append(make([]termbox.Cell, p.width))
}

func (p *Pane) CursorPosition(c *pansi.AnsiEscapeCode) {
	if len(c.Values) == 0 {
		p.cursor.Set(0, 0)
	} else {
		p.cursor.Set(c.Values[1]-1, c.Values[0]-1)
	}
}
