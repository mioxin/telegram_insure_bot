package handlers

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/resources"
)

func (h *HandlerCommands) get(input_message *tgapi.Message) (string, int) {
	log.Printf("get: [%s] %s", input_message.From.UserName, input_message.Text)

	msg := tgapi.NewMessage(input_message.Chat.ID, resources.GET)

	buttons := make([][]tgapi.InlineKeyboardButton, 0)
	for _, k := range h.FilesId.ListUsers() {
		keyboardRow := tgapi.NewInlineKeyboardRow(tgapi.NewInlineKeyboardButtonData(k, "user:"+k))
		buttons = append(buttons, keyboardRow)
	}
	msg.ReplyMarkup = tgapi.NewInlineKeyboardMarkup(buttons...)

	m, _ := h.Bot.Send(msg)
	return "get", m.MessageID
}

func init() {
	registered_commands["get"] = RegisteredCommand{
		Description: "Работа с файлами, отправленными боту клиентами.",
		Worker:      (*HandlerCommands).get,
		ShowInHelp:  false,
		Adm:         true,
	}
}
