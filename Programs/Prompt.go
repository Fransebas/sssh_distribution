package Programs

import (
	"fmt"
	"net/http"
	"sssh_server/CustomUtils"
	"sssh_server/Modules/Configuration"
	"strings"
)

func Prompt(config Configuration.Configuration) {
	_, e := http.Post(fmt.Sprintf("http://localhost:2000/newcommand?SSSH_USER=%v", config.UserId), "text/html", strings.NewReader(config.History))
	CustomUtils.CheckPrint(e)
	updateVariables(config)
}

// I don't like this solution but time will tell
func updateVariables(config Configuration.Configuration) {
	data := CustomUtils.ExecuteCommand("env")
	_, e := http.Post(fmt.Sprintf("http://localhost:2000/variables?SSSH_USER=%v", config.UserId), "text/html", strings.NewReader(data))
	CustomUtils.CheckPrint(e)
}
