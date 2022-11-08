package main

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

/* Returns a `Command` type
Args:
`name` - name used to call command
`minNumOfArgs` - minimum number of arguments needed for command to run (will error if doesn't have enough args)
`handler` - the handler to call for the command*/
func NewCommand(name string, minNumOfArgs int, handler handlerType) Command {
	return Command{
		name:         name,
		minNumOfArgs: minNumOfArgs,
		handler:      handler,
	}
}

/* The command handlerType function type
takes in a `*Shell` and `[]string`*/
type handlerType func(shell *Shell, stdout io.Writer, args []string) error

/* The Command type
Use `NewCommand` function to create a command*/
type Command struct {
	name         string
	minNumOfArgs int
	handler      handlerType
}

func NewShell(cmds []Command) (*Shell, error) {
	s := &Shell{}

	s.prefix = "> "
	s.currentInput = ""
	// s.quitChannel = make(chan bool)
	s.History = []string{}
	s.specialKeyChannel = make(chan keyboard.Key)
	s.quitChannel = make(chan bool, 1)
	// // s.runCmdChannel = make(chan bool)

	s.Path = make(map[string]Command)

	exit := NewCommand("exit", 0, func(shell *Shell, stdout io.Writer, args []string) error {
		fmt.Fprintln(stdout, "Bye bye...")
		s.quitChannel <- true
		return nil
	})
	s.Path[exit.name] = exit

	for _, c := range cmds {
		if _, exists := s.Path[c.name]; !exists {
			s.Path[c.name] = c
		} else { // command with that name already exists
			return nil, fmt.Errorf("duplicate command name: %v", c.name)
		}
	}

	return s, nil
}

type Shell struct {
	Path              map[string]Command // Similar to unix $PATH
	History           []string           // Similar to unix-terminal history
	prefix            string
	currentInput      string
	specialKeyChannel chan keyboard.Key
	quitChannel       chan bool
}

// Handle errors by printing the error to  the terminal
func (s *Shell) errorHandle(err error) {
	fmt.Println()
	fmt.Printf("Error: %v\n", err)
}

// Listen for keyboard input and update the s.currentInput or send the keys to s.specialKeyChannel
func (s *Shell) listenToKeyboardInput(keysChannel <-chan keyboard.KeyEvent) {
	for {
		// char, key, err := keyboard.GetKey()
		event := <-keysChannel
		key, char, err := event.Key, event.Rune, event.Err
		// fmt.Printf("--- %v, %v, %v ---\n", char, key, err)
		if err != nil {
			s.errorHandle(err)
		}
		switch key {
		case keyboard.KeyBackspace, keyboard.KeyBackspace2:
			inputLen := len(s.currentInput)
			if inputLen >= 1 {
				s.currentInput = s.currentInput[:inputLen-1]
			}
		case keyboard.KeySpace:
			s.currentInput += " "
		case keyboard.Key(0): // if is not a special key, its a letter or a symbol
			s.currentInput += string(char)
		default:
			s.specialKeyChannel <- key
			// // case keyboard.KeyDelete, keyboard.KeyEnd, keyboard.KeyEnter, keyboard.KeyHome, keyboard.KeyArrowLeft, keyboard.KeyArrowRight, keyboard.KeyArrowUp, keyboard.KeyArrowDown
		}
	}
}

func (s *Shell) runCommand() error {
	args := strings.Split(s.currentInput, " ")
	cmdName := args[0]

	cmd, exists := s.Path[cmdName]
	if !exists {
		return fmt.Errorf("command named: \"%v\", not found", cmdName)
	}

	cmdStdout := new(strings.Builder)

	err := cmd.handler(s, cmdStdout, args)
	if err != nil {
		return err
	}

	if cmdStdout.Len() > 0 { // only if has something to print
		// should print "\n" so text appears on newline and not `> cmdOutput`
		fmt.Println()
		out := cmdStdout.String()
		if strings.HasSuffix(out, "\n") {
			fmt.Print(out)
		} else {
			fmt.Println(out)
		}
	}
	s.currentInput = ""

	return nil
}

func (s *Shell) handleSpecialKey(key keyboard.Key) error {
	switch key {
	case keyboard.KeyCtrlC, keyboard.KeyCtrlZ:
		s.quitChannel <- true
	case keyboard.KeyEnter:
		return s.runCommand()
	}
	return nil
}

func (s *Shell) defaultDisplay() {
	fmt.Print("\033[2K\r")
	fmt.Print(s.prefix)
	fmt.Print(s.currentInput)
}

func (s *Shell) mainLoop() {
	fmt.Println("Enter \"exit\" to quit")

	ticker := time.Tick(100 * time.Millisecond)
	for {
		select {
		case <-s.quitChannel:
			return
		case key := <-s.specialKeyChannel:
			err := s.handleSpecialKey(key)
			if err != nil {
				s.errorHandle(err)
			}
		default:
			s.defaultDisplay()
		}

		<-ticker
	}
}

func (s *Shell) close() error {
	return keyboard.Close()
}

func (s *Shell) Run() error {
	keyEvents, err := keyboard.GetKeys(1)
	if err != nil {
		return err
	}

	go s.listenToKeyboardInput(keyEvents)

	s.mainLoop()

	err = s.close()
	if err != nil {
		return err
	}

	fmt.Println()

	return nil
}
