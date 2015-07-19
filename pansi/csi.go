package pansi

import (
	"strconv"
	"strings"
)

type (
	csiEntry        struct{ clearEntry }
	csiIgnore       struct{ nullState }
	csiParam        struct{ nullState }
	csiIntermediate struct{ nullState }
	csiDispatch     struct{ nullState }
)

func (p *Parser) csiDispatch(b byte) *AnsiEscapeCode {
	var err error
	var t AnsiEscapeType
	var values []int
	if len(p.params) > 1 {
		spl := strings.Split(string(p.params), ";")
		values = make([]int, len(spl))
		for i, x := range spl {
			values[i], err = strconv.Atoi(x)
			if err != nil {
				return nil
			}
		}
	} else {
		values = []int{}
	}
	switch b {
	case 0x6D:
		t = SetGraphicMode
	case 0x66:
		fallthrough
	case 0x48:
		t = CursorPosition
	case 0x41:
		t = CursorUp
	case 0x42:
		t = CursorDown
	case 0x43:
		t = CursorForward
	case 0x44:
		t = CursorBackward
	case 0x4B:
		t = EraseLine
	default:
		return nil
	}
	return &AnsiEscapeCode{t, values}
}

func (s *csiEntry) Execute(p *Parser, b byte) (State, *AnsiEscapeCode) {
	switch {
	case b == 0x3A:
		return &csiIgnore{}, nil
	case b >= 0x20 && b <= 0x2F:
		p.Collect(b)
		return &csiIntermediate{}, nil
	case b == 0x3B:
		fallthrough
	case b >= 0x30 && b <= 0x39:
		p.Param(b)
		return &csiParam{}, nil
	case b >= 0x3C && b <= 0x3F:
		p.Collect(b)
		return &csiParam{}, nil
	case b >= 0x40 && b <= 0x7E:
		return nil, p.csiDispatch(b)
	default:
		return nil, nil
	}
}

func (s *csiParam) Execute(p *Parser, b byte) (State, *AnsiEscapeCode) {
	switch {
	case b >= 0x20 && b <= 0x2F:
		return &csiIntermediate{}, nil
	case b == 0x3B:
		fallthrough
	case b >= 0x30 && b <= 0x39:
		p.Param(b)
		return &csiParam{}, nil
	case b == 0x3A:
		fallthrough
	case b >= 0x3C && b <= 0x3F:
		return &csiIgnore{}, nil
	case b >= 0x40 && b <= 0x7E:
		return nil, p.csiDispatch(b)
	default:
		return nil, nil
	}
}

func (s *csiIntermediate) Execute(p *Parser, b byte) (State, *AnsiEscapeCode) {
	switch {
	case b >= 0x20 && b <= 0x2F:
		p.Collect(b)
		return s, nil
	case b >= 0x30 && b <= 0x3F:
		return &csiIgnore{}, nil
	case b >= 0x40 && b <= 0x7E:
		return nil, nil
	default:
		return nil, nil
	}
}

func (s *csiIgnore) Execute(p *Parser, b byte) (State, *AnsiEscapeCode) {
	switch {
	case b == 0x7F:
		fallthrough
	case b >= 0x20 && b <= 0x3F:
		return s, nil
	case b >= 0x40 && b <= 0x7E:
		return nil, nil
	default:
		return nil, nil
	}
}
