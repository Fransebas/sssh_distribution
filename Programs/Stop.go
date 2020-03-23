package Programs

import (
	"fmt"
	"net/http"
	"sssh_server/CustomUtils"
	"sssh_server/Modules/Configuration"
	"strings"
)

func Stop(config Configuration.Configuration) {
	_, e := http.Post(fmt.Sprintf("http://localhost:%v/stop", config.HTTPPort), "text/html", strings.NewReader(""))
	CustomUtils.CheckPanic(e, "Couldn't stop server")
	fmt.Println("server stopped")
}
