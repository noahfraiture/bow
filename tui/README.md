# TUI Library

A simple and flexible Terminal User Interface (TUI) library for Go, enabling developers to build interactive text-based applications with ease. It provides a modular architecture for creating panels, arranging layouts, and handling user input in the terminal.

## Features

- **Modular Panels**: Create reusable UI components that display content and respond to input.
- **Flexible Layouts**: Arrange panels using splits and nodes for complex screen arrangements.
- **Builtin Panels**: Ready-to-use components like lists, text inputs, and info displays.
- **Input Handling**: Comprehensive support for keyboard input, including special keys, characters, and modifiers.
- **Automatic Rendering**: Handles drawing, borders, titles, colors, and cursor management.
- **Terminal Management**: Manages raw mode, resizing, and cleanup automatically.
- **Cross-Platform**: Works on Unix-like systems with terminal support.

## Installation

Add the library to your Go project:

```bash
go get bow/tui
```

Then import it in your code:

```go
import "bow/tui"
```

## Quick Start

Here's a minimal example to get you started:

```go
package main

import (
    "bow/tui"
)

func main() {
    // Create a simple info panel
    info := &tui.InfoPanel{
        Lines: []string{"Hello, World!", "This is a TUI app."},
    }
    info.Title = "Welcome"
    info.Border = true

    // Define a layout with a single panel
    layout := &tui.PanelNode{Panel: info}

    // Create and run the app
    app := tui.NewApp(layout)
    app.Run()
}
```

Run this code, and you'll see a bordered panel with the title "Welcome" displaying the lines. Press `q` or `Ctrl-C` to quit.

## Creating Panels

Panels are the core UI components. To create a custom panel, embed `tui.PanelBase` and implement the `tui.Panel` interface.

### Panel Interface

```go
type Panel interface {
    GetBase() *PanelBase
    Update(msg InputMessage) bool
    Draw(active bool) string
}
```

- `GetBase()`: Returns the embedded `PanelBase` (usually no need to override).
- `Update(msg InputMessage) bool`: Handles input. Return `true` if the panel needs redrawing.
- `Draw(active bool) string`: Returns the panel's content as a string with `\n` for new lines.

### Example Custom Panel

```go
import (
    "bow/tui"
    "strings"
)

type CounterPanel struct {
    tui.PanelBase
    Count int
}

func (cp *CounterPanel) Update(msg tui.InputMessage) bool {
    if msg.IsChar('+') {
        cp.Count++
        return true
    } else if msg.IsChar('-') {
        cp.Count--
        return true
    }
    return false
}

func (cp *CounterPanel) Draw(active bool) string {
    return fmt.Sprintf("Count: %d\nPress + or - to change", cp.Count)
}
```

Usage:

```go
counter := &CounterPanel{}
counter.Title = "Counter"
counter.Border = true
```

### PanelBase Fields

- `Title string`: Sets the panel's title (displayed in the border).
- `Border bool`: Enables/disables the border around the panel.

The library handles positioning (`x`, `y`, `w`, `h`) internally.

## Builtin Panels

The library provides three ready-to-use panels:

### ListPanel[T]

Displays a selectable list of items. Generic type allows any item type.

```go
list := &tui.ListPanel[string]{
    Items:    []string{"Option 1", "Option 2", "Option 3"},
    Selected: 0,
}
list.Title = "Menu"
list.Border = true
```

- **Input**: `j`/`k` or `↓`/`↑` to navigate.
- **Display**: Highlights the selected item (reversed when active, yellow when inactive).

### TextPanel

Provides text input and editing capabilities.

```go
text := &tui.TextPanel{
    Text:   []rune("Initial text"),
    Cursor: 0,
}
text.Title = "Input"
text.Border = true
```

- **Input**: Arrow keys to move cursor, backspace to delete, enter to clear, printable characters to insert.
- **Display**: Shows the current text with cursor positioning.

### InfoPanel

Simple display panel for static information.

```go
info := &tui.InfoPanel{
    Lines: []string{"Line 1", "Line 2", "Line 3"},
}
info.Title = "Info"
info.Border = true
```

- **Input**: No input handling (static display).
- **Display**: Joins lines with newlines.

## Layouts

Layouts define how panels are arranged on the screen. They implement the internal `layout` interface.

### PanelNode

Represents a single panel filling the entire allocated area.

```go
node := &tui.PanelNode{Panel: myPanel}
```

### HorizontalSplit

Divides the area into left and right sections (50/50 split).

```go
split := &tui.HorizontalSplit{
    Left:  &tui.PanelNode{Panel: panel1},
    Right: &tui.PanelNode{Panel: panel2},
}
```

### VerticalSplit

Divides the area into top and bottom sections (50/50 split).

```go
split := &tui.VerticalSplit{
    Top:    &tui.PanelNode{Panel: panel1},
    Bottom: &tui.PanelNode{Panel: panel2},
}
```

Layouts can be nested for complex arrangements:

```go
layout := &tui.VerticalSplit{
    Top: &tui.HorizontalSplit{
        Left:  &tui.PanelNode{Panel: list},
        Right: &tui.PanelNode{Panel: info},
    },
    Bottom: &tui.PanelNode{Panel: text},
}
```

## Running the App

Create an `App` instance with your layout and call `Run()`:

```go
app := tui.NewApp(layout)
app.Run()
```

The `App` automatically:
- Enables raw mode for direct input handling.
- Positions panels based on the layout and terminal size.
- Handles input parsing and dispatching to the active panel.
- Manages rendering, including borders, colors, and cursor.
- Responds to terminal resizing (`SIGWINCH`).
- Cleans up on exit (disables raw mode, shows cursor).

## Key Bindings

Global key bindings (handled by the app):

- `Tab`: Switch to the next panel.
- `q` / `Q` / `Ctrl-C`: Quit the application.

Panel-specific bindings (handled by individual panels):

- `↑` / `↓` / `j` / `k`: Navigate in lists.
- `←` / `→`: Move cursor in text panels.
- `Backspace`: Delete character in text panels.
- `Enter`: Confirm/clear in text panels.
- Printable characters: Insert in text panels.

The status bar at the bottom shows available actions.

## Best Practices

- **Embed PanelBase**: Always embed `tui.PanelBase` in custom panels for positioning and border support.
- **Efficient Updates**: Only return `true` from `Update` when the display actually changes to avoid unnecessary redraws.
- **Handle Input Carefully**: Use `InputMessage` methods like `IsKey()`, `IsChar()`, and `HasModifier()` for robust input handling.
- **Test Panels**: Create simple test apps to verify panel behavior before integrating into larger layouts.
- **Layout Design**: Plan your layout hierarchy to ensure panels have adequate space.
- **Error Handling**: The library handles most terminal errors gracefully, but check for issues in custom panel logic.
- **Performance**: Keep `Draw` methods fast, as they're called on every render cycle.

## API Reference

### Types

- `Panel`: Interface for all panels.
- `PanelBase`: Base struct providing common panel functionality.
- `App`: Main application struct.
- `InputMessage`: Represents user input events.
- `ListPanel[T]`: Generic selectable list panel.
- `TextPanel`: Text input panel.
- `InfoPanel`: Static information display panel.
- `PanelNode`: Layout node for a single panel.
- `HorizontalSplit`: Layout for left-right splits.
- `VerticalSplit`: Layout for top-bottom splits.

### Key Functions

- `NewApp(layout layout) *App`: Creates a new app instance.
- `(*App) Run()`: Starts the main event loop.

For detailed method signatures and more examples, run `go doc bow/tui` or refer to the source code.
