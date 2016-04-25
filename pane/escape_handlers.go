package pane

import (
	"github.com/chrisseto/sux/pansi"
	"github.com/nsf/termbox-go"
	"log"
	"runtime/debug"
)

type EscapeCodeHandler func(p *Pane, c *pansi.AnsiEscapeCode)

var ESCAPE_HANDLERS = map[pansi.AnsiEscapeType]EscapeCodeHandler{
	pansi.CursorUp:       (*Pane).CursorUp,
	pansi.CursorDown:     (*Pane).CursorDown,
	pansi.CursorForward:  (*Pane).CursorForward,
	pansi.CursorBackward: (*Pane).CursorBackward,

	pansi.EraseLine:      (*Pane).EraseLine,
	pansi.EraseDisplay:   (*Pane).EraseDisplay,
	pansi.ReverseIndex:   (*Pane).ReverseIndex,
	pansi.SetGraphicMode: (*Pane).SetGraphicMode,
	pansi.CursorPosition: (*Pane).CursorPosition,
}

func (p *Pane) handleEscapeCode(c *pansi.AnsiEscapeCode) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error while handling escape code %+v. \n%v", c, r)
			log.Printf("%s", debug.Stack())
		}
	}()

	if handler, ok := ESCAPE_HANDLERS[c.Type]; ok {
		handler(p, c)
	} else {
		p.defaultEscapeCodeHandler(c)
	}
}

func (p *Pane) defaultEscapeCodeHandler(c *pansi.AnsiEscapeCode) {
	log.Printf("Got unhandled escape code %+v\n", c)
}

func (p *Pane) CursorUp(c *pansi.AnsiEscapeCode) {
	p.cursor.Up(1)
}

func (p *Pane) CursorDown(c *pansi.AnsiEscapeCode) {
	p.cursor.Down(1)
}

func (p *Pane) CursorForward(c *pansi.AnsiEscapeCode) {
	p.cursor.Right(1)
}

func (p *Pane) CursorBackward(c *pansi.AnsiEscapeCode) {
	p.cursor.Left(1)
}

func (p *Pane) SetGraphicMode(c *pansi.AnsiEscapeCode) {
	switch {
	case len(c.Values) == 0:
		fallthrough
	case c.Values[0] == 0:
		// Reset to default colors
		p.fg, p.bg = 8, 1

	case c.Values[0] == 1:
		// Set text to bold
		p.fg |= termbox.AttrBold

	case c.Values[0] == 7:
		// Inverse fg/bg
		p.fg, p.bg = p.bg, p.fg

	case 30 <= c.Values[0] && c.Values[0] <= 37:
		// Set fg to Black/Red/Green/Yellow/Blue/Magenta/Cyan
		p.fg = termbox.Attribute(c.Values[0] - 29) //-30 + 1 for termbox offset

	case c.Values[0] == 38:
		if c.Values[1] == 5 {
			p.fg = termbox.Attribute(c.Values[2] + 1)
		} else if c.Values[1] == 2 {
			//TODO Parse RGB color code
		} else {
			panic("Malformed SGR Code")
		}

	case c.Values[0] == 39:
		p.fg = termbox.ColorWhite

	case 40 <= c.Values[0] && c.Values[0] <= 47:
		// Set bg to Black/Red/Green/Yellow/Blue/Magenta/Cyan
		p.bg = termbox.Attribute(c.Values[0] - 29) //-30 + 1 for termbox offset

	case c.Values[0] == 48:
		if c.Values[1] == 5 {
			p.bg = termbox.Attribute(c.Values[2] + 1)
		} else if c.Values[1] == 2 {
			//TODO Parse RGB color code
		} else {
			panic("Malformed SGR Code")
		}

	case c.Values[0] == 49:
		p.bg = termbox.ColorBlack
	}

}

func (p *Pane) EraseDisplay(c *pansi.AnsiEscapeCode) {
	// p.buffer.Clear()
	p.screen.Clear()
	p.screen.NewLine()
	// p.buffer.Append(make([]termbox.Cell, p.width))
}

func (p *Pane) CursorPosition(c *pansi.AnsiEscapeCode) {
	if len(c.Values) == 0 {
		p.cursor.Set(0, 0)
	} else {
		p.cursor.Set(c.Values[1]-1, c.Values[0]-1)
	}
}

func (p *Pane) EraseLine(c *pansi.AnsiEscapeCode) {
	row := p.screen.Row(p.cursor.Y())
	for i := p.cursor.X(); i < len(row); i++ {
		row[i] = termbox.Cell{' ', p.fg, p.bg}
	}
}

func (p *Pane) ReverseIndex(c *pansi.AnsiEscapeCode) {
	if p.cursor.Y() > 0 {
		p.cursor.Up(1)
	} else {
		p.screen.RollBack(1)
	}
}
