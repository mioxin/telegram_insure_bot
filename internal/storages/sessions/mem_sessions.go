package sessions

import (
	"fmt"
	"log"
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
	defer mses.mtx.Unlock()
	if _, ok := mses.sessions[id]; !ok {
		return fmt.Errorf("error UpdateSession: session id=%v not found. created new session", id)
	}
	if ses == nil {
		return fmt.Errorf("error in updateSession: session id=%v is <nil>", id)
	}
	mses.sessions[id] = *ses
	return err

}

func (mses *MemSessions) AddSession(id int64, ses *Session) {
	mses.mtx.Lock()
	defer mses.mtx.Unlock()
	mses.sessions[id] = *ses
}

func (mses *MemSessions) GetIdsByUser(user string) []int64 {
	aId := make([]int64, 0)
	mses.mtx.Lock()
	defer mses.mtx.Unlock()

	for k, v := range mses.sessions {
		if v.UserName == user {
			aId = append(aId, k)
		}
	}
	return aId
}

func (mses *MemSessions) Close() {
	log.Printf("End Close MemSessions.\n")

}
