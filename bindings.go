package main

import "github.com/nsf/termbox-go"

var (
	leader                 = termbox.KeyCtrlB
	outstandingLeaderPress = false
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
		}
	case termbox.EventError:
		panic(ev.Err)
	}
	return false
}

func InputLoop() {
	var raw = make([]byte, 1)
	for {
		switch ev := termbox.PollRawEvent(raw); ev.Type {
		case termbox.EventError:
			panic(ev.Err)

		case termbox.EventRaw:
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
					SelectedPane.Pty.Write(raw)
				}
			}
		}
	}
}
