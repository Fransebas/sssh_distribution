package SSH

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sssh_server/Services/SSH/LimitlessChannel"
	"strings"
)

type Handler func(s *SSHSession, w io.Writer, r io.Reader)

type SSSHServer struct {
	config            *ssh.ServerConfig
	handlers          map[string]Handler
	NewSessionHandler func(session *SSHSession)
	authorizedKeysMap map[string]bool
}

type SSHSession struct {
	ssh.Channel
	conn     *ssh.ServerConn
	channels <-chan ssh.NewChannel
	reqs     <-chan *ssh.Request
}

func (session *SSHSession) GetSessionID() string {
	return hex.EncodeToString(session.conn.SessionID())
}

func (server *SSSHServer) OnNewSession(f func(session *SSHSession)) {
	server.NewSessionHandler = f
}

const AUTHORIZED_KEYS_FILE = "/Users/fransebas/.ssh/authorized_keys"

func (server *SSSHServer) ReadAuthorizedKeys(path string) {
	server.authorizedKeysMap = map[string]bool{}
	file, err := os.Open(path)
	if err != nil {
		// If no file no problem just no keys
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "#") {
			pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(line))
			if err != nil {
				log.Fatal(err)
			}
			server.authorizedKeysMap[string(pubKey.Marshal())] = true
		}
	}

}

func (server *SSSHServer) AddHostKeys() {
	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key: ", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}
	server.config.AddHostKey(private)
}

func (server *SSSHServer) HandleFunc(msgType string, handler Handler) {
	if server.handlers == nil {
		server.handlers = make(map[string]Handler)
	}
	server.handlers[msgType] = handler
}

func (server *SSSHServer) initAuthCallbacks() {
	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	server.config = &ssh.ServerConfig{
		// Remove to disable password auth.
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			// Should use constant-time compare (or better, salt+hash) in
			// a production setting.
			fmt.Println("Password auth")
			if c.User() == "testuser" && string(pass) == "tiger" {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
		// Remove to disable public key auth.
		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			fmt.Println("Key auth")
			if server.authorizedKeysMap[string(pubKey.Marshal())] {
				return &ssh.Permissions{
					// Record the public key used for authentication.
					Extensions: map[string]string{
						"pubkey-fp": ssh.FingerprintSHA256(pubKey),
					},
				}, nil
			}
			return nil, fmt.Errorf("unknown public key for %q", c.User())
		},
	}
}

func (server *SSSHServer) initServer() {
	server.ReadAuthorizedKeys(AUTHORIZED_KEYS_FILE)
	if server.handlers == nil {
		server.handlers = make(map[string]Handler)
	}
	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	server.initAuthCallbacks()
	server.AddHostKeys()
}

func (server *SSSHServer) ListenAndServe() {
	server.initServer()
	server.serve()
}

func (server *SSSHServer) serve() {
	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", "localhost:2222")
	if err != nil {
		log.Fatal("failed to listen for connection: ", err)
	}

	for {

		// For now we will accept all incoming request and reject them later on based on credential or permissions
		nConn, err := listener.Accept()
		if err != nil {
			log.Fatal("failed to accept incoming connection: ", err)
		}

		// Before use, a handshake must be performed on the incoming
		// net.Conn.
		conn, chans, reqs, err := ssh.NewServerConn(nConn, server.config)
		if err != nil {
			log.Fatal("failed to handshake: ", err)
		}

		// Create a session for each incoming request
		session := SSHSession{
			conn:     conn,
			channels: chans,
			reqs:     reqs,
		}

		server.NewSessionHandler(&session)

		// The incoming Request channel must be serviced.
		// Dont know what the fuck does that means but I'm going to accept all request
		// go AcceptRequests(reqs)

		// Service the incoming Channel channel.
		for newChannel := range chans {
			// Channels have a type, depending on the application level
			// I don't really care about any of that because the client cant issue new channel types
			// that's why I'm going to implement my own multiplex on top of a channel

			// Accept all channels, we don't care about their type,
			// Could be bananas for all I know
			channel, requests, err := newChannel.Accept()
			if err != nil {
				log.Fatalf("Could not accept channel: %v", err)
			}

			go server.AcceptRequests(requests, &channel, &session)

		}
		fmt.Println("hello")
	}
}

func (server *SSSHServer) AcceptRequests(in <-chan *ssh.Request, channel *ssh.Channel, session *SSHSession) {
	for req := range in {
		switch req.Type {
		case "subsystem":
			// handle ftp here
			if string(req.Payload[4:]) == "sftp" {
				// Requesting a SFTP
				_ = req.Reply(true, nil)
				go server.startSFTP(channel)
			}
		default:
			// our custom protocol here
			if req.WantReply {
				_ = req.Reply(true, nil)
				msgType := string(req.Payload[4:])
				go func() { server.handleChannel(channel, session, msgType) }()
			}

		}
	}
}

func (s *SSSHServer) startSFTP(channel *ssh.Channel) {
	serverOptions := []sftp.ServerOption{
		sftp.WithDebug(os.Stderr),
	}
	server, err := sftp.NewServer(
		*channel,
		serverOptions...,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Serve(); err == io.EOF {
		server.Close()
		log.Print("sftp client exited session.")
	} else if err != nil {
		log.Fatal("sftp server completed with error:", err)
	}
}

func (server *SSSHServer) handleChannel(channel *ssh.Channel, session *SSHSession, msgType string) {
	// We don't care of the channel type, we will read the data
	// and based on the first line we will multiplex the connection
	// The first message should signal the type/mod they want to use

	// msgType := ReadChannel(channel)

	// TODO: deal with err
	// _, _ = (*channel).Write([]byte("ack")) // Send the string ack signaling that the other side can start sending and receiving

	f := server.handlers[msgType]

	// I wonder if this will copy something by passing the value rather than the pointer, but Go wont let me pass the pointer as a io.Writer/Reader ...
	// bitch
	limitlessChannel := LimitlessChannel.LimitlessChannel{
		Writer: *channel,
	}
	if f != nil {
		f(session, limitlessChannel, *channel)
	} else {
		fmt.Println("No handler for this message " + msgType)
	}

}

func ReadChannel(channel *ssh.Channel) string {
	data := make([]byte, 1024)
	// TODO: deal with err
	l, _ := (*channel).Read(data)
	return string(data[:l])
}
