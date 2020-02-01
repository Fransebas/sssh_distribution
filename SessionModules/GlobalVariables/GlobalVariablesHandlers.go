package GlobalVariables

import (
	"encoding/json"
	"io"
	"sssh_server/CustomUtils"
	"sssh_server/SessionModules/API"
)

func (gv *GlobalVariables) GetHandlers() []*API.RequestHandler {
	// Return an array of vars with the global variables
	globalVars := API.RequestHandler{
		RequestHandler: func(w io.Writer, r io.Reader) {
			_, _ = w.Write([]byte(gv.getVariables()))
		},
		Name: "globalVars",
	}

	// Set a given var
	setVar := API.RequestHandler{
		RequestHandler: func(w io.Writer, r io.Reader) {
			b, _ := CustomUtils.Read(r)
			var bashVar BashVar
			_ = json.Unmarshal(b, &bashVar)
			_ = gv.storeVariable(bashVar)
		},
		Name: "setVar",
	}

	return []*API.RequestHandler{&globalVars, &setVar}
}
