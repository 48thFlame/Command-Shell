package shell

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

// NewCommand returns a `Command` type,
// Args:
// 	name // name used to call command
// 	minNumOfArgs // minimum number of arguments needed for command to run (the `Shell` won't call the command without enough args)
// 	handler // the handler to call for the command
func NewCommand(name string, minNumOfArgs int, handler HandlerType) Command {
	return Command{
		Name:         name,
		MinNumOfArgs: minNumOfArgs,
		Handler:      handler,
	}
}

type CommandHandlerInput struct {
	Shell  *Shell    // the `Shell` struct
	Stdout io.Writer // write any command output to here, avoid printing to `os.Stdout`
	Args   []string  // command-line arguments
	Cmd    *Command  // the `Command` struct
}

// HandlerType is the type for
// 	Command.Handler
// Its a function called when running its command
//	*CommandHandlerInput -> error
type HandlerType func(*CommandHandlerInput) error

// Command type, use
// 	NewCommand()
// to create a new Command
type Command struct {
	Name         string      // commands name
	MinNumOfArgs int         // minium number of arguments needed to run command
	Handler      HandlerType // the `HandlerType` for this command
}

// Return a new `Shell` and any `error` that occurred.
// Takes as input a `[]Command` that the shell should know how to run,
// Use the `NewCommand` function to create new `Command`s

// NewShell is the constructor for `Shell` type
//	[]Command -> (*Shell, error)
func NewShell(cmds ...Command) (*Shell, error) {
	s := &Shell{}

	s.prefix = "> "
	s.currentInput = ""

	s.Path = make(map[string]Command)
	s.History = []string{}

	s.specialKeyChannel = make(chan keyboard.Key)
	s.quitChannel = make(chan bool, 1)

	exit := NewCommand("exit", 0, func(i *CommandHandlerInput) error {
		i.Shell.quitChannel <- true
		return nil
	})
	s.Path[exit.Name] = exit

	for _, c := range cmds {
		if _, exists := s.Path[c.Name]; !exists {
			s.Path[c.Name] = c
		} else { // command with that name already exists
			return nil, fmt.Errorf("duplicate command name: \"%v\"", c.Name)
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

// errorHandle handles errors by printing the error to the terminal
func (s *Shell) errorHandle(err error) {
	fmt.Println()
	fmt.Printf("ERROR: %v\n", err)
}

// listenToKeyboardInput listens for keyboard input and update the
//	s.currentInput // (user typing)
// or send the keys to
//	s.specialKeyChannel // (Enter, CtrlC etc)
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

// runCommand runs the command at `s.currentInput`, returns any errors
func (s *Shell) runCommand() error {
	input := strings.TrimSpace(s.currentInput)
	args := strings.Split(input, " ")
	cmdName, args := args[0], args[1:] // pop(0)

	cmd, exists := s.Path[cmdName]
	if !exists {
		return fmt.Errorf("command \"%v\" not found", cmdName)
	}

	numOfArgs := len(args) // -1 minus the command itself
	if numOfArgs < cmd.MinNumOfArgs {
		return fmt.Errorf("command \"%v\" needs %v arguments, but %v where provided", cmd.Name, cmd.MinNumOfArgs, numOfArgs)
	}

	cmdStdout := new(strings.Builder)

	cmdInput := &CommandHandlerInput{
		Shell:  s,
		Stdout: cmdStdout,
		Args:   args,
		Cmd:    &cmd,
	}

	err := cmd.Handler(cmdInput)
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

// handlerSpecialKey handles a special key (one that isn't plain text, Enter, CtrlC etc)
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

// Run runs the `Shell` and return any `error`s,
// its a blocking function only returns when user exits the `Shell` (CtrlC, etc)
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
