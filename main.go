package main

import (
	"fmt"
	"io"
	"log"
)

func hello(s *Shell, stdout io.Writer, args []string) error {
	fmt.Fprintf(stdout, "Args: %v", args)

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
