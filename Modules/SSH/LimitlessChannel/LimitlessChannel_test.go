package LimitlessChannel

import (
	"fmt"
	"testing"
)

func TestSendingSmallMessage(t *testing.T) {
	lw := LimitlessChannel{
		Writer: ConsoleWriter{},
	}
	lw.Write([]byte("123456789abcdefghijkmnopqrstuvwx"))
}

type ConsoleWriter struct {
}

func (ConsoleWriter) Write(b []byte) (int, error) {
	fmt.Print(string(b))
	return len(b), nil
}
