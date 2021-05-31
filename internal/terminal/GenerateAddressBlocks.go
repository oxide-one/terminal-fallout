package Terminal

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

// func GeneratePasswordBlock(lineCount int, lineWidth int, columnCount int, passwordCount int, passwords map[string]passStruct, passwordList []string) [][][]string {
func GenerateAddressBlocks(AddressBlocks []SimpleBlock, settings TerminalSettings) []SimpleBlock {
	// The deviation between the minimum and max is 396, always.
	const standardDeviation int = 396

	// Declare the minimum address values (000-FFF)
	const addrMin int = 0
	const addrMax int = 4095

	// Generate the start and end ranges
	var startAddress int = rand.Intn(addrMax-addrMin) + addrMin
	var endAddress int = startAddress + standardDeviation

	// Pull the settings
	addressWidth := settings.AddressWidth
	lineCount := settings.Lines
	columnCount := settings.Columns

	// Determine the total line count
	totalLineCount := (lineCount * columnCount) - 2

	// Create a checkmap to see if the value already exists
	chkMap := make(map[int]bool)

	// Create a list of addresses
	var addrList []int
	addrList = append(addrList, startAddress)
	// Iterate over the number of addresses min 2, to create an address map
	for line := 0; line < totalLineCount; line++ {
		for {
			randAddress := startAddress + 1 + rand.Intn(standardDeviation)
			if _, ok := chkMap[randAddress]; !ok {
				addrList = append(addrList, randAddress)
				chkMap[randAddress] = true
				break
			}
		}
	}
	addrList = append(addrList, endAddress)
	sort.Ints(addrList)

	// Create a slice of addresses that are formatted correctly
	for column := 0; column < columnCount; column++ {
		for line := 0; line < lineCount; line++ {
			currentLine := addrList[line] // The current line
			// The remaining number of letters to fill
			formattedString := strings.ToUpper(strconv.FormatInt(int64(currentLine), 16)) // Convert to hex
			formattedStringLength := len(formattedString)                                 // Calculate the length of that hex
			lineWidthRemaining := addressWidth - formattedStringLength - 3                // Figure out how many zeros to fill
			zeroFill := strings.Repeat("0", lineWidthRemaining)                           // Create the zero Fill string
			output := fmt.Sprintf("0x%sF%s", zeroFill, formattedString)

			// Figure out how many 0x00... to add to make it to the line width
			AddressBlocks[column].Content[line].Content = output
			AddressBlocks[column].Content[line].Position.Start.X = AddressBlocks[column].Position.Start.X
			AddressBlocks[column].Content[line].Position.Start.Y = AddressBlocks[column].Position.Start.Y + line
			AddressBlocks[column].Content[line].Position.End.X = AddressBlocks[column].Position.End.X
			AddressBlocks[column].Content[line].Position.End.Y = AddressBlocks[column].Content[line].Position.Start.Y

		}
		addrList = addrList[lineCount:]
	}

	return AddressBlocks
}
