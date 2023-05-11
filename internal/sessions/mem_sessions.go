package sessions

import "fmt"

type memSessions map[int64]Session

func NewMemSessions() *memSessions {
	s := memSessions(make(map[int64]Session))
	return &s
}

func (mses memSessions) GetSession(id int64) (*Session, error) {
	if ses, ok := (mses)[id]; ok {
		return &ses, nil
	} else {
		return nil, fmt.Errorf("error in getSession: session id=%v not found", id)
	}
}

func (mses memSessions) UpdateSession(id int64, ses *Session) error {
	var err error
	if _, ok := mses[id]; !ok {
		err = fmt.Errorf("error in updateSession: session id=%v not found. created new session", id)
	}
	mses[id] = *ses
	return err

}

func (mses memSessions) AddSession(id int64, ses *Session) {
	mses[id] = *ses
}

func (mses memSessions) GetIdsByUser(user string) []int64 {
	aId := make([]int64, 0)
	for k, v := range mses {
		if v.UserName == user {
			aId = append(aId, k)
		}
	}
	return aId
}
