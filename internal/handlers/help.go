package handlers

import (
	"fmt"
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const HELP string = `Команды Вы можете набрать с помощью клавиатуры или выбрать в меню.
`

func (h *HandlerCommands) help(input_message *tgapi.Message) (string, int) {
	log.Printf("help: [%s] %s", input_message.From.UserName, input_message.Text)
	str := ""
	for key, cmnd := range registered_commands {
		if cmnd.ShowInHelp {
			str += fmt.Sprintf("/%s - %s\n", key, cmnd.Description)
		}
	}
	msg := tgapi.NewMessage(input_message.Chat.ID, str+HELP)
	m, _ := h.Bot.Send(msg)
	return "", m.MessageID
}

func init() {
	registered_commands["help"] = RegisteredCommand{Description: "Подсказка по командам бота.", Worker: (*HandlerCommands).help, ShowInHelp: true}
}
