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
	DeleteCharacter //DCH (default 1)
	// CHA—Cursor Horizontal Absolute
	// Move the active position to the n-th character of the active line.
	// Format: CSI Pn G
	// Parameters
	// 	Pn (default 1): is the number of active positions to the n-th character of the active line.
	// Description
	// 	The active position is moved to the n-th character position of the active line.
	CursorHorizontalAbsolute
)

type AnsiEscapeCode struct {
	Type   AnsiEscapeType
	Values []int
}
