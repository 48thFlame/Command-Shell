# Command-shell

A [go](https://go.dev) library for creating simple command-based `tui`'s (terminal user interfaces).
Built on top of this [keyboard package](https://github.com/eiannone/keyboard).

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

import (
    "fmt"
    "strconv"

    "github.com/48thFlame/Command-shell"
)

func addCommand(input *shell.CommandHandlerInput) error {
    args := input.Args
    a, err := strconv.Atoi(args[0])
    if err != nil {
        return err
    }
    b, err := strconv.Atoi(args[1])
    if err != nil {
        return err
    }

    fmt.Fprintf(input.Stdout, "%v", a+b)

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
