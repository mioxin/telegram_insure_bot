package storages

import "time"

type FileInfo struct {
	Time     time.Time `json:"time,omitempty"`
	UserName string    `json:"user_name,omitempty"`
	FileName string    `json:"name"`
	FileId   string    `json:"id"`
}
