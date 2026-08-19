package config

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/mrmioxin/gak_telegram_bot/resources"
)

type rawConfig struct {
	Token   string `name:"token" short:"t" env:"TELEGRAM_BOT_TOKEN"`
	Company string `name:"company" env:"COMPANY"`
	LogFile string `name:"log-file" env:"LOG_FILE"`
	AccWord string `name:"acc-word" env:"ACC_WORD"`
	Deny    string `name:"deny" env:"DENY"`
	Allow   string `name:"allow" env:"ALLOW"`
	Admins  string `name:"adm-group" env:"ADM_GROUP"`
}

func newFromKong(values rawConfig, configName string) *Config {
	company := strings.TrimSpace(values.Company)
	if company == "" {
		company = resources.COMPANY
	}
	logFile := strings.TrimSpace(values.LogFile)
	if logFile == "" {
		logFile = resources.DEFAULT_LOG_FILE
	}
	conf := &Config{
		Token:    strings.TrimSpace(values.Token),
		Company:  company,
		LogFile:  logFile,
		ConfFile: configName,
		AccWord:  strings.TrimSpace(values.AccWord),
		Deny:     str_map(values.Deny),
		Allow:    str_map(values.Allow),
		Admins:   str_map(values.Admins),
	}
	resources.COMPANY = conf.Company
	return conf
}

func kongConfigLoader(reader io.Reader) (kong.Resolver, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	trimmed := bytes.TrimSpace(content)
	if len(trimmed) == 0 {
		return kong.ResolverFunc(func(*kong.Context, *kong.Path, *kong.Flag) (any, error) {
			return nil, nil
		}), nil
	}
	withoutComments := bytes.TrimSpace(stripConfigComments(trimmed))
	if len(withoutComments) > 0 && withoutComments[0] == '{' {
		return kong.JSON(bytes.NewReader(withoutComments))
	}
	return dotenvResolver(bytes.NewReader(trimmed))
}

func stripConfigComments(content []byte) []byte {
	lines := bytes.Split(content, []byte("\n"))
	filtered := make([][]byte, 0, len(lines))
	for _, line := range lines {
		trimmed := bytes.TrimSpace(line)
		if bytes.HasPrefix(trimmed, []byte("***REMOVED***")) || bytes.HasPrefix(trimmed, []byte("//")) {
			continue
		}
		filtered = append(filtered, line)
	}
	return bytes.Join(filtered, []byte("\n"))
}

func dotenvResolver(reader io.Reader) (kong.Resolver, error) {
	values := make(map[string]string)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "***REMOVED***") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, value, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(key) == "" {
			return nil, fmt.Errorf("invalid .env line %q", line)
		}
		value = strings.TrimSpace(value)
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}
		values[strings.TrimSpace(key)] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return kong.ResolverFunc(func(_ *kong.Context, _ *kong.Path, flag *kong.Flag) (any, error) {
		for _, envName := range flag.Envs {
			if value, ok := values[envName]; ok {
				return value, nil
			}
		}
		return nil, nil
	}), nil
}

func loadConfig(configName string) (rawConfig, error) {
	return loadWithKong(configName, ".env")
}

func loadWithKong(configName, envFile string) (rawConfig, error) {
	values := rawConfig{}
	parser, err := kong.New(&values, kong.Configuration(kongConfigLoader, configName, envFile))
	if err != nil {
		return values, err
	}
	_, err = parser.Parse(nil)
	return values, err
}

type LoadOptions struct {
	Verbose bool
	Debug   bool
}

type cliConfig struct {
	rawConfig
	ConfigFile string `name:"config" short:"c" type:"path" default:"bot.cfg" help:"Configuration JSON file."`
	Verbose    bool   `name:"verbose" short:"v" help:"Write logs to stdout."`
	Debug      bool   `name:"debug" short:"d" help:"Enable Telegram debug mode."`
}

func Load(args []string) (*Config, LoadOptions, error) {
	first, err := parseCLI(args, false, "")
	if err != nil {
		return nil, LoadOptions{}, err
	}
	second, err := parseCLI(args, true, first.ConfigFile)
	if err != nil {
		return nil, LoadOptions{}, err
	}
	if strings.TrimSpace(second.Token) == "" {
		return nil, LoadOptions{}, fmt.Errorf("Telegram token is required; use -t, TELEGRAM_BOT_TOKEN, or .env")
	}
	return newFromKong(second.rawConfig, second.ConfigFile), LoadOptions{Verbose: second.Verbose, Debug: second.Debug}, nil
}

func parseCLI(args []string, withConfig bool, configFile string) (cliConfig, error) {
	cli := cliConfig{}
	options := []kong.Option{}
	if withConfig {
		options = append(options, kong.Configuration(kongConfigLoader, configFile, ".env"))
	}
	parser, err := kong.New(&cli, options...)
	if err != nil {
		return cli, err
	}
	_, err = parser.Parse(args)
	return cli, err
}
