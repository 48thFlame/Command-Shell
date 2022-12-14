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
	}
}

// CommandInput is the input type for `CommandHandler` function
type CommandInput struct {
	Shell  *Shell    // the `Shell` struct
	Stdout io.Writer // write any command output to here, avoid printing to `os.Stdout`
	Args   []string  // command-line arguments
	Cmd    *Command  // the `Command` struct
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
	Name         string      // commands name
	MinNumOfArgs int         // minium number of arguments needed to run command
	Handler      HandlerType // the `HandlerType` for this command
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
// TODO: fix bug when has 2 spaces!!
// TODO: fix bug when in strange historyIndex place
func (s *Shell) runCommand() error {
	input := strings.TrimSpace(s.getInput())

	if input == "" {
		return nil
	}

	args := strings.Split(input, " ")
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

	// when done should insert at beginning of history empty slot to be new input
	s.history = prependString(s.history, "")

	return nil
}
