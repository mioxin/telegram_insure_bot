package commands

import (
	"fmt"
	"log"
	"time"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/sessions"
)

const WRONG_INPUT string = `Введите команду или выберите ёё из меню.
`

type reguest struct {
	ok_text    string
	wrong_text string
	worker     func(c *Commander, mes *tgapi.Message) error
}

var requests_list = make([]reguest, 0)

var registered_commands = map[string]func(c *Commander, mes *tgapi.Message) string{}

type Service interface {
	Calculate() (string, error)
}
type SessionsI interface {
	GetSession(id int64) (*sessions.Session, error)
	UpdateSession(id int64, ses *sessions.Session) error
	AddSession(id int64, ses *sessions.Session)
}
type Commander struct {
	bot             *tgapi.BotAPI
	Product_service Service
	Sessions        SessionsI
	// Idx             int
	// Handler         string
}

func NewCommander(bot *tgapi.BotAPI, serv Service, ses SessionsI) *Commander {
	return &Commander{bot, serv, ses}
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
			ses, err := cmder.Sessions.GetSession(update.Message.Chat.ID)
			if err != nil {
				continue
			}
			switch ses.ActionName {
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
		ses, err := cmder.Sessions.GetSession(update.Message.Chat.ID)
		if err != nil {
			ses = sessions.NewSession(update.Message.Chat.UserName)
			cmder.Sessions.AddSession(update.Message.Chat.ID, ses)
		}
		ses.ActionName = command(cmder, update.Message)
		cmder.Sessions.UpdateSession(update.Message.Chat.ID, ses)
		log.Println(ses)
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
	ses, err := cmder.Sessions.GetSession(update.Message.Chat.ID)
	if err != nil {
		ses = sessions.NewSession(update.Message.Chat.UserName)
		cmder.Sessions.AddSession(update.Message.Chat.ID, ses)
	}

	strIdx := fmt.Sprintf(" (Idx=%d)", ses.IdxRequest)
	err = requests_list[ses.IdxRequest].worker(cmder, update.Message)
	fmt.Println(err, strIdx, requests_list[ses.IdxRequest])

	if err != nil {
		log.Printf("error: Idx=%v %v", ses.IdxRequest, err)
		cmder.bot.Send(tgapi.NewMessage(update.Message.Chat.ID, requests_list[ses.IdxRequest].wrong_text+strIdx))
	} else {
		mes := tgapi.NewMessage(update.Message.Chat.ID, requests_list[ses.IdxRequest].ok_text+strIdx)
		mes.ParseMode = "Markdown"
		cmder.bot.Send(mes)
		ses.IdxRequest++
		if ses.IdxRequest >= len(requests_list) {
			ses.ResetSession()
		}
		cmder.Sessions.UpdateSession(update.Message.Chat.ID, ses)
	}
}
