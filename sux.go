package main

import (
  "io"
  "os"
  "fmt"
  "flag"
  "os/exec"
  "os/signal"
  "github.com/jroimartin/gocui"
)

type Cmd struct {
  command exec.Cmd
  stdin io.WriteCloser
  stderr io.ReadCloser
  stdout io.ReadCloser
}

func parseCommands(commands []string) []Cmd {
  cmds := make([]Cmd, len(commands))

  for i, cmd := range commands {
    command := exec.Command(cmd)

    stdin, err := command.StdinPipe()
    if err != nil {
      panic(err)
    }
    stdout, err := command.StdoutPipe()
    if err != nil {
      panic(err)
    }
    stderr, err := command.StderrPipe()
    if err != nil {
      panic(err)
    }

    cmds[i] = Cmd {
      command: *command,
      stdin: stdin,
      stderr: stderr,
      stdout: stdout,
    }

    command.Start()
  }

  return cmds
}

func setLayout() {

}


func main() {
  // sigchan := make(chan os.Signal, 1)
  // x := make(chan os.Signal, 1)
  // signal.Notify(x, os.Interrupt)
  // go func() {
  //   for _ = range sigchan {
  //     // <-x
  //     fmt.Println("Got interrupt")
  //   }
  // }()

  if len(flag.Args()) == 0 {
    fmt.Println("Given 0 commands to run.")
    return
  }

  cmds := parseCommands(flag.Args())

  g := gocui.NewGui()
  if err := g.Init(); err != nil {
    panic(err)
  }
  defer g.Close()

  g.SetLayout(layout)

  if err := initKeybindings(g); err != nil {
    panic(err)
  }

  err = g.MainLoop()
  if err != nil && err != gocui.Quit {
    panic(err)
  }

  // go io.Copy(stdin, os.Stdin)
  // go io.Copy(os.Stdout, stdout)
  // go io.Copy(os.Stderr, stderr)

  // for _, cmd := range cmds {
  //   if err := cmd.command.Wait(); err != nil {
  //     panic(err)
  //   }
}
