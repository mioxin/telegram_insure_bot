package config

import (
	"bufio"
	"io"
	"log"
	"strings"
)

type Config struct {
	Token   string
	LogFile string
	AccWord string
	Deny    map[string]struct{}
	Allow   map[string]struct{}
}

func NewConfig(config_file io.Reader) (*Config, error) {
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
		arr_str := strings.Split(strings.ToLower(strings.TrimSpace(str)), " ")
		if len(arr_str) < 2 || arr_str[0][:2] == "//" || arr_str[0][:1] == "#" {
			continue
		}
		switch strings.TrimSpace(arr_str[0]) {
		case "secure_token":
			token = strings.TrimSpace(arr_str[1])
		case "log-file":
			logFile = strings.TrimSpace(arr_str[1])
		case "deny":
			arr_deny := strings.Split(arr_str[1], ",")
			for _, d := range arr_deny {
				deny[strings.TrimSpace(d)] = struct{}{}
			}
		case "allow":
			arr_allow := strings.Split(arr_str[1], ",")
			for _, d := range arr_allow {
				allow[strings.TrimSpace(d)] = struct{}{}
			}
		case "acc_word":
			acc_word = strings.TrimSpace(arr_str[1])

		default:
			log.Println("config: unrecognized param in the line: ", str)
		}
	}
	return &Config{Token: token, LogFile: logFile, Deny: deny, Allow: allow, AccWord: acc_word}, nil
}

func (conf *Config) IsAccess(user string) bool {
	user = strings.ToLower(user)
	res := false
	if _, ok := conf.Allow["all"]; ok && len(conf.Deny) == 0 {
		return ok
	} else {
		res = true
	}

	if _, ok := conf.Deny["all"]; ok {
		res = false
	}
	_, res = conf.Allow[user]
	if _, ok := conf.Deny[user]; ok {
		res = !ok
	}

	return res
}
