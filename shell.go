package main

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

// Returns a `Command` type
// Args:
// `name` - name used to call command
// `minNumOfArgs` - minimum number of arguments needed for command to run (the `Shell` won't call the command without enough args)
// `handler` - the handler to call for the command
func NewCommand(name string, minNumOfArgs int, handler handlerType) Command {
	return Command{
		name:         name,
		minNumOfArgs: minNumOfArgs,
		handler:      handler,
	}
}

// The command handlerType function type
// takes in a `*Shell`, `io.Writer`, `[]string` and returns an error
type handlerType func(shell *Shell, stdout io.Writer, args []string) error

// The Command type
// Use `NewCommand` function to create a command
type Command struct {
	name         string
	minNumOfArgs int
	handler      handlerType
}

// Return a new `Shell` and any `error` that occurred.
// Takes as input a `[]Command` that the shell should know how to run,
// Use the `NewCommand` function to create new `Command`s
func NewShell(cmds []Command) (*Shell, error) {
	s := &Shell{}

	s.prefix = "> "
	s.currentInput = ""

	s.Path = make(map[string]Command)
	s.History = []string{}

	s.specialKeyChannel = make(chan keyboard.Key)
	s.quitChannel = make(chan bool, 1)

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
			return nil, fmt.Errorf("duplicate command name: \"%v\"", c.name)
		}
	}

	return s, nil
}

// The `Shell` type
type Shell struct {
	Path              map[string]Command // Similar to unix $PATH
	History           []string           // Similar to unix-terminal history
	prefix            string
	currentInput      string
	specialKeyChannel chan keyboard.Key
	quitChannel       chan bool
}

// Handle errors by printing the error to the terminal
func (s *Shell) errorHandle(err error) {
	fmt.Println()
	fmt.Printf("ERROR: %v\n", err)
}

// Listen for keyboard input and update the `s.currentInput` (user typing) or send the keys to `s.specialKeyChannel`` (Enter, CtrlC etc)
func (s *Shell) listenToKeyboardInput(keysChannel <-chan keyboard.KeyEvent) {
	for {
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

// Run the command using s.currentInput, return any errors
func (s *Shell) runCommand() error {
	input := strings.TrimSpace(s.currentInput)
	args := strings.Split(input, " ")
	cmdName := args[0]

	cmd, exists := s.Path[cmdName]
	if !exists {
		return fmt.Errorf("command \"%v\" not found", cmdName)
	}

	numOfArgs := len(args) - 1 // -1 minus the command itself
	if numOfArgs < cmd.minNumOfArgs {
		return fmt.Errorf("command \"%v\" needs %v arguments, but %v where provided", cmd.name, cmd.minNumOfArgs, numOfArgs)
	}

	cmdStdout := new(strings.Builder)

	err := cmd.handler(s, cmdStdout, args)
	if err != nil {
		return err
	}

	if cmdStdout.Len() > 0 { // only if has something to print
		// should print "\n" so text appears on newline and not `> cmdOutput`
		// but instead `> cmd\nOutput`
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

// Handle a special key (one that isn't plain text, Enter, CtrlC etc)
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

// Run the `Shell` return any `error`s, its a blocking function only returns when user exits the `*Shell`
func (s *Shell) Run() error {
	keyEvents, err := keyboard.GetKeys(1)
	if err != nil {
		return err
	}

	go s.listenToKeyboardInput(keyEvents)

	s.mainLoop()

	err = keyboard.Close()
	if err != nil {
		return err
	}

	fmt.Println()

	return nil
}
