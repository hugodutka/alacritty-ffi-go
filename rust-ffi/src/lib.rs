use std::os::raw::{c_int, c_uint};
use std::slice;

use alacritty_terminal::{Term, event::VoidListener, grid::Dimensions};
use alacritty_terminal::term::{Config, cell::{Cell, Flags}};
use alacritty_terminal::vte::ansi::{Color, NamedColor, Rgb, Handler};
use alacritty_terminal::index::{Point, Line, Column};

/// C-compatible cell structure
#[repr(C)]
#[derive(Debug, Clone, Copy)]
pub struct CCell {
    pub c: u32,           // Unicode codepoint
    pub fg_r: u8,         // Foreground color RGB
    pub fg_g: u8,
    pub fg_b: u8,
    pub bg_r: u8,         // Background color RGB  
    pub bg_g: u8,
    pub bg_b: u8,
    pub flags: u16,       // Cell flags (bold, italic, etc.)
}

impl Default for CCell {
    fn default() -> Self {
        CCell {
            c: b' ' as u32,
            fg_r: 255, fg_g: 255, fg_b: 255,  // White foreground
            bg_r: 0, bg_g: 0, bg_b: 0,        // Black background
            flags: 0,
        }
    }
}

/// C-compatible terminal size structure
#[repr(C)]
#[derive(Debug, Clone, Copy)]
pub struct CTermSize {
    pub columns: c_uint,
    pub screen_lines: c_uint,
}

impl Dimensions for CTermSize {
    fn columns(&self) -> usize {
        self.columns as usize
    }

    fn screen_lines(&self) -> usize {
        self.screen_lines as usize
    }

    fn total_lines(&self) -> usize {
        self.screen_lines as usize
    }
}

/// Opaque terminal handle
pub struct CTerminal {
    term: Term<VoidListener>,
    size: CTermSize,
}

/// Convert Alacritty Cell to CCell
fn cell_to_ccell(cell: &Cell) -> CCell {
    let c = cell.c as u32;
    
    // Extract colors
    let (fg_r, fg_g, fg_b) = match cell.fg {
        Color::Named(NamedColor::Foreground) => (255, 255, 255),
        Color::Named(NamedColor::Background) => (0, 0, 0),
        Color::Named(NamedColor::Black) => (0, 0, 0),
        Color::Named(NamedColor::Red) => (255, 0, 0),
        Color::Named(NamedColor::Green) => (0, 255, 0),
        Color::Named(NamedColor::Yellow) => (255, 255, 0),
        Color::Named(NamedColor::Blue) => (0, 0, 255),
        Color::Named(NamedColor::Magenta) => (255, 0, 255),
        Color::Named(NamedColor::Cyan) => (0, 255, 255),
        Color::Named(NamedColor::White) => (255, 255, 255),
        Color::Named(NamedColor::BrightBlack) => (128, 128, 128),
        Color::Named(NamedColor::BrightRed) => (255, 128, 128),
        Color::Named(NamedColor::BrightGreen) => (128, 255, 128),
        Color::Named(NamedColor::BrightYellow) => (255, 255, 128),
        Color::Named(NamedColor::BrightBlue) => (128, 128, 255),
        Color::Named(NamedColor::BrightMagenta) => (255, 128, 255),
        Color::Named(NamedColor::BrightCyan) => (128, 255, 255),
        Color::Named(NamedColor::BrightWhite) => (255, 255, 255),
        Color::Spec(Rgb { r, g, b }) => (r, g, b),
        Color::Indexed(idx) => {
            // Simple mapping for indexed colors
            match idx {
                0..=7 => {
                    let colors = [(0,0,0), (255,0,0), (0,255,0), (255,255,0), 
                                 (0,0,255), (255,0,255), (0,255,255), (255,255,255)];
                    colors[idx as usize]
                }
                8..=15 => {
                    let colors = [(128,128,128), (255,128,128), (128,255,128), (255,255,128),
                                 (128,128,255), (255,128,255), (128,255,255), (255,255,255)];
                    colors[(idx - 8) as usize]
                }
                _ => (128, 128, 128), // Default gray for other indexed colors
            }
        }
        _ => (255, 255, 255), // Default white for unhandled colors
    };

    let (bg_r, bg_g, bg_b) = match cell.bg {
        Color::Named(NamedColor::Background) => (0, 0, 0),
        Color::Named(NamedColor::Foreground) => (255, 255, 255),
        Color::Named(NamedColor::Black) => (0, 0, 0),
        Color::Named(NamedColor::Red) => (255, 0, 0),
        Color::Named(NamedColor::Green) => (0, 255, 0),
        Color::Named(NamedColor::Yellow) => (255, 255, 0),
        Color::Named(NamedColor::Blue) => (0, 0, 255),
        Color::Named(NamedColor::Magenta) => (255, 0, 255),
        Color::Named(NamedColor::Cyan) => (0, 255, 255),
        Color::Named(NamedColor::White) => (255, 255, 255),
        Color::Named(NamedColor::BrightBlack) => (128, 128, 128),
        Color::Named(NamedColor::BrightRed) => (255, 128, 128),
        Color::Named(NamedColor::BrightGreen) => (128, 255, 128),
        Color::Named(NamedColor::BrightYellow) => (255, 255, 128),
        Color::Named(NamedColor::BrightBlue) => (128, 128, 255),
        Color::Named(NamedColor::BrightMagenta) => (255, 128, 255),
        Color::Named(NamedColor::BrightCyan) => (128, 255, 255),
        Color::Named(NamedColor::BrightWhite) => (255, 255, 255),
        Color::Spec(Rgb { r, g, b }) => (r, g, b),
        Color::Indexed(idx) => {
            match idx {
                0..=7 => {
                    let colors = [(0,0,0), (255,0,0), (0,255,0), (255,255,0), 
                                 (0,0,255), (255,0,255), (0,255,255), (255,255,255)];
                    colors[idx as usize]
                }
                8..=15 => {
                    let colors = [(128,128,128), (255,128,128), (128,255,128), (255,255,128),
                                 (128,128,255), (255,128,255), (128,255,255), (255,255,255)];
                    colors[(idx - 8) as usize]
                }
                _ => (0, 0, 0), // Default black for other indexed colors
            }
        }
        _ => (0, 0, 0), // Default black for unhandled colors
    };

    // Convert flags
    let mut flags = 0u16;
    if cell.flags.contains(Flags::BOLD) {
        flags |= 1;
    }
    if cell.flags.contains(Flags::ITALIC) {
        flags |= 2;
    }
    if cell.flags.contains(Flags::UNDERLINE) {
        flags |= 4;
    }
    if cell.flags.contains(Flags::INVERSE) {
        flags |= 8;
    }

    CCell {
        c,
        fg_r, fg_g, fg_b,
        bg_r, bg_g, bg_b,
        flags,
    }
}

/// Create a new terminal instance
#[no_mangle]
pub extern "C" fn terminal_new(cols: c_uint, rows: c_uint) -> *mut CTerminal {
    let size = CTermSize {
        columns: cols,
        screen_lines: rows,
    };
    
    let config = Config::default();
    let term = Term::new(config, &size, VoidListener);
    
    let terminal = Box::new(CTerminal { term, size });
    Box::into_raw(terminal)
}

/// Free a terminal instance
#[no_mangle]
pub extern "C" fn terminal_free(terminal: *mut CTerminal) {
    if !terminal.is_null() {
        unsafe {
            let _ = Box::from_raw(terminal);
        }
    }
}

/// Process input bytes and return number of changed lines
#[no_mangle]
pub extern "C" fn terminal_process_bytes(
    terminal: *mut CTerminal,
    input: *const u8,
    input_len: usize,
) -> c_int {
    if terminal.is_null() || input.is_null() {
        return -1;
    }

    unsafe {
        let terminal = &mut *terminal;
        let input_slice = slice::from_raw_parts(input, input_len);
        
        // Process the input through VTE parser
        for &byte in input_slice {
            terminal.term.input(byte as char);
        }
        
        // For simplicity, assume all lines might have changed
        // In a real implementation, you'd track damage more precisely
        terminal.size.screen_lines as c_int
    }
}

/// Get a cell at the specified position
#[no_mangle]
pub extern "C" fn terminal_get_cell(
    terminal: *const CTerminal,
    x: c_uint,
    y: c_uint,
) -> CCell {
    if terminal.is_null() {
        return CCell::default();
    }

    unsafe {
        let terminal = &*terminal;
        
        if x >= terminal.size.columns || y >= terminal.size.screen_lines {
            return CCell::default();
        }

        let point = Point::new(Line(y as i32), Column(x as usize));
        let cell = &terminal.term.grid()[point];
        cell_to_ccell(cell)
    }
}

/// Get all cells for a specific line
#[no_mangle]
pub extern "C" fn terminal_get_line(
    terminal: *const CTerminal,
    y: c_uint,
    output_cells: *mut CCell,
    max_cells: usize,
) -> c_int {
    if terminal.is_null() || output_cells.is_null() {
        return -1;
    }

    unsafe {
        let terminal = &*terminal;
        
        if y >= terminal.size.screen_lines {
            return -1;
        }

        let cols = std::cmp::min(terminal.size.columns as usize, max_cells);
        let output_slice = slice::from_raw_parts_mut(output_cells, cols);
        
        for x in 0..cols {
            let point = Point::new(Line(y as i32), Column(x as usize));
            let cell = &terminal.term.grid()[point];
            output_slice[x] = cell_to_ccell(cell);
        }
        
        cols as c_int
    }
}

/// Resize the terminal
#[no_mangle]
pub extern "C" fn terminal_resize(
    terminal: *mut CTerminal,
    cols: c_uint,
    rows: c_uint,
) -> c_int {
    if terminal.is_null() {
        return -1;
    }

    unsafe {
        let terminal = &mut *terminal;
        terminal.size.columns = cols;
        terminal.size.screen_lines = rows;
        terminal.term.resize(terminal.size);
        0
    }
}

/// Get terminal size
#[no_mangle]
pub extern "C" fn terminal_get_size(
    terminal: *const CTerminal,
    cols: *mut c_uint,
    rows: *mut c_uint,
) -> c_int {
    if terminal.is_null() || cols.is_null() || rows.is_null() {
        return -1;
    }

    unsafe {
        let terminal = &*terminal;
        *cols = terminal.size.columns;
        *rows = terminal.size.screen_lines;
        0
    }
}

/// Get cursor position
#[no_mangle]
pub extern "C" fn terminal_get_cursor(
    terminal: *const CTerminal,
    x: *mut c_uint,
    y: *mut c_uint,
) -> c_int {
    if terminal.is_null() || x.is_null() || y.is_null() {
        return -1;
    }

    unsafe {
        let terminal = &*terminal;
        let cursor = terminal.term.grid().cursor.point;
        *x = cursor.column.0 as c_uint;
        *y = cursor.line.0 as c_uint;
        0
    }
}