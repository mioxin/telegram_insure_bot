package commands

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const HELP string = `Для ввода данных необходимых, чтобы расчитать стоимость выберите команду /calc.
Вы можете так же выбрать команду /send чтобы отправить сканы документов менеджеру для расчета.
После расчета менеджер сообщит вам результат и размер скидки или кэшбека.
Информация о компании - команда /about.
Все команды Вы можете набрать с помощью клавиатуры или выбрать в меню.
`

func (c *Commander) help(input_message *tgapi.Message) string {
	log.Printf("help: [%s] %s", input_message.From.UserName, input_message.Text)

	msg := tgapi.NewMessage(input_message.Chat.ID, HELP)
	c.bot.Send(msg)
	return ""
}

func init() {
	registered_commands["help"] = (*Commander).help
}
