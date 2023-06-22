package handlers

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/resources"
)

func (h *HandlerCommands) send(input_message *tgapi.Message) (string, int) {
	log.Printf("send: [%s] %s", input_message.From.UserName, input_message.Text)

	msg := tgapi.NewMessage(input_message.Chat.ID, resources.SEND)
	m, _ := h.Bot.Send(msg)
	return "send", m.MessageID
}

func init() {
	registered_commands["send"] = RegisteredCommand{Description: "Отправка информации для расчета страховки менеджером.", Worker: (*HandlerCommands).send, ShowInHelp: true}
}
