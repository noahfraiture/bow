# TUI Library

This is a simple Terminal User Interface (TUI) library for Go, designed for building text-based applications with panels, layouts, and input handling.

## Overview

The library provides:
- **Panels**: UI components that can display content and handle input.
- **Layouts**: Structures to arrange panels on the screen (e.g., splits, nodes).
- **App**: The main application runner that manages the terminal, input, and rendering.

## Getting Started

1. Import the library.

2. Create panels by embedding `PanelBase` and implementing the `Panel` interface.

3. Define a layout to arrange your panels.

4. Create an `App` with the layout and run it.

## Creating Custom Panels

To create a custom panel, embed `PanelBase` (or another existing panel) in your struct. This provides common fields like `Title` and `Border`.

```go
type MyPanel struct {
    tui.PanelBase
    // Your custom fields
    Data []string
}

// Implement the Panel interface
func (mp *MyPanel) Update(input byte) bool {
    // Handle input, return true if redraw needed
    return false
}

func (mp *MyPanel) Draw(active bool) string {
    // Return content as string with \n for lines
    content := strings.Join(mp.Data, "\n")
    return mp.WrapWithBorder(content, active)
}
```

- **Embedding**: Use an unnamed field to embed `PanelBase`. This allows access to its fields and methods.
- **Override Methods**: Implement `Update` and `Draw`. Keep `GetBase` as-is (it's for internal use).
- **Borders and Titles**: Use `Base: true` as field to add borders. Set `Title: <title>` as field to set the title.

## Builtin Panels

The library includes ready-to-use panels:
- `ListPanel`: Displays a list with selection.
- `TextPanel`: For text input and editing.
- `InfoPanel`: Simple display of lines.

Example:
```go
list := &tui.ListPanel{
    Items: []string{"Item 1", "Item 2"},
}
list.Title = "My List"
list.Border = true
```

## Layouts

Arrange panels using layouts:
- `PanelNode`: Single panel filling the area.
- `HorizontalSplit`: Left-right split.
- `VerticalSplit`: Top-bottom split.

Example:
```go
layout := &tui.HorizontalSplit{
    Left:  &tui.PanelNode{Panel: list},
    Right: &tui.PanelNode{Panel: textPanel},
}
```

## Running the App

```go
app := tui.NewApp(layout)
app.Run()
```

The app handles:
- Terminal setup (raw mode).
- Input processing (arrows, keys).
- Panel switching (Tab).
- Resizing (SIGWINCH).
- Rendering and cleanup.

## Key Bindings

- `Tab`: Switch active panel.
- Arrow keys: Navigate (handled by panels).
- `q` / `Ctrl-C`: Quit.

## Best Practices

- Always embed `PanelBase` for custom panels to ensure compatibility.
- Test your panels by implementing the interface and running in a layout.
- Handle redraws efficiently in `Update`.

## API Reference

See the Go documentation for detailed method signatures and examples. Run `go doc` on the package for more info.
