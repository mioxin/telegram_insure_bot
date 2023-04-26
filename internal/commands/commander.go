package commands

import (
	"log"
	"time"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const WRONG_INPUT string = `Введите команду или выберите ёё из меню.
`

var registered_commands = map[string]func(c *Commander, mes *tgapi.Message){}

type Commander struct {
	bot *tgapi.BotAPI
	//product_service Service
}

func NewCommander(bot *tgapi.BotAPI) *Commander {
	return &Commander{bot}
}

func (cmder *Commander) Run() error {
	u := tgapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := cmder.bot.GetUpdatesChan(u)
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

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
	} else {
		msg := tgapi.NewMessage(update.Message.Chat.ID, WRONG_INPUT)
		cmder.bot.Send(msg)
	}

}
