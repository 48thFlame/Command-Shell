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
	cmds := []Command{NewCommand("hello", 0, hello)}
	s, err := NewShell(cmds)
	if err != nil {
		log.Fatal(err)
	}

	s.Run()
}
