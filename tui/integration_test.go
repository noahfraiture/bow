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

// newTestApp creates a new App instance for testing with pipe input.
func newTestApp(layout layout) *App {
	app := NewApp(layout)
	app.term.cols = 80
	app.term.rows = 24
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
		r.Close()
		w.Close()
	}()

	// Create test app
	app := newTestApp(layout)

	// Start the app in goroutine
	go app.Run()

	w.Write([]byte{'+'})

	// Wait a bit for processing
	time.Sleep(10 * time.Millisecond)

	// Send 'q' to quit
	w.Write([]byte{'q'})

	// Wait for app to quit
	time.Sleep(10 * time.Millisecond)

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
		r.Close()
		w.Close()
	}()

	// Create test app
	app := newTestApp(layout)

	// Start the app in goroutine
	go app.Run()

	// Send '+' to increment shared count
	w.Write([]byte{'+'})

	// Wait a bit for processing
	time.Sleep(10 * time.Millisecond)

	// Send 'q' to quit
	w.Write([]byte{'q'})

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
