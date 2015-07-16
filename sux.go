package main

import (
	"flag"
	"fmt"
	"github.com/nsf/termbox-go"
)

func MaybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	Quit   chan struct{}
	Redraw chan struct{}
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("Given 0 commands to run.")
		return
	}

	DefaultMode = InputMode
	CurrentMode = DefaultMode

	Quit = make(chan struct{})
	Redraw = make(chan struct{})

	MaybePanic(termbox.Init())

	termbox.SetInputMode(termbox.InputEsc)

	defer termbox.Close()
	defer EndPanes()

	go InputLoop()
	go OutputLoop()

	MaybePanic(RunPanes())

	<-Quit
}
