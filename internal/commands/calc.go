package commands

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	TXT string = `Расчет страховой суммы, страховой премии.
	
Введите общее количество работников с учетом работников филиалов (одно число).`
)

func (c *Commander) calc(input_message *tgapi.Message) string {
	log.Printf("calc: [%s] %s", input_message.From.UserName, input_message.Text)

	msg := tgapi.NewMessage(input_message.Chat.ID, TXT)
	c.bot.Send(msg)

	return "calc"
}

func init() {
	registered_commands["calc"] = (*Commander).calc
}
