package main

import (
	"flag"
	"github.com/nsf/termbox-go"
	"os"
	"strings"
  "fmt"
	// "strconv"
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
	uwidth, uheight := uint16(width), uint16(height-1)
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
  f, _ := os.Create("log.log")
  defer f.Close()
	WriteStatusBar(selected.Prog)
	termbox.Flush()
	for {
		select {
		case cells := <-selected.CellOutput:
      fmt.Fprintf(f, "%v", cells)
			for _, cell := range cells {
        fmt.Fprintf(f, "[%v, %v, %#U]\n", cell.x, cell.y, cell.Ch)
				termbox.SetCell(cell.x, cell.y, cell.Ch, cell.Fg, cell.Bg)
			}
			WriteStatusBar(selected.Prog)
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
	for y, line := range SelectedPane.Cells() {
		for x, cell := range line {
			termbox.SetCell(x, y, cell.Ch, cell.Fg, cell.Bg)
		}
	}
	WriteStatusBar(SelectedPane.Prog)
	// termbox.Sync()
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
	for _, cmd := range RunningPanes {
		cmd.Close()
	}
}

func WriteStatusBar(prog string) {
	width, height := termbox.Size()
	i := 0
	for _, char := range prog {
		i++
		termbox.SetCell(i, height-1, char, termbox.ColorBlack, termbox.ColorGreen)
	}
	for ; i < width; i++ {
		termbox.SetCell(i, height-1, ' ', termbox.ColorBlack, termbox.ColorGreen)
	}
	// termbox.HideCursor()
	// termbox.Flush()
}
