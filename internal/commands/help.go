package commands

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const HELP string = `Уважаемый предприниматель!
Наш бот предлагает услуги менеджера Государственной Аннуитетной Компании в г. Петропавловск (Казахстан).
Вы можете получить расчет стоимости Обязательного страхования работника от несчестного случая (ОСНС).
Предусмотрены скидки или кэшбек.

Для ввода данных необходимых, чтобы расчитать стоимость выберите команду /calc.
Вы можете так же выбрать команду /send чтобы отправить сканы документов менеджеру для расчета.
После расчета менеджер сообщит вам результат и размер скидки или кэшбека.
Для прямой сязи с менеджером номер +7(701)172-67-88 (Whatsapp)
`

func (c *Commander) help(input_message *tgapi.Message) {
	log.Printf("help: [%s] %s", input_message.From.UserName, input_message.Text)

	msg := tgapi.NewMessage(input_message.Chat.ID, HELP)
	c.bot.Send(msg)
}

func init() {
	registered_commands["help"] = (*Commander).help
}
