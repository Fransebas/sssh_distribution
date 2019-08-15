package CommandList

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"sssh_server/CustomUtils"
	"sssh_server/Services/RecentCommands/Models"
	"strings"
)

type CommandList struct {
	GetCommandsScriptPath string
}

func getContent(path string) string {
	data, err := ioutil.ReadFile(path)
	CustomUtils.CheckPrint(err)
	return string(data)
}

func NewCommandList() (cl *CommandList) {
	cl = new(CommandList)
	basePath, err := filepath.Abs("Assets/Scripts")
	cl.GetCommandsScriptPath = fmt.Sprintf("%v/%v", basePath, "get_commands")
	CustomUtils.CheckPrint(err)
	return
}

func processRaw(commandsRaw string) []Models.Command {
	commandsStr := strings.Split(commandsRaw, "\n")
	commands := []Models.Command{}
	for _, command := range commandsStr {
		commands = append(commands, Models.NewCommand(command))
	}
	return commands
}

func (cl *CommandList) GetList(arg string) []Models.Command {
	c := exec.Command("compgen", strings.Split(arg, " ")...)
	b, e := c.Output()
	CustomUtils.CheckPrint(e)
	commandsRaw := string(b)
	return processRaw(commandsRaw)
}

func (cl *CommandList) GetCommandsList() []Models.Command {
	data, err := ioutil.ReadFile(cl.GetCommandsScriptPath)
	CustomUtils.CheckPrint(err)
	//fmt.Printf("The data ist %v \n", string(data))
	c := exec.Command("bash", "-c", string(data))
	b, e := c.Output()
	CustomUtils.CheckPrint(e)
	return processRaw(string(b))
}
