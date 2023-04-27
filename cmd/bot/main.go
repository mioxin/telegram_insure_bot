package main

import (
	"flag"
	"log"
	"os"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/commands"
	"github.com/mrmioxin/gak_telegram_bot/internal/config"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
)

var (
	verbouse bool
	conf     *config.Config
)

func init() {
	flag.BoolVar(&verbouse, "v", false, "Output fool log to StdOut (shorthand)")
	flag.BoolVar(&verbouse, "verbouse", false, "Output fool log to StdOut")

}

func parsConf() {
	var err error
	conf, err = config.NewConfig("bot.cfg")
	if err != nil {
		log.Panic("error reading config file:", err)
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
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	insure := services.NewInsurence("ОСНС", 1000.00)

	c := commands.NewCommander(bot, insure)
	if err := (*c).Run(); err != nil {
		log.Panic(err)
	}
}
