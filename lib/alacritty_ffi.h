#ifndef ALACRITTY_FFI_H
#define ALACRITTY_FFI_H

#include <stdint.h>
#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

// C-compatible cell structure
typedef struct {
    uint32_t c;        // Unicode codepoint
    uint8_t fg_r;      // Foreground color RGB
    uint8_t fg_g;
    uint8_t fg_b;
    uint8_t bg_r;      // Background color RGB
    uint8_t bg_g;
    uint8_t bg_b;
    uint16_t flags;    // Cell flags (bold, italic, etc.)
} CCell;

// Opaque terminal handle
typedef struct CTerminal CTerminal;

// Cell flag constants
#define CELL_FLAG_BOLD      (1 << 0)
#define CELL_FLAG_ITALIC    (1 << 1)
#define CELL_FLAG_UNDERLINE (1 << 2)
#define CELL_FLAG_INVERSE   (1 << 3)

// Function declarations
CTerminal* terminal_new(uint32_t cols, uint32_t rows);
void terminal_free(CTerminal* terminal);
int terminal_process_bytes(CTerminal* terminal, const uint8_t* input, size_t input_len);
CCell terminal_get_cell(const CTerminal* terminal, uint32_t x, uint32_t y);
int terminal_get_line(const CTerminal* terminal, uint32_t y, CCell* output_cells, size_t max_cells);
int terminal_resize(CTerminal* terminal, uint32_t cols, uint32_t rows);
int terminal_get_size(const CTerminal* terminal, uint32_t* cols, uint32_t* rows);
int terminal_get_cursor(const CTerminal* terminal, uint32_t* x, uint32_t* y);

#ifdef __cplusplus
}
#endif

#endif // ALACRITTY_FFI_H