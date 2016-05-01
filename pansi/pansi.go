//go:generate python generate.py
//go:generate gofmt -w codes.go
//go:generate stringer -type=AnsiEscapeType

//pansi provides functions to aid in parsing ansi escape code
//All rules for parsing are pulled from:
//http://ascii-table.com/ansi-escape-sequences.php
//http://www.vt100.net/emu/dec_ansi_parser
package pansi

type AnsiEscapeType int
type AnsiEscapeCode struct {
	Type   AnsiEscapeType
	Values []int
}

var globalLexer = NewLexer()

func Result() *AnsiEscapeCode {
	return globalLexer.Result()
}

// func State() State {
// 	return l.State()
// }

func Clear() {
	globalLexer.Clear()
}

func Parse(s string) *AnsiEscapeCode {
	return globalLexer.Parse(s)
}

func FeedRune(r rune) {
	globalLexer.FeedRune(r)
}

func Feed(b byte) {
	globalLexer.Feed(b)
}
