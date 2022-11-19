package main

import (
	"github.com/48thFlame/CommandShell"
	"fmt"
	"log"
)

func hello(i *shell.CommandHandlerInput) error {
	fmt.Fprintf(i.Stdout, "Args: %v", i.Args)

	return nil
}

func main() {
	cmds := []shell.Command{shell.NewCommand("hello", 1, hello)}
	s, err := shell.NewShell(cmds)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
