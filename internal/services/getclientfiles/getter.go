package getclientfiles

import (
	"fmt"
	"log"
	"strings"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/storages"
	"github.com/mrmioxin/gak_telegram_bot/internal/storages/sessions"
	"github.com/mrmioxin/gak_telegram_bot/resources"
)

type IFileStorage interface {
	GetFileId(name, user string) (string, error)
	ListFiles(user string) []*storages.FileInfo
}

type Getter struct {
	FileStore IFileStorage
}

func NewGetter(store IFileStorage) *Getter {
	return &Getter{store}
}

func (g *Getter) Execute(bot *tgapi.BotAPI, ses *sessions.Session, update tgapi.Update) {
	var mId int
	var err error
	ses.LastRequestIsError = false

	key := ""
	if update.CallbackQuery != nil {
		chatId := update.CallbackQuery.Message.Chat.ID
		arr := strings.Split(update.CallbackQuery.Data, ":")
		log.Printf("getter Execute: slice CallbackQuery.Data %s", arr)

		if len(arr) > 1 {
			key = strings.TrimSpace(arr[0])
		} else {
			log.Printf("error getter Execute: CallbackQuery.Data %s", update.CallbackQuery.Data)
			mId, _ = sendSysError(bot, ses, chatId)
		}
		switch key {
		case "user":
			mId, err = g.UploadListFiles(bot, arr[1], chatId)
			if err != nil {
				log.Printf("error getter Execute: %s", err)
				mId, _ = sendSysError(bot, ses, chatId)
			}
		case "f": //file
			f := strings.Split(arr[1], ";")
			fileID, _ := g.FileStore.GetFileId(f[0], f[1])
			mId, _ = g.UploadFile(bot, fileID, chatId)
			if err != nil {
				log.Printf("error getter Execute: %s", err)
			}
		default:
		}
	}
	ses.LastMessageID = mId

}

func (g *Getter) UploadListFiles(bot *tgapi.BotAPI, user string, chatId int64) (int, error) {
	msg := tgapi.NewMessage(chatId, fmt.Sprintf(resources.LIST_FILES, user))

	buttons := make([][]tgapi.InlineKeyboardButton, 0)
	for _, k := range g.FileStore.ListFiles(user) {
		keyboardRow := tgapi.NewInlineKeyboardRow(tgapi.NewInlineKeyboardButtonData(k.FileName, fmt.Sprintf("f:%s;%d", user, k.Id)))
		buttons = append(buttons, keyboardRow)
	}

	msg.ReplyMarkup = tgapi.NewInlineKeyboardMarkup(buttons...)

	m, err := bot.Send(msg)
	if err != nil {
		log.Printf("getter UploadListFile: buttons: %#v", buttons)
		return 0, err
	}
	return m.MessageID, nil
}

func (g *Getter) UploadFile(bot *tgapi.BotAPI, fileID string, chatId int64) (int, error) {
	log.Printf("getter UploadFile: id: %s", fileID)

	m, err := bot.Send(tgapi.NewDocumentShare(chatId, fileID))

	if err != nil { // file isnt document
		log.Printf("error getter UploadFile: %s, chat ID:%v fileID:%v", err, chatId, fileID)

		m, err = bot.Send(tgapi.NewPhotoShare(chatId, fileID))
		if err != nil {
			log.Printf("error getter UploadFile: %s, chat ID:%v fileID:%v", err, chatId, fileID)
			return 0, err
		}
	}
	return m.MessageID, nil

}

func sendSysError(bot *tgapi.BotAPI, ses *sessions.Session, chatId int64) (int, error) {
	var msg tgapi.EditMessageTextConfig
	if ses.LastRequestIsError {
		msg = tgapi.NewEditMessageText(chatId, ses.LastMessageID, resources.WRONG_AGAIN+resources.WRONG_SYS)
	} else {
		msg = tgapi.NewEditMessageText(chatId, ses.LastMessageID, resources.WRONG_SYS)
	}
	m, err := bot.Send(msg)
	if err != nil {
		return 0, err
	}
	ses.LastRequestIsError = true
	return m.MessageID, nil
}
