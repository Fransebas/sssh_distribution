package SessionLayer

import (
	"sssh_server/Services/API"
	"sssh_server/Services/CommandExecuter"
	"sssh_server/Services/CommandList"
	"sssh_server/Services/EchoService"
	"sssh_server/Services/GlobalVariables"
	"sssh_server/Services/History"
	"sssh_server/Services/TerminalService"
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
var Services []API.Service = []API.Service{}

// List of the services that implement the API.OnCommandService interface
var CommandServices []API.OnCommandService = []API.OnCommandService{}

func clasify(i interface{}) {
	if s, ok := i.(API.Service); ok {
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
