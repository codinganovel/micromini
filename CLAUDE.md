# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is **micromini**, a strategically stripped-down version of the micro text editor. The project reduces the original micro editor from 25,000+ lines to ~22,500 lines by removing the plugin system, color scheme complexity, and most syntax definitions while adding AutoCD functionality for enhanced workflow.

**Key Architectural Changes:**
- **Plugin System Removed**: Complete removal of Lua VM integration and plugin architecture
- **Hardcoded Color Scheme**: Single dark theme in `internal/config/colorscheme.go` instead of runtime theme system
- **Reduced Syntax Support**: Only 7 essential languages (Go, JavaScript, Python, HTML, CSS, Markdown, default)
- **AutoCD Integration**: Uses `github.com/codinganovel/autocd-go v0.1.1` for directory inheritance on exit

## Development Commands

### Building
```bash
# Quick build (recommended for development)
make build-quick
# Outputs binary as 'micro'

# Build with debug information
make build-dbg

# Full build with generation step
make build

# Install to GOPATH/bin
make install
```

### Testing
```bash
# Run all tests
make test

# Run specific package tests
go test ./internal/buffer/
go test ./internal/action/

# Run benchmarks (includes comprehensive buffer performance tests)
make bench

# Run single benchmark
go test -bench=BenchmarkEdit1000Lines1Cursor ./internal/buffer/
```

### Runtime Generation
```bash
# Generate embedded runtime files (required after syntax changes)
make generate
```

## Architecture Overview

Micromini follows a layered architecture with clear separation of concerns:

### Core Architecture Layers

**Terminal Layer** (`internal/display/`, `internal/screen/`)
- `bufwindow.go` - Main editor window rendering and cursor line highlighting
- `tabwindow.go` - Tab management and display
- Uses `github.com/micro-editor/tcell/v2` for terminal manipulation

**Action Layer** (`internal/action/`)
- `bufpane.go` - Buffer pane management and editor actions
- `command.go` - Command system (plugin commands stubbed out)
- `tab.go` - Tab operations and management
- Event handling and keybinding system

**Buffer Layer** (`internal/buffer/`)
- `buffer.go` - Core text buffer implementation with multi-cursor support
- `cursor.go` - Cursor management and operations
- Line-based storage with efficient operations for large files
- Comprehensive benchmark suite for performance validation

**Configuration Layer** (`internal/config/`)
- `colorscheme.go` - **Hardcoded dark theme** (82 color definitions)
- `rtfiles.go` - Runtime file embedding and syntax loading
- No dynamic theme loading (simplified from original micro)

**Syntax Layer** (`pkg/highlight/`, `runtime/syntax/`)
- `highlighter.go` - Syntax highlighting engine
- `parser.go` - Syntax definition parsing
- **Only 7 syntax files retained**: go.yaml, javascript.yaml, python3.yaml, html.yaml, css.yaml, markdown.yaml, default.yaml

### Key Implementation Details

**AutoCD Integration**
- Located in `cmd/micro/micro.go` (lines 245-258)
- Uses `--autocd` flag for opt-in behavior
- Implementation: `autocd.ExitWithDirectoryOrFallback(targetDir, fallback)`
- Changes to file's directory on editor exit when flag is present

**Plugin System Removal**
- All plugin-related code replaced with stub functions
- `internal/action/command.go` contains inactive plugin command stubs
- No Lua VM dependency, significantly reduced memory footprint

**Hardcoded Color Scheme**
- `InitColorscheme()` in `internal/config/colorscheme.go` sets up all colors
- Cursor line highlighting uses navy blue background (`tcell.ColorNavy`)
- No runtime theme switching or configuration

## Important Behavioral Notes

**Exit Behavior:**
- Default: Traditional exit (stays in current directory)
- With `--autocd` flag: Changes to file's directory on exit
- AutoCD only activates when editing files with paths

**Syntax Highlighting:**
- Limited to 7 essential languages
- Hardcoded theme with consistent dark appearance
- Syntax files are embedded at build time

**Performance Characteristics:**
- 7ms startup time (significantly faster than original micro)
- Sub-microsecond buffer operations for typical file sizes
- Excellent multi-cursor performance up to 100 cursors

## Build System Notes

**Makefile Variables:**
- `VERSION`, `HASH`, `DATE` embedded in binary via ldflags
- `CGO_ENABLED=0` except on macOS (for dynamic linking requirements)
- Runtime files embedded via `go generate ./runtime`

**Binary Output:**
- Default binary name: `micro`
- Can be built as `micromini` with: `go build -o micromini cmd/micro/*.go`
- Single static binary with no external dependencies

## Testing Architecture

**Buffer Tests** (`internal/buffer/buffer_test.go`):
- Comprehensive benchmarks from 10 lines to 1M lines
- Multi-cursor performance testing (1, 10, 100, 1000 cursors)
- File creation, reading, and editing benchmarks

**Test Execution:**
- Buffer tests include performance validation
- Cross-platform compatibility testing
- Syntax highlighting verification for retained languages

## Key Files for Development

**Entry Points:**
- `cmd/micro/micro.go` - Main application entry point with AutoCD integration
- `cmd/micro/clean.go` - Configuration cleanup utilities

**Core Components:**
- `internal/action/bufpane.go` - Primary editor functionality
- `internal/buffer/buffer.go` - Text buffer implementation
- `internal/display/bufwindow.go` - Main editor rendering
- `internal/config/colorscheme.go` - Hardcoded theme system

**Runtime Assets:**
- `runtime/syntax/*.yaml` - 7 retained syntax definitions
- `runtime/help/*.md` - Editor help documentation
- Embedded via `runtime/runtime.go` generation

This codebase represents a successful strategic simplification that maintains full text editing functionality while achieving significant code reduction and performance improvements.