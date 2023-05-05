package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/commands"
	"github.com/mrmioxin/gak_telegram_bot/internal/config"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
	"github.com/mrmioxin/gak_telegram_bot/internal/sessions"
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
	configFile, err := os.ReadFile("bot.cfg")
	if err != nil {
		log.Panic("error reading config file:", err)
	}
	if conf, err = config.NewConfig(strings.NewReader(string(configFile))); err != nil {
		log.Panic(err)
	}
	if conf.Token == "" {
		log.Panic("error config file: secure token expected")
	}

	if conf.LogFile == "" {
		conf.LogFile = "bot.log"
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
	flag.Parse()

	parsConf()

	bot, err := tgapi.NewBotAPI(conf.Token)
	if err != nil {
		log.Fatal(">>>", conf.Token, ">>> ", err)
	}

	bot.Debug = debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	insure := services.NewInsurence("ОСНС", 1000.00)
	ses := make(map[int64]sessions.Session)

	c := commands.NewCommander(bot, conf, insure, sessions.MemSessions(ses))
	u := tgapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		c.HandlerMain(update)
	}
}
