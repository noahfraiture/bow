package tui

import (
	"strings"
	"testing"
)

func TestPanelBaseDraw(t *testing.T) {
	t.Run("with_border", func(t *testing.T) {
		pb := &PanelBase{
			x: 0, y: 0, w: 10, h: 5, Title: "Test", Border: true,
		}

		// Test with active = true (cyan border)
		result := pb.Draw(true)
		full := pb.wrapWithBorder(result, true)
		lines := strings.Split(full, "\n")
		if len(lines) != 5 {
			t.Errorf("Expected 5 lines, got %d", len(lines))
		}

		// Check top border with title
		var expectedTop = clrCyan + "┌ [Test] ┐" + reset
		if lines[0] != expectedTop {
			t.Errorf("Top line mismatch: got %q, want %q", lines[0], expectedTop)
		}

		// Check middle lines (empty)
		for i := 1; i < 4; i++ {
			var expected = clrCyan + "│" + reset + clrWhite + "        " + reset + clrCyan + "│" + reset
			if lines[i] != expected {
				t.Errorf("Line %d mismatch: got %q, want %q", i, lines[i], expected)
			}
		}

		// Check bottom border
		var expectedBottom = clrCyan + "└────────┘" + reset
		if lines[4] != expectedBottom {
			t.Errorf("Bottom line mismatch: got %q, want %q", lines[4], expectedBottom)
		}

		// Test with active = false (white border)
		result = pb.Draw(false)
		full = pb.wrapWithBorder(result, false)
		lines = strings.Split(full, "\n")
		expectedTop = clrWhite + "┌ [Test] ┐" + reset
		if lines[0] != expectedTop {
			t.Errorf("Inactive top line mismatch: got %q, want %q", lines[0], expectedTop)
		}

		// Test without title
		pb.Title = ""
		result = pb.Draw(true)
		full = pb.wrapWithBorder(result, true)
		lines = strings.Split(full, "\n")
		expectedTop = clrCyan + "┌────────┐" + reset
		if lines[0] != expectedTop {
			t.Errorf("No title top line mismatch: got %q, want %q", lines[0], expectedTop)
		}

		// Check middle lines (empty)
		for i := 1; i < 4; i++ {
			expected := clrCyan + "│" + reset + clrWhite + "        " + reset + clrCyan + "│" + reset
			if lines[i] != expected {
				t.Errorf("Line %d mismatch: got %q, want %q", i, lines[i], expected)
			}
		}

		// Check bottom border
		expectedBottom = clrCyan + "└────────┘" + reset
		if lines[4] != expectedBottom {
			t.Errorf("Bottom line mismatch: got %q, want %q", lines[4], expectedBottom)
		}

		// Test small panel
		pb.w = 2
		pb.h = 2
		result = pb.Draw(true)
		if result != "" {
			t.Errorf("Small panel should return empty, got %q", result)
		}
	})

	t.Run("without_border", func(t *testing.T) {
		pb := &PanelBase{
			x: 0, y: 0, w: 10, h: 5, Title: "Test", Border: false,
		}

		// Test with active = true
		result := pb.Draw(true)
		full := pb.wrapWithBorder(result, true)
		lines := strings.Split(full, "\n")
		if len(lines) != 5 {
			t.Errorf("Expected 5 lines, got %d", len(lines))
		}

		// Check top with title
		expectedTop := "Test      "
		if lines[0] != expectedTop {
			t.Errorf("Top line mismatch: got %q, want %q", lines[0], expectedTop)
		}

		// Check middle lines (empty)
		for i := 1; i < 5; i++ {
			expected := "          "
			if lines[i] != expected {
				t.Errorf("Line %d mismatch: got %q, want %q", i, lines[i], expected)
			}
		}

		// Test with active = false
		result = pb.Draw(false)
		full = pb.wrapWithBorder(result, false)
		lines = strings.Split(full, "\n")
		expectedTop = "Test      "
		if lines[0] != expectedTop {
			t.Errorf("Inactive top line mismatch: got %q, want %q", lines[0], expectedTop)
		}

		// Test without title
		pb.Title = ""
		result = pb.Draw(true)
		full = pb.wrapWithBorder(result, true)
		lines = strings.Split(full, "\n")
		expectedTop = "          "
		if lines[0] != expectedTop {
			t.Errorf("No title top line mismatch: got %q, want %q", lines[0], expectedTop)
		}

		// Check all lines (empty)
		for i := range 5 {
			expected := "          "
			if lines[i] != expected {
				t.Errorf("Line %d mismatch: got %q, want %q", i, lines[i], expected)
			}
		}

		// Test small panel
		pb.w = 2
		pb.h = 2
		result = pb.Draw(true)
		full = pb.wrapWithBorder(result, true)
		if full != "" {
			t.Errorf("Small panel should return empty, got %q", full)
		}
	})
}

func TestListPanelDraw(t *testing.T) {
	t.Run("with_border", func(t *testing.T) {
		lp := &ListPanel[string]{
			PanelBase: PanelBase{w: 10, h: 5, Title: "List"},
			Items:     []string{"item1", "item2", "item3"},
			Selected:  0,
		}
		lp.Border = true

		res := lp.Draw(true)
		full := (&lp.PanelBase).wrapWithBorder(res, true)
		lines := strings.Split(full, "\n")
		if len(lines) != 5 {
			t.Errorf("Expected 5 lines, got %d", len(lines))
		}

		// Title is handled by app, not in panel Draw

		// Check first item (selected, active)
		var expectedItem1 = clrCyan + "│" + reset + clrWhite + Reverse + "item1" + reset + "   " + reset + clrCyan + "│" + reset
		if lines[1] != expectedItem1 {
			t.Errorf("Item1 mismatch: got %q, want %q", lines[1], expectedItem1)
		}

		// Check second item
		var expectedItem2 = clrCyan + "│" + reset + clrWhite + "item2" + reset + "   " + reset + clrCyan + "│" + reset
		if lines[2] != expectedItem2 {
			t.Errorf("Item2 mismatch: got %q, want %q", lines[2], expectedItem2)
		}

		// Test truncation
		lp.Items = []string{"very long item name"}
		res = lp.Draw(true)
		full = lp.wrapWithBorder(res, true)
		lines = strings.Split(full, "\n")
		if !strings.Contains(lines[1], "very l..") {
			t.Errorf("Truncation failed: %q", lines[1])
		}

		// Test overflow (more items than height)
		lp.Items = []string{"1", "2", "3", "4", "5"}
		lp.h = 4
		res = lp.Draw(true)
		full = lp.wrapWithBorder(res, true)
		lines = strings.Split(full, "\n")
		if len(lines) != 4 {
			t.Errorf("Expected 4 lines, got %d", len(lines))
		}
		// Last content line should be "..."
		if !strings.Contains(lines[2], "...") {
			t.Errorf("Overflow not handled: %q", lines[2])
		}
	})

	t.Run("without_border", func(t *testing.T) {
		lp := &ListPanel[string]{
			PanelBase: PanelBase{w: 10, h: 5, Title: "List"},
			Items:     []string{"item1", "item2", "item3"},
			Selected:  0,
		}
		lp.Border = false

		res := lp.Draw(true)
		full := lp.wrapWithBorder(res, true)
		lines := strings.Split(full, "\n")
		if len(lines) != 5 {
			t.Errorf("Expected 5 lines, got %d", len(lines))
		}

		// Check first item (selected, active)
		expectedItem1 := Reverse + "item1" + reset + "     "
		if lines[1] != expectedItem1 {
			t.Errorf("Item1 mismatch: got %q, want %q", lines[1], expectedItem1)
		}

		// Check second item
		expectedItem2 := "item2" + reset + "     "
		if lines[2] != expectedItem2 {
			t.Errorf("Item2 mismatch: got %q, want %q", lines[2], expectedItem2)
		}

		// Test truncation
		lp.Items = []string{"very long item name"}
		res = lp.Draw(true)
		full = lp.wrapWithBorder(res, true)
		lines = strings.Split(full, "\n")
		if !strings.Contains(lines[1], "very lon..") {
			t.Errorf("Truncation failed: %q", lines[1])
		}

		// Test overflow (more items than height)
		lp.Items = []string{"1", "2", "3", "4", "5"}
		lp.h = 4
		res = lp.Draw(true)
		full = lp.wrapWithBorder(res, true)
		lines = strings.Split(full, "\n")
		if len(lines) != 4 {
			t.Errorf("Expected 4 lines, got %d", len(lines))
		}
		// Last content line should be "..."
		if !strings.Contains(lines[3], "...") {
			t.Errorf("Overflow not handled: %q", lines[3])
		}
	})
}

func TestTextPanelDraw(t *testing.T) {
	t.Run("with_border", func(t *testing.T) {
		tp := &TextPanel{
			PanelBase: PanelBase{w: 10, h: 4, Title: "Text"},
			Text:      []rune("hello"),
			Cursor:    2,
		}
		tp.Border = true

		res := tp.Draw(true)
		full := (&tp.PanelBase).wrapWithBorder(res, true)
		lines := strings.Split(full, "\n")
		if len(lines) != 4 {
			t.Errorf("Expected 4 lines, got %d", len(lines))
		}

		// Check text content
		var expectedText = clrCyan + "│" + reset + clrWhite + "hello   " + reset + clrCyan + "│" + reset
		if lines[1] != expectedText {
			t.Errorf("Text line mismatch: got %q, want %q", lines[1], expectedText)
		}

		// Test truncation
		tp.Text = []rune("very long text")
		res = tp.Draw(true)
		full = (&tp.PanelBase).wrapWithBorder(res, true)
		lines = strings.Split(full, "\n")
		if !strings.Contains(lines[1], "very l..") {
			t.Errorf("Text truncation failed: %q", lines[1])
		}
	})

	t.Run("without_border", func(t *testing.T) {
		tp := &TextPanel{
			PanelBase: PanelBase{w: 10, h: 4, Title: "Text"},
			Text:      []rune("hello"),
			Cursor:    2,
		}
		tp.Border = false

		res := tp.Draw(true)
		full := (&tp.PanelBase).wrapWithBorder(res, true)
		lines := strings.Split(full, "\n")
		if len(lines) != 4 {
			t.Errorf("Expected 4 lines, got %d", len(lines))
		}

		// Check text content
		expectedText := "hello     "
		if lines[1] != expectedText {
			t.Errorf("Text line mismatch: got %q, want %q", lines[1], expectedText)
		}

		// Test truncation
		tp.Text = []rune("very long text")
		res = tp.Draw(true)
		full = (&tp.PanelBase).wrapWithBorder(res, true)
		lines = strings.Split(full, "\n")
		if !strings.Contains(lines[1], "very lon..") {
			t.Errorf("Text truncation failed: %q", lines[1])
		}
	})

	t.Run("without_border", func(t *testing.T) {
		tp := &TextPanel{
			PanelBase: PanelBase{w: 10, h: 4, Title: "Text"},
			Text:      []rune("hello"),
			Cursor:    2,
		}
		tp.Border = false

		res := tp.Draw(true)
		full := (&tp.PanelBase).wrapWithBorder(res, true)
		lines := strings.Split(full, "\n")
		if len(lines) != 4 {
			t.Errorf("Expected 4 lines, got %d", len(lines))
		}

		// Check text content
		expectedText := "hello     "
		if lines[1] != expectedText {
			t.Errorf("Text line mismatch: got %q, want %q", lines[1], expectedText)
		}

		// Test truncation
		tp.Text = []rune("very long text")
		res = tp.Draw(true)
		full = (&tp.PanelBase).wrapWithBorder(res, true)
		lines = strings.Split(full, "\n")
		if !strings.Contains(lines[1], "very lon..") {
			t.Errorf("Text truncation failed: %q", lines[1])
		}
	})
}

func TestInfoPanelDraw(t *testing.T) {
	t.Run("with_border", func(t *testing.T) {
		ip := &InfoPanel{
			PanelBase: PanelBase{w: 10, h: 5, Title: "Info"},
			Lines:     []string{"line1", "line2", "line3"},
		}
		ip.Border = true

		res := ip.Draw(true)
		full := (&ip.PanelBase).wrapWithBorder(res, true)
		lines := strings.Split(full, "\n")
		if len(lines) != 5 {
			t.Errorf("Expected 5 lines, got %d", len(lines))
		}

		// Check lines
		var expectedLine1 = clrCyan + "│" + reset + clrWhite + "line1   " + reset + clrCyan + "│" + reset
		if lines[1] != expectedLine1 {
			t.Errorf("Line1 mismatch: got %q, want %q", lines[1], expectedLine1)
		}

		// Test overflow
		ip.Lines = []string{"1", "2", "3", "4", "5"}
		res = ip.Draw(true)
		full = ip.wrapWithBorder(res, true)
		lines = strings.Split(full, "\n")
		if !strings.Contains(lines[3], "...") {
			t.Errorf("Overflow not handled: %q", lines[3])
		}
	})

	t.Run("without_border", func(t *testing.T) {
		ip := &InfoPanel{
			PanelBase: PanelBase{w: 10, h: 5, Title: "Info"},
			Lines:     []string{"line1", "line2", "line3"},
		}
		ip.Border = false

		res := ip.Draw(true)
		full := (&ip.PanelBase).wrapWithBorder(res, true)
		lines := strings.Split(full, "\n")
		if len(lines) != 5 {
			t.Errorf("Expected 5 lines, got %d", len(lines))
		}

		// Check lines
		expectedLine1 := "line1     "
		if lines[1] != expectedLine1 {
			t.Errorf("Line1 mismatch: got %q, want %q", lines[1], expectedLine1)
		}

		// Test overflow
		ip.Lines = []string{"1", "2", "3", "4", "5"}
		res = ip.Draw(true)
		full = (&ip.PanelBase).wrapWithBorder(res, true)
		lines = strings.Split(full, "\n")
		if lines[3] != "3         " {
			t.Errorf("Line3 mismatch: got %q, want %q", lines[3], "3         ")
		}
		if !strings.Contains(lines[4], "...") {
			t.Errorf("Overflow not handled: %q", lines[4])
		}
		if !strings.Contains(lines[4], "...") {
			t.Errorf("Overflow not handled: %q", lines[4])
		}
	})
}

// Mock panel for testing
type mockPanel struct {
	PanelBase
	id string
}

func (mp *mockPanel) Draw(active bool) string {
	return mp.id
}

func TestWeightedLayouts(t *testing.T) {
	t.Run("PanelNode_GetWeight", func(t *testing.T) {
		tests := []struct {
			name     string
			weight   int
			expected int
		}{
			{"positive weight", 5, 5},
			{"zero weight", 0, 1},
			{"negative weight", -1, 1},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				pn := &PanelNode{Weight: tt.weight}
				if got := pn.GetWeight(); got != tt.expected {
					t.Errorf("PanelNode.GetWeight() = %v, want %v", got, tt.expected)
				}
			})
		}
	})

	t.Run("HorizontalSplit_GetWeight", func(t *testing.T) {
		tests := []struct {
			name     string
			weight   int
			expected int
		}{
			{"positive weight", 3, 3},
			{"zero weight", 0, 1},
			{"negative weight", -2, 1},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				hs := &HorizontalSplit{Weight: tt.weight}
				if got := hs.GetWeight(); got != tt.expected {
					t.Errorf("HorizontalSplit.GetWeight() = %v, want %v", got, tt.expected)
				}
			})
		}
	})

	t.Run("VerticalSplit_GetWeight", func(t *testing.T) {
		tests := []struct {
			name     string
			weight   int
			expected int
		}{
			{"positive weight", 4, 4},
			{"zero weight", 0, 1},
			{"negative weight", -3, 1},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				vs := &VerticalSplit{Weight: tt.weight}
				if got := vs.GetWeight(); got != tt.expected {
					t.Errorf("VerticalSplit.GetWeight() = %v, want %v", got, tt.expected)
				}
			})
		}
	})

	t.Run("HorizontalSplit_position_equal_weights", func(t *testing.T) {
		p1 := &mockPanel{PanelBase: PanelBase{}, id: "p1"}
		p2 := &mockPanel{PanelBase: PanelBase{}, id: "p2"}
		p3 := &mockPanel{PanelBase: PanelBase{}, id: "p3"}

		hs := &HorizontalSplit{
			Panels: []Layout{
				&PanelNode{Panel: p1, Weight: 1},
				&PanelNode{Panel: p2, Weight: 1},
				&PanelNode{Panel: p3, Weight: 1},
			},
		}

		panels := hs.position(0, 0, 90, 10)

		if len(panels) != 3 {
			t.Errorf("Expected 3 panels, got %d", len(panels))
		}

		// Check positions and sizes (90 width / 3 panels = 30 each)
		expectedWidths := []int{30, 30, 30}
		expectedX := []int{0, 30, 60}

		for i, panel := range panels {
			pb := panel.GetBase()
			if pb.w != expectedWidths[i] {
				t.Errorf("Panel %d width = %d, want %d", i, pb.w, expectedWidths[i])
			}
			if pb.x != expectedX[i] {
				t.Errorf("Panel %d x = %d, want %d", i, pb.x, expectedX[i])
			}
			if pb.y != 0 {
				t.Errorf("Panel %d y = %d, want 0", i, pb.y)
			}
			if pb.h != 10 {
				t.Errorf("Panel %d h = %d, want 10", i, pb.h)
			}
		}
	})

	t.Run("HorizontalSplit_position_unequal_weights", func(t *testing.T) {
		p1 := &mockPanel{PanelBase: PanelBase{}, id: "p1"}
		p2 := &mockPanel{PanelBase: PanelBase{}, id: "p2"}
		p3 := &mockPanel{PanelBase: PanelBase{}, id: "p3"}

		hs := &HorizontalSplit{
			Panels: []Layout{
				&PanelNode{Panel: p1, Weight: 1}, // 1/7 of space
				&PanelNode{Panel: p2, Weight: 2}, // 2/7 of space
				&PanelNode{Panel: p3, Weight: 4}, // 4/7 of space
			},
		}

		panels := hs.position(0, 0, 70, 10)

		if len(panels) != 3 {
			t.Errorf("Expected 3 panels, got %d", len(panels))
		}

		// Check positions and sizes (70 width, weights 1+2+4=7)
		// Panel 1: (1/7)*70 = 10, Panel 2: (2/7)*70 = 20, Panel 3: (4/7)*70 = 40
		// But last panel gets remainder, so: 10, 20, 40 (total 70)
		expectedWidths := []int{10, 20, 40}
		expectedX := []int{0, 10, 30}

		for i, panel := range panels {
			pb := panel.GetBase()
			if pb.w != expectedWidths[i] {
				t.Errorf("Panel %d width = %d, want %d", i, pb.w, expectedWidths[i])
			}
			if pb.x != expectedX[i] {
				t.Errorf("Panel %d x = %d, want %d", i, pb.x, expectedX[i])
			}
		}
	})

	t.Run("HorizontalSplit_position_zero_weights", func(t *testing.T) {
		p1 := &mockPanel{PanelBase: PanelBase{}, id: "p1"}
		p2 := &mockPanel{PanelBase: PanelBase{}, id: "p2"}
		p3 := &mockPanel{PanelBase: PanelBase{}, id: "p3"}

		hs := &HorizontalSplit{
			Panels: []Layout{
				&PanelNode{Panel: p1, Weight: 0}, // Should get weight 1
				&PanelNode{Panel: p2, Weight: 0}, // Should get weight 1
				&PanelNode{Panel: p3, Weight: 0}, // Should get weight 1
			},
		}

		panels := hs.position(0, 0, 90, 10)

		if len(panels) != 3 {
			t.Errorf("Expected 3 panels, got %d", len(panels))
		}

		// All panels should have equal width (90 / 3 = 30)
		for i, panel := range panels {
			pb := panel.GetBase()
			if pb.w != 30 {
				t.Errorf("Panel %d width = %d, want 30", i, pb.w)
			}
		}
	})

	t.Run("VerticalSplit_position_equal_weights", func(t *testing.T) {
		p1 := &mockPanel{PanelBase: PanelBase{}, id: "p1"}
		p2 := &mockPanel{PanelBase: PanelBase{}, id: "p2"}

		vs := &VerticalSplit{
			Panels: []Layout{
				&PanelNode{Panel: p1, Weight: 1},
				&PanelNode{Panel: p2, Weight: 1},
			},
		}

		panels := vs.position(0, 0, 20, 40)

		if len(panels) != 2 {
			t.Errorf("Expected 2 panels, got %d", len(panels))
		}

		// Check positions and sizes (40 height / 2 panels = 20 each)
		expectedHeights := []int{20, 20}
		expectedY := []int{0, 20}

		for i, panel := range panels {
			pb := panel.GetBase()
			if pb.h != expectedHeights[i] {
				t.Errorf("Panel %d height = %d, want %d", i, pb.h, expectedHeights[i])
			}
			if pb.y != expectedY[i] {
				t.Errorf("Panel %d y = %d, want %d", i, pb.y, expectedY[i])
			}
			if pb.x != 0 {
				t.Errorf("Panel %d x = %d, want 0", i, pb.x)
			}
			if pb.w != 20 {
				t.Errorf("Panel %d w = %d, want 20", i, pb.w)
			}
		}
	})

	t.Run("VerticalSplit_position_unequal_weights", func(t *testing.T) {
		p1 := &mockPanel{PanelBase: PanelBase{}, id: "p1"}
		p2 := &mockPanel{PanelBase: PanelBase{}, id: "p2"}
		p3 := &mockPanel{PanelBase: PanelBase{}, id: "p3"}

		vs := &VerticalSplit{
			Panels: []Layout{
				&PanelNode{Panel: p1, Weight: 1}, // 1/4 of space
				&PanelNode{Panel: p2, Weight: 2}, // 2/4 of space
				&PanelNode{Panel: p3, Weight: 1}, // 1/4 of space
			},
		}

		panels := vs.position(0, 0, 20, 40)

		if len(panels) != 3 {
			t.Errorf("Expected 3 panels, got %d", len(panels))
		}

		// Check positions and sizes (40 height, weights 1+2+1=4)
		// Panel 1: (1/4)*40 = 10, Panel 2: (2/4)*40 = 20, Panel 3: (1/4)*40 = 10
		expectedHeights := []int{10, 20, 10}
		expectedY := []int{0, 10, 30}

		for i, panel := range panels {
			pb := panel.GetBase()
			if pb.h != expectedHeights[i] {
				t.Errorf("Panel %d height = %d, want %d", i, pb.h, expectedHeights[i])
			}
			if pb.y != expectedY[i] {
				t.Errorf("Panel %d y = %d, want %d", i, pb.y, expectedY[i])
			}
		}
	})

	t.Run("Nested_weighted_layouts", func(t *testing.T) {
		// Create a complex nested layout:
		// VerticalSplit (weight: 1)
		// ├── HorizontalSplit (weight: 2) with 3 panels (weights: 1,1,1)
		// └── PanelNode (weight: 1)

		p1 := &mockPanel{PanelBase: PanelBase{}, id: "p1"}
		p2 := &mockPanel{PanelBase: PanelBase{}, id: "p2"}
		p3 := &mockPanel{PanelBase: PanelBase{}, id: "p3"}
		p4 := &mockPanel{PanelBase: PanelBase{}, id: "p4"}

		innerHS := &HorizontalSplit{
			Panels: []Layout{
				&PanelNode{Panel: p1, Weight: 1},
				&PanelNode{Panel: p2, Weight: 1},
				&PanelNode{Panel: p3, Weight: 1},
			},
			Weight: 2,
		}

		vs := &VerticalSplit{
			Panels: []Layout{
				innerHS,
				&PanelNode{Panel: p4, Weight: 1},
			},
			Weight: 1,
		}

		panels := vs.position(0, 0, 60, 40)

		if len(panels) != 4 {
			t.Errorf("Expected 4 panels, got %d", len(panels))
		}

		// Vertical split: innerHS gets 2/3 of height (26px), p4 gets 1/3 (13px)
		// But due to rounding, last panel gets remainder

		// Check that panels are positioned correctly
		pb1 := panels[0].GetBase() // p1 from innerHS
		pb2 := panels[1].GetBase() // p2 from innerHS
		pb3 := panels[2].GetBase() // p3 from innerHS
		pb4 := panels[3].GetBase() // p4

		// All panels from innerHS should have same Y position and height
		if pb1.y != pb2.y || pb2.y != pb3.y {
			t.Errorf("Inner panels should have same Y position")
		}
		if pb1.h != pb2.h || pb2.h != pb3.h {
			t.Errorf("Inner panels should have same height")
		}

		// p4 should be below the inner panels
		if pb4.y < pb1.y+pb1.h {
			t.Errorf("p4 should be below inner panels")
		}

		// Check horizontal distribution in inner layout (60 width / 3 = 20 each)
		if pb1.w != 20 || pb2.w != 20 || pb3.w != 20 {
			t.Errorf("Inner panels should each have width 20, got %d, %d, %d", pb1.w, pb2.w, pb3.w)
		}
	})

	t.Run("Empty_layouts", func(t *testing.T) {
		hs := &HorizontalSplit{Panels: []Layout{}}
		panels := hs.position(0, 0, 100, 50)
		if len(panels) != 0 {
			t.Errorf("Empty layout should return 0 panels, got %d", len(panels))
		}

		vs := &VerticalSplit{Panels: []Layout{}}
		panels = vs.position(0, 0, 100, 50)
		if len(panels) != 0 {
			t.Errorf("Empty layout should return 0 panels, got %d", len(panels))
		}
	})

	t.Run("Single_panel_layouts", func(t *testing.T) {
		p1 := &mockPanel{PanelBase: PanelBase{}, id: "p1"}

		hs := &HorizontalSplit{
			Panels: []Layout{&PanelNode{Panel: p1, Weight: 5}},
		}
		panels := hs.position(0, 0, 100, 50)

		if len(panels) != 1 {
			t.Errorf("Expected 1 panel, got %d", len(panels))
		}

		pb := panels[0].GetBase()
		if pb.x != 0 || pb.y != 0 || pb.w != 100 || pb.h != 50 {
			t.Errorf("Single panel should fill entire space: x=%d, y=%d, w=%d, h=%d", pb.x, pb.y, pb.w, pb.h)
		}
	})
}

func TestTruncateToWidth(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		width    int
		expected string
	}{
		{"no truncation", "hello", 10, "hello"},
		{"exact width", "hello", 5, "hello"},
		{"truncation", "hello world", 8, "hello .."},
		{"short width", "hello", 2, "he"},
		{"empty string", "", 5, ""},
		{"width zero", "hello", 0, ""},
		{"with ANSI", "\x1b[31mred\x1b[0m text", 6, "\x1b[31mred\x1b[0m .."},
		{"only ANSI", "\x1b[31m\x1b[0m", 5, "\x1b[31m\x1b[0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateToWidth(tt.input, tt.width)
			if result != tt.expected {
				t.Errorf("truncateToWidth(%q, %d) = %q, want %q", tt.input, tt.width, result, tt.expected)
			}
		})
	}
}

func TestDisplayWidth(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"simple", "hello", 5},
		{"with spaces", "hello world", 11},
		{"empty", "", 0},
		{"with ANSI", "\x1b[31mred\x1b[0m", 3},
		{"only ANSI", "\x1b[31m\x1b[0m", 0},
		{"mixed", "a\x1b[31mb\x1b[0mc", 3},
		{"unicode", "café", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := displayWidth(tt.input)
			if result != tt.expected {
				t.Errorf("displayWidth(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestWrapWithBorder(t *testing.T) {
	t.Run("small dimensions", func(t *testing.T) {
		pb := &PanelBase{w: 2, h: 2}
		result := pb.wrapWithBorder("content", true)
		if result != "" {
			t.Errorf("Expected empty string for small dimensions, got %q", result)
		}
	})

	t.Run("without border", func(t *testing.T) {
		pb := &PanelBase{w: 10, h: 5, Title: "Test"}
		pb.Border = false
		result := pb.wrapWithBorder("line1\nline2", true)
		lines := strings.Split(result, "\n")
		if len(lines) != 5 {
			t.Errorf("Expected 5 lines, got %d", len(lines))
		}
		if lines[0] != "Test      " {
			t.Errorf("Title line mismatch: got %q", lines[0])
		}
		if lines[1] != "line1     " {
			t.Errorf("Content line1 mismatch: got %q", lines[1])
		}
		if lines[2] != "line2     " {
			t.Errorf("Content line2 mismatch: got %q", lines[2])
		}
	})

	t.Run("with border", func(t *testing.T) {
		pb := &PanelBase{w: 10, h: 5, Title: "Test"}
		pb.Border = true
		result := pb.wrapWithBorder("line1\nline2", true)
		lines := strings.Split(result, "\n")
		if len(lines) != 5 {
			t.Errorf("Expected 5 lines, got %d", len(lines))
		}
		expectedTop := clrCyan + "┌ [Test] ┐" + reset
		if lines[0] != expectedTop {
			t.Errorf("Top border mismatch: got %q, want %q", lines[0], expectedTop)
		}
		expectedContent := clrCyan + "│" + reset + clrWhite + "line1   " + reset + clrCyan + "│" + reset
		if lines[1] != expectedContent {
			t.Errorf("Content line mismatch: got %q, want %q", lines[1], expectedContent)
		}
	})

	t.Run("content overflow", func(t *testing.T) {
		pb := &PanelBase{w: 10, h: 4, Title: "Test"}
		pb.Border = true
		result := pb.wrapWithBorder("1\n2\n3\n4\n5", true)
		lines := strings.Split(result, "\n")
		if len(lines) != 4 {
			t.Errorf("Expected 4 lines, got %d", len(lines))
		}
		// Last content line should be truncated with "..."
		if !strings.Contains(lines[2], "...") {
			t.Errorf("Overflow not handled: %q", lines[2])
		}
	})
}
