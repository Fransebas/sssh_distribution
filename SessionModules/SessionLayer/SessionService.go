/*
This module is responsible for providing sessions and user for the SSSH server. Authentication is done at the SSH level.
*/
package SessionLayer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"sssh_server/CustomUtils"
	"sssh_server/Modules/SSH"
	"sssh_server/Modules/SSH/LimitlessChannel"
	"sssh_server/SessionModules/API"
	"strings"
	"time"
)

/*
Holds all the sessions that exist in the server and maps SSH IDs to Sessions ID's
*/
type SessionService struct {
	Server            *SSH.SSSHServer
	Sessions          map[string]*TerminalSession
	SSHidToTerminalID map[string]string
	KeyPath           string
	port              int
}

type PubKeyShare struct {
	Key      string
	Hash     string
	Mnemonic string
}

var HISTORY_FILE_NAME = "history"

// Creates a new SessionService
func Constructor(KeyPath string, port int) (s *SessionService) {
	s = new(SessionService)

	s.port = port
	s.KeyPath = KeyPath
	s.Server = &SSH.SSSHServer{}
	s.Sessions = make(map[string]*TerminalSession)
	s.SSHidToTerminalID = make(map[string]string)

	s.ChannelHandler()

	s.Server.OnNewSession(func(anonymousSession *SSH.SSHSession) {
		//
	})

	return s
}

// Return the pubKey associated with KeyPath
func (s *SessionService) GetPubKey() (*PubKeyShare, error) {
	key, e := ioutil.ReadFile(s.KeyPath + ".pub")
	if e != nil {
		return nil, e
	}
	hash, e := SSH.GetKeyHash(key)
	if e != nil {
		return nil, e
	}
	mnemonic, e := SSH.MakeMnemonic(key)
	if e != nil {
		return nil, e
	}

	pk := PubKeyShare{
		Key:      string(key),
		Hash:     base64.RawStdEncoding.EncodeToString(hash),
		Mnemonic: mnemonic,
	}
	return &pk, nil
}

// Returns the TerminalSession associated with the SSH.SSHSession
func (s *SessionService) SSHSessionToTerminalSession(sshSession *SSH.SSHSession) *TerminalSession {
	return s.Sessions[s.SSHidToTerminalID[sshSession.GetSessionID()]]
}

// Creates a new session when a client connects
func (s *SessionService) createSession(msgType string, sshSession *SSH.SSHSession, w io.Writer, r io.Reader) {
	// TODO: handle the null terminalSession number as a new terminalSession
	b, _ := CustomUtils.Read(r)
	sessionID := string(b)

	CustomUtils.Print("Terminal session started " + sessionID)

	var terminalSession *TerminalSession

	if s.Sessions[sessionID] != nil {
		// Terminal Session already exists, add sshSession to it
		terminalSession = s.Sessions[sessionID]
		terminalSession.SSHSessions[sshSession.GetSessionID()] = sshSession
	} else {
		// Terminal Session doesn't exist, create one
		terminalSession = newSession(sessionID, sshSession.Conn.User())
		s.Sessions[sessionID] = terminalSession
		terminalSession.SSHSessions[sshSession.GetSessionID()] = sshSession

		// Lifecycle hook
		terminalSession.OnNewSessionLifecycleHook()
	}
	terminalSession.LastConnected = time.Now().Unix()

	// Lifecycle hook
	terminalSession.OnNewConnectionLifecycleHook(sshSession)

	s.SSHidToTerminalID[sshSession.GetSessionID()] = sessionID

	bytesSession, _ := json.Marshal(terminalSession)
	_, _ = w.Write(bytesSession)
}

// Renames the session
func (s *SessionService) changeSessionName(name, id string) *TerminalSession {
	s.Sessions[id].Name = name
	return s.Sessions[id]
}

// Renames the session
func (s *SessionService) Close() {

}

// Deletes the session
func (s *SessionService) deleteSession(id string) *TerminalSession {
	var session = s.Sessions[id]
	session.Close()
	delete(s.Sessions, id)
	return session
}

func (s *SessionService) openSessions(w io.Writer, r io.Reader) {
	// Return all the open sessions for the client to choose
	var sessions = []*TerminalSession{}
	for _, session := range s.Sessions {
		sessions = append(sessions, session)
	}
	b, e := json.Marshal(sessions)
	CustomUtils.CheckPanic(e, "Unable to parse sessions, this should never happen")
	_, e = w.Write(b)
	CustomUtils.CheckPrint(e)
}

// Make the mapping for any kind of messages
// Here is where the multiplexing of channels happen
// Also if it's of type request, it will be passed through the limitless channel
// because by default each message has a default size but with this, you can send message the size you wanted to
func (s *SessionService) ChannelHandler() {
	s.Server.SetAnyHandler(func(fullMsgType string, sshSession *SSH.SSHSession, w io.Writer, r io.Reader) {
		parts := strings.Split(fullMsgType, "&")
		mode := parts[0]
		msgType := parts[1]
		limitlessChannel := LimitlessChannel.LimitlessChannel{
			Writer: w,
		}
		if msgType == "session" {
			// On the session message create a new session
			s.createSession(msgType, sshSession, limitlessChannel, r)
		} else if msgType == "open.sessions" {
			s.openSessions(limitlessChannel, r)
		} else if msgType == "session.name" {
			b, _ := CustomUtils.Read(r)
			var auxSession TerminalSession
			_ = json.Unmarshal(b, &auxSession)
			newSessions := s.changeSessionName(auxSession.Name, auxSession.ID)
			b2, _ := json.Marshal(newSessions)
			_, _ = limitlessChannel.Write(b2)
		} else if msgType == "session.delete" {
			b, _ := CustomUtils.Read(r)
			var auxSession TerminalSession
			_ = json.Unmarshal(b, &auxSession)
			s.deleteSession(auxSession.ID)
			s.openSessions(limitlessChannel, r)
		} else {
			terminalSession := s.SSHSessionToTerminalSession(sshSession)

			if mode == "request" {

				if handler, ok := terminalSession.HandlersMap[msgType]; ok {
					handler(limitlessChannel, r)
				} else {
					CustomUtils.CheckPrint(fmt.Errorf("Message type " + msgType + " doesn't exist"))
				}
			} else if mode == "channel" {
				if handler, ok := terminalSession.HandlersMap[msgType]; ok {
					handler(w, r)
				} else {
					CustomUtils.CheckPrint(fmt.Errorf("Message type " + msgType + " doesn't exist"))
				}
			}
		}
	})
}

// Serve the server
func (s *SessionService) Serve() {
	s.Server.InitServer(s.KeyPath, s.port)
	s.Server.ListenAndServe()
}

// Adds a new command to the recently used commands
func (ss *SessionService) AddCommand(data string, id string) {
	if _, ok := ss.Sessions[id]; ok {
		session := ss.Sessions[id]
		for _, service := range session.Modules {
			if historyService, ok := service.(API.HistoryService); ok {
				historyService.OnNewCommand(data)
			}
		}
	}
}

// Updates the variables
func (ss *SessionService) UpdateVariables(data string, id string) {
	if _, ok := ss.Sessions[id]; ok {
		session := ss.Sessions[id]
		for _, service := range session.Modules {
			if variablesService, ok := service.(API.VariablesService); ok {
				variablesService.OnUpdateVariables(data)
			}
		}
	}
}

// Updates pwd
func (ss *SessionService) UpdatePWD(path string, id string) {
	if _, ok := ss.Sessions[id]; ok {
		session := ss.Sessions[id]
		for _, service := range session.Modules {
			if pwdService, ok := service.(API.PWDService); ok {
				pwdService.OnPWD(path)
			}
		}
	}
}
