package main

import (
	"fmt"
	"os"
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

var s string

func main() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		err := keyboard.Close()
		if err != nil {
			fmt.Println("!!!ERROR:", err)
		}
	}()

	go func() {
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			switch key {
			case keyboard.KeyEsc, keyboard.KeyCtrlC:
				fmt.Println()
				os.Exit(0)
			case keyboard.KeyBackspace, keyboard.KeyBackspace2:
				if sLen := len(s); sLen >= 1 {
					s = s[:sLen-1]
				}
			case keyboard.KeySpace:
				s += " "
			// case keyboard.KeyDelete, keyboard.KeyEnd, keyboard.KeyEnter, keyboard.KeyHome
			default:
				s += string(char)
			}
		}
	}()
	fmt.Println("Press ESC to quit")
	ticker := time.Tick(100 * time.Millisecond)
	for {
		fmt.Print("\033[2K\r")
		fmt.Print(s)

		<-ticker
	}
}
