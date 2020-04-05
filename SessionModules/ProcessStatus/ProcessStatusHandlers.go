package ProcessStatus

import (
	"encoding/json"
	"io"
	"sssh_server/CustomUtils"
	"sssh_server/Modules/Logging"
	"sssh_server/SessionModules/API"
)

func (cl *ProcessStatusModule) GetHandlers() []*API.RequestHandler {
	processStatusRequest := API.RequestHandler{
		RequestHandler: func(w io.Writer, r io.Reader) {

			CustomUtils.Logger.Println(Logging.INFO, "Hit")

			mapp := cl.getProcessStatus()
			b, e := json.Marshal(mapp)
			CustomUtils.CheckPrint(e)
			_, e = w.Write(b)
			CustomUtils.CheckPrint(e)
		},
		Name: "psaux",
	}

	return []*API.RequestHandler{&processStatusRequest}
}
