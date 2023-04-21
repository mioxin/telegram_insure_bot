package main

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/commands"
	"github.com/mrmioxin/gak_telegram_bot/internal/config"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
)

func main() {
	conf, err := config.NewConfig("bot.cfg")
	if err != nil {
		log.Panic("error reading config file:", err)
	}
	if conf.Token == "" {
		log.Panic("error config file: secure token expected")
	}
	bot, err := tgapi.NewBotAPI(conf.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	insure := services.NewInsurence("ОСНС", 1000.00)
	c := commands.NewCommander(bot, insure)
	if err := (*c).Run(); err != nil {
		log.Panic(err)
	}
}
