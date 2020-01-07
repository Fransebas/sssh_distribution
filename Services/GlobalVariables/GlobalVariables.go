package GlobalVariables

import (
	"encoding/json"
	"fmt"
	"os"
	"sssh_server/Services/API"
	"sssh_server/Services/CommandExecuter"
	"sssh_server/Services/SSH"
	"strings"
)

type GlobalVariables struct {
	vars string
}

type BashVar struct {
	Name  string
	Value string
}

func (g *GlobalVariables) getVariables() string {
	//executer := CommandExecuter.CommandExecuter{}
	//res := executer.ExecuteCommand("env")
	lines := strings.Split(g.vars, "\n")
	variables := []BashVar{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		r := strings.SplitAfterN(line, "=", 2)
		variable := BashVar{
			Name:  r[0][:len(r[0])-1],
			Value: r[1],
		}
		variables = append(variables, variable)
	}
	b, e := json.Marshal(variables)

	if e != nil {
		panic(e)
	}
	return string(b)
}

func (*GlobalVariables) storeVariable(bashVar BashVar) error {
	//executer := CommandExecuter.CommandExecuter{}
	// TODO : create function
	//executer.ExecuteCommand(fmt.Sprintf( "export %v=%v", bashVar.Name, bashVar.Value))

	return os.Setenv(bashVar.Name, bashVar.Value)
}

func (*GlobalVariables) OnNewSession(s API.TerminalSessionInterface) {}
func (*GlobalVariables) OnNewConnection(sshSession *SSH.SSHSession)  {}

func (*GlobalVariables) GetClientCode() API.ClientCode {
	return func(cmnd string) string {
		exec := CommandExecuter.CommandExecuter{}
		return exec.ExecuteCommand("env")
	}
}

func (*GlobalVariables) GetName() string {
	return "GlobalVariables"
}

func (g *GlobalVariables) ClientResponse(res string) {
	fmt.Println("response!!!!" + res)
	g.vars = res
}
