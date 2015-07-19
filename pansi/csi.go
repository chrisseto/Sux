package pansi

// var (
// csiEntry = &State{
// 	Entry:   actionClear,
// 	Execute: csiEntryExecute,
// 	Exit:    actionNull,
// }

// csiIntermediate = &State{
// 	Entry:   actionNull,
// 	Execute: csiIntermediateExecute,
// 	Exit:    actionNull,
// }

// csiIgnore = &State{
// 	Entry:   actionNull,
// 	Execute: csiIgnoreExecute,
// 	Exit:    actionNull,
// }

// csiParam = &State{
// 	Entry:   actionNull,
// 	Execute: csiParamExecute,
// 	Exit:    actionNull,
// }
// )

type (
	csiEntry        struct{ clearEntry }
	csiIgnore       struct{ nullState }
	csiParam        struct{ nullState }
	csiIntermediate struct{ nullState }
)

func (p *Parser) csiDispatch(b byte) *AnsiEscapeCode {
	switch b {
	case 0x6D:
		return &AnsiEscapeCode{SetGraphicMode, p.params}
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
		return nil, nil
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

// func csiIgnoreExecute(p *Parser, b byte) (*State, *AnsiEscapeCode) {
// 	switch {
// 	case b == 0x7F:
// 		fallthrough
// 	case b >= 0x20 && b <= 0x3F:
// 		return csiIgnore, nil
// 	case b >= 0x40 && b <= 0x7E:
// 		return nil, nil
// 	default:
// 		return nil, nil
// 	}
// }
