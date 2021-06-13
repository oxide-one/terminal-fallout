package Terminal

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
	"github.com/oxide-one/terminal-fallout/pkg/clear"
	"github.com/oxide-one/terminal-fallout/pkg/sleeper"
)

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)

		s.Show()
		x += w
	}

}

func refreshSelection(s tcell.Screen, cell String, style tcell.Style) {
	emitStr(s, cell.Position.Start.X, cell.Position.Start.Y, style, cell.Content)
}

func displayHeader(terminal Terminal, s tcell.Screen) {
	//position := terminal.Header.Position
	content := terminal.Header.Content
	//x_start := position.Start.X
	//y_start := position.Start.Y
	for _, line := range content {
		linePosition := line.Position
		lineContent := line.Content
		emitStr(s, linePosition.Start.X, linePosition.Start.Y, terminal.Style.Default, lineContent)
	}
}

func displayBlocks(terminal Terminal, s tcell.Screen) {
	// Display the blocks
	addressBlocks := terminal.AddressBlocks
	memoryBlocks := terminal.MemoryBlocks
	for column := 0; column < terminal.Settings.Columns; column++ {
		addressBlock := addressBlocks[column]
		memoryBlock := memoryBlocks[column]
		_ = memoryBlock
		for line := 0; line < terminal.Settings.Lines; line++ {
			// Pulling address information
			addressLine := addressBlock.Content[line]
			addressLinePosition := addressLine.Position
			// Display the Addressblock
			emitStr(s, addressLinePosition.Start.X, addressLinePosition.Start.Y, terminal.Style.Default, addressLine.Content)
			// Display the Line block
			memoryLine := memoryBlock.Content[line]
			for stringSet := range memoryLine.Content {
				myString := memoryLine.Content[stringSet]
				//myStringPosition := myString.Position
				refreshSelection(s, myString, terminal.Style.Default)
				//emitStr(s, myStringPosition.Start.X, myStringPosition.Start.Y, tcell.StyleDefault, myString.Content)
			}
		}
	}
}

func Display(terminal Terminal) bool {
	// Clear the TTY
	clear.ClearTTY()
	//vaultTec.EarlyBoot()

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(terminal.Style.Default)
	encoding.Register()
	// Init the Cursor
	cursor := Cursor{}
	cursor.Init(terminal)
	// Init the AttemptBox
	attemptBox := AttemptBox{}
	attemptBox.Init(terminal)
	// Init the InputBox
	inputBox := InputBox{}
	inputBox.Init(terminal)
	for {

		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			displayHeader(terminal, s)
			attemptBox.flash(terminal, s)
			displayBlocks(terminal, s)
			cursor.display(terminal, s, true)
		case *tcell.EventKey:

			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				return true
			}
			//refreshSelection(s, terminal, false)
			switch ev.Key() {

			case tcell.KeyUp:
				cursor.moveUp(terminal, s)
			case tcell.KeyDown:
				cursor.moveDown(terminal, s)
			case tcell.KeyRight:
				cursor.moveRight(terminal, s)
			case tcell.KeyLeft:
				cursor.moveLeft(terminal, s)
			case tcell.KeyEnter:
				if cursor.Cell.Attempted {
					// var messages = []string{
					// 	">" + cursor.Cell.Content,
					// 	">Already Tried",
					// }
					// inputBox.addMessage(messages)
				} else if cursor.Cell.StringType == "password" {
					if cursor.Cell.Content == terminal.Passwords.CorrectPassword {
						var messages = []string{
							">" + cursor.Cell.Content,
							">Exact match!",
							">Please wait",
							">while system",
							">is accessed.",
						}
						inputBox.addMessage(messages)
						attemptBox.flash(terminal, s)
						inputBox.flash(terminal, s, cursor.Cell)
						inputBox.pushList(terminal, s)
						sleeper.Sleep(3000)
						s.Fini()
						return true
					} else {
						wrongPassword := cursor.Cell.Content
						wrongPasswordInfo := terminal.Passwords.Content[wrongPassword]
						var messages = []string{
							">" + cursor.Cell.Content,
							">Entry denied.",
							fmt.Sprintf(">Likeness=%d/%d", wrongPasswordInfo.Similarity, wrongPasswordInfo.Length),
						}
						inputBox.addMessage(messages)
						attemptBox.RemainingAttempts -= 1
					}

				} else if cursor.Cell.StringType == "match" {
					var messages []string
					effect := (rand.Intn(2) == 1)
					// Random choice
					if effect {
						messages = []string{
							">" + cursor.Cell.Content,
							">Allowance",
							">Replenished.",
						}
						attemptBox.RemainingAttempts = attemptBox.TotalAttempts
					} else {
						messages = []string{
							">" + cursor.Cell.Content,
							">Dud removed.",
						}
						for _, password := range terminal.Passwords.Listing {
							passwordMap := terminal.Passwords.Content[password]
							passwordCell := &terminal.MemoryBlocks[passwordMap.Column].Content[passwordMap.Line].Content[passwordMap.LinePosition]
							if password != terminal.Passwords.CorrectPassword && !passwordCell.Attempted {
								passwordCell.Attempted = true
								passwordX := passwordCell.Position.Start.X
								passwordY := passwordCell.Position.Start.Y
								passwordLength := passwordCell.Length
								emitStr(s, passwordX, passwordY, terminal.Style.LowDefault, strings.Repeat(" ", passwordLength))
								break
							}
						}
					}

					inputBox.addMessage(messages)
				} else {
					var messages = []string{
						">" + cursor.Cell.Content,
						">Entry denied.",
						fmt.Sprintf(">Likeness=%d/%d", 0, len(terminal.Passwords.CorrectPassword)),
					}
					inputBox.addMessage(messages)
					attemptBox.RemainingAttempts -= 1
				}
				terminal.MemoryBlocks[cursor.ColumnNumber].Content[cursor.LineNumber].Content[cursor.LinePosition].Attempted = true
				cursor.Cell.Attempted = true

				inputBox.pushList(terminal, s)

			}
			attemptBox.flash(terminal, s)
			inputBox.flash(terminal, s, cursor.Cell)
		}
	}

}
