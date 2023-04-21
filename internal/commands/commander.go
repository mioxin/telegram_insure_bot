package commands

import (
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Service interface {
	Calculate(data any) (string, error)
}

type Commander struct {
	bot             *tgapi.BotAPI
	product_service Service
}

func NewCommander(bot *tgapi.BotAPI, product_service Service) *Commander {
	return &Commander{bot, product_service}
}

func (cmder *Commander) Run() error {
	u := tgapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := cmder.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		// If we got a message
		switch update.Message.Command() {
		case "help":
			cmder.help(update.Message)
		case "about":
			cmder.about(update.Message)
		case "calc":
			cmder.calc(update.Message)
		default:
		}
	}
	return nil
}
