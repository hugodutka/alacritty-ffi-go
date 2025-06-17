package main

import (
	"fmt"
	"log"

	alacritty "github.com/example/alacritty-go"
)

func main() {
	// Create a new terminal
	term := alacritty.NewTerminal(80, 24)
	if term == nil {
		log.Fatal("Failed to create terminal")
	}
	defer term.Close()

	// Test basic functionality
	fmt.Println("=== Terminal Created ===")
	cols, rows, err := term.GetSize()
	if err != nil {
		log.Fatal("Failed to get size:", err)
	}
	fmt.Printf("Terminal size: %dx%d\n", cols, rows)

	// Test writing some text
	fmt.Println("\n=== Writing Text ===")
	testText := "Hello, World!\nThis is a test of the Alacritty FFI.\n"
	changedLines, err := term.Write([]byte(testText))
	if err != nil {
		log.Fatal("Failed to write:", err)
	}
	fmt.Printf("Changed lines: %d\n", changedLines)

	// Get cursor position
	x, y, err := term.GetCursor()
	if err != nil {
		log.Fatal("Failed to get cursor:", err)
	}
	fmt.Printf("Cursor position: (%d, %d)\n", x, y)

	// Test getting individual cells
	fmt.Println("\n=== Reading Cells ===")
	for i := uint32(0); i < 13; i++ { // "Hello, World!" length
		cell, err := term.GetCell(i, 0)
		if err != nil {
			log.Printf("Failed to get cell (%d, 0): %v", i, err)
			continue
		}
		fmt.Printf("Cell (%d, 0): '%c' fg:(%d,%d,%d) bg:(%d,%d,%d) bold:%v\n",
			i, cell.Char, cell.FgColor.R, cell.FgColor.G, cell.FgColor.B,
			cell.BgColor.R, cell.BgColor.G, cell.BgColor.B, cell.Bold)
	}

	// Test getting a full line
	fmt.Println("\n=== Reading Line ===")
	line, err := term.GetLine(0)
	if err != nil {
		log.Fatal("Failed to get line:", err)
	}
	fmt.Printf("Line 0 (%d cells): ", len(line))
	for _, cell := range line {
		if cell.Char == 0 || cell.Char == ' ' {
			fmt.Print(" ")
		} else {
			fmt.Printf("%c", cell.Char)
		}
	}
	fmt.Println()

	// Test ANSI sequences
	fmt.Println("\n=== Testing ANSI Sequences ===")
	ansiText := "\x1b[31mRed text\x1b[0m \x1b[1mBold text\x1b[0m \x1b[4mUnderlined\x1b[0m"
	_, err = term.Write([]byte(ansiText))
	if err != nil {
		log.Fatal("Failed to write ANSI:", err)
	}

	// Check if colors were applied
	cell, err := term.GetCell(0, 2) // First character of "Red text"
	if err != nil {
		log.Fatal("Failed to get colored cell:", err)
	}
	fmt.Printf("Colored cell: '%c' fg:(%d,%d,%d) bg:(%d,%d,%d)\n",
		cell.Char, cell.FgColor.R, cell.FgColor.G, cell.FgColor.B,
		cell.BgColor.R, cell.BgColor.G, cell.BgColor.B)

	// Test terminal resize
	fmt.Println("\n=== Testing Resize ===")
	err = term.Resize(100, 30)
	if err != nil {
		log.Fatal("Failed to resize:", err)
	}
	cols, rows, err = term.GetSize()
	if err != nil {
		log.Fatal("Failed to get new size:", err)
	}
	fmt.Printf("New terminal size: %dx%d\n", cols, rows)

	// Display full terminal content
	fmt.Println("\n=== Full Terminal Content ===")
	content := term.String()
	fmt.Printf("Terminal content:\n%s\n", content)

	fmt.Println("\n=== Test Complete ===")
}