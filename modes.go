package main

import "github.com/nsf/termbox-go"

type InputTrigger int
type InputHandler func([]byte, termbox.Event)

type Mode struct {
	Name        string
	Trigger     InputTrigger
	HandleInput InputHandler
	SubModes    []Mode
}

var (
	//Populated in main() to avoid
	//Cyclic references
	CurrentMode Mode
	DefaultMode Mode

	// ScrollMode = Mode{
	// 	Name:        "Scroll",
	// 	Trigger:     InputTrigger('['),
	// 	HandleInput: ScrollModeHandler,
	// 	SubModes:    nil,
	// }

	CommandMode = Mode{
		Name:        "Command",
		Trigger:     InputTrigger(termbox.KeyCtrlB),
		HandleInput: CommandModeHandler,
		SubModes:    nil,
		// SubModes:    []Mode{ScrollMode},
	}

	InputMode = Mode{
		Name: "Input",
		// Trigger:     nil,
		HandleInput: InputModeHandler,
		SubModes:    []Mode{CommandMode},
	}
)

func SetMode(newMode *Mode) {
	if newMode == nil {
		CurrentMode = DefaultMode
	} else {
		CurrentMode = *newMode
	}
	Redraw()
}

func InputModeHandler(raw []byte, ev termbox.Event) {
	SelectedPane.Send(raw)
}

// func ScrollModeHandler(raw []byte, ev termbox.Event) {
// 	switch ev.Key {
// 	case termbox.KeyEsc:
// 		SetMode(&DefaultMode)
// 	case termbox.KeyArrowUp:
// 		SelectedPane.Scroll(-1)
// 	case termbox.KeyArrowDown:
// 		SelectedPane.Scroll(1)
// 	}
// }

func CommandModeHandler(raw []byte, ev termbox.Event) {
	switch ev.Key {
	case termbox.KeyCtrlC:
		Quit <- struct{}{}
	case termbox.KeyArrowRight:
		NextPane()
	case termbox.KeyArrowLeft:
		PrevPane()
	}
	SetMode(&DefaultMode)
}
