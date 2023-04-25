package commands

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const ABOUT string = "100% государственная компания. Единственным акционером АО \"КСЖ ГАК\" является государство в лице Правительства Республики Казахстан."

func (c *Commander) about(input_message *tgapi.Message) {
	log.Printf("ABOUT: [%s] %s", input_message.From.UserName, input_message.Text)

	msg := tgapi.NewMessage(input_message.Chat.ID, ABOUT)
	c.bot.Send(msg)
}

func init() {
	registered_commands["about"] = (*Commander).calc
}
