# Command-shell

A go library for creating simple command-based `tui`'s (terminal user interfaces).
Built on top of this [keyboard package](https://github.com/eiannone/keyboard).

## Features

- Default `exit` command
- Minimal number of arguments safety
- Customizable `LinePrefix`
- Up-arrow for previous command
- History
- And more...

## TODO

- Commands can store data
- Default & automatic `help` command
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

// Imports
import (
    "fmt"
    "log"
    "strconv"

    shell "github.com/48thFlame/Command-Shell"
)

// addCommand represents the `add` command you can see in the output section
func addCommand(input *shell.CommandInput) error {
    args := input.Args

    // Convert arguments to numbers
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

    // Output the addition
    fmt.Fprint(input.Stdout, a+b)

    // No errors occurred
    return nil
}

func main() {
    s, err := shell.NewShell(shell.NewCommand("add", 2, addCommand))
    if err != nil {
        panic(err)
    }

    err = s.Run()
    if err != nil {
        panic(err)
    }
}
```
