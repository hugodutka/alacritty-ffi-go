.PHONY: all clean build-rust build-go test example

# Default target
all: build-rust build-go

# Build the Rust FFI library
build-rust:
	cd rust-ffi && cargo build --release
	mkdir -p lib
	cp rust-ffi/target/release/libalacritty_ffi.a lib/
	cp rust-ffi/alacritty_ffi.h lib/

# Build the Go bindings
build-go: build-rust
	cd go-bindings && go build

# Build and run the example
example: build-rust
	cd examples && go build && ./example

# Test the implementation
test: build-rust
	cd go-bindings && go test -v

# Clean build artifacts
clean:
	cd rust-ffi && cargo clean
	rm -rf lib/
	cd examples && rm -f example
	cd go-bindings && go clean

# Check Rust code
check-rust:
	cd rust-ffi && cargo check

# Format Rust code
fmt-rust:
	cd rust-ffi && cargo fmt

# Format Go code
fmt-go:
	cd go-bindings && go fmt
	cd examples && go fmt