package session

import "time"

type Session struct {
	ActionName  string
	Idx_request int
	LastTime    time.Time
	UserName    string
}

func NewSession(user string) *Session {
	return &Session{UserName: user}
}

func (ses *Session) ResetSession() {
	ses.Idx_request = 0
	ses.ActionName = ""
}
