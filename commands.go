package main

import (
	"flag"
	"github.com/nsf/termbox-go"
	"os"
	"strings"
)

var (
	quitChan        chan bool
	selectChan      chan *Cmd
	selectedIndex   = 0
	SelectedCommand *Cmd
	RunningCommands []*Cmd
)

func StartCommands() error {
	quitChan = make(chan bool)
	selectChan = make(chan *Cmd)
	cmds := strings.Split(strings.Join(flag.Args(), " "), ",")

	RunningCommands = make([]*Cmd, len(cmds))

	for i, cmd := range cmds {
		cmdsp := strings.Split(strings.Trim(cmd, " "), " ")
		RunningCommands[i] = Command(cmdsp[0], cmdsp[1:]...)
		err := RunningCommands[i].Start()

		if err != nil {
			return err
		}
	}

	SelectedCommand = RunningCommands[selectedIndex]
	selectChan <- SelectedCommand
	return nil
}

func OutputLoop() {
	selected := <-selectChan
	for {
		select {
		case buf := <-selected.Output:
			os.Stdout.Write(buf)
		case selected = <-selectChan:
		case <-quitChan:
			return
		}
	}
}

func setCommand() {
	SelectedCommand = RunningCommands[selectedIndex]
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Sync()
	selectChan <- SelectedCommand
}

func NextCommand() {
	selectedIndex++
	if selectedIndex >= len(RunningCommands) {
		selectedIndex = 0
	}
	setCommand()
}

func PrevCommand() {
	selectedIndex--
	if selectedIndex < 0 {
		selectedIndex = len(RunningCommands) - 1
	}
	setCommand()
}

func EndCommands() {
	quitChan <- true
	for _, cmd := range RunningCommands {
		cmd.Close()
	}
}
