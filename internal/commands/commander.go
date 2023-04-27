package commands

import (
	"log"
	"time"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const WRONG_INPUT string = `Введите команду или выберите ёё из меню.
`

type reguest struct {
	text   string
	worker func(c *Commander, mes *tgapi.Message)
}

var requests_list = make([]reguest, 0)

var registered_commands = map[string]func(c *Commander, mes *tgapi.Message) string{}

type Service interface {
	Calculate() (string, error)
}

type Commander struct {
	bot             *tgapi.BotAPI
	Product_service Service
	handler         *string
	idx_request     *int
}

func NewCommander(bot *tgapi.BotAPI, serv Service) *Commander {
	h := ""
	i := 0
	return &Commander{bot, serv, &h, &i}
}

func (cmder *Commander) Run() error {
	u := tgapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := cmder.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		if update.Message == nil {
			continue
		}
		switch *cmder.handler {
		case "calc":
			cmder.HandlerRequest(update)
		default:
			cmder.HandlerCommand(update)
		}
	}
	return nil
}

func (cmder *Commander) HandlerCommand(update tgapi.Update) {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			log.Printf("recover panic in HandlerCommand %v:", panicVal)
		}
	}()
	// If we got a message
	if command, ok := registered_commands[update.Message.Command()]; ok {
		(*cmder.handler) = command(cmder, update.Message)
	} else {
		msg := tgapi.NewMessage(update.Message.Chat.ID, WRONG_INPUT)
		cmder.bot.Send(msg)
	}

}

func (cmder *Commander) HandlerRequest(update tgapi.Update) {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			log.Printf("recover panic in HandlerRequest%v:", panicVal)
		}
	}()

	idx := *cmder.idx_request
	msg := tgapi.NewMessage(update.Message.Chat.ID, requests_list[idx].text)
	cmder.bot.Send(msg)

	requests_list[idx].worker(cmder, update.Message)
	*cmder.idx_request++
}

func (cmder *Commander) resetCommander() {
	*cmder.handler = ""
	*cmder.idx_request = 0
}
