package pansi

type Lexer struct {
	state  State
	mode   byte
	params []byte
	result *AnsiEscapeCode
}

func NewLexer() Lexer {
	return Lexer{state: Ground}
}

func (l *Lexer) Result() *AnsiEscapeCode {
	return l.result
}

func (l *Lexer) State() State {
	return l.state
}

func (l *Lexer) Clear() {
	l.state = Ground
	l.mode = 0
	l.params = make([]byte, 0, 15)
	l.result = nil
}

func (l *Lexer) Parse(s string) *AnsiEscapeCode {
	for _, r := range s {
		l.Feed(byte(r))
	}
	return l.Result()
}

func (l *Lexer) FeedRune(r rune) {
	l.Feed(byte(r))
}

func (l *Lexer) Feed(b byte) {
	rule, ok := globalRules[b]
	if ok {
		l.state, l.result = rule.State, rule.Transition(l, b)
		return
	}

	for _, rule := range states[l.state] {
		if b >= rule.Start && b <= rule.End {
			l.state, l.result = rule.State, rule.Transition(l, b)
			return
		}
	}
}
