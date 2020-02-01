package History

import (
	"encoding/json"
	"io"
	"sssh_server/SessionModules/API"
)

func (h *History) GetHandlers() []*API.RequestHandler {
	connection := API.RequestHandler{
		RequestHandler: func(w io.Writer, r io.Reader) {
			// Here a new connection is establish, should send all the current commands
			historyJSON, _ := json.Marshal(h.Cmnds)
			_, _ = w.Write(historyJSON)
			h.OnCommandsUpdateF = func(commands []Command) {
				// TODO: properly handle error
				historyJSON, _ := json.Marshal(commands)
				_, _ = w.Write(historyJSON)
			}
		},
		Name: "commands",
	}

	return []*API.RequestHandler{&connection}
}
