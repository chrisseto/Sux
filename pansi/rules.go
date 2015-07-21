package pansi

type State int

const (
	Ground State = iota
	Escape
	CsiEntry
	CsiIgnore
	CsiParam
	CsiIntermediate
)

type Rule struct {
	Start, End byte
	Transition Transition
	State      State
}

var globalRules = map[byte]Rule{
	0x1B: Rule{0, 0, clear, Escape},
}

var states = map[State][]Rule{
	Ground: []Rule{},
	Escape: []Rule{
		Rule{0x5B, 0x5B, noTransition, CsiEntry},
	},
	CsiEntry: []Rule{
		Rule{0x3A, 0x3A, noTransition, CsiIgnore},
		Rule{0x20, 0x2F, collect, CsiIntermediate},
		Rule{0x3B, 0x3B, param, CsiParam},
		Rule{0x30, 0x39, param, CsiParam},
		Rule{0x40, 0x7E, csiDispatch, Ground},
	},
	CsiParam: []Rule{
		Rule{0x20, 0x2F, noTransition, CsiIntermediate},
		Rule{0x3B, 0x3B, param, CsiParam},
		Rule{0x30, 0x39, param, CsiParam},
		Rule{0x3A, 0x3A, noTransition, CsiIgnore},
		Rule{0x40, 0x7E, csiDispatch, Ground},
	},
}
