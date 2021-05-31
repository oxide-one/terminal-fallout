package Terminal

import (
	"fmt"
	"math/rand"
)

func generateBracketPairs(lineCount int, columnCount int) map[string]bool {
	matchCount := rand.Intn(lineCount-10) + 10

	bracketPairMap := make(map[string]bool)
	// Iterate through the number of passwords needed
	for i := 0; i < matchCount; i++ {
		for {
			// Generate a random number between 0 and the lineCount
			lineNo := rand.Intn(lineCount)

			// Generate a random number between 0 and the columns
			colNo := rand.Intn(columnCount)

			// Generate a 'hash' of sorts to ensure no column and line collissions
			colLineHash := fmt.Sprintf("%d%d", colNo, lineNo)

			if _, ok := bracketPairMap[colLineHash]; !ok {
				// Set the hash value
				bracketPairMap[colLineHash] = true
				break
			}
		}
	}
	// Sort the matches
	// sort.Ints(bracketPairLines)
	return bracketPairMap
}

func generatePasswordMatches(lineCount int, columnCount int, passwordCount int, matchMap map[string]bool) map[string]bool {
	// Checkmap so we don't get conflicts
	passMap := make(map[string]bool)

	// Iterate through the number of passwords needed
	for i := 0; i < passwordCount; i++ {
		for {
			// Generate a random number between 0 and the lineCount
			lineNo := rand.Intn(lineCount)

			// Generate a random number between 0 and the columns
			colNo := rand.Intn(columnCount)

			// Generate a 'hash' of sorts to ensure no column and line collissions
			colLineHash := fmt.Sprintf("%d%d", colNo, lineNo)

			// Evaluate the hashes
			_, passMapChk := passMap[colLineHash]
			_, matchMapChk := matchMap[colLineHash]

			// If There are no collissions, add it in.
			if !matchMapChk && !passMapChk {
				// Append the line number and column number

				// Set the hash value
				passMap[colLineHash] = true
				break
			}
		}
	}
	// Sort the arrays
	// sort.Ints(passwordLines)
	return passMap
}

func GenerateMemoryBlock(terminal Terminal, passStruct PassStruct) ([]Block, PassStruct) {
	passwordList := passStruct.Listing
	// List of allowed punctuation
	var punctuationList = []string{",", ",", ".", "/", "?", "@", "'", ":", ";", "~", "#", "{", "[", "+", "=", "-", "_", "(", "*", "&", "^", "%", "$", "\"", "!", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	var punctuationListLen int = len(punctuationList)

	// List of matched pairs
	var bracketPairs = []string{"{}", "[]", "()", "<>"}
	var bracketPairsLen int = len(bracketPairs)

	// Pull some needed variables
	lineCount := terminal.Settings.Lines
	lineWidth := terminal.Settings.MemoryWidth
	columnCount := terminal.Settings.Columns
	passwordCount := terminal.Settings.Passwords

	// Generate the Bracket Matches
	bracketPairMap := generateBracketPairs(lineCount, columnCount)
	passwordMap := generatePasswordMatches(lineCount, columnCount, passwordCount, bracketPairMap)

	// Columns
	memoryBlocks := terminal.MemoryBlocks

	for column := 0; column < columnCount; column++ {

		// One Column
		memoryBlock := memoryBlocks[column]
		// Iterate over each line
		for line := 0; line < lineCount; line++ {
			// The line in the current memory block
			lineInMemoryBlock := memoryBlock.Content[line]

			// Generate a 'hash' of sorts to ensure no column and line collissions
			colLineHash := fmt.Sprintf("%d%d", column, line)

			// Evaluate the hashes
			_, passwordMapCheck := passwordMap[colLineHash]
			_, bracketPairCheck := bracketPairMap[colLineHash]

			lineCharsLeft := lineWidth
			var currentBuffer []String

			// If the current line and column is where a password should be
			if passwordMapCheck {
				myPassword := String{
					Content:    passwordList[0],
					StringType: "password",
				}
				// Append the first password to the buffer
				currentBuffer = append(currentBuffer, myPassword)
				// Remove the length of that password from the remaining characters left
				lineCharsLeft -= len(passwordList[0])
				// Pop the first password from the passwordList
				passwordList = passwordList[1:]

			}

			// If the current line and column is where a bracket pair should be
			if bracketPairCheck {
				// Grab a random index of the bracketPairs
				matchIndex := rand.Intn(bracketPairsLen)
				matchIndexStart := string(bracketPairs[matchIndex][0])
				matchIndexEnd := string(bracketPairs[matchIndex][1])
				// Generate a random rune count
				runeCount := rand.Intn(4)
				inbetween := ""
				for i := 0; i < runeCount; i++ {
					inbetween += string(punctuationList[rand.Intn(len(punctuationList))])
				}
				var fullMatch string = matchIndexStart + inbetween + matchIndexEnd
				// Append to the current Buffer
				myMatch := String{
					Content:    fullMatch,
					StringType: "match",
				}
				currentBuffer = append(currentBuffer, myMatch)
				lineCharsLeft -= len(fullMatch)
			}

			// Fill the remaining matches with random punctuation
			for fillChars := lineCharsLeft; fillChars > 0; fillChars-- {
				puncIndex := rand.Intn(punctuationListLen)
				myPunc := String{
					Content:    punctuationList[puncIndex],
					StringType: "punctuation",
				}
				currentBuffer = append(currentBuffer, myPunc)
			}

			// Finally, Randomize the list and populate the StringSet
			currentBufferLength := len(currentBuffer)
			currentBufferPerms := rand.Perm(currentBufferLength)

			// Set the array length of the line in memory
			lineInMemoryBlock.Content = make([]String, currentBufferLength)

			// Set the first element to be the StartX position
			lineInMemoryBlock.Content[0].Position.Start = lineInMemoryBlock.Position.Start

			for i, v := range currentBufferPerms {
				currentStringInLine := lineInMemoryBlock.Content[i]

				currentString := currentBuffer[v]
				currentStringLength := len(currentString.Content)

				currentStringInLine.Content = currentString.Content
				currentStringInLine.Length = currentStringLength
				currentStringInLine.StringType = currentString.StringType

				// Set line, position and line position
				if currentString.StringType == "password" {
					currentPass := passStruct.Content[currentString.Content]
					currentPass.Column = column
					currentPass.Line = line
					currentPass.LinePosition = i
					passStruct.Content[currentPass.Content] = currentPass
				}

				if i == 0 {
					currentStringInLine.Position.Start = lineInMemoryBlock.Position.Start
				} else {
					lastStringInLine := lineInMemoryBlock.Content[i-1]
					currentStringInLine.Position.Start.X = lastStringInLine.Position.End.X + 1
					currentStringInLine.Position.Start.Y = lastStringInLine.Position.End.Y
				}
				if currentStringInLine.Length == 1 {
					currentStringInLine.Position.End = currentStringInLine.Position.Start
					currentStringInLine.Position.Mid = currentStringInLine.Position.Start
				} else {
					currentStringInLine.Position.End.X = currentStringInLine.Position.Start.X + currentStringLength - 1
					currentStringInLine.Position.Mid.X = currentStringInLine.Position.Start.X + (currentStringLength / 2) - 1
					currentStringInLine.Position.End.Y = lineInMemoryBlock.Position.Start.Y
					currentStringInLine.Position.Start.Y = lineInMemoryBlock.Position.Start.Y
					currentStringInLine.Position.Mid.Y = lineInMemoryBlock.Position.Start.Y
				}

				// Save back
				lineInMemoryBlock.Content[i] = currentStringInLine
				lineInMemoryBlock.Length = len(lineInMemoryBlock.Content)
			}

			// Writeback the content of the current line
			memoryBlock.Content[line] = lineInMemoryBlock

		}
		// Writeback the content of the current block
		memoryBlocks[column] = memoryBlock
	}

	//fmt.Println(passwordBlockLines)

	//fmt.Println(passwordLines, passwordCols, matchLines, matchCols)
	//fmt.Printf("%+v", memoryBlocks[0].Content[1].Content)
	return memoryBlocks, passStruct
}
