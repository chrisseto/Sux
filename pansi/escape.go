package pansi

type escape struct {
	clearEntry
}

// var escape = State{
// 	Entry:   actionClear,
// 	Execute: csiEntryExecute,
// 	Exit:    actionNull,
// }

func (s escape) Execute(p *Parser, b byte) (State, *AnsiEscapeCode) {
	switch b {
	case 0x5B:
		return &csiEntry{}, nil
	default:
		return nil, nil
	}
}
