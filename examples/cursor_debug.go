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

	tests := []string{
		"",
		"Hello",
		"Hello\n",
		"Hello\nWorld",
		"A\nB\nC",
	}

	for i, input := range tests {
		// Create fresh terminal for each test
		term := alacritty.NewTerminal(80, 24)
		defer term.Close()

		if input != "" {
			_, err := term.Write([]byte(input))
			if err != nil {
				log.Printf("Failed to write input: %v", err)
				continue
			}
		}

		x, y, err := term.GetCursor()
		if err != nil {
			log.Printf("Failed to get cursor: %v", err)
			continue
		}

		fmt.Printf("Test %d: input=%q -> cursor=(%d,%d)\n", i, input, x, y)
	}
}