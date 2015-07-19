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
	switch b {
	case 0x6D:
		spl := strings.Split(string(p.params), ";")
		// var ok bool
		var err error
		values := make([]int, len(spl))
		for i, x := range spl {
			values[i], err = strconv.Atoi(x)
			if err != nil {
				return nil
			}
		}
		return &AnsiEscapeCode{SetGraphicMode, values}
	default:
		return nil
	}
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
