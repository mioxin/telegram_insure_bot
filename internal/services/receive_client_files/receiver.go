package receive_client_files

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/storages/sessions"
	"github.com/mrmioxin/gak_telegram_bot/resources"
)

type IFileStorage interface {
	GetFileId(name, user string) (string, error)
	SetFileId(name, user, id string) error
}

type Recever struct {
	FileStore IFileStorage
}

func NewReceiver(store IFileStorage) *Recever {
	return &Recever{store}
}

func (r *Recever) Execute(bot *tgapi.BotAPI, ses *sessions.Session, update tgapi.Update) {
	var msg tgapi.Chattable
	var err error
	switch { //expected doc or img
	case update.Message.Photo != nil || update.Message.Document != nil:
		ses.LastRequestIsError = false
		if msg, err = r.saveFile(ses, update); err != nil {
			log.Printf("error in Receiver Execute: %v\n", err)
			ses.LastRequestIsError = true
		}
		ses.ActionName = ""
	default: //not doc or img
		log.Printf("error in Receiver Execute: expected attached a doc or img files\n")
		if ses.LastRequestIsError { //error was been in last time
			msg = tgapi.NewEditMessageText(update.Message.Chat.ID, ses.LastMessageID, resources.WRONG_AGAIN+" "+resources.WRONG_GET_FILES)
		} else {
			msg = tgapi.NewMessage(update.Message.Chat.ID, resources.WRONG_GET_FILES)
		}
		ses.LastRequestIsError = true
	}

	m, _ := bot.Send(msg)
	ses.LastMessageID = m.MessageID
}

func (r *Recever) saveFile(ses *sessions.Session, update tgapi.Update) (tgapi.Chattable, error) {
	var msg tgapi.Chattable
	var err error
	var userName, fileName, fileId string
	strMsg := resources.TXT_GET_FILES

	userName = update.Message.Chat.UserName
	switch {
	case update.Message.Photo != nil:
		len_arr := len(*update.Message.Photo)
		len_id := len((*update.Message.Photo)[len_arr-1].FileID)
		addname := (*update.Message.Photo)[len_arr-1].FileID
		if len_id > 10 {
			addname = (*update.Message.Photo)[len_arr-1].FileID[len_id-10:]
		}
		fileName = "typePhoto_" + addname
		fileId = (*update.Message.Photo)[len_arr-1].FileID
	case update.Message.Document != nil:
		fileName = update.Message.Document.FileName
		fileId = update.Message.Document.FileID
	default:
	}
	//save file info to map & json
	r.FileStore.SetFileId(fileName, userName, fileId)

	if ses.LastRequestIsError { //error was been in last time
		msg = tgapi.NewEditMessageText(update.Message.Chat.ID, ses.LastMessageID, resources.WRONG_AGAIN+" "+strMsg)
	} else {
		msg = tgapi.NewMessage(update.Message.Chat.ID, strMsg)
	}

	return msg, err
}
