package sessions

import "time"

type Session struct {
	ActionName         string    `json:"action_name"`
	IdxRequest         int       `json:"idx_request"`
	LastTime           time.Time `json:"last_time"`
	UserName           string    `json:"user_name"`
	LastRequestIsError bool
	LastMessageID      int `json:"last_message_id"`
	AccessCommand      map[string]struct{}
}

func NewSession(user string) *Session {
	return &Session{UserName: user, LastTime: time.Now(), AccessCommand: map[string]struct{}{}}
}

func (ses *Session) ResetSession() {
	ses.IdxRequest = 0
	ses.ActionName = ""
	ses.LastRequestIsError = false
	ses.LastMessageID = 0
}
