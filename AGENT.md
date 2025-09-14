# AGENT.md

## Project Overview

Bow is a utility for Arcanist, the tool to manage diffs (PRs) for Phabricator. It provides an interactive terminal interface to select diffs, commits, and base commits.

The project uses a TUI (Terminal User Interface) library located in the `tui/` directory for building the text-based interface.

## TUI Library

The TUI library provides components for building terminal applications:
- **Panels**: UI components with `Update()` and `Draw()` methods
- **Layouts**: Arrangements like splits and nodes
- **Input Handling**: Structured message system for keyboard input
- **App Management**: Terminal setup, event loop, and rendering

## Go Best Practices

Follow idiomatic Go code for the latest release (Go 1.25.1). Use proper naming conventions, interfaces for abstraction, and standard error handling patterns.

## Testing Instructions

- Run `go test ./...` to execute all tests
- Build with `go build` to verify compilation
- Never run the executable during testing
- Ensure tests pass and build succeeds before changes

## Development Workflow

- Test and build after any changes
- Follow established patterns for panels and layouts
- Keep code modular and well-documented