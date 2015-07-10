package main

import (
  // "io"
  "os"
	"flag"
	"strings"
	"github.com/nsf/termbox-go"
)

var (
  quitChan chan bool
	selectedIndex   = 0
	SelectedCommand *Cmd
	RunningCommands []*Cmd
)

func StartCommands() error {
  quitChan = make(chan bool)
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
  go pipeSelected()
	return nil
}


func pipeSelected() {
    for {
      select {
      case <- quitChan:
        return
      case buf := <-SelectedCommand.Output:
        os.Stdout.Write(buf)

      // default:
      //   io.CopyN(os.Stdout, SelectedCommand.Pty, 1)
      }
    }
}

func setCommand() {
  quitChan <- true
  SelectedCommand = RunningCommands[selectedIndex]
  termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
  termbox.Sync()
  go pipeSelected()
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


