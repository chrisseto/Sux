package main

import (
	"flag"
	"fmt"
	"github.com/nsf/termbox-go"
	"log"
	"os"
)

func MaybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	Quit   chan struct{}
	redraw chan struct{}
)

func main() {
	f, err := os.OpenFile("sux.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)

	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("sux: no commands given")
		fmt.Println("Usage sux [command ...]")
		return
	}

	DefaultMode = InputMode
	CurrentMode = DefaultMode

	Quit = make(chan struct{})
	redraw = make(chan struct{})

	MaybePanic(termbox.Init())

	termbox.SetInputMode(termbox.InputEsc)
	termbox.SetOutputMode(termbox.Output256)

	defer termbox.Close()
	defer EndPanes()

	go InputLoop()
	go OutputLoop()

	MaybePanic(RunPanes())

	<-Quit
}
