package main

import (
	"fmt"
	"log"

	alacritty "github.com/example/alacritty-go"
)

func main() {
	term := alacritty.NewTerminal(80, 24)
	if term == nil {
		log.Fatal("Failed to create terminal")
	}
	defer term.Close()

	// Test "Hello\nWorld"
	_, err := term.Write([]byte("Hello\nWorld"))
	if err != nil {
		log.Fatal("Failed to write:", err)
	}

	x, y, err := term.GetCursor()
	if err != nil {
		log.Fatal("Failed to get cursor:", err)
	}
	fmt.Printf("Cursor: (%d, %d)\n", x, y)

	// Check line contents
	for i := uint32(0); i < 3; i++ {
		line, err := term.GetLine(i)
		if err != nil {
			continue
		}
		
		var content string
		for j, cell := range line {
			if j > 20 { // Only show first 20 chars
				break
			}
			if cell.Char == 0 || cell.Char == ' ' {
				content += "."
			} else {
				content += string(cell.Char)
			}
		}
		fmt.Printf("Line %d: '%s'\n", i, content)
	}
}