package config

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	Token   string
	LogFile string
	AccWord string
	Deny    map[string]struct{}
	Allow   map[string]struct{}
	ModTime time.Time
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

	return config, err
}

func ParseConfigFile(config_file io.Reader) (*Config, error) {
	var token, logFile, acc_word string
	deny := make(map[string]struct{})
	allow := make(map[string]struct{})
	in := bufio.NewReader(config_file)
	for {
		str, err := in.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		arr_str := strings.Split(strings.TrimSpace(str), "=")
		if len(arr_str) < 2 || arr_str[0][:2] == "//" || arr_str[0][:1] == "#" {
			continue
		}
		switch strings.TrimSpace(arr_str[0]) {
		case "secure_token":
			token = strings.TrimSpace(arr_str[1])
			if token == "" {
				return nil, fmt.Errorf("error ParseConfigFile: in the config file secure token expected")
			}
		case "log-file":
			logFile = strings.TrimSpace(arr_str[1])
		case "deny":
			arr_deny := strings.Split(strings.ToLower(arr_str[1]), ",")
			for _, d := range arr_deny {
				deny[strings.TrimSpace(d)] = struct{}{}
			}
		case "allow":
			arr_allow := strings.Split(strings.ToLower(arr_str[1]), ",")
			for _, d := range arr_allow {
				allow[strings.TrimSpace(d)] = struct{}{}
			}
		case "acc_word":
			acc_word = strings.TrimSpace(arr_str[1])

		default:
			log.Println("ParseConfigFile: unrecognized param in the line: ", str)
		}
	}
	if logFile == "" {
		logFile = "bot.log"
	}

	return &Config{Token: token, LogFile: logFile, Deny: deny, Allow: allow, AccWord: acc_word}, nil
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

func (conf *Config) Watch(configFile string, watchTime time.Duration, ok chan any) {
	var isModify bool
	for {
		if flInfo, err := os.Stat(configFile); err != nil {
			log.Printf("error in conf.Watch: error get FileInfo for  %v.", configFile)

			time.Sleep(watchTime)
			continue
		} else {
			isModify = flInfo.ModTime() != conf.ModTime
		}

		if isModify {
			conf, _ = NewConfig(configFile)
			log.Printf("Config file %v was modify. Config is updated.", configFile)
			ok <- []string{}
		}
		time.Sleep(watchTime)
	}
}

func (conf *Config) Close() {
}
