package shell

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
)

// Return a new `Shell` and any `error` that occurred.
// Takes as input a `[]Command` that the shell should know how to run,
// Use the `NewCommand` function to create new `Command`s
// NewShell is the constructor for `Shell` type
func NewShell(cmds ...*Command) (*Shell, error) {
	s := &Shell{}

	s.LinePrefix = "> "

	s.Path = make(map[string]*Command)
	s.history = []string{}
	s.historyIndex = defaultHistoryIndex
	s.quitChannel = make(chan bool, 1)

	cmds = append(defaultCommands(), cmds...)
	for _, c := range cmds {
		if _, exists := s.Path[c.Name]; !exists {
			s.Path[c.Name] = c
		} else { // command with that name already exists
			return nil, fmt.Errorf("duplicate command name: \"%v\"", c.Name)
		}
	}

	return s, nil
}

// Shell type the main part of the program, use
//
//	NewShell()
//
// to create a new Shell
type Shell struct {
	LinePrefix   string              // Prefix printed at beginning of every command line
	Path         map[string]*Command // Similar to unix $PATH
	history      []string            // Similar to unix-terminal history
	historyIndex int
	currentInput string
	quitChannel  chan bool
}

// errorHandle handles errors by printing the error to the terminal
func (s *Shell) errorHandle(err error) {
	fmt.Println()
	fmt.Printf("ERROR: %v\n", err)
}

func (s *Shell) defaultDisplay() {
	fmt.Print("\033[2K\r")
	fmt.Print(s.LinePrefix)
	fmt.Print(s.getInput())
}

func (s *Shell) mainLoop() {
	fmt.Println("Enter \"exit\" to quit")

	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-s.quitChannel:
			return
		default:
			s.defaultDisplay()
		}

		<-ticker.C
	}
}

// GetHistory returns command history similar to Unix history,
// user can use default `history` command to display history
func (s *Shell) GetHistory() []string {
	return s.history
}

// addToHistory add element `i` to beginning of slice `s.history`
func (s *Shell) addToHistory(i string) {
	s.history = append(s.history, "")
	copy(s.history[1:], s.history)
	s.history[0] = i
	s.historyIndex = defaultHistoryIndex
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
