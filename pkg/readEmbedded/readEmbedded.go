package readembedded

import (
	"embed"
	"log"
	"strings"
)

func File(fileObj embed.FS, path string) []string {
	byteArr, err := fileObj.ReadFile(path)
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	strArr := strings.Split(string(byteArr), "\n")
	return strArr
}
