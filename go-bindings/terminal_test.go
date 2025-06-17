package alacritty

import (
	"strings"
	"testing"
)

func TestTerminalCreation(t *testing.T) {
	term := NewTerminal(80, 24)
	if term == nil {
		t.Fatal("Failed to create terminal")
	}
	defer term.Close()

	cols, rows, err := term.GetSize()
	if err != nil {
		t.Fatalf("Failed to get terminal size: %v", err)
	}

	if cols != 80 {
		t.Errorf("Expected 80 columns, got %d", cols)
	}
	if rows != 24 {
		t.Errorf("Expected 24 rows, got %d", rows)
	}
}

func TestBasicTextInput(t *testing.T) {
	term := NewTerminal(80, 24)
	if term == nil {
		t.Fatal("Failed to create terminal")
	}
	defer term.Close()

	// Write simple text
	text := "Hello, World!"
	_, err := term.Write([]byte(text))
	if err != nil {
		t.Fatalf("Failed to write text: %v", err)
	}

	// Check each character
	for i, expectedChar := range text {
		cell, err := term.GetCell(uint32(i), 0)
		if err != nil {
			t.Fatalf("Failed to get cell (%d, 0): %v", i, err)
		}
		if cell.Char != expectedChar {
			t.Errorf("Cell (%d, 0): expected '%c', got '%c'", i, expectedChar, cell.Char)
		}
	}
}

func TestNewlineProcessing(t *testing.T) {
	term := NewTerminal(80, 24)
	if term == nil {
		t.Fatal("Failed to create terminal")
	}
	defer term.Close()

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:  "Single newline",
			input: "Line1\nLine2",
			expected: []string{
				"Line1",
				"     Line2", // Line2 starts at column 5 (after Line1)
			},
		},
		{
			name:  "Multiple newlines",
			input: "A\nB\nC",
			expected: []string{
				"A",
				" B", // B starts at column 1 (after A)
				"  C", // C starts at column 2 (after B)
			},
		},
		{
			name:  "Empty lines",
			input: "First\n\nThird",
			expected: []string{
				"First",
				"", // Empty line
				"     Third", // Third starts at column 5 (after First)
			},
		},
		{
			name:  "Trailing newline",
			input: "Text\n",
			expected: []string{
				"Text",
				"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh terminal for each test
			term := NewTerminal(80, 24)
			defer term.Close()

			_, err := term.Write([]byte(tt.input))
			if err != nil {
				t.Fatalf("Failed to write input: %v", err)
			}

			// Check each expected line
			for lineNum, expectedLine := range tt.expected {
				line, err := term.GetLine(uint32(lineNum))
				if err != nil {
					t.Fatalf("Failed to get line %d: %v", lineNum, err)
				}

				// Extract text from line, trimming trailing spaces
				var actualLine string
				for _, cell := range line {
					if cell.Char == 0 {
						break // Stop at first null character
					}
					actualLine += string(cell.Char)
				}
				// Trim trailing spaces for comparison
				actualLine = strings.TrimRight(actualLine, " ")

				if actualLine != expectedLine {
					t.Errorf("Line %d: expected '%s', got '%s'", lineNum, expectedLine, actualLine)
				}
			}
		})
	}
}

func TestCursorMovement(t *testing.T) {
	term := NewTerminal(80, 24)
	if term == nil {
		t.Fatal("Failed to create terminal")
	}
	defer term.Close()

	tests := []struct {
		name     string
		input    string
		expectedX uint32
		expectedY uint32
	}{
		{
			name:     "Initial position",
			input:    "",
			expectedX: 0,
			expectedY: 0,
		},
		{
			name:     "Simple text",
			input:    "Hello",
			expectedX: 5,
			expectedY: 0,
		},
		{
			name:     "Single newline",
			input:    "Hello\n",
			expectedX: 5,
			expectedY: 1,
		},
		{
			name:     "Text after newline",
			input:    "Hello\nWorld",
			expectedX: 10,
			expectedY: 1,
		},
		{
			name:     "Multiple newlines",
			input:    "A\nB\nC",
			expectedX: 3,
			expectedY: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh terminal for each test
			term := NewTerminal(80, 24)
			defer term.Close()

			if tt.input != "" {
				_, err := term.Write([]byte(tt.input))
				if err != nil {
					t.Fatalf("Failed to write input: %v", err)
				}
			}

			x, y, err := term.GetCursor()
			if err != nil {
				t.Fatalf("Failed to get cursor position: %v", err)
			}

			if x != tt.expectedX {
				t.Errorf("Cursor X: expected %d, got %d", tt.expectedX, x)
			}
			if y != tt.expectedY {
				t.Errorf("Cursor Y: expected %d, got %d", tt.expectedY, y)
			}
		})
	}
}

func TestANSIColorProcessing(t *testing.T) {
	term := NewTerminal(80, 24)
	if term == nil {
		t.Fatal("Failed to create terminal")
	}
	defer term.Close()

	tests := []struct {
		name        string
		input       string
		checkPos    uint32
		expectedFg  RGB
		expectedBg  RGB
		expectedChar rune
	}{
		{
			name:        "Red foreground",
			input:       "\x1b[31mR",
			checkPos:    0,
			expectedFg:  RGB{R: 255, G: 0, B: 0},
			expectedBg:  RGB{R: 0, G: 0, B: 0},
			expectedChar: 'R',
		},
		{
			name:        "Green foreground",
			input:       "\x1b[32mG",
			checkPos:    0,
			expectedFg:  RGB{R: 0, G: 255, B: 0},
			expectedBg:  RGB{R: 0, G: 0, B: 0},
			expectedChar: 'G',
		},
		{
			name:        "Blue foreground",
			input:       "\x1b[34mB",
			checkPos:    0,
			expectedFg:  RGB{R: 0, G: 0, B: 255},
			expectedBg:  RGB{R: 0, G: 0, B: 0},
			expectedChar: 'B',
		},
		{
			name:        "Yellow foreground",
			input:       "\x1b[33mY",
			checkPos:    0,
			expectedFg:  RGB{R: 255, G: 255, B: 0},
			expectedBg:  RGB{R: 0, G: 0, B: 0},
			expectedChar: 'Y',
		},
		{
			name:        "Reset to default",
			input:       "\x1b[31mR\x1b[0mD",
			checkPos:    1,
			expectedFg:  RGB{R: 255, G: 255, B: 255}, // Default white
			expectedBg:  RGB{R: 0, G: 0, B: 0},
			expectedChar: 'D',
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh terminal for each test
			term := NewTerminal(80, 24)
			defer term.Close()

			_, err := term.Write([]byte(tt.input))
			if err != nil {
				t.Fatalf("Failed to write input: %v", err)
			}

			cell, err := term.GetCell(tt.checkPos, 0)
			if err != nil {
				t.Fatalf("Failed to get cell (%d, 0): %v", tt.checkPos, err)
			}

			if cell.Char != tt.expectedChar {
				t.Errorf("Character: expected '%c', got '%c'", tt.expectedChar, cell.Char)
			}

			if cell.FgColor != tt.expectedFg {
				t.Errorf("Foreground color: expected RGB(%d,%d,%d), got RGB(%d,%d,%d)",
					tt.expectedFg.R, tt.expectedFg.G, tt.expectedFg.B,
					cell.FgColor.R, cell.FgColor.G, cell.FgColor.B)
			}

			if cell.BgColor != tt.expectedBg {
				t.Errorf("Background color: expected RGB(%d,%d,%d), got RGB(%d,%d,%d)",
					tt.expectedBg.R, tt.expectedBg.G, tt.expectedBg.B,
					cell.BgColor.R, cell.BgColor.G, cell.BgColor.B)
			}
		})
	}
}

func TestANSITextFormatting(t *testing.T) {
	term := NewTerminal(80, 24)
	if term == nil {
		t.Fatal("Failed to create terminal")
	}
	defer term.Close()

	tests := []struct {
		name           string
		input          string
		checkPos       uint32
		expectedBold   bool
		expectedItalic bool
		expectedUnder  bool
		expectedChar   rune
	}{
		{
			name:         "Bold text",
			input:        "\x1b[1mB",
			checkPos:     0,
			expectedBold: true,
			expectedChar: 'B',
		},
		{
			name:           "Italic text",
			input:          "\x1b[3mI",
			checkPos:       0,
			expectedItalic: true,
			expectedChar:   'I',
		},
		{
			name:          "Underlined text",
			input:         "\x1b[4mU",
			checkPos:      0,
			expectedUnder: true,
			expectedChar:  'U',
		},
		{
			name:         "Bold + Underline",
			input:        "\x1b[1;4mBU",
			checkPos:     0,
			expectedBold: true,
			expectedUnder: true,
			expectedChar: 'B',
		},
		{
			name:         "Reset formatting",
			input:        "\x1b[1mB\x1b[0mN",
			checkPos:     1,
			expectedBold: false,
			expectedChar: 'N',
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh terminal for each test
			term := NewTerminal(80, 24)
			defer term.Close()

			_, err := term.Write([]byte(tt.input))
			if err != nil {
				t.Fatalf("Failed to write input: %v", err)
			}

			cell, err := term.GetCell(tt.checkPos, 0)
			if err != nil {
				t.Fatalf("Failed to get cell (%d, 0): %v", tt.checkPos, err)
			}

			if cell.Char != tt.expectedChar {
				t.Errorf("Character: expected '%c', got '%c'", tt.expectedChar, cell.Char)
			}

			if cell.Bold != tt.expectedBold {
				t.Errorf("Bold: expected %v, got %v", tt.expectedBold, cell.Bold)
			}

			if cell.Italic != tt.expectedItalic {
				t.Errorf("Italic: expected %v, got %v", tt.expectedItalic, cell.Italic)
			}

			if cell.Underline != tt.expectedUnder {
				t.Errorf("Underline: expected %v, got %v", tt.expectedUnder, cell.Underline)
			}
		})
	}
}

func TestTerminalResize(t *testing.T) {
	term := NewTerminal(80, 24)
	if term == nil {
		t.Fatal("Failed to create terminal")
	}
	defer term.Close()

	// Write some text first
	_, err := term.Write([]byte("Hello, World!"))
	if err != nil {
		t.Fatalf("Failed to write text: %v", err)
	}

	// Resize terminal
	err = term.Resize(100, 30)
	if err != nil {
		t.Fatalf("Failed to resize terminal: %v", err)
	}

	// Check new size
	cols, rows, err := term.GetSize()
	if err != nil {
		t.Fatalf("Failed to get terminal size: %v", err)
	}

	if cols != 100 {
		t.Errorf("Expected 100 columns after resize, got %d", cols)
	}
	if rows != 30 {
		t.Errorf("Expected 30 rows after resize, got %d", rows)
	}

	// Verify text is still there
	cell, err := term.GetCell(0, 0)
	if err != nil {
		t.Fatalf("Failed to get cell after resize: %v", err)
	}
	if cell.Char != 'H' {
		t.Errorf("Expected 'H' at (0,0) after resize, got '%c'", cell.Char)
	}
}

func TestComplexANSISequences(t *testing.T) {
	term := NewTerminal(80, 24)
	if term == nil {
		t.Fatal("Failed to create terminal")
	}
	defer term.Close()

	// Test complex sequence with multiple formatting and colors
	input := "Normal \x1b[1;31mBold Red\x1b[0m \x1b[4;32mUnder Green\x1b[0m Normal"
	_, err := term.Write([]byte(input))
	if err != nil {
		t.Fatalf("Failed to write complex ANSI sequence: %v", err)
	}

	tests := []struct {
		pos      uint32
		char     rune
		fg       RGB
		bold     bool
		underline bool
	}{
		{0, 'N', RGB{255, 255, 255}, false, false}, // "Normal "
		{7, 'B', RGB{255, 0, 0}, true, false},      // "Bold Red"
		{16, 'U', RGB{0, 255, 0}, false, true},     // "Under Green"
		{28, 'N', RGB{255, 255, 255}, false, false}, // " Normal"
	}

	for i, tt := range tests {
		cell, err := term.GetCell(tt.pos, 0)
		if err != nil {
			t.Fatalf("Test %d: Failed to get cell (%d, 0): %v", i, tt.pos, err)
		}

		if cell.Char != tt.char {
			t.Errorf("Test %d: Character at pos %d: expected '%c', got '%c'", i, tt.pos, tt.char, cell.Char)
		}

		if cell.FgColor != tt.fg {
			t.Errorf("Test %d: Foreground at pos %d: expected RGB(%d,%d,%d), got RGB(%d,%d,%d)",
				i, tt.pos, tt.fg.R, tt.fg.G, tt.fg.B, cell.FgColor.R, cell.FgColor.G, cell.FgColor.B)
		}

		if cell.Bold != tt.bold {
			t.Errorf("Test %d: Bold at pos %d: expected %v, got %v", i, tt.pos, tt.bold, cell.Bold)
		}

		if cell.Underline != tt.underline {
			t.Errorf("Test %d: Underline at pos %d: expected %v, got %v", i, tt.pos, tt.underline, cell.Underline)
		}
	}
}

func TestCarriageReturn(t *testing.T) {
	term := NewTerminal(80, 24)
	if term == nil {
		t.Fatal("Failed to create terminal")
	}
	defer term.Close()

	// Test carriage return behavior
	_, err := term.Write([]byte("Hello\rWorld"))
	if err != nil {
		t.Fatalf("Failed to write text with carriage return: %v", err)
	}

	// After "Hello\rWorld", cursor should be at position 5 on line 0
	// and the line should show "World" (overwriting "Hello")
	_, y, err := term.GetCursor()
	if err != nil {
		t.Fatalf("Failed to get cursor position: %v", err)
	}

	if y != 0 {
		t.Errorf("Cursor Y: expected 0, got %d", y)
	}

	// Check that "World" overwrote "Hello"
	expectedText := "World"
	for i, expectedChar := range expectedText {
		cell, err := term.GetCell(uint32(i), 0)
		if err != nil {
			t.Fatalf("Failed to get cell (%d, 0): %v", i, err)
		}
		if cell.Char != expectedChar {
			t.Errorf("Cell (%d, 0): expected '%c', got '%c'", i, expectedChar, cell.Char)
		}
	}
}

func TestUnicodeSupport(t *testing.T) {
	term := NewTerminal(80, 24)
	if term == nil {
		t.Fatal("Failed to create terminal")
	}
	defer term.Close()

	// Test various Unicode characters
	unicodeTests := []struct {
		name string
		text string
		char rune
	}{
		{"ASCII", "A", 'A'},
		{"Latin", "é", 'é'},
		{"Greek", "α", 'α'},
		{"Emoji", "🚀", '🚀'},
		{"CJK", "中", '中'},
	}

	for _, tt := range unicodeTests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh terminal for each test
			term := NewTerminal(80, 24)
			defer term.Close()

			_, err := term.Write([]byte(tt.text))
			if err != nil {
				t.Fatalf("Failed to write Unicode text '%s': %v", tt.text, err)
			}

			cell, err := term.GetCell(0, 0)
			if err != nil {
				t.Fatalf("Failed to get cell (0, 0): %v", err)
			}

			if cell.Char != tt.char {
				t.Errorf("Unicode character: expected '%c' (U+%04X), got '%c' (U+%04X)",
					tt.char, tt.char, cell.Char, cell.Char)
			}
		})
	}
}

func BenchmarkTerminalWrite(b *testing.B) {
	term := NewTerminal(80, 24)
	if term == nil {
		b.Fatal("Failed to create terminal")
	}
	defer term.Close()

	text := []byte("Hello, World! This is a benchmark test.\n")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := term.Write(text)
		if err != nil {
			b.Fatalf("Failed to write text: %v", err)
		}
	}
}

func BenchmarkTerminalGetCell(b *testing.B) {
	term := NewTerminal(80, 24)
	if term == nil {
		b.Fatal("Failed to create terminal")
	}
	defer term.Close()

	// Fill terminal with some text
	for i := 0; i < 24; i++ {
		term.Write([]byte("This is line " + string(rune('0'+i)) + " with some text content.\n"))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := term.GetCell(uint32(i%80), uint32(i%24))
		if err != nil {
			b.Fatalf("Failed to get cell: %v", err)
		}
	}
}