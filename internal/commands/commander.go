package commands

import (
	"fmt"
	"log"
	"time"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const WRONG_INPUT string = `Введите команду или выберите ёё из меню.
`

type reguest struct {
	ok_text    string
	wrong_text string
	worker     func(c *Commander, mes *tgapi.Message) error
}

var BackKeyboard = tgapi.NewReplyKeyboard(
	tgapi.NewKeyboardButtonRow(
		tgapi.NewKeyboardButton("<<< Назад"),
	))

var requests_list = make([]reguest, 0)

var registered_commands = map[string]func(c *Commander, mes *tgapi.Message) string{}

type Service interface {
	Calculate() (string, error)
}

type Commander struct {
	bot             *tgapi.BotAPI
	Product_service Service
	Idx             int
	Handler         string
}

func NewCommander(bot *tgapi.BotAPI, serv Service) *Commander {
	return &Commander{bot, serv, 0, ""}
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
			//processing callbacks...
			continue
		}
		if update.Message.IsCommand() {
			cmder.HandlerCommand(update)
		} else {
			switch cmder.Handler {
			case "calc":
				cmder.HandlerRequest(update)
			default:
			}
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
		// ses, ok := Sessions[update.Message.Chat.ID]
		// if !ok {
		// 	main.Sessions[update.Message.Chat.ID] = *session.NewSession(update.Message.From.UserName)
		// 	log.Printf("handlerCommand: create new session on %v", update)
		// }
		// ses.LastTime = time.Now()
		cmder.Handler = command(cmder, update.Message)
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
	// ses, err := cmder.getSession(update.Message.Chat.ID)
	// if err != nil {
	// 	log.Printf("error commander.Run: while calc operation %v\n get update: %v", err, update.Message)
	// 	return
	// }
	// ses.LastTime = time.Now()
	// idx := cmder.getSessionIdx(update.Message.Chat.ID)
	// if !cmder.ErrorInput {
	// 	msg := tgapi.NewMessage(update.Message.Chat.ID, requests_list[cmder.Idx].text)
	// 	msg.ReplyMarkup = BackKeyboard
	// 	cmder.bot.Send(msg)
	// }
	strIdx := fmt.Sprintf(" (Idx=%d)", cmder.Idx)
	err := requests_list[cmder.Idx].worker(cmder, update.Message)
	fmt.Println(err, strIdx, requests_list[cmder.Idx])

	if err != nil {
		log.Printf("error: Idx=%v %v", cmder.Idx, err)
		cmder.bot.Send(tgapi.NewMessage(update.Message.Chat.ID, requests_list[cmder.Idx].wrong_text+strIdx))
	} else {
		mes := tgapi.NewMessage(update.Message.Chat.ID, requests_list[cmder.Idx].ok_text+strIdx)
		mes.ParseMode = "Markdown"
		cmder.bot.Send(mes)
		cmder.Idx++
		if cmder.Idx >= len(requests_list) {
			cmder.ResetSession()
		}
	}
}

func (cmder *Commander) ResetSession() {
	cmder.Idx = 0
	cmder.Handler = ""
}

// func (cmder *Commander) getSession(sesID int64) (*session.Session, error) {
// 	if ses, ok := main.Sessions[sesID]; ok {
// 		return &ses, nil
// 	}
// 	return nil, fmt.Errorf("error getSession: session %v not found", sesID)
// }

// func (cmder *Commander) getSessionHandler(sesID int64) string {
// 	res := ""
// 	if ses, err := cmder.getSession(sesID); err == nil {
// 		return ses.Handler
// 	}
// 	return res
// }

// func (cmder *Commander) getSessionIdx(sesID int64) int {
// 	res := -1
// 	if ses, err := cmder.getSession(sesID); err == nil {
// 		return ses.Idx_request
// 	} else {
// 		log.Panicf("handlerCommand:  %v", err)
// 	}
// 	return res
// }

// func (cmder *Commander) setSessionLtime(sesID int64) {
// 	if ses, ok := main.Sessions[sesID]; ok {
// 		ses.LastTime = time.Now()
// 	} else {
// 		log.Panicf("handlerCommand: not found session on %v", sesID)
// 	}
// }
