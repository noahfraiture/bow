package tui

import (
	"strings"
	"testing"
)

func TestPanelBaseDraw(t *testing.T) {
	pb := &PanelBase{
		x: 0, y: 0, w: 10, h: 5, Title: "Test", Border: true,
	}

	// Test with active = true (cyan border)
	result := pb.Draw(true)
	full := pb.WrapWithBorder(result, true)
	lines := strings.Split(full, "\n")
	if len(lines) != 5 {
		t.Errorf("Expected 5 lines, got %d", len(lines))
	}

	// Check top border with title
	var expectedTop = ClrCyan + "┌ [Test] ┐" + Reset
	if lines[0] != expectedTop {
		t.Errorf("Top line mismatch: got %q, want %q", lines[0], expectedTop)
	}

	// Check middle lines (empty)
	for i := 1; i < 4; i++ {
		var expected = ClrCyan + "│" + Reset + ClrWhite + "        " + Reset + ClrCyan + "│" + Reset
		if lines[i] != expected {
			t.Errorf("Line %d mismatch: got %q, want %q", i, lines[i], expected)
		}
	}

	// Check bottom border
	var expectedBottom = ClrCyan + "└────────┘" + Reset
	if lines[4] != expectedBottom {
		t.Errorf("Bottom line mismatch: got %q, want %q", lines[4], expectedBottom)
	}

	// Test with active = false (white border)
	result = pb.Draw(false)
	full = pb.WrapWithBorder(result, false)
	lines = strings.Split(full, "\n")
	expectedTop = ClrWhite + "┌ [Test] ┐" + Reset
	if lines[0] != expectedTop {
		t.Errorf("Inactive top line mismatch: got %q, want %q", lines[0], expectedTop)
	}

	// Test without title
	pb.Title = ""
	result = pb.Draw(true)
	full = pb.WrapWithBorder(result, true)
	lines = strings.Split(full, "\n")
	expectedTop = ClrCyan + "┌────────┐" + Reset
	if lines[0] != expectedTop {
		t.Errorf("No title top line mismatch: got %q, want %q", lines[0], expectedTop)
	}

	// Check middle lines (empty)
	for i := 1; i < 4; i++ {
		expected := ClrCyan + "│" + Reset + ClrWhite + "        " + Reset + ClrCyan + "│" + Reset
		if lines[i] != expected {
			t.Errorf("Line %d mismatch: got %q, want %q", i, lines[i], expected)
		}
	}

	// Check bottom border
	expectedBottom = ClrWhite + "└────────┘" + Reset
	if lines[4] != expectedBottom {
		t.Errorf("Bottom line mismatch: got %q, want %q", lines[4], expectedBottom)
	}

	// Test without title
	pb.Title = ""
	result = pb.Draw(true)
	full = pb.WrapWithBorder(result, true)
	lines = strings.Split(full, "\n")
	expectedTop = ClrCyan + "┌────────┐" + Reset
	if lines[0] != expectedTop {
		t.Errorf("No title top line mismatch: got %q, want %q", lines[0], expectedTop)
	}

	// Test small panel
	pb.w = 2
	pb.h = 2
	result = pb.Draw(true)
	if result != "" {
		t.Errorf("Small panel should return empty, got %q", result)
	}
}

func TestListPanelDraw(t *testing.T) {
	lp := &ListPanel{
		PanelBase: PanelBase{w: 10, h: 5, Title: "List"},
		Items:     []string{"item1", "item2", "item3"},
		Selected:  0,
	}

	res := lp.Draw(true)
	full := lp.WrapWithBorder(res, true)
	lines := strings.Split(full, "\n")
	if len(lines) != 5 {
		t.Errorf("Expected 5 lines, got %d", len(lines))
	}

	// Title is handled by app, not in panel Draw

	// Check first item (selected, active)
	var expectedItem1 = ClrCyan + "│" + Reset + ClrWhite + Reverse + "item1   " + Reset + ClrCyan + "│" + Reset
	if lines[1] != expectedItem1 {
		t.Errorf("Item1 mismatch: got %q, want %q", lines[1], expectedItem1)
	}

	// Check second item
	var expectedItem2 = ClrCyan + "│" + Reset + ClrWhite + "item2   " + Reset + ClrCyan + "│" + Reset
	if lines[2] != expectedItem2 {
		t.Errorf("Item2 mismatch: got %q, want %q", lines[2], expectedItem2)
	}

	// Test truncation
	lp.Items = []string{"very long item name"}
	res = lp.Draw(true)
	full = lp.WrapWithBorder(res, true)
	lines = strings.Split(full, "\n")
	if !strings.Contains(lines[1], "very l..") {
		t.Errorf("Truncation failed: %q", lines[1])
	}

	// Test overflow (more items than height)
	lp.Items = []string{"1", "2", "3", "4", "5"}
	lp.h = 4
	res = lp.Draw(true)
	full = lp.WrapWithBorder(res, true)
	lines = strings.Split(full, "\n")
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(lines))
	}
	// Last content line should be "..."
	if !strings.Contains(lines[2], "...") {
		t.Errorf("Overflow not handled: %q", lines[2])
	}
}

func TestTextPanelDraw(t *testing.T) {
	tp := &TextPanel{
		PanelBase: PanelBase{w: 10, h: 4, Title: "Text"},
		Text:      []rune("hello"),
		Cursor:    2,
	}

	res := tp.Draw(true)
	full := tp.WrapWithBorder(res, true)
	lines := strings.Split(full, "\n")
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(lines))
	}

	// Check text content
	var expectedText = ClrCyan + "│" + Reset + ClrWhite + "hello   " + Reset + ClrCyan + "│" + Reset
	if lines[1] != expectedText {
		t.Errorf("Text line mismatch: got %q, want %q", lines[1], expectedText)
	}

	// Test truncation
	tp.Text = []rune("very long text")
	res = tp.Draw(true)
	full = tp.WrapWithBorder(res, true)
	lines = strings.Split(full, "\n")
	if !strings.Contains(lines[1], "very l..") {
		t.Errorf("Text truncation failed: %q", lines[1])
	}
}

func TestInfoPanelDraw(t *testing.T) {
	ip := &InfoPanel{
		PanelBase: PanelBase{w: 10, h: 5, Title: "Info"},
		Lines:     []string{"line1", "line2", "line3"},
	}

	res := ip.Draw(true)
	full := ip.WrapWithBorder(res, true)
	lines := strings.Split(full, "\n")
	if len(lines) != 5 {
		t.Errorf("Expected 5 lines, got %d", len(lines))
	}

	// Check lines
	var expectedLine1 = ClrCyan + "│" + Reset + ClrWhite + "line1   " + Reset + ClrCyan + "│" + Reset
	if lines[1] != expectedLine1 {
		t.Errorf("Line1 mismatch: got %q, want %q", lines[1], expectedLine1)
	}

	// Test overflow
	ip.Lines = []string{"1", "2", "3", "4", "5"}
	res = ip.Draw(true)
	full = ip.WrapWithBorder(res, true)
	lines = strings.Split(full, "\n")
	if !strings.Contains(lines[3], "...") {
		t.Errorf("Overflow not handled: %q", lines[3])
	}
}
