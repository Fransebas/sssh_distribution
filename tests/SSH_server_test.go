package tests

import (
	"fmt"
	gossh "golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"sssh_server/CustomUtils"
	"testing"
)

func TestGliderBallSSHServer(t *testing.T) {

	//handler :=  func(srv *ssh.Server, conn *gossh.ServerConn, newChan gossh.NewChannel, ctx ssh.Context) {
	//	ch, reqs, err := newChan.Accept()
	//	if err != nil {
	//		// TODO: trigger event callback
	//		return
	//	}
	//	sess := &session{
	//		Channel:   ch,
	//		conn:      conn,
	//		handler:   srv.Handler,
	//		ptyCb:     srv.PtyCallback,
	//		sessReqCb: srv.SessionRequestCallback,
	//		ctx:       ctx,
	//	}
	//	sess.handleRequests(reqs)
	//}
	//server := ssh.Server{
	//	Addr: ":2222",
	//	Handler: handler,
	//}
	//
	//// server.ChannelHandlers["session"] = handler.(ssh.ChannelHandler)
	//
	//log.Fatal(server.ListenAndServe())
}

func TestGoSSHServer(t *testing.T) {
	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &gossh.ServerConfig{
		// Remove to disable password auth.
		PasswordCallback: func(c gossh.ConnMetadata, pass []byte) (*gossh.Permissions, error) {
			// Should use constant-time compare (or better, salt+hash) in
			// a production setting.
			if c.User() == "testuser" && string(pass) == "tiger" {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}

	privateBytes, err := ioutil.ReadFile("../id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key: ", err)
	}

	private, err := gossh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}

	config.AddHostKey(private)

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", "0.0.0.0:2222")
	if err != nil {
		log.Fatal("failed to listen for connection: ", err)
	}
	nConn, err := listener.Accept()
	if err != nil {
		log.Fatal("failed to accept incoming connection: ", err)
	}

	// Before use, a handshake must be performed on the incoming
	// net.Conn.
	conn, chans, reqs, err := gossh.NewServerConn(nConn, config)
	if err != nil {
		log.Fatal("failed to handshake: ", err)
	}

	// The incoming Request channel must be serviced.
	go gossh.DiscardRequests(reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans {
		// Channels have a type, depending on the application level
		// protocol intended. In the case of a shell, the type is
		// "session" and ServerShell may be used to present a simple
		// terminal interface.
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(gossh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Fatalf("Could not accept channel: %v", err)
		}

		go func() { handleRequests(requests, channel) }()
	}

}

func handleRequests(reqs <-chan *gossh.Request, channel gossh.Channel) {
	// This for loop will run while the channel is open
	for req := range reqs {
		switch req.Type {
		default:
			// TODO: debug log
			e := req.Reply(true, nil)
			CustomUtils.CheckPrint(e)
			channel.Write([]byte("Hola Nena"))
		}
	}
}
