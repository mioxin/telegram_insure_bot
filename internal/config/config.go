package config

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

type Config struct {
	Token   string
	LogFile string
	AccWord string
	Deny    map[string]struct{}
	Allow   map[string]struct{}
}

func NewConfig(config_file string) (*Config, error) {
	var token, logFile, acc_word string
	deny := make(map[string]struct{})
	allow := make(map[string]struct{})
	conf, err := os.ReadFile(config_file)
	if err != nil {
		return nil, err
	}
	in := bufio.NewReader(strings.NewReader(string(conf)))
	for {
		str, err := in.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		arr_str := strings.Split(strings.TrimSpace(str), " ")
		if len(arr_str) < 2 || arr_str[0][:2] == "//" || arr_str[0][:1] == "#" {
			continue
		}
		switch strings.TrimSpace(arr_str[0]) {
		case "secure_token":
			token = strings.TrimSpace(arr_str[1])
		case "log-file":
			logFile = strings.TrimSpace(arr_str[1])
		case "deny":
			deny[arr_str[1]] = struct{}{}
		case "allow":
			allow[arr_str[1]] = struct{}{}
		case "acc_word":
			acc_word = strings.TrimSpace(arr_str[1])

		default:
			log.Println("config: unrecognized param in the line: ", str)
		}
	}
	return &Config{Token: token, LogFile: logFile, Deny: deny, Allow: allow, AccWord: acc_word}, nil
}
