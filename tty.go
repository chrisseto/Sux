package main

const (
	MODE_CURSOR = 1 << iota
	MODE_INSERT
	MODE_KCURSOR
	MODE_KKEYPAD
	MODE_WRAP
	MODE_MOUSE_STANDARD
	MODE_MOUSE_BUTTON
	MODE_BLINKING
	MODE_MOUSE_UTF8
	MODE_MOUSE_SGR
	MODE_BRACKETPASTE
	MODE_FOCUSON
)

type Tty struct {
	Pty   *os.File
	modes int
}

func (t *Tty) AddMode(mode int) {
	t.modes &= mode
}

func (t *Tty) RemoveMode(mode int) {
	t.modes &= ^mode
}
