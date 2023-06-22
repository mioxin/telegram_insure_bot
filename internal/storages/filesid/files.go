package filesid

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mrmioxin/gak_telegram_bot/resources"
)

type File struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}
type mapFilesId struct {
	store map[string]string
	//file  *os.File
}

func NewMapFilesId() *mapFilesId {
	//open file
	mapf, err := loadMap(resources.FILES_ID)
	if err != nil {
		log.Printf("error in NewMapFilesId: %v\n", err)
		mapf = make(map[string]string)
	}
	return &mapFilesId{mapf}
}

func (mf *mapFilesId) GetFileId(name string) (string, error) {
	if id, ok := mf.store[name]; ok {
		return id, nil
	} else {
		return "", fmt.Errorf("error in GetFileId: file_id \"%v\" in the map not found\n", name)
	}
}

func (mf *mapFilesId) SetFileId(name, id string) {
	mf.store[name] = id
	fl, err := os.OpenFile(resources.FILES_ID, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Printf("error in SetFileId: can not open file \"%v\"\n", resources.FILES_ID)
		return
	}
	defer fl.Close()

	jsonFinfo, err := json.Marshal(File{name, id})
	if err != nil {
		log.Printf("error in SetFileId: can not marshal %v\n", File{name, id})
		return

	}
	_, err = fl.WriteString(string(jsonFinfo))
	if err != nil {
		log.Printf("error in SetFileId: can not write the file info to %v\n", resources.FILES_ID)
		return
	}
}

func loadMap(fileName string) (map[string]string, error) {
	var finfo File
	m := make(map[string]string)
	fl, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	defer fl.Close()

	jdec := json.NewDecoder(bufio.NewReader(fl))

	for {
		err := jdec.Decode(&finfo)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("error in loadMap: %v\n", err)
		} else {
			m[finfo.Name] = finfo.Id
		}
	}
	return m, nil
}

// func (mf *mapFilesId) Close() {
// 	if mf.file != nil {
// 		mf.file.Close()
// 	}
// }
