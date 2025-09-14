package tui

import (
	"os"
	"testing"
	"time"
)

// newTestApp creates a new App instance for testing with pipe input.
func newTestApp(layout layout) *App {
	app := NewApp(layout)
	app.term.cols = 80
	app.term.rows = 24
	return app
}

func TestCounterPanelIntegration(t *testing.T) {
	// Create a simple layout with CounterPanel
	counter := &CounterPanel{
		PanelBase: PanelBase{Title: "Counter", Border: true},
		Count:     0,
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

	// Send '+' to increment counter
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
		if cp.Count != 1 {
			t.Errorf("Expected count 1, got %d", cp.Count)
		}
	} else {
		t.Errorf("Expected CounterPanel, got %T", app.panels[0])
	}
}
