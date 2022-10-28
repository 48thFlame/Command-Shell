package main

import (
	"fmt"
	"io"
)

func hello(args ArgsType, w io.Writer) error {
	fmt.Fprintln(w, "yoooo")

	return nil
}

func main() {
	// cmds := []command{hello}
	s, err := NewShell(hello)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	s.Run()
}
