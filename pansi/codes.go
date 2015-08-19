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
	VPA        //Vertical Line Position Absolute
	VPR        //Vertical Line Position Relative
	DeleteLine //DL (default 1)
	InsertLine //IL (default 1)
	Index
	ReverseIndex
)

type AnsiEscapeCode struct {
	Type   AnsiEscapeType
	Values []int
}
