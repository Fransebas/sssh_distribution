package CommandExecuter

import (
	"os/exec"
	"rgui/CustomUtils"
)

type CommandExecuter struct {
}

func (cme *CommandExecuter) ExecuteCommand(cmmnd string) string {
	//fmt.Println("comnd : " + cmmnd)
	// s := strings.Split(cmmnd, " ")
	// c := exec.Command(s[0], s[1:]...)
	c := exec.Command("bash", "-c", cmmnd)
	b, e := c.Output()
	CustomUtils.CheckPrint(e)

	return string(b)
}
