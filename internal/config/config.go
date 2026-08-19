package config

import (
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	Token    string
	Company  string
	LogFile  string
	ConfFile string
	AccWord  string
	Deny     map[string]struct{}
	Allow    map[string]struct{}
	Admins   map[string]struct{}
	ModTime  time.Time
}

func NewConfig(configName string) (*Config, error) {
	values, err := loadConfig(configName)
	if err != nil {
		return nil, err
	}

	if flInfo, err := os.Stat(configName); err != nil {
		return nil, err
	} else {
		config := newFromKong(values, configName)
		config.ModTime = flInfo.ModTime()
		return config, nil
	}
}
func str_map(s string) map[string]struct{} {
	m := make(map[string]struct{})
	if strings.TrimSpace(s) != "" {
		arr_str := strings.Split(strings.ToLower(s), ",")
		for _, d := range arr_str {
			m[strings.TrimSpace(d)] = struct{}{}
		}
	}
	return m
}

func (conf *Config) Update() {
	values, err := loadConfig(conf.ConfFile)
	if err != nil {
		log.Printf("error in Update config: %v", err)
		return
	}
	loaded := newFromKong(values, conf.ConfFile)
	*conf = *loaded

	if flInfo, err := os.Stat(conf.ConfFile); err != nil {
		log.Printf("error in Update config: %***REMOVED***v\n%***REMOVED***v\n", err, conf)
	} else {
		conf.ModTime = flInfo.ModTime()
	}

}

func (conf *Config) IsAdmin(user string) bool {
	_, ok := conf.Admins[strings.ToLower(user)]
	log.Println("Is admin:", user, ok)
	return ok
}

func (conf *Config) IsAccess(user string) bool {
	user = strings.ToLower(user)
	res := false
	if len(conf.Allow) == 0 && len(conf.Deny) == 0 {
		return true
	}
	_, okall := conf.Allow["all"]
	_, noall := conf.Deny["all"]
	if okall && len(conf.Deny) == 0 {
		return okall
	}

	if noall && okall {
		return false
	}

	if _, ok := conf.Allow[user]; ok || okall {
		res = true
	}
	if _, ok := conf.Deny[user]; ok {
		res = !ok
	}
	return res
}

// func (conf *Config) Watch(configFile string, watchTime time.Duration) chan any {
// 	var isModify bool
// 	ok := make(chan any)
// 	go func() {
// 		for {
// 			if flInfo, err := os.Stat(configFile); err != nil {
// 				log.Printf("error in conf.Watch: error get FileInfo for  %v.", configFile)

// 				time.Sleep(watchTime)
// 				continue
// 			} else {
// 				isModify = flInfo.ModTime() != conf.ModTime
// 			}

// 			if isModify {
// 				// c, _ := NewConfig(configFile)
// 				// log.Printf("Config file before %v \n%***REMOVED***v", &conf, conf)
// 				log.Printf("Config file %v was modify.", configFile)
// 				ok <- struct{}{}
// 			}
// 			time.Sleep(watchTime)
// 		}
// 	}()
// 	return ok
// }

func (conf *Config) Close() {
	log.Printf("End Close Config.\n")

}

func (conf *Config) IsAccWord(str string) bool {
	return conf.AccWord == str
}
