package commands

import (
	"fmt"
	"log"
	"time"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/session"
)

const WRONG_INPUT string = `Введите команду или выберите ёё из меню.
`

type reguest struct {
	text   string
	worker func(c *Commander, mes *tgapi.Message)
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
	Sessions        map[int64]session.Session
}

func NewCommander(bot *tgapi.BotAPI, serv Service) *Commander {
	s := make(map[int64]session.Session)
	return &Commander{bot, serv, s}
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

		switch cmder.getSessionHandler(update.Message.Chat.ID) {
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
		if ses, ok := cmder.Sessions[update.Message.Chat.ID]; ok {
			ses.Handler = command(cmder, update.Message)
		} else {
			log.Panicf("handlerCommand: not found session on %v", update)
		}
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
	ses, err := cmder.getSession(update.Message.Chat.ID)
	if err != nil {
		log.Printf("error commander.Run: while calc operation %v\n get update: %v", err, update.Message)
		return
	}

	idx := cmder.getSessionIdx(update.Message.Chat.ID)
	if ses.ErrorInput {
		msg := tgapi.NewMessage(update.Message.Chat.ID, requests_list[idx].text)
		cmder.bot.Send(msg)
	} else {
		requests_list[idx].worker(cmder, update.Message)
	}
}

func (cmder *Commander) getSession(sesID int64) (*session.Session, error) {
	if ses, ok := cmder.Sessions[sesID]; ok {
		return &ses, nil
	} else {
		log.Panicf("handlerCommand: not found session on %v", sesID)
	}
	return nil, fmt.Errorf("error getSession: session %v not found", sesID)
}

func (cmder *Commander) getSessionHandler(sesID int64) string {
	res := ""
	if ses, err := cmder.getSession(sesID); err == nil {
		return ses.Handler
	} else {
		log.Panicf("handlerCommand:  %v", err)
	}
	return res
}

func (cmder *Commander) getSessionIdx(sesID int64) int {
	res := -1
	if ses, err := cmder.getSession(sesID); err == nil {
		return ses.Idx_request
	} else {
		log.Panicf("handlerCommand:  %v", err)
	}
	return res
}

func (cmder *Commander) setSessionLtime(sesID int64) {
	if ses, ok := cmder.Sessions[sesID]; ok {
		ses.LastTime = time.Now()
	} else {
		log.Panicf("handlerCommand: not found session on %v", sesID)
	}
}
