package Terminal

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type AttemptBox struct {
	// The current lolumn
	TotalAttempts int
	// The current line
	RemainingAttempts int
	// The current line position
	Position MultiCoordinates
}

func (a *AttemptBox) Init(t Terminal) {
	a.TotalAttempts = t.Settings.TotalAttempts
	a.RemainingAttempts = t.Settings.TotalAttempts
	a.Position = t.Header.Content[len(t.Header.Content)-1].Position
}

func (a *AttemptBox) flash(t Terminal, s tcell.Screen) {
	boxes := strings.Repeat("▪ ", a.RemainingAttempts)
	emptyBoxes := strings.Repeat("  ", a.TotalAttempts-a.RemainingAttempts)
	fullString := fmt.Sprintf("%d ATTEMPTS REMAINING: %s%s", a.RemainingAttempts, boxes, emptyBoxes)
	emitStr(s, a.Position.Start.X, a.Position.Start.Y, t.Style.Default, fullString)
}
