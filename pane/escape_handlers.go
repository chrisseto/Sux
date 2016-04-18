package pane

import (
	"github.com/chrisseto/sux/pansi"
	"log"
)

type EscapeCodeHandler func(p *Pane)

var ESCAPE_HANDLERS = map[pansi.AnsiEscapeType]EscapeCodeHandler{}

func (p *Pane) handleEscapeCode(c *pansi.AnsiEscapeCode) {
	if handler, ok := ESCAPE_HANDLERS[c.Type]; ok {
		handler(p)
	} else {
		p.defaultEscapeCodeHandler(c)
	}
}

func (p *Pane) defaultEscapeCodeHandler(c *pansi.AnsiEscapeCode) {
	log.Printf("Go unhandled escape code %+v\n", c)
}

func (p *Pane) setGraphicMode(c *pansi.AnsiEscapeCode) {
}
