package main

import (
	"flag"
	"log"
	"os"
	"time"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/commands"
	"github.com/mrmioxin/gak_telegram_bot/internal/config"
	"github.com/mrmioxin/gak_telegram_bot/internal/httpclient"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
	"github.com/mrmioxin/gak_telegram_bot/internal/sessions"
)

const (
	CONFIG_FILE_NAME      string        = "bot.cfg"
	DURATION_WATCH_CONFIG time.Duration = 3 * time.Second
)

var (
	verbouse bool
	debug    bool
	conf     *config.Config
)

func init() {
	flag.BoolVar(&verbouse, "v", false, "Output fool log to StdOut (shorthand)")
	flag.BoolVar(&verbouse, "verbouse", false, "Output fool log to StdOut")
	flag.BoolVar(&debug, "d", false, "Output debug info to StdOut (shorthand)")
	// registers.Sessions = *sessions.NewMemSessions()
	// registers.RegisteredCommands = make(map[string]func(mes *tgapi.Message) string)
	// registers.RegisteredActions = make(map[string]actions.Action)

}

func parsConf() {
	var err error

	flag.Parse()

	if conf, err = config.NewConfig(CONFIG_FILE_NAME); err != nil {
		log.Panic(err)
	}

	if !verbouse {
		output_log, err := os.OpenFile(conf.LogFile, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal("Cant open ouput file for loging bot.log.\n", err)
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
	log.Printf("Authorized on account %s", bot.Self.UserName)

	isModifyConfig := make(chan any)
	go conf.Watch(CONFIG_FILE_NAME, DURATION_WATCH_CONFIG, isModifyConfig)

	insure := services.NewInsurence("ОСНС", 1000.00)
	//srvs := make(map[string]sessions.Services)

	c := commands.NewCommander(bot, conf, insure, sessions.NewMemSessions())
	//go c.WatchConfig(isModifyConfig)

	u := tgapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(">>> ", err)
	}
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		go c.HandlerMain(update)
	}
}
