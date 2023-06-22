package handlers

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/storages/filesid"
	"github.com/mrmioxin/gak_telegram_bot/internal/storages/sessions"
	"github.com/mrmioxin/gak_telegram_bot/resources"
)

type IFilesID interface {
	GetFileId(name string) (string, error)
	SetFileId(name, id string)
}

type HandlerCommands struct {
	Handler
	FilesId IFilesID
}

type RegisteredCommand struct {
	Description string
	Worker      func(c *HandlerCommands, mes *tgapi.Message) (string, int)
	ShowInHelp  bool
}

var registered_commands = map[string]RegisteredCommand{}

func NewHandlerCommand(bot *tgapi.BotAPI, ses *sessions.Session, update tgapi.Update) *HandlerCommands {
	files_id := filesid.NewMapFilesId()
	return &HandlerCommands{Handler{bot, ses, update}, files_id}
}

func (h *HandlerCommands) Execute() {
	if command, ok := registered_commands[h.Update.Message.Command()]; ok {
		_, okAll := h.Ses.AccessCommand["all"]
		_, okCmd := h.Ses.AccessCommand[h.Update.Message.Command()]
		if okAll || okCmd {
			h.Ses.ActionName, h.Ses.LastMessageID = command.Worker(h, h.Update.Message)
		} else {
			msg := tgapi.NewMessage(h.Update.Message.Chat.ID, resources.WRONG_ACCESS)
			h.Bot.Send(msg)
			log.Printf("HandlerCommand: deny acces fo @%s on /%s\n", h.Update.Message.Chat.UserName, h.Update.Message.Command())
		}
	}
}
