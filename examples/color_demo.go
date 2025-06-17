package main

import (
	"fmt"
	"log"

	alacritty "github.com/example/alacritty-go"
)

func main() {
	// Create a new terminal
	term := alacritty.NewTerminal(80, 10)
	if term == nil {
		log.Fatal("Failed to create terminal")
	}
	defer term.Close()

	fmt.Println("=== Testing ANSI Color Processing ===")

	// Test basic colors
	colorTests := []struct {
		name string
		ansi string
	}{
		{"Red text", "\x1b[31mRed\x1b[0m"},
		{"Green text", "\x1b[32mGreen\x1b[0m"},
		{"Blue text", "\x1b[34mBlue\x1b[0m"},
		{"Bold text", "\x1b[1mBold\x1b[0m"},
		{"Underlined", "\x1b[4mUnder\x1b[0m"},
	}

	for i, test := range colorTests {
		// Write test on separate lines
		testStr := fmt.Sprintf("Line %d: %s\n", i, test.ansi)
		_, err := term.Write([]byte(testStr))
		if err != nil {
			log.Printf("Failed to write %s: %v", test.name, err)
			continue
		}

		// Check the first character of the colored text
		cell, err := term.GetCell(8, uint32(i)) // Position after "Line X: "
		if err != nil {
			log.Printf("Failed to get cell for %s: %v", test.name, err)
			continue
		}

		fmt.Printf("%s: char='%c' fg:(%d,%d,%d) bg:(%d,%d,%d) bold:%v underline:%v\n",
			test.name, cell.Char,
			cell.FgColor.R, cell.FgColor.G, cell.FgColor.B,
			cell.BgColor.R, cell.BgColor.G, cell.BgColor.B,
			cell.Bold, cell.Underline)
	}

	// Test cursor position after all writes
	x, y, err := term.GetCursor()
	if err != nil {
		log.Fatal("Failed to get cursor:", err)
	}
	fmt.Printf("\nFinal cursor position: (%d, %d)\n", x, y)

	// Show full terminal content
	fmt.Println("\n=== Full Terminal Content ===")
	for i := uint32(0); i < 5; i++ {
		line, err := term.GetLine(i)
		if err != nil {
			continue
		}
		fmt.Printf("Line %d: ", i)
		for _, cell := range line {
			if cell.Char == 0 {
				fmt.Print(" ")
			} else {
				fmt.Printf("%c", cell.Char)
			}
		}
		fmt.Println()
	}
}