package Bifrost

import (
	"net/rpc"
	"sssh_server/CustomUtils"
)

const PORT = "8888"

type Bifrost struct {
}

func (*Bifrost) RunCmnd(cmnd string) (string, error) {
	client, err := rpc.DialHTTP("tcp", "localhost:"+PORT)
	CustomUtils.CheckPanic(err, "Couldn't start client")
	var response string
	err = client.Call("Symbiont.RunCmnd", cmnd, &response)
	return response, err
}
