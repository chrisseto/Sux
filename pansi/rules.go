package pansi

type State int

const (
	ground State = iota
	escape
	csiEntry
	csiIgnore
	csiParam
	csiIntermediate
)

type Rule struct {
	Start, End byte
	Transition Transition
	State      State
}

var globalRules = map[byte]Rule{
	0x1B: Rule{0, 0, clear, escape},
}

var states = map[State][]Rule{
	ground: []Rule{},
	escape: []Rule{
		Rule{0x5B, 0x5B, noTransition, csiEntry},
	},
	csiEntry: []Rule{
		Rule{0x3A, 0x3A, noTransition, csiIgnore},
		Rule{0x20, 0x2F, collect, csiIntermediate},
		Rule{0x3B, 0x3B, param, csiParam},
		Rule{0x30, 0x39, param, csiParam},
		Rule{0x40, 0x7E, csiDispatch, ground},
	},
	csiParam: []Rule{
		Rule{0x20, 0x2F, noTransition, csiIntermediate},
		Rule{0x3B, 0x3B, param, csiParam},
		Rule{0x30, 0x39, param, csiParam},
		Rule{0x3A, 0x3A, noTransition, csiIgnore},
		Rule{0x40, 0x7E, csiDispatch, ground},
	},
}
