package SFTP

import (
	"fmt"
	"testing"
)

func TestSimple(t *testing.T) {
	h := New("fransebas")
	filePath := "/Users/fransebas"

	if h.CanWrite(filePath) {
		fmt.Println("Can write file " + filePath)
	}

	if h.CanRead(filePath) {
		fmt.Println("Can read file " + filePath)
	}

	if h.CanExecute(filePath) {
		fmt.Println("Can execute file " + filePath)
	}
}
