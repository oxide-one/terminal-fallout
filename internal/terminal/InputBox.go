package Terminal

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type InputBox struct {
	// The current line position
	Position SimpleCoordinates
	// The last line length
	LastLength     int
	Messages       [][]string
	LongestMessage int
}

func (i *InputBox) Init(t Terminal) {
	i.Position = t.MemoryBlocks[len(t.MemoryBlocks)-1].Position.End
	i.Position.X += 2
	i.Position.Y -= 1
	i.LastLength = 0
	i.LongestMessage = 2
}

func (i *InputBox) flash(t Terminal, s tcell.Screen, cell String) {
	// Wipe the current line
	emptyBoxes := strings.Repeat(" ", i.LastLength)
	emitStr(s, i.Position.X, i.Position.Y, t.Style.Default, emptyBoxes)

	// Compile the current string
	var fullString string
	// if cell.Attempted {
	// 	fullString = fmt.Sprintf("> %s", "Already Attempted.")
	// } else {
	fullString = fmt.Sprintf("> %s", cell.Content)
	//}
	emitStr(s, i.Position.X, i.Position.Y, t.Style.Default, fullString)
	i.LastLength = len(fullString)
}

func (i *InputBox) pushList(t Terminal, s tcell.Screen) {
	// Print the new buffer
	Y := i.Position.Y - 1
	for messagePos := len(i.Messages) - 1; messagePos >= 0; messagePos-- {
		MessageSet := i.Messages[messagePos]
		Y -= len(MessageSet)
		for _, Message := range MessageSet {
			emitStr(s, i.Position.X, Y, t.Style.Default, strings.Repeat(" ", 15))
			emitStr(s, i.Position.X, Y, t.Style.Default, Message)
			Y++
		}
		Y -= len(MessageSet)
	}
}

func (i *InputBox) addMessage(message []string) {
	for _, msg := range message {
		if len(msg) > i.LongestMessage {
			i.LongestMessage = len(msg)
		}
	}
	i.Messages = append(i.Messages, message)
}
