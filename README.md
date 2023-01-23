# Command-shell

A go library for creating simple command-based `tui`'s (terminal user interfaces).
Built on top of this [keyboard package](https://github.com/eiannone/keyboard).

## Features

- Default `exit` command
- Minimal number of arguments safety
- Commands can store data
- Up-arrow for previous command
- History
- Default & automatic `help` command
- Customizable `LinePrefix`
- And more...

## TODO

- Tab autocomplete
- Cursor moving around line

## Example

### output

```terminal
> add 2 6
8
>
```

### code

```go
package main

// imports
import (
    "fmt"
    "strconv"

    shell "github.com/48thFlame/Command-Shell"
)

// addCommand represents the `add` command you can see in the output section
func addCommand(input *shell.CommandInput) error {
    args := input.Args

    // convert arguments to numbers
    // can safely assume there will be at lease 2 arguments
    // we specify this in the `shell.NewCommand` constructor
    a, err := strconv.Atoi(args[0])
    if err != nil {
        return err
    }
    b, err := strconv.Atoi(args[1])
    if err != nil {
        return err
    }

    // output the addition
    fmt.Fprint(input.Stdout, a+b)

    // no errors occurred
    return nil
}

func timesRanCommand(input *shell.CommandInput) error {
    // get number of times ran
    // if doesn't exist the error will be thrown away `_` and val will be the default int -> 0
    val, _ := input.Cmd.Data["times"].(int)

    // increment `times` by 1
    input.Cmd.Data["times"] = val + 1

    // output times ran
    fmt.Fprint(input.Stdout, val)

    // no errors
    return nil
}

func main() {
    s, err := shell.NewShell(
        shell.NewCommand("add", 2, addCommand),
        shell.NewCommand("times", 0, timesRanCommand),
    )
    if err != nil {
        panic(err)
    }

    err = s.Run()
    if err != nil {
        panic(err)
    }
}

```

## Install

`go get https://pkg.go.dev/github.com/48thFlame/Command-Shell`
