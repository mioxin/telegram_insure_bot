package commands

import (
	"log"
	"runtime/debug"
	"strings"
	"time"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/handlers"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
	srv_cfls "github.com/mrmioxin/gak_telegram_bot/internal/services/clientfiles"
	"github.com/mrmioxin/gak_telegram_bot/internal/services/product1"
	"github.com/mrmioxin/gak_telegram_bot/internal/storages/clientfiles"
	"github.com/mrmioxin/gak_telegram_bot/internal/storages/sessions"
)

type IHandler interface {
	Execute()
}

type IFilesStorage interface {
	GetFileId(name string) (string, error)
	SetFileId(name, user, id string) error
	Close()
}

type IConfig interface {
	IsAccess(user string) bool
	IsAccWord(word string) bool
	Close()
}
type ISessions interface {
	GetSession(id int64) (*sessions.Session, error)
	GetIdsByUser(user string) []int64
	UpdateSession(id int64, ses *sessions.Session) error
	AddSession(id int64, ses *sessions.Session)
	Close()
}
type Commander struct {
	bot      *tgapi.BotAPI
	Config   IConfig
	Products map[string]services.IService
	Sessions ISessions
	Files    IFilesStorage
	// done     chan struct{}
}

func NewCommander(bot *tgapi.BotAPI, conf IConfig) *Commander {
	prods := make(map[string]services.IService)
	files := clientfiles.NewMapStorage()

	prods["calc"] = product1.NewInsurence("ОСНС") //processing of calc command (get data and colculate the cost of osns Insurance)
	prods["send"] = srv_cfls.NewReceiver(files)   //receive and store a files was send to bot from users

	s := sessions.NewMemSessions()
	// done := make(chan struct{})

	return &Commander{bot, conf, prods, s, files}
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
	log.Println(">>> End commander.")

}

func (cmder *Commander) handl(update tgapi.Update) error {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			log.Printf("recover panic in HandlerMain %v:\n%v", panicVal, string(debug.Stack()))
		}
	}()

	var h IHandler
	var ses *sessions.Session
	var chatID int64
	log.Printf("Sessions in HandlerMain: %v", cmder.Sessions)

	switch {
	case update.CallbackQuery != nil:
		log.Printf("CallbackQuery in HandlerMain: %v", update.CallbackQuery)
		chatID = update.CallbackQuery.Message.Chat.ID
		ses, _ = cmder.Sessions.GetSession(chatID)
		h = handlers.NewHandlerMessage(cmder.bot, ses, cmder.Products, update)

	case update.Message.IsCommand():
		log.Printf("Command in HandlerMain: %#v", update.Message)
		ses = sessions.NewSession(update.Message.Chat.UserName)

		if cmder.Config.IsAccess(update.Message.Chat.UserName) {
			ses.AccessCommand["all"] = struct{}{}
		} else if update.Message.CommandArguments() != "" {
			s := strings.Split(update.Message.CommandArguments(), " ")
			if cmder.Config.IsAccWord(strings.TrimSpace(s[0])) {
				ses.AccessCommand[update.Message.Command()] = struct{}{}
			}
		} else {
			ses.AccessCommand["about"] = struct{}{}
		}
		chatID = update.Message.Chat.ID
		cmder.Sessions.AddSession(chatID, ses)
		h = handlers.NewHandlerCommand(cmder.bot, cmder.Files, ses, update)

	case update.Message != nil:
		log.Printf("Message in HandlerMain: %#v", update.Message)
		chatID = update.Message.Chat.ID
		ses, _ = cmder.Sessions.GetSession(chatID)
		h = handlers.NewHandlerMessage(cmder.bot, ses, cmder.Products, update)

	default:
		log.Printf("error HandlerMain: invalid Message %#v (user %v)", update.Message, update.Message.Chat.UserName)
	}

	h.Execute()

	if err := cmder.Sessions.UpdateSession(chatID, ses); err != nil {
		log.Printf("After session update in HandlerMain: %v\n %v", err, cmder.Sessions)
	}

	return nil
}

func (cmder *Commander) Stop() {
	cmder.Sessions.Close()
	cmder.Config.Close()
	cmder.Files.Close()
	cmder.bot.StopReceivingUpdates()
	log.Println("Commander stoped.")

}
