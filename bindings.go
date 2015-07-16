package main

import "github.com/nsf/termbox-go"

func InputLoop() {
	var trigger InputTrigger
	var raw = make([]byte, 5)

Loop:
	for {
		raw = make([]byte, 5)

		ev := termbox.PollRawEvent(raw)
		if ev.Type == termbox.EventError {
			panic(ev.Err)
		}

		raw = raw[:ev.N] // Truncate raw
		ev = termbox.ParseEvent(raw)
		if ev.Type == termbox.EventError {
			panic(ev.Err)
		}

		if ev.Ch == 0x0 {
			trigger = InputTrigger(ev.Key)
		} else {
			trigger = InputTrigger(ev.Ch)
		}

		for _, submode := range CurrentMode.SubModes {
			if submode.Trigger == trigger {
				SetMode(&submode)
				continue Loop
			}
		}

		CurrentMode.HandleInput(raw, ev)
	}
}
