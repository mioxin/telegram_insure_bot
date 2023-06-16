package handlers

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const START string = `Уважаемый предприниматель!
Наш бот предлагает услуги менеджера Государственной Аннуитетной Компании в г. Петропавловск (Казахстан).
Он поможет Вам расчитать стоимость Обязательного страхования работника от несчастного случая (ОСНС).
Вы можете так же отправить сканы документов менеджеру для расчета. После чего менеджер сообщит вам результат.
Предусмотрены скидки или кэшбек.

Для дополниельной подсказки выберите команду /help.

Для прямой сязи с менеджером: номер +7(701)172-67-88 (Whatsapp)
`

func (h *HandlerCommands) start(input_message *tgapi.Message) (string, int) {
	log.Printf("start: [%s] %s", input_message.From.UserName, input_message.Text)

	msg := tgapi.NewMessage(input_message.Chat.ID, START)
	m, _ := h.bot.Send(msg)
	return "", m.MessageID
}

func init() {
	registered_commands["start"] = RegisteredCommand{Description: "Первоначальная информация при подключении к боту.", Worker: (*HandlerCommands).start, ShowInHelp: true}
}
