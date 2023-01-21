package shell

import (
	"fmt"
	"io"
	"strings"
)

// NewCommand returns a `Command` type,
// Args:
//
//	name // name used to call command
//	minNumOfArgs // minimum number of arguments needed for command to run (the `Shell` won't call the command without enough args)
//	handler // the handler to call for the command
func NewCommand(name string, minNumOfArgs int, handler HandlerType) *Command {
	return &Command{
		Name:         name,
		MinNumOfArgs: minNumOfArgs,
		Handler:      handler,
		Data:         make(map[string]interface{}),
	}
}

// CommandInput is the input type for `CommandHandler` function
type CommandInput struct {
	Shell  *Shell    // the `Shell` type
	Stdout io.Writer // write any command output to here, avoid printing to `os.Stdout`
	Args   []string  // command-line arguments
	Cmd    *Command  // the `Command` type
}

// HandlerType is the type for
//
//	Command.Handler
//
// Its a function called when running its command
type HandlerType func(*CommandInput) error

// Command type, use
//
//	NewCommand()
//
// to create a new Command
type Command struct {
	Name         string                 // commands name
	MinNumOfArgs int                    // minium number of arguments needed to run command
	Handler      HandlerType            // the `HandlerType` for this command
	Data         map[string]interface{} // data
}

func defaultCommands() []*Command {
	var exitCommand = NewCommand("exit", 0, func(ci *CommandInput) error {
		ci.Shell.exit()
		return nil
	})

	var historyCommand = NewCommand("history", 0, func(ci *CommandInput) error {
		// TODO: make prettyprint
		fmt.Fprint(ci.Stdout, ci.Shell.GetHistory())
		return nil
	})

	return []*Command{exitCommand, historyCommand}
}

// runCommand runs the command from input
func (s *Shell) runCommand() error {
	args := parsInput(s.getInput())

	if len(args) == 0 {
		fmt.Println() // Print blank line so put a new >
		return nil
	}

	cmdName, args := args[0], args[1:] // pop(0)

	cmd, exists := s.Path[cmdName]
	if !exists {
		return fmt.Errorf("command \"%v\" not found", cmdName)
	}

	numOfArgs := len(args)
	if numOfArgs < cmd.MinNumOfArgs {
		return fmt.Errorf("command \"%v\" needs %v arguments, but %v where provided", cmd.Name, cmd.MinNumOfArgs, numOfArgs)
	}

	cmdStdout := new(strings.Builder)

	cmdInput := &CommandInput{
		Shell:  s,
		Stdout: cmdStdout,
		Args:   args,
		Cmd:    cmd,
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

	s.addToHistory(s.getInput())
	s.replaceInput("")

	return nil
}

func parsInput(i string) []string {
	var args []string
	input := strings.TrimSpace(i)

	if input == "" {
		return []string{}
	}

	argsTemp := strings.Split(input, " ")

	for _, str := range argsTemp {
		if str != "" {
			args = append(args, str)
		}
	}

	return args
}
