package color

import "testing"

func TestStyle_AddModifier(t *testing.T) {
	s := NewStyle(nil, nil)
	s.AddModifier(Bold)

	if len(s.Modifiers) != 1 {
		t.Error("expected 1 modifier")
	}
}
