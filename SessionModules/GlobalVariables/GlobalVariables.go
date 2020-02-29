package GlobalVariables

import (
	"encoding/json"
	"os"
	"sssh_server/Modules/SSH"
	"sssh_server/SessionModules/API"
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

func (g *GlobalVariables) OnUpdateVariables(vars string) {
	g.vars = vars
}
