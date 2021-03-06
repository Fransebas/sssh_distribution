package SSH

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"sssh_server/Modules/DirectoryManager"
	"strings"

	"sssh_server/CustomUtils"
	"sssh_server/Modules/Authentication"
	"sssh_server/Modules/Logging"

	"sssh_server/Modules/SSH/SFTP"
)

type Handler func(s *SSHSession, w io.Writer, r io.Reader)
type AnyHandler func(msgType string, s *SSHSession, w io.Writer, r io.Reader)

type SSSHServer struct {
	config            *ssh.ServerConfig
	handlers          map[string]Handler
	AnyHandler        AnyHandler
	NewSessionHandler func(session *SSHSession)
	KeyPath           string // Path for the host key i.e. /etc/ssh/id_rsa or a custom path
	hasBeenInit       bool
	port              int
}

type User struct {
	// Later I have to add here the keys
	ID string
}

type SSHSession struct {
	ssh.Channel
	Conn     *ssh.ServerConn
	channels <-chan ssh.NewChannel
	reqs     <-chan *ssh.Request
}

func (session *SSHSession) GetSessionID() string {
	return hex.EncodeToString(session.Conn.SessionID())
}

func (server *SSSHServer) OnNewSession(f func(session *SSHSession)) {
	server.NewSessionHandler = f
}

// TODO: Change this!!!!!!!!
//const AUTHORIZED_KEYS_FILE = "/Users/fransebas/.ssh/authorized_keys"

func (server *SSSHServer) ReadAuthorizedKeys(username string) map[string]bool {
	dm := DirectoryManager.New(username)
	cpath := path.Join(dm.UserDirectory, ".ssh/authorized_keys")
	authorizedKeysMap := map[string]bool{}
	file, err := os.Open(cpath)
	if err != nil {
		// If no file no problem just no keys
		CustomUtils.CheckPrint(err)
		return map[string]bool{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "#") {
			pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(line))
			CustomUtils.CheckPrint(err)
			if err == nil && pubKey != nil {
				authorizedKeysMap[string(pubKey.Marshal())] = true
			}
		}
	}
	return authorizedKeysMap
}

func (server *SSSHServer) AddHostKeys(KeyPath string) {
	privateBytes, err := ioutil.ReadFile(KeyPath)
	if err != nil {
		log.Fatal("Failed to load private key, Do you have the right permissions?: ", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}
	server.config.AddHostKey(private)
}

func (server *SSSHServer) SetAnyHandler(f AnyHandler) {
	server.AnyHandler = f
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
			if val, _ := Authentication.ValidateUser(c.User(), string(pass)); val {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
		// Remove to disable public key auth.
		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {

			if !Authentication.UserExist(c.User()) {
				return nil, fmt.Errorf("invalid user %v", c.User())
			}

			authorizedKeysMap := server.ReadAuthorizedKeys(c.User())

			if val, ok := authorizedKeysMap[string(pubKey.Marshal())]; val && ok {
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

func (server *SSSHServer) InitServer(KeyPath string, port int) {
	server.port = port
	server.hasBeenInit = true

	if server.handlers == nil {
		server.handlers = make(map[string]Handler)
	}
	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	server.initAuthCallbacks()
	server.AddHostKeys(KeyPath)
}

func (server *SSSHServer) ListenAndServe() {
	if !server.hasBeenInit {
		panic("Server hasn't been initialized, call initServer(KeyPath string) first")
	}
	server.serve()
}

func (server *SSSHServer) serve() {
	// Once a ServerConfig has been configured, connections can be
	// accepted.
	url := fmt.Sprintf("%v:%v", "", server.port)
	listener, err := net.Listen("tcp", url)
	if err != nil {
		log.Fatal("failed to listen for connection: ", err)
	}

	for {

		// For now we will accept all incoming request and reject them later on based on credential or permissions
		nConn, err := listener.Accept()
		if err != nil {
			log.Fatal("failed to accept incoming connection: ", err)
		}

		go func() {
			// Before use, a handshake must be performed on the incoming
			// net.Conn.
			conn, chans, reqs, err := ssh.NewServerConn(nConn, server.config)
			if err != nil {
				//log.Fatal("failed to handshake: ", err)
				CustomUtils.CheckPrint(err)
				return
			}

			// Create a session for each incoming request
			session := SSHSession{
				Conn:     conn,
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
				go func() {
					defer func() {
						if err := recover(); err != nil {
							CustomUtils.CheckPrint(err.(error))
						}
					}()

					channel, requests, err := newChannel.Accept()
					if err != nil {
						log.Fatalf("Could not accept channel: %v", err)
					}

					server.AcceptRequests(requests, &channel, &session)
				}()
			}
		}()
	}
}

func (server *SSSHServer) AcceptRequests(in <-chan *ssh.Request, channel *ssh.Channel, session *SSHSession) {

	for req := range in {
		CustomUtils.Logger.Println(Logging.INFO, req.Type)
		switch req.Type {
		//case "session":
		//	session.
		case "subsystem":
			// handle ftp here
			if string(req.Payload[4:]) == "sftp" {
				// Requesting a SFTP
				_ = req.Reply(true, nil)
				go server.startSFTP(channel, session)
			}
		case "exec":
			// our custom protocol here
			if req.WantReply {
				_ = req.Reply(true, nil)
				msgType := string(req.Payload[4:])
				go func() { server.handleChannel(channel, session, msgType) }()
			}
		default:
			CustomUtils.CheckPrint(errors.New("ssh request type " + req.Type + "is not supported"))
		}
	}
}

func (s *SSSHServer) startSFTP(channel *ssh.Channel, session *SSHSession) {
	handlerServer := SFTP.New(session.Conn.User())

	handles := sftp.Handlers{
		FileGet:  handlerServer,
		FileList: handlerServer,
		FileCmd:  handlerServer,
		FilePut:  handlerServer,
	}

	server := sftp.NewRequestServer(*channel, handles)

	if err := server.Serve(); err == io.EOF {
		_ = server.Close()
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

	// I wonder if this will copy something by passing the value rather than the pointer, but Go wont let me pass the pointer as a io.Writer/Reader ...
	// bitch
	server.AnyHandler(msgType, session, *channel, *channel)
}

func ReadChannel(channel *ssh.Channel) string {
	data := make([]byte, 1024)
	// TODO: deal with err
	l, _ := (*channel).Read(data)
	return string(data[:l])
}
