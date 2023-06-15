package commands

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/sessions"
)

const (
	WRONG_AGAIN  string = `Опять ошибка. `
	WRONG_INPUT  string = `Введите команду или выберите её из меню.`
	WRONG_ACCESS string = "Извините, пока доступ закрыт."
	YES          string = "Да"
	NO           string = "Нет"
)

type ITypeOfBusiness interface {
	Get(vid string) (string, error)
}

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
	TypeOfBuseness  ITypeOfBusiness
}

func NewCommander(bot *tgapi.BotAPI, conf IConfig, serv IService, ses ISessions, tob ITypeOfBusiness) *Commander {
	return &Commander{bot, conf, serv, ses, tob}
}

func (cmder *Commander) HandlerMain(update tgapi.Update) error {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			log.Printf("recover panic in HandlerMain %v:", panicVal)
		}
	}()

	if update.CallbackQuery != nil {
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
			cmder.HandlerCalc(update, ses)
		default:
		}
	}
	return nil
}

// func (cmder *Commander) WatchConfig(isModify chan any) {
// 	for range isModify {

// 	}
// }
