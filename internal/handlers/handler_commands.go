package handlers

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/storages/sessions"
)

const WRONG_ACCESS string = "Извините, пока доступ закрыт."

type HandlerCommands struct {
	Handler
}

type ICommander interface {
}

type RegisteredCommand struct {
	Description string
	Worker      func(c *HandlerCommands, mes *tgapi.Message) (string, int)
	ShowInHelp  bool
}

var registered_commands = map[string]RegisteredCommand{}

func NewHandlerCommand(bot *tgapi.BotAPI, ses *sessions.Session, update tgapi.Update) *HandlerCommands {
	return &HandlerCommands{Handler{bot, ses, update}}
}

func (h *HandlerCommands) Execute() {
	if command, ok := registered_commands[h.Update.Message.Command()]; ok {
		_, okAll := h.Ses.AccessCommand["all"]
		_, okCmd := h.Ses.AccessCommand[h.Update.Message.Command()]
		if okAll || okCmd {
			h.Ses.ActionName, h.Ses.LastMessageID = command.Worker(h, h.Update.Message)
		} else {
			msg := tgapi.NewMessage(h.Update.Message.Chat.ID, WRONG_ACCESS)
			h.Bot.Send(msg)
			log.Printf("HandlerCommand: deny acces fo @%s on /%s", h.Update.Message.Chat.UserName, h.Update.Message.Command())
		}
	}
}
