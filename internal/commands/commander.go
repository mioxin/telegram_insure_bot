package commands

import (
	"log"
	"time"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/handlers"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
	"github.com/mrmioxin/gak_telegram_bot/internal/services/product1"
	"github.com/mrmioxin/gak_telegram_bot/internal/storages/sessions"
)

const (
	WRONG_AGAIN  string = `Опять ошибка. `
	WRONG_INPUT  string = `Введите команду или выберите её из меню.`
	WRONG_ACCESS string = "Извините, пока доступ закрыт."
	YES          string = "Да"
	NO           string = "Нет"
)

type IHandler interface {
	Execute()
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
	bot      *tgapi.BotAPI
	Config   IConfig
	Products map[string]services.IService
	Sessions ISessions
}

func NewCommander(bot *tgapi.BotAPI, conf IConfig) *Commander {
	prods := make(map[string]services.IService)
	serv := product1.NewInsurence("ОСНС")
	prods["calc"] = serv

	ses := sessions.NewMemSessions()

	return &Commander{bot, conf, prods, ses}
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

	var h IHandler
	var ses *sessions.Session
	log.Printf("Sessions in HandlerMain: %v", cmder.Sessions)

	switch {
	case update.CallbackQuery != nil:
		log.Printf("CallbackQuery in HandlerMain: %v", update.CallbackQuery)
		h = handlers.NewHandlerMessage(cmder.bot, ses, cmder.Products, update)
		//cmder.HandlerCallback(ses, update)

	case update.Message.IsCommand():
		log.Printf("Command in HandlerMain: %v", update.Message)
		ses = sessions.NewSession(update.Message.Chat.UserName)
		if cmder.Config.IsAccess(update.Message.Chat.UserName) {
			ses.AccessCommand["all"] = struct{}{}
		} else {
			ses.AccessCommand["about"] = struct{}{}
		}
		cmder.Sessions.AddSession(update.Message.Chat.ID, ses)

		h = handlers.NewHandlerCommand(cmder.bot, ses, update)
		//cmder.HandlerCommand(update)

	case update.Message != nil:
		log.Printf("Message in HandlerMain: %v", update.Message)
		if ses, err := cmder.Sessions.GetSession(update.Message.Chat.ID); err == nil {
			h = handlers.NewHandlerMessage(cmder.bot, ses, cmder.Products, update)
		} else {
			log.Printf("error HandlerMain: not found session for Message \"%v\" (user %v)\n%v", update.Message, update.Message.Chat.UserName, err)
		}

	default:
		log.Printf("error HandlerMain: invalid Message \"%v\" (user %v)", update.Message, update.Message.Chat.UserName)

	}

	h.Execute()
	cmder.Sessions.UpdateSession(update.Message.Chat.ID, ses)
	// if update.Message != nil {
	// } else {
	// 	log.Printf("error HandlerMain: cant update session Message \"%v\" (user %v)", update.Message, update.Message.Chat.UserName)
	// }

	return nil
}
