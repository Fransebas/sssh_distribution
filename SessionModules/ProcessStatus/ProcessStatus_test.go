package ProcessStatus

import (
	"testing"
)

func TestGetVariables(t *testing.T) {
	ps := ProcessStatusModule{}
	ps.getProcessStatus()
}
