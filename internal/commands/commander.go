package commands

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var registered_commands = map[string]func(c *Commander, mes *tgapi.Message){}

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
		cmder.HandlerCommand(update)
	}
	return nil
}

func (cmder *Commander) HandlerCommand(update tgapi.Update) {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			log.Printf("recover panic %v:", panicVal)
		}
	}()
	// If we got a message
	if command, ok := registered_commands[update.Message.Command()]; ok {
		command(cmder, update.Message)
	}

}
