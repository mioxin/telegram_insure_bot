package sessions

import (
	"fmt"
	"sync"
)

type MemSessions struct {
	sessions map[int64]Session
	mtx      *sync.Mutex
}

func NewMemSessions() *MemSessions {
	s := make(map[int64]Session)
	return &MemSessions{s, new(sync.Mutex)}
}

func (mses *MemSessions) GetSession(id int64) (*Session, error) {
	if ses, ok := mses.sessions[id]; ok {
		return &ses, nil
	} else {
		return nil, fmt.Errorf("error in getSession: session id=%v not found", id)
	}
}

func (mses *MemSessions) UpdateSession(id int64, ses *Session) error {
	var err error
	mses.mtx.Lock()
	if _, ok := mses.sessions[id]; !ok {
		err = fmt.Errorf("error in updateSession: session id=%v not found. created new session", id)
	}
	mses.sessions[id] = *ses
	mses.mtx.Unlock()
	return err

}

func (mses *MemSessions) AddSession(id int64, ses *Session) {
	mses.mtx.Lock()
	mses.sessions[id] = *ses
	mses.mtx.Unlock()
}

func (mses MemSessions) GetIdsByUser(user string) []int64 {
	aId := make([]int64, 0)
	mses.mtx.Lock()

	for k, v := range mses.sessions {
		if v.UserName == user {
			aId = append(aId, k)
		}
	}
	mses.mtx.Unlock()
	return aId
}
