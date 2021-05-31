package clear

func ClearTTY() {
	// Echos a very specific string to clear the terminal
	println("\033[;H\033[2J")
}
