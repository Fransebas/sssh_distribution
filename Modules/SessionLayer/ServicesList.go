package SessionLayer

import (
	"sssh_server/Modules/API"
	"sssh_server/Modules/CommandExecuter"
	"sssh_server/Modules/CommandList"
	"sssh_server/Modules/EchoService"
	"sssh_server/Modules/GlobalVariables"
	"sssh_server/Modules/History"
	"sssh_server/Modules/TerminalService"
)

var rawServices []interface{} = []interface{}{
	new(TerminalService.TerminalService),
	new(History.History),
	new(CommandList.CommandListService),
	new(CommandExecuter.CommandExecuter),
	new(EchoService.EchoService),
	new(GlobalVariables.GlobalVariables),
}

// List of all the available services
var Services []API.Module = []API.Module{}

// List of the services that implement the API.OnCommandService interface
var CommandServices []API.OnCommandService = []API.OnCommandService{}

func clasify(i interface{}) {
	if s, ok := i.(API.Module); ok {
		Services = append(Services, s)
	}
	if s, ok := i.(API.OnCommandService); ok {
		CommandServices = append(CommandServices, s)
	}
}

func init() {
	for _, service := range rawServices {
		clasify(service)
	}
}
