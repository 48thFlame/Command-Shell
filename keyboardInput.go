package shell

import (
	"github.com/eiannone/keyboard"
)

// listenToKeyboardInput listens for keyboard input and
// calls the corresponding function (ie del - `s.del`, space - `s.space` ec)
func (s *Shell) listenToKeyboardInput(keysChannel <-chan keyboard.KeyEvent) {
	for {
		event := <-keysChannel
		key, char, err := event.Key, event.Rune, event.Err
		// fmt.Printf("--- %v, %v, %v ---\n", char, key, err)
		if err != nil {
			s.errorHandle(err)
		}

		switch key {
		case keyboard.KeyCtrlC, keyboard.KeyCtrlZ:
			s.exit()
		case keyboard.KeyEnter:
			s.enter()
		case keyboard.KeyArrowUp:
			s.upArrow()
		case keyboard.KeyArrowDown:
			s.downArrow()
		case keyboard.KeyBackspace, keyboard.KeyBackspace2:
			s.del()
		case keyboard.KeySpace:
			s.space()
		case keyboard.Key(0): // if is not a special key, its a letter or a symbol
			s.typing(char)
			// case keyboard.KeyDelete, keyboard.KeyEnd, keyboard.KeyHome, keyboard.KeyArrowLeft, keyboard.KeyArrowRight, keyboard.KeyArrowDown
		}
	}
}

const defaultHistoryIndex = -1

// getInput returns the input that should be displayed
func (s *Shell) getInput() string {
	if s.historyIndex == defaultHistoryIndex {
		return s.currentInput
	} else {
		historyItem := s.history[s.historyIndex]
		s.currentInput = historyItem
		return historyItem
	}
}

// replaceInput replaces the input that is displayed
func (s *Shell) replaceInput(newInput string) {
	if s.historyIndex == defaultHistoryIndex {
		s.currentInput = newInput
	} else {
		s.currentInput = newInput
		s.historyIndex = defaultHistoryIndex
	}
}

// updateInput updates the input py adding string to end of input
func (s *Shell) updateInput(str string) {
	input := s.getInput()
	newInput := input + str

	s.replaceInput(newInput)
}

func (s *Shell) exit() {
	s.quitChannel <- true
}

func (s *Shell) enter() {
	err := s.runCommand()
	if err != nil {
		s.errorHandle(err)
	}
}

func (s *Shell) del() {
	i := s.getInput()
	inputLen := len(i)

	if inputLen >= 1 {
		s.replaceInput(i[:inputLen-1])
	}
}

func (s *Shell) space() {
	s.updateInput(" ")
}

func (s *Shell) typing(char rune) {
	s.updateInput(string(char))
}

func (s *Shell) upArrow() {
	s.historyIndex += 1
	if s.historyIndex >= len(s.history) {
		s.historyIndex = defaultHistoryIndex
		s.replaceInput("")
	}
}

func (s *Shell) downArrow() {
	s.historyIndex -= 1
	if s.historyIndex < defaultHistoryIndex {
		s.historyIndex = len(s.history) - 1
	}
}
