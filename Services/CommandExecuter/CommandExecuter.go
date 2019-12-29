package CommandExecuter

import (
	"os/exec"
	"sssh_server/CustomUtils"
	"sssh_server/Services/API"
	"sssh_server/Services/SSH"
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

func (cme *CommandExecuter) OnNewSession(session API.TerminalSessionInterface) {

}
func (cme *CommandExecuter) OnNewConnection(sshSession *SSH.SSHSession) {

}
