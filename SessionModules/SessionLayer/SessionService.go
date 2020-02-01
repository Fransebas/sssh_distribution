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
	"sssh_server/SessionModules/API"
	"sssh_server/SessionModules/SSH"
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
}

type PubKeyShare struct {
	Key      string
	Hash     string
	Mnemonic string
}

var HISTORY_FILE_NAME = "history"

// Creates a new SessionService
func Constructor(KeyPath string) (s *SessionService) {
	s = new(SessionService)

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

	fmt.Println("Terminal session started " + sessionID)

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
func (s *SessionService) changeSessionName(name, id string) {
	s.Sessions[id].Name = name
}

// Make the mapping for any kind of messages
// Here is where the multiplexing of channels happen
func (s *SessionService) ChannelHandler() {
	s.Server.SetAnyHandler(func(msgType string, sshSession *SSH.SSHSession, w io.Writer, r io.Reader) {

		if msgType == "session" {
			// On the session message create a new session
			s.createSession(msgType, sshSession, w, r)
		} else if msgType == "open.sessions" {
			// Return all the open sessions for the client to choose
			var sessions = []*TerminalSession{}
			for _, session := range s.Sessions {
				sessions = append(sessions, session)
			}
			b, e := json.Marshal(sessions)
			CustomUtils.CheckPanic(e, "Unable to parse sessions, this should never happen")
			_, e = w.Write(b)
			CustomUtils.CheckPrint(e)
		} else if msgType == "session.name" {
			b, _ := CustomUtils.Read(r)
			var auxSession TerminalSession
			_ = json.Unmarshal(b, &auxSession)
			s.changeSessionName(auxSession.Name, auxSession.ID)

		} else {
			terminalSession := s.SSHSessionToTerminalSession(sshSession)
			if handler, ok := terminalSession.HandlersMap[msgType]; ok {
				handler(w, r)
			} else {
				CustomUtils.CheckPrint(fmt.Errorf("Message type " + msgType + " doesn't exist"))
			}
		}
	})
}

// Serve the server
func (s *SessionService) Serve() {
	s.Server.InitServer(s.KeyPath)
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
