# AGENT.md

## Project Overview

Bow is an interactive Terminal User Interface (TUI) utility for managing Arcanist diffs in Phabricator. It simplifies selecting diffs, commits, and base commits for creating or updating diffs. The application fetches commit history from the local Git repository and lists existing diffs via `arc list`, presenting them in a text-based interface.

Written in Go, it includes a custom TUI library in `tui/` for building terminal applications. Integrates with Git using `go-git` and relies on Arcanist for Phabricator operations.

## High-Level Architecture

- **Main Application**: Entry point in `main.go` sets up panels for commit selection, diff management, and commands. Uses the TUI library for UI rendering and input handling.
- **TUI Library (`tui/`)**: Custom framework for terminal UIs. Core components include app management, panel interfaces, layout systems, and low-level terminal operations. Supports generics for reusable panels like lists and text inputs.
- **Integration Modules**: Separate modules for Git commits (`commit.go`), Arcanist diffs (`arc.go`), and command execution (`command.go`). Handle external tool interactions and data parsing.

## Key Concepts

- **Panel Pattern**: UI components implement the `Panel` interface with `Update()` for input handling and `Draw()` for rendering. Embed `PanelBase` for shared properties.
- **Layout System**: Dynamic arrangement of panels using splits (horizontal/vertical) and nodes.
- **Input Handling**: Structured messages for keyboard input, enabling key detection and character processing.
- **Modularity**: Code organized into focused modules; avoid global state by passing dependencies explicitly.
- **Error Handling**: Standard Go practices with error wrapping and fatal logging in main.

## Coding Conventions

- **Idiomatic Go**: Follow Go 1.25.1 best practices, including naming conventions, interfaces, and error handling with `fmt.Errorf`.
- **Generics**: Used for reusable components (e.g., `ListPanel[T]`).
- **Colorized Output**: ANSI escape codes for terminal formatting.
- **No Comments**: Rely on self-documenting code; avoid unnecessary comments.
- **Testing**: Unit tests with mock data for external dependencies.
- **Minimal public interface**: keep number of public element minimal when working on a library.

## Tools and Dependencies

- **Go 1.25.1**: Required for modern features like generics.
- **go-git**: For Git repository operations.
- **Arcanist**: External tool for Phabricator diff management.
- **golangci-lint**: For code quality checks.
- **Custom TUI Library**: Internal UI framework.

## Instructions for AI Agent

- **Development Workflow**: Test and build after changes. Run `go test ./...`, `go build`, and `golangci-lint run ./...`. Ensure tests pass before completion.
- **Panel Implementation**: Follow `Panel` interface; embed `PanelBase`; use generics for lists.
- **Layout Usage**: Arrange panels with splits and nodes.
- **Integration**: Use `go-git` for commits, parse `arc list` output for diffs.
- **Input/Rendering**: Handle via `InputMessage`; redraw when `Update()` returns true; use ANSI codes.
- **Error Handling**: Wrap errors with `%w`; fatal logs in main.
- **Code Style**: Match existing patterns; modular, no comments.
- **Testing**: Add tests with mocks.
- **Security**: Avoid secrets.
- **Proactivity**: Act only on explicit requests.
- **File Editing**: Read before editing; prefer modifications over new files.
- **Tool Usage**: Use search tools extensively; batch calls; run lint after changes.

This provides essential context for coding assistance while maintaining project consistency.
