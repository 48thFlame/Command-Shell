package main

import (
	"fmt"
	"log"
)

func hello(i *CommandHandlerInput) error {
	fmt.Fprintf(i.Stdout, "Args: %v", i.Args)

	return nil
}

func main() {
	cmds := []Command{NewCommand("hello", 1, hello)}
	s, err := NewShell(cmds)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
