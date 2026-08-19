package config

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mrmioxin/gak_telegram_bot/resources"
)

type Config struct {
	Token    string
	LogFile  string
	ConfFile string
	AccWord  string
	Deny     map[string]struct{}
	Allow    map[string]struct{}
	Admins   map[string]struct{}
	ModTime  time.Time
}

func NewConfig(configName string) (*Config, error) {
	configFile, err := os.ReadFile(configName)
	if err != nil {
		return nil, err
	}

	config, err := ParseConfigFile(strings.NewReader(string(configFile)))
	if err != nil {
		return nil, err
	}

	if flInfo, err := os.Stat(configName); err != nil {
		return nil, err
	} else {
		config.ModTime = flInfo.ModTime()
	}

	config.ConfFile = configName

	return &config, err
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

func ParseConfigFile(config_file io.Reader) (Config, error) {
	var token, logFile, acc_word string
	deny := make(map[string]struct{})
	allow := make(map[string]struct{})
	adm := make(map[string]struct{})
	config := Config{}

	in := bufio.NewReader(config_file)
	for {
		str, err := in.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return config, err
			}
		}
		line := strings.TrimSpace(str)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "***REMOVED***") {
			if err == io.EOF {
				break
			}
			continue
		}
		arr_str := strings.SplitN(line, "=", 2)
		if len(arr_str) < 2 {
			if err == io.EOF {
				break
			}
			continue
		}
		switch strings.TrimSpace(arr_str[0]) {
		case "secure_token":
			token = strings.TrimSpace(arr_str[1])
			if token == "" {
				return config, fmt.Errorf("error ParseConfigFile: in the config file secure token expected")
			}
		case "log_file":
			logFile = strings.TrimSpace(arr_str[1])
		case "deny":
			deny = str_map(strings.ToLower(arr_str[1]))
		case "allow":
			allow = str_map(strings.ToLower(arr_str[1]))
		case "adm_group":
			adm = str_map(strings.ToLower(arr_str[1]))
		case "acc_word":
			acc_word = strings.TrimSpace(arr_str[1])

		default:
			log.Println("ParseConfigFile: unrecognized param in the line: ", str)
		}
		if err == io.EOF {
			break
		}
	}
	if logFile == "" {
		logFile = resources.DEFAULT_LOG_FILE
	}

	return Config{Token: token, LogFile: logFile, Deny: deny, Allow: allow, Admins: adm, AccWord: acc_word}, nil
}

func (conf *Config) Update() {
	configFile, err := os.ReadFile(conf.ConfFile)
	if err != nil {
		log.Printf("error in Update config: %v", err)
		return
	}
	c, err := ParseConfigFile(strings.NewReader(string(configFile)))
	if err != nil {
		log.Printf("error in Update config: %v", err)
		return
	}
	c.ConfFile = conf.ConfFile
	*conf = c

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
