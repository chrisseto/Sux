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
	SetUsg0
	DECKPAM
	DesignateG0CharacterSet
)

type AnsiEscapeCode struct {
	Type   AnsiEscapeType
	Values []int
}
