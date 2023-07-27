package clientfiles

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/mrmioxin/gak_telegram_bot/internal/storages"
	"github.com/mrmioxin/gak_telegram_bot/resources"
)

type MapStorage struct {
	//map[user] map[filename]FileInfo
	store map[string]map[string]*storages.FileInfo
	maxid int
	mtx   *sync.Mutex
}

func NewMapStorage() *MapStorage {
	//open file
	mapf, maxid, err := loadMap(resources.FILES_ID)
	if err != nil {
		log.Printf("error in NewMapStorage: the data not load from json, created empty map. %v\n", err)
		mapf = make(map[string]map[string]*storages.FileInfo)
	}
	return &MapStorage{mapf, maxid, new(sync.Mutex)}
}

func (mf *MapStorage) Close() {
	// save map to files_id.json
	fl, err := os.OpenFile(resources.FILES_ID, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("map Storage Close: error while open files_id.json. %v\n", err)
	}
	defer fl.Close()
	wr := bufio.NewWriter(fl)
	defer wr.Flush()

	jenc := json.NewEncoder(wr)

	for _, flMap := range mf.store {
		for _, file := range flMap {
			err := jenc.Encode(file)
			if err != nil {
				log.Printf("map Storage Close: error while save map to files_id.json. %v\n%#v", err, file)
			}
		}
	}
	log.Printf("End Close MapStorage.\n")
}

func (mf *MapStorage) GetFileId(user, ids string) (string, error) {
	mf.mtx.Lock()
	defer mf.mtx.Unlock()
	id, err := strconv.Atoi(ids)
	if err != nil {
		return "", fmt.Errorf("error in GetFileId: id \"%v\" not converted to int\n %#v", ids, err)
	}

	f, ok := mf.store[user]
	if ok {
		for _, fl := range f {
			if fl.Id == id {
				return fl.FileId, nil
			}
		}
	} else {
		log.Printf("GetFileId: %v, %v\n %#v.\n", id, user, f)
		return "", fmt.Errorf("error in GetFileId: file_id \"%v\" in the map %#v not found", id, f)
	}
	// if f, ok := mf.store[user][name]; ok {
	// 	return f.FileId, nil
	// } else {
	// 	log.Printf("GetFileId: %v, %v\n %#v.\n", name, user, mf.store)
	// 	return "", fmt.Errorf("error in GetFileId: file_id \"%v\" in the map not found", name)
	// }
	return "", fmt.Errorf("error in GetFileId: file_id \"%v\" in the map %#v not found", id, f)
}

func (mf *MapStorage) SetFileId(name, user, id string) error {
	mf.maxid++
	f := storages.FileInfo{Id: mf.maxid, Time: time.Now(), UserName: user, FileName: name, FileId: id}
	//save to map
	mf.mtx.Lock()
	defer mf.mtx.Unlock()
	if _, ok := mf.store[user]; !ok {
		mf.store[user] = make(map[string]*storages.FileInfo)
	}

	mf.store[user][name] = &f

	return nil
}

func (mf *MapStorage) ListFiles(user string) []*storages.FileInfo {
	slFiles := make([]*storages.FileInfo, 0)
	for _, f := range mf.store[user] {
		slFiles = append(slFiles, f) //.FileName+"; "+f.Time.Format("2/1/2006 15:16:17"))
	}
	return slFiles
}

func (mf *MapStorage) ListUsers() []string {
	slUser := make([]string, 0)
	for k, _ := range mf.store {
		slUser = append(slUser, k)
	}
	return slUser
}

func loadMap(filesIdJson string) (map[string]map[string]*storages.FileInfo, int, error) {
	//var finfo *storages.FileInfo
	m := make(map[string]map[string]*storages.FileInfo)
	maxid := 0
	fl, err := os.OpenFile(filesIdJson, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, maxid, err
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
			log.Printf("LoadMap: error while load map from files_id.json.\n%#v", err)
		} else {
			if maxid < finfo.Id {
				maxid = finfo.Id
			}
			if _, ok := m[finfo.UserName]; !ok {
				m[finfo.UserName] = make(map[string]*storages.FileInfo)
			}
			m[finfo.UserName][finfo.FileName] = finfo
		}
	}
	return m, maxid, nil
}
