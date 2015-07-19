//pansi provides functions to aid in parsing ansi escape code
//All rules for parsing are pulled from:
//http://ascii-table.com/ansi-escape-sequences.php
//http://www.vt100.net/emu/dec_ansi_parser
package pansi

// import (
// 	"strconv"
// 	"strings"
// )

type AnsiEscapeType int

const (
	Invalid AnsiEscapeType = iota
	CursorPosition
	CursorUp
	CursorDown
	CursorForward
	CursorBackward
	SaveCursorPosition
	RestoreCursorPosition
	EraseDisplay
	EraseLine
	SetGraphicMode
	SetMode
	ResetMode
	SetKeyboardStrings
)

type State interface {
	Entry(*Parser)
	Execute(*Parser, byte) (State, *AnsiEscapeCode)
	Exit(*Parser)
}

type (
	nullState  struct{}
	clearEntry struct{ nullState }

	AnsiEscapeCode struct {
		Type   AnsiEscapeType
		Params []byte
	}

	Parser struct {
		state  State
		mode   byte
		params []byte
		result *AnsiEscapeCode
	}
)

var (
	actionNull  = func(_ *Parser) {}
	actionClear = func(p *Parser) { p.Clear() }
)

func (s nullState) Entry(_ *Parser)  {}
func (s nullState) Exit(_ *Parser)   {}
func (s clearEntry) Entry(p *Parser) { p.Clear() }

func (p *Parser) Clear() {
	p.params = nil
	p.mode = 0x0
	p.state = nil
}

func (p *Parser) Feed(r rune) {
	b := byte(r)
	switch b {
	case 0x1B:
		p.state, p.result = &escape{}, nil
	default:
		if p.state != nil {
			p.state, p.result = p.state.Execute(p, b)
		}
	}
}

func (p *Parser) Collect(b byte) {
	p.mode = b
}

func (p *Parser) Param(b byte) {
	p.params = append(p.params, b)
}

func (p *Parser) Result() *AnsiEscapeCode {
	return p.result
}

func NewParser() Parser {
	return Parser{nil, 0x0, make([]byte, 0, 16), nil}
}

func Parse(s string) *AnsiEscapeCode {
	p := NewParser()
	for _, ch := range s {
		p.Feed(ch)
		// if p.state == Invalid {
		// 	break
		// }
	}
	return p.Result()
}
