package sessions

import "fmt"

type MemSessions map[int64]Session

func NewMemSessions() *MemSessions {
	s := MemSessions(make(map[int64]Session))
	return &s
}

func (mses MemSessions) GetSession(id int64) (*Session, error) {
	if ses, ok := (mses)[id]; ok {
		return &ses, nil
	} else {
		return nil, fmt.Errorf("error in getSession: session id=%v not found", id)
	}
}

func (mses MemSessions) UpdateSession(id int64, ses *Session) error {
	var err error
	if _, ok := mses[id]; !ok {
		err = fmt.Errorf("error in updateSession: session id=%v not found. created new session", id)
	}
	(mses)[id] = *ses
	return err

}
