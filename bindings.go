package main

import "github.com/nsf/termbox-go"

type Mode int

const (
	NormalMode Mode = iota
	ScrollMode
)

var (
	leader      = termbox.KeyCtrlB
	currentMode = NormalMode
)

func leaderPress() bool {
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyCtrlC:
			InputChan <- nil
			return true
		case termbox.KeyArrowRight:
			NextPane()
		case termbox.KeyArrowLeft:
			PrevPane()
		default:
			switch ev.Ch {
			case '[':
				currentMode = ScrollMode
			}
		}
	case termbox.EventError:
		panic(ev.Err)
	}
	return false
}

func InputLoop() {
	var raw = make([]byte, 5)
	for {
		raw = make([]byte, 5)
		switch ev := termbox.PollRawEvent(raw); ev.Type {
		case termbox.EventError:
			panic(ev.Err)

		case termbox.EventRaw:
			raw = raw[:ev.N]
			switch ev := termbox.ParseEvent(raw); ev.Type {
			case termbox.EventError:
				panic(ev.Err)

			case termbox.EventKey:
				switch ev.Key {
				case leader:
					if leaderPress() {
						return
					}

				default:
					switch currentMode {
					case NormalMode:
						NormalModeHandler(raw, ev)
					case ScrollMode:
						ScrollModeHandler(raw, ev)
					}
				}
			}
		}
	}
}

//TODO Make more plugable
func NormalModeHandler(raw []byte, ev termbox.Event) {
	SelectedPane.Pty.Write(raw)
}

func ScrollModeHandler(raw []byte, ev termbox.Event) {
	switch ev.Key {
	case termbox.KeyEsc:
		currentMode = NormalMode
	case termbox.KeyArrowUp:
		SelectedPane.Scroll(-1)
	case termbox.KeyArrowDown:
		SelectedPane.Scroll(1)
	}
}
