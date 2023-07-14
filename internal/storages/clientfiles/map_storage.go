package clientfiles

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mrmioxin/gak_telegram_bot/internal/storages"
	"github.com/mrmioxin/gak_telegram_bot/resources"
)

type MapStorage struct {
	store map[string]*storages.FileInfo
	mtx   *sync.Mutex
}

func NewMapStorage() *MapStorage {
	//open file
	mapf, err := loadMap(resources.FILES_ID)
	if err != nil {
		log.Printf("error in NewMapStorage: the dat anot load from json, created empty map. %v\n", err)
		mapf = make(map[string]*storages.FileInfo)
	}
	return &MapStorage{mapf, new(sync.Mutex)}
}

func (mf *MapStorage) Close() {
	// save map to files_id.json
	fl, err := os.OpenFile(resources.FILES_ID, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("map Storage Close: error while open files_id.json. Error %v\n", err)
	}
	defer fl.Close()

	jenc := json.NewEncoder(bufio.NewWriter(fl))

	for _, v := range mf.store {
		err := jenc.Encode(v)
		if err != nil {
			log.Printf("map Storage Close: error while save map to files_id.json. Error %v\n%#v", err, v)
		}
	}
	log.Printf("End Close MapStorage.\n")
}

func (mf *MapStorage) GetFileId(name string) (string, error) {
	mf.mtx.Lock()
	defer mf.mtx.Unlock()
	if f, ok := mf.store[name]; ok {
		return f.FileId, nil
	} else {
		return "", fmt.Errorf("error in GetFileId: file_id \"%v\" in the map not found", name)
	}
}

func (mf *MapStorage) SetFileId(name, user, id string) error {
	f := storages.FileInfo{Time: time.Now(), UserName: user, FileName: name, FileId: id}
	//save to map
	mf.mtx.Lock()
	defer mf.mtx.Unlock()
	mf.store[filepath.Join(user, name)] = &f

	//save to files_id.json
	// fl, err := os.OpenFile(resources.FILES_ID, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	// if err != nil {
	// 	log.Printf("error in SetFileId: can not open file \"%v\"\n", resources.FILES_ID)
	// 	return err
	// }
	// defer fl.Close()

	// jsonFinfo, err := json.Marshal(f)
	// if err != nil {
	// 	log.Printf("error in SetFileId: can not marshal %v. %v\n", f, err)
	// 	return err
	// }
	// _, err = fl.WriteString(string(jsonFinfo) + "\n")
	// if err != nil {
	// 	log.Printf("error in SetFileId: can not write the file info to %v. %v\n", resources.FILES_ID, err)
	// 	return err
	// }
	//log.Printf("SetFileId: Map %#v\n", mf.store)

	return nil
}

func loadMap(filesIdJson string) (map[string]*storages.FileInfo, error) {
	//var finfo *storages.FileInfo
	m := make(map[string]*storages.FileInfo)
	fl, err := os.OpenFile(filesIdJson, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	defer fl.Close()

	jdec := json.NewDecoder(bufio.NewReader(fl))

	for {
		finfo := new(storages.FileInfo)
		err := jdec.Decode(finfo)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("LoadMap: error while load map from files_id.json. Error %v\n%#v", err, finfo)
		} else {
			//key is a path ./username/filename
			fileName := filepath.Join(finfo.UserName, finfo.FileName)
			m[fileName] = finfo
		}
	}
	return m, nil
}
