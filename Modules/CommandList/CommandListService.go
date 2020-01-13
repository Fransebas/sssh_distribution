package CommandList

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"sssh_server/CustomUtils"
	"sssh_server/Modules/API"
	"sssh_server/Modules/RecentCommands/Models"
	"sssh_server/Modules/SSH"
	"strings"
)

type CommandListService struct {
	GetCommandsScriptPath string
}

func processRaw(commandsRaw string) []Models.Command {
	commandsStr := strings.Split(commandsRaw, "\n")
	commands := []Models.Command{}
	for _, command := range commandsStr {
		commands = append(commands, Models.NewCommand(command))
	}
	return commands
}

func (cl *CommandListService) OnNewSession(session API.TerminalSessionInterface) {
	basePath, err := filepath.Abs("Assets/Scripts")
	cl.GetCommandsScriptPath = fmt.Sprintf("%v/%v", basePath, "get_commands")
	CustomUtils.CheckPrint(err)
}
func (cl *CommandListService) OnNewConnection(sshSession *SSH.SSHSession) {}

func (cl *CommandListService) getList(arg string) []Models.Command {
	c := exec.Command("compgen", strings.Split(arg, " ")...)
	b, e := c.Output()
	CustomUtils.CheckPrint(e)
	commandsRaw := string(b)
	return processRaw(commandsRaw)
}

func (cl *CommandListService) getCommandsList() []Models.Command {
	data, err := ioutil.ReadFile(cl.GetCommandsScriptPath)
	CustomUtils.CheckPrint(err)
	//fmt.Printf("The data ist %v \n", string(data))
	c := exec.Command("bash", "-c", string(data))
	b, e := c.Output()
	CustomUtils.CheckPrint(e)
	return processRaw(string(b))
}

var _ API.Module = (*CommandListService)(nil) // Verify that *T implements I.
