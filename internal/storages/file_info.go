package storages

import "time"

type FileInfo struct {
	Id       int       `json:"id"`
	Time     time.Time `json:"time,omitempty"`
	UserName string    `json:"user_name,omitempty"`
	FileName string    `json:"name"`
	FileId   string    `json:"fileid"`
}
