package DirectoryManager

import (
	"fmt"
	"testing"
)

func TestHomeDir(t *testing.T) {

	dm := New("fransebas")
	fmt.Println(dm.UserDirectory)
	fmt.Println(dm.GetConfigFile("abc"))
	fmt.Println(dm.GetVariableFile("abc"))
}
