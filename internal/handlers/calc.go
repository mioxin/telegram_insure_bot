package handlers

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	TXT_BIN string = "Ваш БИН или ИИН."
	TXT     string = `Расчет страховой суммы, страховой премии.
`
)

func (h *HandlerCommands) calc(input_message *tgapi.Message) (string, int) {
	log.Printf("calc: [%s] %s", input_message.From.UserName, input_message.Text)

	msg := tgapi.NewMessage(input_message.Chat.ID, TXT+TXT_BIN)
	m, _ := h.Bot.Send(msg)
	//requests_list[c.Idx].worker(c, input_message)
	return "calc", m.MessageID
}

func init() {
	registered_commands["calc"] = RegisteredCommand{Description: "Расчет страховки ОСНС.", Worker: (*HandlerCommands).calc, ShowInHelp: true}
}
