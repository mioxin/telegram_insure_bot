package commands

import (
	"log"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/config"
	"github.com/mrmioxin/gak_telegram_bot/internal/handlers"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
	"github.com/mrmioxin/gak_telegram_bot/internal/services/product1"
	srv_cfls "github.com/mrmioxin/gak_telegram_bot/internal/services/receive_client_files"
	"github.com/mrmioxin/gak_telegram_bot/internal/services/uploadclientfiles"
	"github.com/mrmioxin/gak_telegram_bot/internal/storages/clientfiles"
	"github.com/mrmioxin/gak_telegram_bot/internal/storages/sessions"
	"github.com/mrmioxin/gak_telegram_bot/resources"
)

type IHandler interface {
	Execute()
}

type IFilesStorage interface {
	GetFileId(name, user string) (string, error)
	SetFileId(name, user, id string) error
	ListUsers() []string
	// ListFiles(user string) []*storages.FileInfo
	Close()
}

//	type IConfig interface {
//		IsAdmin(user string) bool
//		IsAccess(user string) bool
//		IsAccWord(word string) bool
//		Close()
//	}
type ISessions interface {
	GetSession(id int64) (*sessions.Session, error)
	GetIdsByUser(user string) []int64
	UpdateSession(id int64, ses *sessions.Session) error
	AddSession(id int64, ses *sessions.Session)
	Close()
}
type Commander struct {
	bot         *tgapi.BotAPI
	Config      *config.Config
	Products    map[string]services.IService
	Sessions    ISessions
	Files       IFilesStorage
	updateLocks sync.Map
	// newConf  chan any
}

func NewCommander(bot *tgapi.BotAPI, conf *config.Config) *Commander {
	prods := make(map[string]services.IService)
	files := clientfiles.NewMapStorage()

	prods["calc"] = product1.NewInsurence("ОСНС")     //processing of calc command (get data and colculate the cost of osns Insurance)
	prods["send"] = srv_cfls.NewReceiver(files)       //receive and store a files was send to bot from users
	prods["get"] = uploadclientfiles.NewGetter(files) //get client files that was sended to bot from users from the store

	s := sessions.NewMemSessions()
	// done := make(chan struct{})
	return &Commander{bot: bot, Config: conf, Products: prods, Sessions: s, Files: files}
}

func (cmder *Commander) Start() {
	go func() {
		for range cmder.Watch(resources.DURATION_WATCH_CONFIG) {
			cmder.Config.Update()
			log.Printf("Update config")
		}
	}()

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
	if update.CallbackQuery != nil {
		if update.CallbackQuery.Message == nil || update.CallbackQuery.Message.Chat == nil {
			log.Printf("error HandlerMain: callback without message: %***REMOVED***v", update.CallbackQuery)
			return nil
		}
		chatID = update.CallbackQuery.Message.Chat.ID
	} else if update.Message != nil && update.Message.Chat != nil {
		chatID = update.Message.Chat.ID
	} else {
		log.Printf("error HandlerMain: update has no supported chat: %***REMOVED***v", update)
		return nil
	}
	lockValue, _ := cmder.updateLocks.LoadOrStore(chatID, &sync.Mutex{})
	chatLock := lockValue.(*sync.Mutex)
	chatLock.Lock()
	defer chatLock.Unlock()

	log.Printf("Sessions in HandlerMain: %v", cmder.Sessions)

	switch {
	case update.CallbackQuery != nil:
		log.Printf("CallbackQuery in HandlerMain: %***REMOVED***v", update.CallbackQuery)
		ses, _ = cmder.Sessions.GetSession(chatID)
		h = handlers.NewHandlerMessage(cmder.bot, ses, cmder.Products, update)

	case update.Message != nil && update.Message.IsCommand():
		log.Printf("Command in HandlerMain: %***REMOVED***v", update.Message.Command())
		ses = sessions.NewSession(update.Message.Chat.UserName)

		cmder.setAvailableCommands(ses, update.Message)

		cmder.Sessions.AddSession(chatID, ses)
		h = handlers.NewHandlerCommand(cmder.bot, cmder.Files, ses, update)

	case update.Message != nil:
		log.Printf("Message in HandlerMain: %***REMOVED***v", update.Message)
		ses, _ = cmder.Sessions.GetSession(chatID)
		h = handlers.NewHandlerMessage(cmder.bot, ses, cmder.Products, update)

	default:
		log.Printf("error HandlerMain: update has no supported message: %***REMOVED***v", update)
		return nil
	}

	if h == nil {
		return nil
	}
	h.Execute()

	if err := cmder.Sessions.UpdateSession(chatID, ses); err != nil {
		log.Printf("After session update in HandlerMain: %v\n %v", err, cmder.Sessions)
	}

	return nil
}

func (cmder *Commander) setAvailableCommands(ses *sessions.Session, msg *tgapi.Message) {
	ses.AccessCommand["about"] = struct{}{}

	if cmder.Config.IsAdmin(msg.Chat.UserName) {
		ses.AccessCommand["adm"] = struct{}{}
		return
	}

	if cmder.Config.IsAccess(msg.Chat.UserName) {
		ses.AccessCommand["all"] = struct{}{}
	}

	if msg.CommandArguments() != "" {
		s := strings.Split(msg.CommandArguments(), " ")
		if cmder.Config.IsAccWord(strings.TrimSpace(s[0])) {
			ses.AccessCommand[msg.Command()] = struct{}{}
		}
	}
}

func (cmder *Commander) Stop() {
	cmder.Sessions.Close()
	cmder.Config.Close()
	cmder.Files.Close()
	cmder.bot.StopReceivingUpdates()
	log.Println("Commander stoped.")

}

func (cmder *Commander) Watch(watchTime time.Duration) chan any {
	ok := make(chan any)
	configFile := cmder.Config.ConfFile
	go func() {
		ticker := time.NewTicker(watchTime)
		for {
			<-ticker.C
			if flInfo, err := os.Stat(configFile); err != nil {
				log.Printf("error in conf.Watch: error get FileInfo for  %v.", configFile)
			} else if flInfo.ModTime() != cmder.Config.ModTime {
				log.Printf("Config file %v was modify.", configFile)
				ok <- struct{}{}
			}
		}
	}()
	return ok
}
