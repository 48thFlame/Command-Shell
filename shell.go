package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

// TODO:
// arrow-keys cursor movement
// https://tldp.org/HOWTO/Bash-Prompt-HOWTO/x361.html
/*
case keyboard.KeyArrowLeft:
	fmt.Print("\033[1D")
case keyboard.KeyArrowRight:
	fmt.Print("\033[1C")
*/

func NewShell(cmd Command) (*Shell, error) {
	s := &Shell{}

	s.linePrefix = "> "
	s.history = []string{}
	s.commands = make(map[string]Command)

	s.commands["hello"] = cmd

	err := keyboard.Open()
	if err != nil {
		return nil, err
	}

	return s, nil
}

type ArgsType []string
type Command func(ArgsType, io.Writer) error

type Shell struct {
	linePrefix  string
	currentLine string
	history     []string
	commands    map[string]Command
}

func (s *Shell) Run() {
	quitChannel := make(chan bool)

	go s.listenToKeyInput(quitChannel)

	fmt.Println("Press ESC to quit")
	ticker := time.Tick(100 * time.Millisecond)
	for {
		select {
		case <-quitChannel:
			return
		default:
			fmt.Print("\033[2K\r")
			fmt.Print(s.linePrefix)
			fmt.Print(s.currentLine)

		}

		<-ticker
	}
}

func (s *Shell) listenToKeyInput(quitChannel chan bool) {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		switch key {
		case keyboard.KeyEsc, keyboard.KeyCtrlC:
			fmt.Println()
			quitChannel <- true
			close(quitChannel)
			return
		case keyboard.KeyBackspace, keyboard.KeyBackspace2:
			if lineLen := len(s.currentLine); lineLen >= 1 {
				s.currentLine = s.currentLine[:lineLen-1]
			}
		case keyboard.KeySpace:
			s.currentLine += " "
		case keyboard.KeyEnter:
			s.runCommand()
		// case keyboard.KeyDelete, keyboard.KeyEnd, keyboard.KeyEnter, keyboard.KeyHome, keyboard.KeyArrowLeft, keyboard.KeyArrowRight, keyboard.KeyArrowUp, keyboard.KeyArrowDown
		case keyboard.Key(0): // if is not a special key, its a letter or a symbol
			s.currentLine += string(char)
		}
	}
}

func (s *Shell) runCommand() {
	args := strings.Split(s.currentLine, " ")
	cmd := args[0]
	// args = args[1:]
	fmt.Println()
	fmt.Println(args, cmd)
	b := new(bytes.Buffer)
	s.commands[cmd](args, b)
	fmt.Print(b)
	s.currentLine = ""
	// fmt.Println("command")

}

func (s *Shell) Close() error {
	fmt.Println("Bye bye...")
	return keyboard.Close()
}
