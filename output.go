package main

import (
	"flag"
	"fmt"
	"github.com/chrisseto/sux/pane"
	"github.com/nsf/termbox-go"
	"strings"
)

var (
	quitChan      chan bool
	selectChan    chan *pane.Pane
	selectedIndex = 0
	SelectedPane  *pane.Pane
	RunningPanes  []*pane.Pane
)

func RunPanes() error {
	quitChan = make(chan bool)
	selectChan = make(chan *pane.Pane)
	width, height := termbox.Size()
	cmds := strings.Split(strings.Join(flag.Args(), " "), ",")

	RunningPanes = make([]*pane.Pane, len(cmds))

	for i, cmd := range cmds {
		cmdsp := strings.Split(strings.Trim(cmd, " "), " ")
		RunningPanes[i] = pane.CreatePane(cmdsp[0], cmdsp[1:], width, height-1)
		err := RunningPanes[i].Start()

		if err != nil {
			return err
		}
	}

	SelectedPane = RunningPanes[selectedIndex]
	selectChan <- SelectedPane
	return nil
}

func Redraw() {
	select {
	case redraw <- struct{}{}:
	default: //Failed to send, a redraw is already happening
	}
}

func OutputLoop() {
	selected := <-selectChan
	selected.Draw(0, 0)
	WriteStatusBar(selected)
	termbox.Flush()
	for {
		select {
		case <-selected.ShouldRedraw:
			selected.Draw(0, 0)
			WriteStatusBar(selected)
			termbox.Flush()
		case <-redraw:
			selected.Draw(0, 0)
			WriteStatusBar(selected)
			termbox.Flush()
		case selected = <-selectChan:
		case <-quitChan:
			return
		}
	}
}

func setPane(index int) {
	SelectedPane = RunningPanes[index]
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	SelectedPane.Draw(0, 0)
	WriteStatusBar(SelectedPane)
	termbox.Flush()
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
	for _, pane := range RunningPanes {
		pane.Stop()
	}
}

func WriteStatusBar(pane *pane.Pane) {
	width, height := termbox.Size()
	statusString := fmt.Sprintf("Pane #%d Command %s Args %v %s Mode", selectedIndex, pane.Prog(), pane.Args(), CurrentMode.Name)
	i := 0
	for _, char := range statusString {
		termbox.SetCell(i, height-1, char, termbox.ColorBlack, termbox.ColorGreen)
		i++
		if i > width {
			return
		}
	}
	for ; i < width; i++ {
		termbox.SetCell(i, height-1, ' ', termbox.ColorBlack, termbox.ColorGreen)
	}
}
