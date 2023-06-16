package handlers

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/sessions"
)

const WRONG_ACCESS string = "Извините, пока доступ закрыт."

type HandlerCommands struct {
	bot    *tgapi.BotAPI
	ses    *sessions.Session
	update tgapi.Update
}

type ICommander interface {
}

type RegisteredCommand struct {
	Description string
	Worker      func(c *HandlerCommands, mes *tgapi.Message) (string, int)
	ShowInHelp  bool
}

var registered_commands = map[string]RegisteredCommand{}

func NewHandlerCommand(bot *tgapi.BotAPI, ses *sessions.Session, upd tgapi.Update) *HandlerCommands {
	return &HandlerCommands{bot, ses, upd}
}

func (h *HandlerCommands) Execute() {
	// If we got a message
	if command, ok := registered_commands[h.update.Message.Command()]; ok {
		_, okAll := h.ses.AccessCommand["all"]
		_, okCmd := h.ses.AccessCommand[h.update.Message.Command()]
		if okAll || okCmd {
			h.ses.ActionName, h.ses.LastMessageID = command.Worker(h, h.update.Message)
			//h.Sessions.UpdateSession(h.update.Message.Chat.ID, ses)
			//log.Println(ses)
		} else {
			msg := tgapi.NewMessage(h.update.Message.Chat.ID, WRONG_ACCESS)
			h.bot.Send(msg)
			log.Printf("HandlerCommand: deny acces fo @%s on /%s", h.update.Message.Chat.UserName, h.update.Message.Command())
		}
	}
}
