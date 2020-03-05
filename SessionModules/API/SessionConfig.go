package API

type SessionConfig struct {
	SessionID       string
	HistoryFilePath string
	BashrcFilePath  string
}

func (sc *SessionConfig) GetSessionID() string {
	return sc.SessionID
}

func (sc *SessionConfig) GetHistoryFilePath() string {
	return sc.HistoryFilePath
}

func (sc *SessionConfig) GetBashrcFilePath() string {
	return sc.BashrcFilePath
}
