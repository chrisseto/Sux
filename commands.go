package main

import (
	"flag"
	"github.com/nsf/termbox-go"
	"os"
	"strings"
)

var (
	quitChan      chan bool
	selectChan    chan *Pane
	selectedIndex = 0
	SelectedPane  *Pane
	RunningPanes  []*Pane
)

func RunPanes() error {
	quitChan = make(chan bool)
	selectChan = make(chan *Pane)
	width, height := termbox.Size()
	uwidth, uheight := uint16(width), uint16(height)
	cmds := strings.Split(strings.Join(flag.Args(), " "), ",")

	RunningPanes = make([]*Pane, len(cmds))

	for i, cmd := range cmds {
		cmdsp := strings.Split(strings.Trim(cmd, " "), " ")
		RunningPanes[i] = CreatePane(uwidth, uheight, cmdsp[0], cmdsp[1:]...)
		err := RunningPanes[i].Start()

		if err != nil {
			return err
		}
	}

	SelectedPane = RunningPanes[selectedIndex]
	selectChan <- SelectedPane
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

func setPane(index int) {
	SelectedPane = RunningPanes[index]
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Sync()
	selectChan <- SelectedPane
}

func NextPane() {
	selectedIndex++
	if selectedIndex >= len(RunningPanes) {
		selectedIndex = 0
	}
	setPane(selectedIndex)
}

func PrevPane() {
	selectedIndex--
	if selectedIndex < 0 {
		selectedIndex = len(RunningPanes) - 1
	}
	setPane(selectedIndex)
}

func EndPanes() {
	quitChan <- true
	for _, cmd := range RunningPanes {
		cmd.Close()
	}
}
