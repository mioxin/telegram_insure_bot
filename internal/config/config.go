package config

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Config struct {
	Token   string
	LogFile string
}

func NewConfig(config_file string) (*Config, error) {
	var token, logFile string
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
		switch arr_str[0] {
		case "secure_token":
			token = arr_str[1]
		case "log-file":
			logFile = arr_str[1]
		default:
		}
	}
	return &Config{Token: token, LogFile: logFile}, nil
}
