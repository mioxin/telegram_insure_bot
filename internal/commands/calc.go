package commands

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/requestdata"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
)

const (
	TXT_TOTAL string = `Расчет страховой суммы, страховой премии.

Введите общее количество работников с учетом работников филиалов (одно число).`
	TXT_VID     string = "Введите Основной вид экономической деятельности (одно 5-ти значное число)."
	TXT_WORKER1 string = `Количество работников с ежемесячым окладом более 10 минимальных заработных плат (далее МЗП) и ГФОТ по этим работникам в тенге через пробел.

Например: у вас 10 работников и ГФОТ 3 000 000. Нужно ввести "10 3000000".`
	TXT_WORKER2 string = "Количество работников с ежемесячым окладом менее или равного 10 МЗП и ГФОТ по этим работникам в тенге через пробел (аналогично ранее введенным данным)."
)

// type DataForCalc struct {
// 	Total_work, Vid, Workers1, Workers2 int
// 	gfot1, gfot2                        float64
// }

func (c *Commander) calc(input_message *tgapi.Message) {
	log.Printf("calc: [%s] %s", input_message.From.UserName, input_message.Text)

	insure := services.NewInsurence("ОСНС", 1000.00)
	reqer := requestdata.NewRequester(c.bot, input_message.Chat.ID, insure)

	if err := reqer.Run(); err != nil {
		log.Printf("calc: error run requester %v", err)
	}

	if sum, err := (*reqer).Product_service.Calculate(); err != nil {
		msg := tgapi.NewMessage(input_message.Chat.ID, sum)
		c.bot.Send(msg)
	}
}

func init() {
	registered_commands["calc"] = (*Commander).calc
}
