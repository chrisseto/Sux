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

var InputChan chan error

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("Given 0 commands to run.")
		return
	}

	InputChan = make(chan error)

	MaybePanic(termbox.Init())

	defer termbox.Close()
	defer EndCommands()

	go InputLoop()
	go OutputLoop()

	MaybePanic(StartCommands())

	<-InputChan
}
