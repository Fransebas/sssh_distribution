package SessionLayer

import (
	"sssh_server/SessionModules/API"
	"sssh_server/SessionModules/CommandExecuter"
	"sssh_server/SessionModules/CommandList"
	"sssh_server/SessionModules/EchoService"
	"sssh_server/SessionModules/GlobalVariables"
	"sssh_server/SessionModules/History"
	"sssh_server/SessionModules/ProcessStatus"
	"sssh_server/SessionModules/TerminalService"
)

var rawServices []interface{} = []interface{}{
	new(TerminalService.TerminalService),
	new(History.History),
	new(CommandList.CommandListService),
	new(CommandExecuter.CommandExecuter),
	new(EchoService.EchoService),
	new(GlobalVariables.GlobalVariables),
	new(ProcessStatus.ProcessStatusModule),
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
