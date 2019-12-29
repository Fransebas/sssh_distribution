package CommandList

import (
	"encoding/json"
	"io"
	"sssh_server/CustomUtils"
	"sssh_server/Services/API"
)

func (cl *CommandListService) GetHandlers() []*API.RequestHandler {
	commandListRequest := API.RequestHandler{
		RequestHandler: func(w io.Writer, r io.Reader) {
			// Here a new connection is establish, should send all the current commands
			str, err := json.Marshal(cl.getCommandsList())
			CustomUtils.CheckPrint(err)
			_, _ = w.Write([]byte(str))
		},
		Name: "commandlist",
	}

	return []*API.RequestHandler{&commandListRequest}
}
