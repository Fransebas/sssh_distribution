package CommandExecuter

import (
	"fmt"
	"sssh_server/CustomUtils"
	"testing"
)

func TestCommandExecuter(t *testing.T) {
	c := CustomUtils.ExecuteCommand("echo abc", "fransebas")

	c2 := c.Args

	fmt.Println(c2)

	b, e := c.Output()
	s := string(b)
	if e != nil {
		panic(e)
		s2 := e.Error()
		fmt.Println(s2)
	}
	fmt.Sprintln(s + "eft \n")
}
