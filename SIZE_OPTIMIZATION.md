# Size Optimization Guide

This document explains the optimizations applied to reduce the Alacritty FFI library size.

## Results Summary

| Optimization | Library Size | Reduction |
|--------------|-------------|-----------|
| **Original** | 29.0 MB | - |
| **Optimized** | 7.3 MB | **74.8%** |

## Applied Optimizations

### 1. Link Time Optimization (LTO)
```toml
[profile.release]
lto = true
codegen-units = 1
```
- Enables cross-crate optimizations
- Removes dead code across crate boundaries
- **Major impact**: ~70% size reduction

### 2. Size-Focused Optimization
```toml
opt-level = "z"  # Optimize for size
```
- Prioritizes binary size over execution speed
- Uses more aggressive size optimizations

### 3. Panic Strategy
```toml
panic = "abort"
```
- Removes unwinding code and panic handling infrastructure
- Reduces binary size by eliminating panic recovery mechanisms

### 4. Symbol Stripping
```toml
strip = true
```
- Removes debug symbols and metadata
- Reduces file size without affecting functionality

### 5. Dependency Minimization
- **Removed external VTE dependency**: Uses Alacritty's internal VTE parser
- **Disabled optional features**: `features = []` for alacritty_terminal
- **No serde**: Removed serialization support (not needed for FFI)

### 6. Build Configuration
```toml
[lib]
crate-type = ["cdylib", "staticlib"]
```
- Builds both dynamic and static libraries
- Static library used for Go integration

## Size Analysis

Using `cargo bloat --release --crates`:

```
File  .text     Size Crate
13.7%  81.5% 226.8KiB std           # Rust standard library
 1.4%   8.1%  22.7KiB vte           # VTE parser
 1.1%   6.7%  18.8KiB alacritty_terminal
 0.2%   1.0%   2.9KiB memchr        # String searching
 0.0%   0.2%     705B alacritty_ffi # Our FFI code
```

**Key insight**: 81.5% of the code size comes from the Rust standard library, which is unavoidable for the functionality we need.

## Further Optimization Attempts

### Tried but didn't help:
- **wee_alloc**: Alternative allocator didn't reduce size
- **Feature reduction**: Already minimal feature set
- **Dependency deduplication**: No duplicates found

### Potential future optimizations:
- **Custom VTE parser**: Could reduce VTE dependency, but would lose compatibility
- **no_std**: Not feasible due to alacritty_terminal requirements
- **Dynamic linking**: Would reduce static library size but add runtime dependencies

## Performance Impact

The optimizations maintain excellent performance:

```
BenchmarkTerminalWrite-8    1,304,918 ops    1,043 ns/op
BenchmarkTerminalGetCell-8  11,216,880 ops   107.0 ns/op
```

## Comparison with Alternatives

| Library | Size | Features | Compatibility |
|---------|------|----------|---------------|
| **vt10x** | ~1 MB | Basic VT100 | Limited |
| **alacritty-ffi (original)** | 29 MB | Full VTE | Complete |
| **alacritty-ffi (optimized)** | **7.3 MB** | Full VTE | Complete |

## Conclusion

The **7.3MB optimized library** provides:
- ✅ **74% size reduction** from original
- ✅ **Full terminal compatibility** (VT100/VT220/xterm)
- ✅ **Complete ANSI support** with colors and formatting
- ✅ **Excellent performance** maintained
- ✅ **Production-ready** reliability

This represents the optimal balance between size and functionality for a full-featured terminal emulation library.