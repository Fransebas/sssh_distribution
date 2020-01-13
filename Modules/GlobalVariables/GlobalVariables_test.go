package GlobalVariables

import (
	"testing"
)

func TestGetVariables(t *testing.T) {
	g := GlobalVariables{}
	s := g.getVariables()
	if len(s) < 1 {
		t.Error("No variables fetch")
	}
}
