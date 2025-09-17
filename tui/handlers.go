package tui

// GlobalHandler defines the interface for handling global input and panel switches.
// Implementations should embed DefaultGlobalHandler to inherit default behavior.
type GlobalHandler interface {
	UpdateGlobal(app *App, msg InputMessage) (redraw bool)
	OnPanelSwitch(app *App, panelName string)
	GetStatus() string
}

// DefaultGlobalHandler provides default global command handling.
// Embed this in custom handlers to override specific methods.
type DefaultGlobalHandler struct{}

// UpdateGlobal handles default global commands like Tab to switch panels and 'q' to quit.
// Returns true for redraw needed (Tab), false for quit (no redraw) or unhandled.
func (dgh *DefaultGlobalHandler) UpdateGlobal(app *App, msg InputMessage) (redraw bool) {
	switch {
	case msg.HasModifier(ModShift) && msg.IsKey(KeyTab):
		app.SwitchPanel(-1)
		return true
	case msg.IsKey(KeyTab):
		app.SwitchPanel(1)
		return true
	case msg.IsChar('q'), msg.IsChar('Q'), msg.IsChar('\x03'): // 'q' or Ctrl+C
		app.Stop()
		return false // Quit, no redraw
	}
	return false
}

// OnPanelSwitch is a no-op hook for panel switches.
// Override in custom handlers for additional logic.
func (dgh *DefaultGlobalHandler) OnPanelSwitch(app *App, panelName string) {
	// No-op
}

// GetStatus returns the default status line displayed at the bottom.
// Override in custom handlers to provide custom status information.
func (dgh *DefaultGlobalHandler) GetStatus() string {
	return " Tab: switch  •  ↑/↓: navigate  •  ←/→: move cursor  •  Enter: confirm  •  q/Ctrl-C: quit "
}
