package Terminal

import "github.com/gdamore/tcell/v2"

type SimpleCoordinates struct {
	// X position
	X int
	// Y position
	Y int
}

type MultiCoordinates struct {
	//
	Start SimpleCoordinates
	Mid   SimpleCoordinates
	End   SimpleCoordinates
}

type TerminalStyle struct {
	Default      tcell.Style
	Highlight    tcell.Style
	LowDefault   tcell.Style
	LowHighlight tcell.Style
}

// Settings for the terminal window
type TerminalSettings struct {
	// The number of lines to generate PER COLUMN
	Lines int

	// The number of columns to generate
	Columns int

	// The number of passwords to generate
	Passwords int

	// Padding from the left edge
	GeneralPaddingX int
	// The padding between Y0 and the Header
	HeaderPaddingTopY int

	// The padding between the enter password line, and the blocks
	HeaderPaddingBottomY int

	// The width of Memory Lines
	MemoryWidth int

	// The padding between Columns
	MemoryPaddingX int

	// The width of Address lines
	AddressWidth int

	// The padding between Address lines
	AddressPaddingX int

	// Total Attempts
	TotalAttempts int
}

type Block struct {
	// A block contains all the necessary content to render that block
	// The content of that block, consisting of multiple lines
	Content  []Line
	Position MultiCoordinates
}

type SimpleBlock struct {
	// A simpleblock just contains a number of lines with strings
	Content  []SimpleLine
	Position MultiCoordinates
}

type SimpleLine struct {
	Content  string
	Position MultiCoordinates
}

type Line struct {
	// Each line contains multiple sets of strings
	Content  []String
	Position MultiCoordinates
	Length   int
}

type String struct {
	Content    string
	Position   MultiCoordinates
	Length     int
	StringType string
	Attempted  bool
}

type Password struct {
	Content    string
	Correct    bool
	Length     int
	Similarity int

	Line         int
	LinePosition int
	Column       int
}

type PassStruct struct {
	Content         map[string]Password
	CorrectPassword string
	Listing         []string
}

/*X	0
0	[TerminalSettings.HeaderPaddingTop]
1	[Terminal.Header] -> Block Height(5)
2 	[TerminalSettings.HeaderPaddingBottom]
3	[	Block	][ TerminalSettings.AddressPadding 	][ Block 	][	TerminalSettings.AddressPadding	][	Block	][ TerminalSettings.AddressPadding 	][ Block 	][	TerminalSettings.AddressPadding	]
4	[ Address	][									][ Memory	][									][ Address	][									][ Memory	][									]

*/
