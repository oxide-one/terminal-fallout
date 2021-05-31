package main

import (
	"fmt"

	term "github.com/oxide-one/terminal-fallout/internal/terminal"
)

func main() {
	terminal := term.NewTerminal(term.DefaultSettings())
	_ = terminal
	fmt.Println(terminal.Passwords)
	term.Display(terminal)
}
