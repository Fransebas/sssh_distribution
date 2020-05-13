package Programs

import (
	"fmt"
	"net/http"
	"sssh_server/CustomUtils"
	"sssh_server/Modules/Configuration"
	"strings"
)

func Prompt(config Configuration.Configuration) {
	promptConfig := config.PromptConfig
	_, e := http.Post(fmt.Sprintf("http://localhost:2000/newcommand?SSSH_USER=%v", promptConfig.UserId), "text/html", strings.NewReader(promptConfig.History))
	CustomUtils.CheckPrint(e)
	updateVariables(promptConfig)
	updatePWD(promptConfig)
}

// I don't like this solution but time will tell
func updateVariables(config Configuration.PromptConfig) {
	data := CustomUtils.SudoExecuteCommandOnce("env")
	_, e := http.Post(fmt.Sprintf("http://localhost:2000/variables?SSSH_USER=%v", config.UserId), "text/html", strings.NewReader(data))
	CustomUtils.CheckPrint(e)
}

func updatePWD(config Configuration.PromptConfig) {
	_, e := http.Post(fmt.Sprintf("http://localhost:2000/pwd?SSSH_USER=%v", config.UserId), "text/html", strings.NewReader(config.Pwd))
	CustomUtils.CheckPrint(e)
}
