package commands

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const TXT_TOTAL string = "Расчет страховой суммы, страховой премии.\n Введите общее количество работников с учетом работников филиалов (одно число)."
const TXT_VID string = "Введите Основной вид экономической деятельности (одно 5-ти значное число)."
const TXT_WORKER1 string = "Количество работников с ежемесячым окладом более 10 минимальных заработных плат (далее МЗП) и ГФОТ по этим работникам в тенге через пробел. Например: у вас 10 работников и ГФОТ 3 000 000. Нужно ввести \"10 3000000\"."
const TXT_WORKER2 string = "Количество работников с ежемесячым окладом менее или равного 10 МЗП и ГФОТ по этим работникам в тенге через пробел (аналогично ранее введенным данным)."

func (c *Commander) calc(input_message *tgapi.Message) {
	var total_work, workers1, gfot1, workers2, gfot2 int
	log.Printf("calc: [%s] %s", input_message.From.UserName, input_message.Text)

	msg := tgapi.NewMessage(input_message.Chat.ID, TXT_TOTAL)
	c.bot.Send(msg)

	c.product_service.Calculate([]int{total_work, workers1, gfot1, workers2, gfot2})
}
