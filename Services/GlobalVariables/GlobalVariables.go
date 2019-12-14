package GlobalVariables

import (
	"encoding/json"
	"os"
	"sssh_server/Services/CommandExecuter"
	"strings"
)

type GlobalVariables struct {
}

type BashVar struct {
	Name  string
	Value string
}

func (*GlobalVariables) GetVariables() string {
	executer := CommandExecuter.CommandExecuter{}
	res := executer.ExecuteCommand("env")
	lines := strings.Split(res, "\n")
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

func (*GlobalVariables) StoreVariable(bashVar BashVar) error {
	//executer := CommandExecuter.CommandExecuter{}
	// TODO : create function
	//executer.ExecuteCommand(fmt.Sprintf( "export %v=%v", bashVar.Name, bashVar.Value))

	return os.Setenv(bashVar.Name, bashVar.Value)
}
