package pansi

import (
	"strconv"
	"strings"
)

type Transition func(*Lexer, byte) *AnsiEscapeCode

var csiDispatchMap = map[byte]AnsiEscapeType{
	0x6D: SetGraphicMode,
	0x66: CursorPosition,
	0x48: CursorPosition,
	0x41: CursorUp,
	0x42: CursorDown,
	0x43: CursorForward,
	0x44: CursorBackward,
	0x4A: EraseDisplay,
	0x4B: EraseLine,
	0x68: DecPrivateModeSet,
	0x72: SetBottomTopLines,
}

var escDispatchMap = map[byte]AnsiEscapeType{
	0x3D: DECKPAM,
}

func csiDispatch(l *Lexer, b byte) *AnsiEscapeCode {
	t, ok := csiDispatchMap[b]
	if !ok {
		return nil
	}

	var err error
	var values []int
	if len(l.params) < 1 {
		values = []int{}
	} else {
		spl := strings.Split(string(l.params), ";")
		values = make([]int, len(spl))
		for i, x := range spl {
			values[i], err = strconv.Atoi(x)
			if err != nil {
				return nil
			}
		}
	}

	return &AnsiEscapeCode{t, values}
}

func escDispatch(l *Lexer, b byte) *AnsiEscapeCode {
	t, ok := escDispatchMap[b]
	if !ok {
		return nil
	}
	return &AnsiEscapeCode{t, nil}
}

func noTransition(l *Lexer, b byte) *AnsiEscapeCode {
	return nil
}

func collect(l *Lexer, b byte) *AnsiEscapeCode {
	l.mode = b
	return nil
}

func param(l *Lexer, b byte) *AnsiEscapeCode {
	l.params = append(l.params, b)
	return nil
}

func clear(l *Lexer, b byte) *AnsiEscapeCode {
	l.params = nil
	l.mode = 0
	return nil
}
