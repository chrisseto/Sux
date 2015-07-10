package main

import (
	"flag"
	// "io"
	// "os"
	// "syscall"
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

	MaybePanic(StartCommands())
	MaybePanic(termbox.Init())

	defer termbox.Close()
	defer EndCommands()

	go InputLoop()

	<-InputChan
}
