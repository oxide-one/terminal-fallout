package Terminal

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

type Cursor struct {
	// The current lolumn
	ColumnNumber int
	// The current line
	LineNumber int
	// The current line position
	LinePosition int
	// The current Memory cell
	Cell String
	// The Line
	Line Line
}

// Sets the display value of the current position
func (c *Cursor) display(t Terminal, s tcell.Screen, highlight bool) {
	// If we want to highlight the string
	var style tcell.Style
	if highlight {
		style = t.Style.Highlight
	} else {
		style = t.Style.Default
	}
	// If the cell has been attempted
	var content string
	if c.Cell.Attempted {
		content = strings.Repeat(" ", len(c.Cell.Content))
		if highlight {
			style = t.Style.LowHighlight
		} else {
			style = t.Style.LowDefault
		}
	} else {
		content = c.Cell.Content

	}
	emitStr(s, c.Cell.Position.Start.X, c.Cell.Position.Start.Y, style, content)
}

// Flashes the current position blank, changes position, then flashes bright
func (c *Cursor) setCell(t Terminal, s tcell.Screen) {
	c.display(t, s, false)
	c.Cell = t.MemoryBlocks[c.ColumnNumber].Content[c.LineNumber].Content[c.LinePosition]
	c.display(t, s, true)
}

// Initializes the cursor at zero
func (c *Cursor) Init(t Terminal) {
	c.ColumnNumber = 0
	c.LineNumber = 0
	c.LinePosition = 0
	c.Line = t.MemoryBlocks[c.ColumnNumber].Content[c.LineNumber]
	c.Cell = c.Line.Content[c.LinePosition]
}

// Movement of the Cursor UP
func (c *Cursor) moveUp(t Terminal, s tcell.Screen) {
	if c.LineNumber == 0 {
		c.LineNumber = t.Settings.Lines - 1
	} else {
		c.LineNumber--
	}
	// RightSide bounds checking
	c.boundsCheck(t)
	c.smoothenMovement(t, c.LineNumber)
	c.setCell(t, s)
}

// Movement of the Cursor DOWN
func (c *Cursor) moveDown(t Terminal, s tcell.Screen) {
	if c.LineNumber == t.Settings.Lines-1 {
		c.LineNumber = 0
	} else {
		c.LineNumber++
	}
	c.boundsCheck(t)
	c.smoothenMovement(t, c.LineNumber)
	c.setCell(t, s)
}

// Movement of the Cursor LEFT
func (c *Cursor) moveLeft(t Terminal, s tcell.Screen) {
	// If we are at the very left of this column
	if c.LinePosition == 0 {

		// If we are at the left st column
		if c.ColumnNumber == 0 {
			c.ColumnNumber = t.Settings.Columns - 1
		} else {
			c.ColumnNumber -= 1
		}

		c.LinePosition = getLineLength(t, c.ColumnNumber, c.LineNumber) - 1
	} else {
		c.LinePosition -= 1
	}
	c.setCell(t, s)
}

// Movement of the Cursor Right
func (c *Cursor) moveRight(t Terminal, s tcell.Screen) {
	// If we are at the very Right of this column
	if c.LinePosition == getLineLength(t, c.ColumnNumber, c.LineNumber)-1 {
		// If we are at the Right most column
		if c.ColumnNumber == t.Settings.Columns-1 {
			c.ColumnNumber = 0
		} else {
			c.ColumnNumber += 1
		}

		c.LinePosition = 0
	} else {
		c.LinePosition += 1
	}
	c.setCell(t, s)
}

func (c *Cursor) boundsCheck(t Terminal) {
	// RightSide bounds checking
	nextLineLength := getLineLength(t, c.ColumnNumber, c.LineNumber)
	if nextLineLength <= c.LinePosition {
		c.LinePosition = nextLineLength - 1
	}
}

func (c *Cursor) smoothenMovement(t Terminal, nextLineNumber int) {
	// Smoothen the movement up and down
	nextLine := t.MemoryBlocks[c.ColumnNumber].Content[c.LineNumber]
	for linePosition := 0; linePosition < nextLine.Length; linePosition++ {
		potentialCursor := nextLine.Content[linePosition]
		// Grab the start and end bounds of the cursor
		potentialCursorStart := potentialCursor.Position.Start.X
		potentialCursorEnd := potentialCursor.Position.End.X
		if potentialCursorStart <= c.Cell.Position.Mid.X && potentialCursorEnd >= c.Cell.Position.Mid.X {
			c.LinePosition = linePosition
		}
	}
}

func getLineLength(t Terminal, columnNumber int, lineNumber int) int {
	return t.MemoryBlocks[columnNumber].Content[lineNumber].Length
}
