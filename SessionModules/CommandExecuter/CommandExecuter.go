package CommandExecuter

import (
	"os/exec"
	"sssh_server/CustomUtils"
	"sssh_server/Modules/SSH"
	"sssh_server/SessionModules/API"
)

type CommandExecuter struct {
	user string
}

func (cme *CommandExecuter) ExecuteCommand(cmmnd string) string {
	// TODO: Important!!! usse the user to prevent a breach
	//fmt.Println("comnd : " + cmmnd)
	// s := strings.Split(cmmnd, " ")
	// c := exec.Command(s[0], s[1:]...)
	c := exec.Command("bash", "-c", cmmnd)
	b, e := c.Output()
	CustomUtils.CheckPrint(e)

	return string(b)
}

func (cme *CommandExecuter) OnNewSession(session API.TerminalSessionInterface) {
	cme.user = session.GetUsername()
}
func (cme *CommandExecuter) OnNewConnection(sshSession *SSH.SSHSession) {

}

func (cme *CommandExecuter) Close() {

}
