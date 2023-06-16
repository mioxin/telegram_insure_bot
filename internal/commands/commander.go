package commands

import (
	"fmt"
	"log"
	"time"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/handlers"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
	"github.com/mrmioxin/gak_telegram_bot/internal/sessions"
)

const (
	TYPE_OF_BUSNS_FILENAME string = "vid.txt"
	WRONG_AGAIN            string = `Опять ошибка. `
	WRONG_INPUT            string = `Введите команду или выберите её из меню.`
	WRONG_ACCESS           string = "Извините, пока доступ закрыт."
	YES                    string = "Да"
	NO                     string = "Нет"
)

// type IHandler interface {
// 	Execute()
// }

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

func NewCommander(bot *tgapi.BotAPI, conf IConfig) *Commander {
	tob, err := services.NewFileTypeOfBusns(TYPE_OF_BUSNS_FILENAME)
	if err != nil {
		log.Fatal("<<<", TYPE_OF_BUSNS_FILENAME, ">>> ", err)
	}

	serv := services.NewInsurence("ОСНС")
	//srvs := make(map[string]sessions.Services)
	ses := sessions.NewMemSessions()

	return &Commander{bot, conf, serv, ses, tob}
}

func (cmder *Commander) Start() {
	u := tgapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := cmder.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(">>> ", err)
	}
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		go cmder.handl(update)
	}

}

func (cmder *Commander) handl(update tgapi.Update) error {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			log.Printf("recover panic in HandlerMain %v:", panicVal)
		}
	}()

	ses, err := cmder.Sessions.GetSession(update.Message.Chat.ID)
	if err != nil {
		log.Printf("error HandlerMain: not found session for Message \"%v\" (user %v)", update.Message.Text, update.Message.Chat.UserName)
		ses = sessions.NewSession(update.Message.Chat.UserName)
		if cmder.Config.IsAccess(update.Message.Chat.UserName) {
			ses.AccessCommand["all"] = struct{}{}
		} else {
			ses.AccessCommand["about"] = struct{}{}
		}
		cmder.Sessions.AddSession(update.Message.Chat.ID, ses)
	}

	if update.CallbackQuery != nil {
		cmder.HandlerCallback(update)
		return nil
	}
	if update.Message.IsCommand() {
		h := handlers.NewHandlerCommand(cmder.bot, ses, update)
		h.Execute()
		//cmder.HandlerCommand(update)
	} else {
		if ses.ActionName == "" {
			log.Printf("error HandlerMain: not found ActionName in session of Message \"%v\" (user %v)", update.Message.Text, update.Message.Chat.UserName)
			return fmt.Errorf("error HandlerMain: not found ActionName in session of Message \"%v\" (user %v)", update.Message.Text, update.Message.Chat.UserName)
		}
		switch ses.ActionName {
		case "calc":
			cmder.HandlerCalc(update, ses)
		default:
		}
	}
	return nil
}
