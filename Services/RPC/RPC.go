package RPC

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"reflect"
	"sssh_server/CustomUtils"
	"sssh_server/Services/API"
	"sssh_server/Services/SessionLayer"
)

type RPC struct {
	Port           int
	receiversMap   map[string]API.OnCommandService
	client         *rpc.Client
	sessionService *SessionLayer.SessionService
	receiver       Receiver
}

type RPCArgs struct {
	Name      string
	Arg       string
	SessionID string
}

type Receiver struct {
	f func(args *RPCArgs)
}

// On server
func (r *Receiver) Receive(arg *RPCArgs, response *string) error {
	r.f(arg)
	return nil
}

// On client
func (p *RPC) OnCommand(cmnd, sessionID string) {
	client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%v", p.Port))
	CustomUtils.CheckPanic(err, "Couldn't make rpc client")
	p.client = client

	for _, service := range p.receiversMap {
		res := service.GetClientCode()(cmnd)

		rpcArgs := new(RPCArgs)
		rpcArgs.Name = service.GetName()
		rpcArgs.Arg = res
		rpcArgs.SessionID = sessionID

		e := p.client.Call("Receiver.Receive", &rpcArgs, nil)
		CustomUtils.CheckPrint(e)
	}
}

// On server and client
func New(Port int) *RPC {
	nrpc := new(RPC)
	nrpc.Port = Port
	nrpc.receiversMap = make(map[string]API.OnCommandService)
	return nrpc
}

// On server and client
func (p *RPC) AddService(service API.OnCommandService) {
	r := new(Receiver)

	e := rpc.Register(r)
	CustomUtils.CheckPanic(e, "Could't register rpc receiver")
	p.receiversMap[service.GetName()] = service
}

// On server
func (p *RPC) Serve(sessionService *SessionLayer.SessionService) {
	p.sessionService = sessionService
	p.receiver.f = func(args *RPCArgs) {
		//p.sessionService.Sessions[args.SessionID].
		//p.receiversMap[args.Name].ClientResponse(args.Arg)
	}
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf(":%v", p.Port))
	CustomUtils.CheckPanic(e, "Could not register symbiont struct")

	go func() {
		e = http.Serve(l, nil)
		panic(e)
	}()
}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
