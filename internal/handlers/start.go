package handlers

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/resources"
)

func (h *HandlerCommands) start(input_message *tgapi.Message) (string, int) {
	log.Printf("start: [%s] %s", input_message.From.UserName, input_message.Text)

	msg := tgapi.NewMessage(input_message.Chat.ID, resources.START)
	m, _ := h.Bot.Send(msg)
	return "", m.MessageID
}

func init() {
	registered_commands["start"] = RegisteredCommand{Description: "Первоначальная информация при подключении к боту.", Worker: (*HandlerCommands).start, ShowInHelp: true}
}
