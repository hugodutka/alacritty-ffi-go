package alacritty

/*
#cgo CFLAGS: -I../lib
#cgo LDFLAGS: -L../lib -lalacritty_ffi -ldl -lm
#include "alacritty_ffi.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

// Cell represents a terminal cell with character and formatting
type Cell struct {
	Char      rune
	FgColor   RGB
	BgColor   RGB
	Bold      bool
	Italic    bool
	Underline bool
	Inverse   bool
}

// RGB represents an RGB color
type RGB struct {
	R, G, B uint8
}

// Terminal represents a terminal emulator instance
type Terminal struct {
	ptr *C.CTerminal
}

// NewTerminal creates a new terminal with the specified dimensions
func NewTerminal(cols, rows uint32) *Terminal {
	ptr := C.terminal_new(C.uint32_t(cols), C.uint32_t(rows))
	if ptr == nil {
		return nil
	}
	
	term := &Terminal{ptr: ptr}
	runtime.SetFinalizer(term, (*Terminal).Close)
	return term
}

// Close frees the terminal resources
func (t *Terminal) Close() {
	if t.ptr != nil {
		C.terminal_free(t.ptr)
		t.ptr = nil
		runtime.SetFinalizer(t, nil)
	}
}

// Write processes input bytes and returns the number of changed lines
func (t *Terminal) Write(data []byte) (int, error) {
	if t.ptr == nil {
		return 0, fmt.Errorf("terminal is closed")
	}
	
	if len(data) == 0 {
		return 0, nil
	}
	
	result := C.terminal_process_bytes(
		t.ptr,
		(*C.uint8_t)(unsafe.Pointer(&data[0])),
		C.size_t(len(data)),
	)
	
	if result < 0 {
		return 0, fmt.Errorf("failed to process bytes")
	}
	
	return int(result), nil
}

// GetCell returns the cell at the specified position
func (t *Terminal) GetCell(x, y uint32) (Cell, error) {
	if t.ptr == nil {
		return Cell{}, fmt.Errorf("terminal is closed")
	}
	
	cCell := C.terminal_get_cell(t.ptr, C.uint32_t(x), C.uint32_t(y))
	
	return Cell{
		Char:      rune(cCell.c),
		FgColor:   RGB{R: uint8(cCell.fg_r), G: uint8(cCell.fg_g), B: uint8(cCell.fg_b)},
		BgColor:   RGB{R: uint8(cCell.bg_r), G: uint8(cCell.bg_g), B: uint8(cCell.bg_b)},
		Bold:      (cCell.flags & C.CELL_FLAG_BOLD) != 0,
		Italic:    (cCell.flags & C.CELL_FLAG_ITALIC) != 0,
		Underline: (cCell.flags & C.CELL_FLAG_UNDERLINE) != 0,
		Inverse:   (cCell.flags & C.CELL_FLAG_INVERSE) != 0,
	}, nil
}

// GetLine returns all cells for a specific line
func (t *Terminal) GetLine(y uint32) ([]Cell, error) {
	if t.ptr == nil {
		return nil, fmt.Errorf("terminal is closed")
	}
	
	cols, _, err := t.GetSize()
	if err != nil {
		return nil, err
	}
	
	// Allocate C array for cells
	cCells := make([]C.CCell, cols)
	
	result := C.terminal_get_line(
		t.ptr,
		C.uint32_t(y),
		(*C.CCell)(unsafe.Pointer(&cCells[0])),
		C.size_t(cols),
	)
	
	if result < 0 {
		return nil, fmt.Errorf("failed to get line")
	}
	
	// Convert C cells to Go cells
	cells := make([]Cell, result)
	for i := 0; i < int(result); i++ {
		cCell := cCells[i]
		cells[i] = Cell{
			Char:      rune(cCell.c),
			FgColor:   RGB{R: uint8(cCell.fg_r), G: uint8(cCell.fg_g), B: uint8(cCell.fg_b)},
			BgColor:   RGB{R: uint8(cCell.bg_r), G: uint8(cCell.bg_g), B: uint8(cCell.bg_b)},
			Bold:      (cCell.flags & C.CELL_FLAG_BOLD) != 0,
			Italic:    (cCell.flags & C.CELL_FLAG_ITALIC) != 0,
			Underline: (cCell.flags & C.CELL_FLAG_UNDERLINE) != 0,
			Inverse:   (cCell.flags & C.CELL_FLAG_INVERSE) != 0,
		}
	}
	
	return cells, nil
}

// Resize changes the terminal size
func (t *Terminal) Resize(cols, rows uint32) error {
	if t.ptr == nil {
		return fmt.Errorf("terminal is closed")
	}
	
	result := C.terminal_resize(t.ptr, C.uint32_t(cols), C.uint32_t(rows))
	if result != 0 {
		return fmt.Errorf("failed to resize terminal")
	}
	
	return nil
}

// GetSize returns the current terminal size
func (t *Terminal) GetSize() (cols, rows uint32, err error) {
	if t.ptr == nil {
		return 0, 0, fmt.Errorf("terminal is closed")
	}
	
	var cCols, cRows C.uint32_t
	result := C.terminal_get_size(t.ptr, &cCols, &cRows)
	if result != 0 {
		return 0, 0, fmt.Errorf("failed to get terminal size")
	}
	
	return uint32(cCols), uint32(cRows), nil
}

// GetCursor returns the current cursor position
func (t *Terminal) GetCursor() (x, y uint32, err error) {
	if t.ptr == nil {
		return 0, 0, fmt.Errorf("terminal is closed")
	}
	
	var cX, cY C.uint32_t
	result := C.terminal_get_cursor(t.ptr, &cX, &cY)
	if result != 0 {
		return 0, 0, fmt.Errorf("failed to get cursor position")
	}
	
	return uint32(cX), uint32(cY), nil
}

// String returns a string representation of the terminal content
func (t *Terminal) String() string {
	_, rows, err := t.GetSize()
	if err != nil {
		return ""
	}
	
	var result string
	for y := uint32(0); y < rows; y++ {
		line, err := t.GetLine(y)
		if err != nil {
			continue
		}
		
		for _, cell := range line {
			if cell.Char == 0 {
				result += " "
			} else {
				result += string(cell.Char)
			}
		}
		if y < rows-1 {
			result += "\n"
		}
	}
	
	return result
}