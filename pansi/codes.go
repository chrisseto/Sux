package pansi

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
	DecPrivateModeSet
	SetBottomTopLines
)

type AnsiEscapeCode struct {
	Type   AnsiEscapeType
	Values []int
}
