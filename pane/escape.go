package pane

import (
	"github.com/chrisseto/sux/pansi"
	"github.com/nsf/termbox-go"
	"log"
)

func (p *Pane) handleEscapeCode(c *pansi.AnsiEscapeCode) {
	switch c.Type {
	case pansi.SetGraphicMode:
		p.SetGraphicMode(c.Values)
	case pansi.CursorPosition:
		if len(c.Values) == 0 {
			p.cx, p.cy = 0, 0
		} else {
			p.cx, p.cy = c.Values[1]-1, c.Values[0]-1
		}
	case pansi.CursorUp:
		p.cy--
	case pansi.CursorDown:
		p.cy++
	case pansi.CursorBackward:
		p.cx--
	case pansi.CursorForward:
		p.cx++
	case pansi.VPA:
		if len(c.Values) == 0 {
			p.cy = 0
		} else {
			p.cy = c.Values[0] - 1
		}
	case pansi.EraseLine:
		row := p.screen.Row(p.cy)
		for i := p.cx; i < len(*row); i++ {
			(*row)[i] = termbox.Cell{' ', p.fg, p.bg}
		}
	case pansi.EraseDisplay:
		p.Clear()
	case pansi.DeleteLine:
		var val int
		if len(c.Values) == 0 {
			val = 1
		} else {
			val = c.Values[0]
		}
		p.screen.DeleteRows(p.cy, val)
		p.screen.AppendRows(val)
	case pansi.ReverseIndex:
		if p.cy > 0 {
			p.cy--
		} else {
			p.screen.Scroll(-1)
		}
	default:
		log.Printf("Doing nothing with %+v\n", *c)
	}
}

func (p *Pane) SetGraphicMode(vals []int) {
	if len(vals) == 0 {
		p.fg, p.bg = 8, 1
		return
	}
	for i := 0; i < len(vals); i++ {
		switch vals[i] {
		case 0:
			p.fg, p.bg = 8, 1
		case 1:
			p.fg |= termbox.AttrBold
		case 7:
			p.fg, p.bg = p.bg, p.fg
		case 38:
			i++
			switch vals[i] {
			case 5:
				i++
				p.fg = termbox.Attribute(vals[i] + 1)
			case 2:
				i += 3 //TODO
			}
		case 39:
			p.fg = termbox.ColorWhite
		case 48:
			i++
			switch vals[i] {
			case 5:
				i++
				p.bg = termbox.Attribute(vals[i] + 1)
			case 2:
				i += 3 //TODO
			}
		case 49:
			p.bg = termbox.ColorBlack
		}
	}
}
