package handlers

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/storages/sessions"
	"github.com/mrmioxin/gak_telegram_bot/resources"
)

type IFilesID interface {
	GetFileId(name, user string) (string, error)
	SetFileId(name, user, id string) error
	ListUsers() []string
}

type HandlerCommands struct {
	Handler
	FilesId IFilesID
}

type RegisteredCommand struct {
	Description string
	Worker      func(c *HandlerCommands, mes *tgapi.Message) (string, int)
	ShowInHelp  bool
	Adm         bool //command fo admins
}

var registered_commands = map[string]RegisteredCommand{}

func NewHandlerCommand(bot *tgapi.BotAPI, files_id IFilesID, ses *sessions.Session, update tgapi.Update) *HandlerCommands {
	user := update.Message.Chat.UserName
	if user == "" {
		user = update.CallbackQuery.Message.Chat.UserName
	}

	return &HandlerCommands{Handler{bot, ses, update, user}, files_id}
}

func (h *HandlerCommands) Execute() {
	if command, ok := registered_commands[h.Update.Message.Command()]; ok {
		if _, ok := h.Ses.AccessCommand["adm"]; ok {
			h.Ses.ActionName, h.Ses.LastMessageID = command.Worker(h, h.Update.Message)
			return
		}

		if _, ok := h.Ses.AccessCommand["all"]; ok && !command.Adm {
			h.Ses.ActionName, h.Ses.LastMessageID = command.Worker(h, h.Update.Message)
			return
		}

		if _, ok := h.Ses.AccessCommand[h.Update.Message.Command()]; ok {
			h.Ses.ActionName, h.Ses.LastMessageID = command.Worker(h, h.Update.Message)
			return
		}

		msg := tgapi.NewMessage(h.Update.Message.Chat.ID, resources.WRONG_ACCESS)
		h.Bot.Send(msg)
		log.Printf("[%s] HandlerCommand: deny acces on /%s\n", h.User, h.Update.Message.Command())
	}
}
