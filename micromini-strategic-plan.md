# Project Plan: micromini - Stripped-Down Micro Editor

## Executive Summary

micromini is a strategic reduction of the micro editor from 25,000 to ~10,000 lines of Go code by removing three major bloat systems: plugin architecture, color schemes, and syntax definitions. This project preserves micro's excellent editing experience while achieving 60% size reduction, faster startup, and lower memory usage. The project includes autocd integration for enhanced workflow and targets completion as a Month 2 learning project.

**Key Objectives:**
- Remove 15,292 lines of bloat while maintaining core functionality
- Achieve sub-100ms startup time and <10MB memory footprint
- Integrate autocd for seamless directory navigation
- Maintain cross-platform compatibility and terminal rendering quality

## Requirements Analysis

### Functional Requirements

**Core Editor Functions (Must Preserve):**
- File opening, editing, saving with full Unicode support
- Multi-cursor editing and advanced selection modes
- Search and replace with regex support
- Undo/redo with unlimited history
- Split panes and tab management
- Mouse support and keyboard shortcuts
- Terminal integration and command execution
- Basic syntax highlighting (hardcoded, no themes)

**autocd Integration:**
- Automatic directory change on file operations
- Directory history tracking and navigation
- Integration with file browser and search

**Performance Requirements:**
- Startup time: <100ms (vs current ~300ms)
- Memory usage: <10MB for typical files
- File loading: Support files up to 100MB efficiently

### Non-Functional Requirements

**Maintainability:**
- Codebase reduced to ~10,000 lines for better comprehension
- Clear separation between core and auxiliary functions
- Minimal external dependencies

**Compatibility:**
- Cross-platform support (Linux, macOS, Windows)
- Terminal compatibility maintained
- File format preservation

## Architecture Design

### System Overview

micromini follows a simplified architecture with four core layers:

```
┌─────────────────────────────────────┐
│           User Interface            │
│     (Terminal Rendering/Input)      │
├─────────────────────────────────────┤
│          Editor Core               │
│   (Buffer, Cursor, Operations)     │
├─────────────────────────────────────┤
│         File System Layer          │
│    (I/O, autocd Integration)       │
├─────────────────────────────────────┤
│        Platform Abstraction        │
│      (Terminal, OS Interface)      │
└─────────────────────────────────────┤
```

### Component Specifications

**Buffer Management (Preserved)**
- Text buffer with efficient line-based storage
- Multi-cursor state management
- Undo/redo stack with memory optimization
- File encoding detection and conversion

**Rendering Engine (Simplified)**
- Direct terminal output without theme system
- Hardcoded color scheme (dark theme only)
- Basic syntax highlighting for common languages
- Status line and command bar rendering

**Input Handler (Preserved)**
- Keyboard event processing and binding
- Mouse event handling
- Command mode processing
- Search interface

**File Operations (Enhanced)**
- Standard file I/O operations
- autocd integration hooks
- Directory navigation tracking
- File type detection for basic highlighting

### Removed Components

**Plugin System (Complete Removal):**
- Lua VM integration (603 lines)
- Plugin loading and management (956 lines)
- Default plugins (887 lines)
- Plugin API and hooks throughout codebase

**Color Schemes (Complete Removal):**
- keep one defult theme, to keep ui rework simple, just always use that one theme
- Theme loading and parsing (832 lines)
- Runtime color switching
- Theme configuration system

**Syntax Definitions (95% Removal):**
- External syntax file loading (6,656 lines)
- Runtime syntax parsing
- Keep only hardcoded highlighting for: Go, JavaScript, Python, HTML, CSS, Markdown


## Technology Stack

**Core Technologies (Unchanged):**
- Go 1.19+ for cross-platform compatibility
- tcell library for terminal manipulation
- Standard library for file operations

**Dependencies Analysis:**
- Remove: Lua VM dependencies, YAML parsing for themes/syntax
- Keep: tcell, clipboard integration, file watching
- Add: Enhanced path manipulation for autocd

**Build System:**
- Maintain existing Go modules approach
- Simplify build flags (remove plugin build tags)
- Static linking for standalone binary

