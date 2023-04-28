package session

import "time"

type Session struct {
	Handler     string
	Idx_request int
	LastTime    time.Time
	UserName    string
	ErrorInput  bool
}

func NewSession(user string) *Session {
	return &Session{UserName: user}
}

func (ses *Session) ResetSession() {
	ses.Idx_request = 0
	ses.Handler = ""
}
