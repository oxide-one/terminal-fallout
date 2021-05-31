package Terminal

import "github.com/gdamore/tcell/v2"

type Terminal struct {
	// The terminal settings
	Settings TerminalSettings

	// The color styling
	Style TerminalStyle
	// The header block
	Header SimpleBlock

	// The address block
	AddressBlocks []SimpleBlock

	// The memory blocks
	MemoryBlocks []Block

	// A and dict list of Passwords
	Passwords PassStruct
}

func DefaultSettings() TerminalSettings {
	settings := TerminalSettings{
		Lines:                32,
		Columns:              2,
		Passwords:            12,
		GeneralPaddingX:      4,
		HeaderPaddingTopY:    2,
		HeaderPaddingBottomY: 5,

		AddressWidth:    7,
		AddressPaddingX: 2,

		MemoryWidth:    32,
		MemoryPaddingX: 4,

		TotalAttempts: 5,
	}
	return settings
}

func NewTerminal(settings TerminalSettings) Terminal {
	terminal := Terminal{Settings: settings}
	terminal.AddressBlocks = make([]SimpleBlock, settings.Columns)
	terminal.MemoryBlocks = make([]Block, settings.Columns)

	// Construct a list of passwords
	terminal.Passwords = GeneratePasswords(settings.Passwords)

	// Construct the style
	terminal.Style = TerminalStyle{
		Default:      tcell.StyleDefault.Foreground(tcell.ColorGreen.TrueColor()).Background(tcell.ColorBlack.TrueColor()).Bold(true),
		Highlight:    tcell.StyleDefault.Foreground(tcell.ColorWhite.TrueColor()).Background(tcell.ColorGreen.TrueColor()).Bold(true),
		LowDefault:   tcell.StyleDefault.Foreground(tcell.ColorDarkGreen.TrueColor()).Background(tcell.ColorBlack.TrueColor()),
		LowHighlight: tcell.StyleDefault.Foreground(tcell.ColorWhite.TrueColor()).Background(tcell.ColorDarkGreen.TrueColor()),
	}
	// Construct the Header
	headerLines := []string{
		"OKAMIDASH INDUSTRIES (TM) TERMLINK PROTOCOL",
		"ENTER PASSWORD NOW",
		"",
		"",
	}

	terminal.Header.Position.Start.Y = terminal.Settings.HeaderPaddingTopY
	terminal.Header.Position.Start.X = terminal.Settings.GeneralPaddingX
	terminal.Header.Position.End.Y = terminal.Settings.HeaderPaddingTopY + 4
	terminal.Header.Position.End.X = terminal.Settings.GeneralPaddingX + len(headerLines[:1])

	for i, line := range headerLines {
		myLine := SimpleLine{
			Content: line,
			Position: MultiCoordinates{
				Start: SimpleCoordinates{
					X: terminal.Header.Position.Start.X,
					Y: terminal.Header.Position.Start.Y + i,
				},
				End: SimpleCoordinates{
					X: terminal.Header.Position.Start.X + len(line),
					Y: terminal.Header.Position.Start.Y + i,
				},
			},
		}
		terminal.Header.Content = append(terminal.Header.Content, myLine)
	}

	// Construct the Address Block
	for column, AddressBlock := range terminal.AddressBlocks {
		AddressBlock.Content = make([]SimpleLine, settings.Lines)
		// Calculate the vertical height
		AddressBlock.Position.Start.Y = (settings.HeaderPaddingTopY + settings.HeaderPaddingBottomY)
		AddressBlock.Position.End.Y = AddressBlock.Position.Start.Y + settings.Lines

		// Calculate the horizontal width and padding
		AddressBlock.Position.Start.X = settings.GeneralPaddingX + (column * (settings.AddressWidth + settings.AddressPaddingX + settings.MemoryWidth + settings.MemoryPaddingX))
		AddressBlock.Position.End.X = AddressBlock.Position.Start.X + settings.AddressWidth
		// Save back the information
		terminal.AddressBlocks[column] = AddressBlock
	}
	terminal.AddressBlocks = GenerateAddressBlocks(terminal.AddressBlocks, terminal.Settings)

	// Construct the Memory Block
	for column, MemoryBlock := range terminal.MemoryBlocks {
		MemoryBlock.Content = make([]Line, settings.Lines)

		// Calculate the vertical height
		MemoryBlock.Position.Start.Y = (settings.HeaderPaddingTopY + settings.HeaderPaddingBottomY)
		MemoryBlock.Position.End.Y = MemoryBlock.Position.Start.Y + settings.Lines

		// Calculate the horizontal width and padding
		MemoryBlock.Position.Start.X = terminal.AddressBlocks[column].Position.End.X + settings.AddressPaddingX
		MemoryBlock.Position.End.X = MemoryBlock.Position.Start.X + settings.MemoryWidth

		// Iterate over each line to fill the Start and End X and Y pos
		for lineN, Line := range MemoryBlock.Content {
			Line.Position.Start.X = MemoryBlock.Position.Start.X
			Line.Position.End.X = MemoryBlock.Position.End.X
			Line.Position.Start.Y = MemoryBlock.Position.Start.Y + lineN
			Line.Position.End.Y = Line.Position.Start.Y
			MemoryBlock.Content[lineN] = Line
		}
		// Save back the information
		terminal.MemoryBlocks[column] = MemoryBlock
	}
	terminal.MemoryBlocks, terminal.Passwords = GenerateMemoryBlock(terminal, terminal.Passwords)

	return terminal
}
