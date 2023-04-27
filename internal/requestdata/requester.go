package requestdata

import (
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type reguest struct {
	text   string
	worker func(c *Requester, mes *tgapi.Message)
}

var requests_list = make([]reguest, 0)

type Service interface {
	Calculate() (string, error)
}

type Requester struct {
	bot             *tgapi.BotAPI
	chatID          int64
	Product_service Service
}

func NewRequester(bot *tgapi.BotAPI, chatID int64, serv Service) *Requester {
	//t := reflect.TypeOf(data)
	return &Requester{bot, chatID, serv}
}

func (req *Requester) Run() error {
	var msg tgapi.MessageConfig
	u := tgapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := req.bot.GetUpdatesChan(u)
	idx := 0

	for idx < len(requests_list) {
		msg = tgapi.NewMessage(req.chatID, requests_list[idx].text)
		req.bot.Send(msg)

		for update := range updates {
			if update.Message == nil {
				continue
			}
			requests_list[idx].worker(req, update.Message)
		}
	}
	return nil
}

// func (req *Requester) HandlerCommand(update tgapi.Update) {
// 	defer func() {
// 		if panicVal := recover(); panicVal != nil {
// 			log.Printf("recover panic %v:", panicVal)
// 		}
// 	}()
// 	// If we got a message
// 	if command, ok := registered_commands[update.Message.Command()]; ok {
// 		command(req, update.Message)
// 	} else {
// 		msg := tgapi.NewMessage(update.Message.Chat.ID, WRONG_INPUT)
// 		req.bot.Send(msg)
// 	}

// }
