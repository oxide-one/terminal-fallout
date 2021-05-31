package Terminal

import (
	"embed"
	"math/rand"
	"strings"
	"time"

	readembedded "github.com/oxide-one/terminal-fallout/pkg/readEmbedded"
)

//go:embed wordlist
var wordList embed.FS

func calculateSimilarity(chosenPassword string, selectedPassword string) int {
	// Calculate the similarity of words by the number of shared letters in them
	// Generate a checkMap to append each rune to
	matchingLetters := 0
	for i, choChar := range chosenPassword {
		selChar := rune(selectedPassword[i])
		if choChar == selChar {
			matchingLetters++
		}
	}
	// chkMap := make(map[rune]bool)
	// for _, char := range chosenPassword {
	// 	chkMap[char] = true
	// }
	// var matchingLetters int
	// for _, char := range selectedPassword {
	// 	if _, ok := chkMap[char]; ok {
	// 		matchingLetters++
	// 	}
	// }
	return matchingLetters
}

func GeneratePasswords(passwordCount int) PassStruct {
	rand.Seed(time.Now().UnixNano())

	// Grab the embedded Wordlist
	wordList := readembedded.File(wordList, "wordlist")
	// Calculate the length of the wordlist
	wordListLen := len(wordList)

	// Make the passwordList and the checkmap
	passStruct := PassStruct{
		Content: make(map[string]Password),
		Listing: make([]string, 0),
	}
	chkMap := make(map[string]bool)

	// Iterate through the count of passwords and
	for i := 0; i < passwordCount; i++ {
		for {
			selectedPassword := strings.ToUpper(wordList[rand.Intn(wordListLen)])
			if _, ok := chkMap[selectedPassword]; !ok {
				passStruct.Listing = append(passStruct.Listing, selectedPassword)
				chkMap[selectedPassword] = true
				break
			}
		}
	}
	chosenPassword := passStruct.Listing[rand.Intn(passwordCount-1)]

	for _, selectedPassword := range passStruct.Listing {
		var correctPassword bool
		if selectedPassword == chosenPassword {
			correctPassword = true
		} else {
			correctPassword = false
		}

		passStruct.Content[selectedPassword] = Password{
			Content:    selectedPassword,
			Correct:    correctPassword,
			Length:     len(selectedPassword),
			Similarity: calculateSimilarity(chosenPassword, selectedPassword),
		}
	}
	passStruct.CorrectPassword = chosenPassword
	return passStruct
}
