package tui

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

type CounterPanel struct {
	PanelBase
	Count *int
}

func (cp *CounterPanel) Update(msg InputMessage) bool {
	switch {
	case msg.IsChar('+'):
		*cp.Count++
		return true
	case msg.IsChar('-'):
		*cp.Count--
		return true
	}
	return false
}

func (cp *CounterPanel) Draw(active bool) string {
	countStr := fmt.Sprintf("Count: %d", *cp.Count)
	instructions := []string{
		"Use + to increment",
		"Use - to decrement",
	}
	lines := []string{countStr}
	lines = append(lines, instructions...)
	return strings.Join(lines, "\n")
}

type dummyItem struct {
	Name  string
	Value int
}

func (d dummyItem) String() string {
	return fmt.Sprintf("%s: %d", d.Name, d.Value)
}

type dummyPanel struct {
	*ListPanel[dummyItem]
	item *dummyItem
}

func (dp *dummyPanel) Update(msg InputMessage) bool {
	redraw := dp.ListPanel.Update(msg)
	if len(dp.Items) > 0 && dp.Selected >= 0 && dp.Selected < len(dp.Items) {
		*dp.item = dp.Items[dp.Selected]
	}
	return redraw
}

func (dp *dummyPanel) Draw(active bool) string {
	return dp.ListPanel.Draw(active)
}

// newTestApp creates a new App instance for testing with pipe input.
func newTestApp(layout Layout) *App {
	app := NewApp(layout, nil)
	app.term.cols = 80
	app.term.rows = 24
	app.noDraw = true // Disable drawing for tests
	return app
}

func TestCounterPanelIntegration(t *testing.T) {
	// Create shared count
	sharedCount := 0
	counter := &CounterPanel{
		PanelBase: PanelBase{Title: "Counter", Border: true},
		Count:     &sharedCount,
	}

	layout := &PanelNode{Panel: counter}

	// Set up pipe for input
	oldStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = r
	defer func() {
		os.Stdin = oldStdin
		if err := r.Close(); err != nil {
			t.Fatal(err)
		}
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	// Create test app
	app := newTestApp(layout)

	// Channel to signal app has stopped
	done := make(chan bool)

	// Start the app in goroutine
	go func() {
		app.Run()
		done <- true
	}()

	if _, err := w.Write([]byte{'+'}); err != nil {
		t.Fatal(err)
	}

	// Wait a bit for processing
	time.Sleep(10 * time.Millisecond)

	// Send 'q' to quit
	if _, err := w.Write([]byte{'q'}); err != nil {
		t.Fatal(err)
	}

	// Wait for app to quit
	select {
	case <-done:
		// App has stopped
	case <-time.After(100 * time.Millisecond):
		t.Fatal("App did not stop within timeout")
	}

	// Check if count increased
	if len(app.panels) != 1 {
		t.Fatalf("Expected 1 panel, got %d", len(app.panels))
	}

	if cp, ok := app.panels[0].(*CounterPanel); ok {
		if *cp.Count != 1 {
			t.Errorf("Expected count 1, got %d", *cp.Count)
		}
	} else {
		t.Errorf("Expected CounterPanel, got %T", app.panels[0])
	}
}

func TestGlobalHandlerIntegration(t *testing.T) {
	// Test that global handler handles Tab and quit correctly
	app := &App{
		panels:    []Panel{&PanelBase{}},
		activeIdx: 0,
		running:   true,
		handler:   &DefaultGlobalHandler{},
	}

	// Test Tab
	handled, redraw := app.handler.UpdateGlobal(app, newKeyMessage(KeyTab, []byte{}))
	if !handled || !redraw {
		t.Errorf("Tab should be handled and redraw")
	}
	if app.activeIdx != 0 { // Only one panel, should stay 0
		t.Errorf("Active index should be 0, got %d", app.activeIdx)
	}

	// Test quit
	handled, redraw = app.handler.UpdateGlobal(app, newCharMessage('q', []byte{}))
	if !handled || redraw {
		t.Errorf("Quit should be handled and not redraw")
	}
	if app.running {
		t.Errorf("App should not be running after quit")
	}
}

func TestAppPublicMethods(t *testing.T) {
	// Create panels
	p1 := &PanelBase{Title: "Panel1"}
	p2 := &PanelBase{Title: "Panel2"}
	app := &App{
		panels:    []Panel{p1, p2},
		activeIdx: 0,
		running:   true,
		handler:   &DefaultGlobalHandler{},
	}

	// Test SwitchPanel
	app.SwitchPanel(1)
	if app.activeIdx != 1 {
		t.Errorf("Expected activeIdx 1, got %d", app.activeIdx)
	}

	app.SwitchPanel(-1)
	if app.activeIdx != 0 {
		t.Errorf("Expected activeIdx 0, got %d", app.activeIdx)
	}

	// Test FocusPanel
	if !app.FocusPanel("Panel2") {
		t.Errorf("FocusPanel should succeed for Panel2")
	}
	if app.activeIdx != 1 {
		t.Errorf("Expected activeIdx 1, got %d", app.activeIdx)
	}

	if !app.FocusPanel("1") { // By index
		t.Errorf("FocusPanel should succeed for index 1")
	}
	if app.activeIdx != 1 {
		t.Errorf("Expected activeIdx 1, got %d", app.activeIdx)
	}

	if app.FocusPanel("NonExistent") {
		t.Errorf("FocusPanel should fail for non-existent panel")
	}

	// Test Stop
	app.Stop()
	if app.running {
		t.Errorf("App should not be running after Stop")
	}
}

func TestMultiplePanelsSharedDataIntegration(t *testing.T) {
	// Create shared count
	sharedCount := 0
	counter1 := &CounterPanel{
		PanelBase: PanelBase{Title: "Counter 1", Border: true},
		Count:     &sharedCount,
	}
	counter2 := &CounterPanel{
		PanelBase: PanelBase{Title: "Counter 2", Border: true},
		Count:     &sharedCount,
	}

	layout := &HorizontalSplit{
		Left:  &PanelNode{Panel: counter1},
		Right: &PanelNode{Panel: counter2},
	}

	// Set up pipe for input
	oldStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = r
	defer func() {
		os.Stdin = oldStdin
		if err := r.Close(); err != nil {
			t.Fatal(err)
		}
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	// Create test app
	app := newTestApp(layout)

	// Start the app in goroutine
	go app.Run()

	// Send '+' to increment shared count
	if _, err := w.Write([]byte{'+'}); err != nil {
		t.Fatal(err)
	}

	// Wait a bit for processing
	time.Sleep(10 * time.Millisecond)

	// Send 'q' to quit
	if _, err := w.Write([]byte{'q'}); err != nil {
		t.Fatal(err)
	}

	// Wait for app to quit
	time.Sleep(10 * time.Millisecond)

	// Check if both panels have the same count
	if len(app.panels) != 2 {
		t.Fatalf("Expected 2 panels, got %d", len(app.panels))
	}

	for i, panel := range app.panels {
		if cp, ok := panel.(*CounterPanel); ok {
			if *cp.Count != 1 {
				t.Errorf("Panel %d: Expected count 1, got %d", i, *cp.Count)
			}
		} else {
			t.Errorf("Panel %d: Expected CounterPanel, got %T", i, panel)
		}
	}
}

func TestSharedPointersIntegration(t *testing.T) {
	// Create mock items
	mockItem1 := dummyItem{
		Name:  "Item A",
		Value: 10,
	}
	mockItem2 := dummyItem{
		Name:  "Item B",
		Value: 20,
	}
	items := []dummyItem{mockItem1, mockItem2}

	// Shared item pointers
	sharedA := &dummyItem{}
	sharedB := &dummyItem{}

	// Create panels
	panelA := &dummyPanel{
		ListPanel: &ListPanel[dummyItem]{
			PanelBase: PanelBase{Title: "Panel A", Border: true},
			Items:     items,
		},
		item: sharedA,
	}

	panelB := &dummyPanel{
		ListPanel: &ListPanel[dummyItem]{
			PanelBase: PanelBase{Title: "Panel B", Border: true},
			Items:     items,
		},
		item: sharedB,
	}

	// Set initial selection to first item
	*sharedA = items[0]
	*sharedB = items[0]

	// Layout with both panels
	layout := &HorizontalSplit{
		Left:  &PanelNode{Panel: panelA},
		Right: &PanelNode{Panel: panelB},
	}

	// Set up pipe for input
	oldStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = r

	// Suppress stdout to avoid printing TUI to screen
	oldStdout := os.Stdout
	_, pw, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = pw

	defer func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
		_ = r.Close()
		_ = w.Close()
		_ = pw.Close()
	}()

	// Create test app
	app := NewApp(layout, nil)

	// Start the app in goroutine
	go app.Run()

	// Simulate navigating to second item (assuming down arrow or 'j' for next)
	// Note: Adjust based on actual key bindings in ListPanel
	_, _ = w.Write([]byte{'j'}) // Assuming 'j' moves down
	time.Sleep(50 * time.Millisecond)

	// Simulate selecting (enter)
	_, _ = w.Write([]byte{'\n'})
	time.Sleep(50 * time.Millisecond)

	// Send 'q' to quit
	_, _ = w.Write([]byte{'q'})
	time.Sleep(50 * time.Millisecond)

	// Check if shared pointers are updated
	if *sharedA != mockItem2 {
		t.Errorf("Expected sharedA to point to second item, got %v", sharedA)
	}
	if *sharedB != mockItem1 { // panelB wasn't selected, should remain first
		t.Errorf("Expected sharedB to remain at first item, got %v", sharedB)
	}
}
