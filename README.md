# micromini

**micromini** is a strategically stripped-down version of the [micro text editor](https://github.com/zyedidia/micro), optimized for simplicity and enhanced workflow. It reduces complexity by ~60% while adding AutoCD functionality for seamless directory navigation.

## What is micromini?

micromini takes the core editing functionality of micro and removes the bloat:

- ❌ **No plugin system** - No Lua VM, no plugin management overhead
- ❌ **No color scheme complexity** - Single hardcoded dark theme
- ❌ **Minimal syntax support** - Only 7 essential languages (Go, JavaScript, Python, HTML, CSS, Markdown, default)
- ✅ **AutoCD integration** - Changes to file's directory on editor exit
- ✅ **All core editing features preserved** - Multi-cursor, search/replace, splits, tabs, etc.

## Key Features

- **Fast startup** - ~7ms startup time (vs ~300ms for full micro)
- **Small binary** - Single static binary with no dependencies
- **AutoCD workflow** - Automatically cd to file's directory when exiting editor
- **Essential syntax highlighting** - Go, JavaScript, Python, HTML, CSS, Markdown
- **Full text editing** - Multi-cursor, unlimited undo/redo, regex search/replace
- **Terminal native** - Splits, tabs, mouse support, common keybindings
- **Cross-platform** - Works on macOS, Linux, Windows

## Installation

### Build from source

```bash
git clone https://github.com/yourusername/micromini.git
cd micromini
make build-quick
```

This creates a `micro` binary in the current directory.

### Quick build

```bash
go build cmd/micro/*.go
```

## AutoCD Usage

By default, micromini works like any other editor. To enable AutoCD functionality:

```bash
# Edit a file with AutoCD - will cd to file's directory on exit
./micro --autocd /path/to/file.txt

# Normal usage - stays in current directory
./micro /path/to/file.txt
```

**AutoCD behavior:**
- Only activates when `--autocd` flag is used
- Changes to the file's directory when you exit the editor
- Falls back to normal exit if AutoCD fails
- Non-breaking - doesn't interfere with normal workflows

## Syntax Support

micromini includes syntax highlighting for essential languages:

- **Go** (.go)
- **JavaScript** (.js, .jsx, .ts, .tsx)  
- **Python** (.py)
- **HTML** (.html, .htm)
- **CSS** (.css)
- **Markdown** (.md)
- **Default** (plain text, unknown extensions)

## Building

### Development commands

```bash
# Quick build for development
make build-quick

# Build with debug info
make build-dbg

# Run tests
make test

# Run benchmarks
make bench
```

### Build requirements

- Go 1.18+
- No external dependencies

## Architecture

micromini maintains micro's proven architecture while removing complexity:

```
micromini/
├── cmd/micro/          # Main entry point with AutoCD integration
├── internal/
│   ├── action/         # Keybindings and commands (plugin refs removed)
│   ├── buffer/         # Text buffer management (multi-cursor support)
│   ├── config/         # Configuration (hardcoded colorscheme)
│   ├── display/        # Terminal rendering
│   └── ...
├── runtime/
│   ├── help/          # Built-in help documentation
│   └── syntax/        # 7 essential syntax files only
└── pkg/highlight/     # Syntax highlighting engine
```

## Performance

- **Startup time**: ~7ms (93% faster than 100ms target)
- **Binary size**: ~10MB (single static binary)
- **Memory usage**: Significantly reduced due to no Lua VM
- **File operations**: Sub-microsecond for typical file sizes

## Differences from micro

| Feature | micro | micromini |
|---------|--------|-----------|
| Plugin system | ✅ Lua plugins | ❌ Removed |
| Color schemes | ✅ 25+ themes | ❌ Single dark theme |
| Syntax languages | ✅ 150+ languages | ✅ 7 essential languages |
| AutoCD | ❌ None | ✅ Built-in |
| Startup time | ~300ms | ~7ms |
| Binary size | Larger | ~10MB |

## License

micromini is dual-licensed under the Coffee License:

- **Free** for entities with net worth < $500M USD
- **Commercial license** ($50 fee) for entities ≥ $500M USD

See [LICENSE.md](LICENSE.md) and [LICENSE-COMMERCIAL.md](LICENSE-COMMERCIAL.md) for details.

Original micro editor components retain their original MIT license (see [LICENSE-ORIGINAL](LICENSE-ORIGINAL)).

## Contributing

micromini is designed to stay minimal. Contributions should focus on:

- Bug fixes for core functionality
- Performance improvements
- AutoCD enhancements
- Essential syntax language improvements

Please avoid contributions that add back removed complexity (plugins, themes, extensive syntax support).

## Credits

micromini is based on the excellent [micro editor](https://github.com/zyedidia/micro) by Zachary Yedidia and contributors.

AutoCD functionality provided by [autocd-go](https://github.com/codinganovel/autocd-go).