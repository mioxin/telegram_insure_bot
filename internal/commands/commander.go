package commands

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/sessions"
)

const (
	WRONG_INPUT  string = `Введите команду или выберите ёё из меню.`
	WRONG_ACCESS string = "Извините, пока доступ закрыт."
)

type reguest struct {
	ok_text    string
	wrong_text string
	worker     func(c *Commander, mes *tgapi.Message) error
}
type RegisteredCommand struct {
	Description string
	Worker      func(c *Commander, mes *tgapi.Message) string
	ShowInHelp  bool
}

var requestsListCalc = make([]reguest, 0)

var registered_commands = map[string]RegisteredCommand{}

type IService interface {
	Calculate() (string, error)
}

type IConfig interface {
	IsAccess(user string) bool
}
type ISessions interface {
	GetSession(id int64) (*sessions.Session, error)
	GetIdsByUser(user string) []int64
	UpdateSession(id int64, ses *sessions.Session) error
	AddSession(id int64, ses *sessions.Session)
}
type Commander struct {
	bot             *tgapi.BotAPI
	Config          IConfig
	Product_service IService
	Sessions        ISessions
}

func NewCommander(bot *tgapi.BotAPI, conf IConfig, serv IService, ses ISessions) *Commander {
	return &Commander{bot, conf, serv, ses}
}

func (cmder *Commander) HandlerMain(update tgapi.Update) error {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			log.Printf("recover panic in HandlerMain %v:", panicVal)
		}
	}()

	if update.CallbackQuery != nil {
		//processing callbacks...
		cmder.HandlerCallback(update)
		return nil
	}
	if update.Message.IsCommand() {
		cmder.HandlerCommand(update)
	} else {
		ses, err := cmder.Sessions.GetSession(update.Message.Chat.ID)
		if err != nil {
			log.Printf("error HandlerMain: not found session for Message \"%v\" (user %v)", update.Message.Text, update.Message.Chat.UserName)
			return err
		}
		switch ses.ActionName {
		case "calc":
			cmder.HandlerRequest(update)
		default:
		}
	}
	return nil
}

func (cmder *Commander) HandlerCommand(update tgapi.Update) {
	// If we got a message
	if command, ok := registered_commands[update.Message.Command()]; ok {
		ses, err := cmder.Sessions.GetSession(update.Message.Chat.ID)
		if err != nil {
			ses = sessions.NewSession(update.Message.Chat.UserName)
			if cmder.Config.IsAccess(update.Message.Chat.UserName) {
				ses.AccessCommand["all"] = struct{}{}
			} else {
				ses.AccessCommand["about"] = struct{}{}
			}
			cmder.Sessions.AddSession(update.Message.Chat.ID, ses)
		}
		_, okAll := ses.AccessCommand["all"]
		_, okCmd := ses.AccessCommand[update.Message.Command()]
		if okAll || okCmd {
			ses.ActionName = command.Worker(cmder, update.Message)
			cmder.Sessions.UpdateSession(update.Message.Chat.ID, ses)
			//log.Println(ses)
		} else {
			msg := tgapi.NewMessage(update.Message.Chat.ID, WRONG_ACCESS)
			cmder.bot.Send(msg)
			log.Printf("HandlerCommand: deny acces fo @%s on /%s", update.Message.Chat.UserName, update.Message.Command())
		}
	}

}

func (cmder *Commander) HandlerRequest(update tgapi.Update) {
	ses, _ := cmder.Sessions.GetSession(update.Message.Chat.ID)
	err := requestsListCalc[ses.IdxRequest].worker(cmder, update.Message)

	if err != nil {
		log.Printf("error: Idx=%v %v", ses.IdxRequest, err)
		cmder.bot.Send(tgapi.NewMessage(update.Message.Chat.ID, requestsListCalc[ses.IdxRequest].wrong_text))
	} else {
		mes := tgapi.NewMessage(update.Message.Chat.ID, requestsListCalc[ses.IdxRequest].ok_text)
		mes.ParseMode = "Markdown"
		cmder.bot.Send(mes)

		ses.IdxRequest++
		if ses.IdxRequest >= len(requestsListCalc) {
			ses.ResetSession()
		}
		cmder.Sessions.UpdateSession(update.Message.Chat.ID, ses)
	}
}

func (cmder *Commander) HandlerCallback(update tgapi.Update) {

}

// func (cmder *Commander) WatchConfig(isModify chan any) {
// 	for range isModify {

// 	}
// }
