package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/commands"
	"github.com/mrmioxin/gak_telegram_bot/internal/config"
	"github.com/mrmioxin/gak_telegram_bot/internal/httpclient"
	"github.com/mrmioxin/gak_telegram_bot/resources"
)

var (
	verbouse          bool
	debug             bool
	token, configFile string
	conf              *config.Config
)

func init() {
	flag.BoolVar(&verbouse, "v", false, "Output fool log to StdOut (shorthand)")
	flag.BoolVar(&verbouse, "verbouse", false, "Output fool log to StdOut")
	flag.BoolVar(&debug, "d", false, "Output debug info to StdOut (shorthand)")
	flag.StringVar(&token, "t", "", "The security token for connecting to Telegram API")
	flag.StringVar(&configFile, "c", "", "The configuration file")
	// registers.Sessions = *sessions.NewMemSessions()
	// registers.RegisteredCommands = make(map[string]func(mes *tgapi.Message) string)
	// registers.RegisteredActions = make(map[string]actions.Action)

}

func parsConf() {
	var err error
	flag.Parse()

	if configFile == "" {
		configFile = resources.CONFIG_FILE_NAME
	}

	if conf, err = config.NewConfig(configFile); err != nil {
		log.Panic(err)
	}

	if token != "" {
		conf.Token = token
	}

	if !verbouse {
		output_log, err := os.OpenFile(conf.LogFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal("Cant open ouput file for loging ", conf.LogFile, ".\n", err)
		}
		log.SetOutput(output_log)
	}
}

func main() {

	parsConf()

	bot, err := tgapi.NewBotAPIWithClient(conf.Token, httpclient.NewHTTPClient())
	if err != nil {
		log.Fatal("<<<", conf.Token, ">>> ", err)
	}

	bot.Debug = debug
	log.Printf("Authorized on account %s\n", bot.Self.UserName)
	chSignal := make(chan os.Signal)

	isModifyConfig := make(chan any)
	go conf.Watch(resources.CONFIG_FILE_NAME, resources.DURATION_WATCH_CONFIG, isModifyConfig)

	c := commands.NewCommander(bot, conf)
	defer c.Stop()

	signal.Notify(chSignal, os.Interrupt, os.Kill)
	go func() {
		<-chSignal
		c.Stop()
		log.Fatal("Cancel bot by OS interruption.")
	}()

	c.Start()
}
