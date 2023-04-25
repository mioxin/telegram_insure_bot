package commands

import (
	"log"
	"strconv"
	"strings"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const TXT_TOTAL string = `Расчет страховой суммы, страховой премии.

Введите общее количество работников с учетом работников филиалов (одно число).`
const TXT_VID string = "Введите Основной вид экономической деятельности (одно 5-ти значное число)."
const TXT_WORKER1 string = `Количество работников с ежемесячым окладом более 10 минимальных заработных плат (далее МЗП) и ГФОТ по этим работникам в тенге через пробел.

Например: у вас 10 работников и ГФОТ 3 000 000. Нужно ввести "10 3000000".`
const TXT_WORKER2 string = "Количество работников с ежемесячым окладом менее или равного 10 МЗП и ГФОТ по этим работникам в тенге через пробел (аналогично ранее введенным данным)."

func (c *Commander) calc(input_message *tgapi.Message) {
	var total_work, vid, workers1, workers2 int
	var gfot1, gfot2 float64
	var err error
	log.Printf("calc: [%s] %s", input_message.From.UserName, input_message.Text)

	msg := tgapi.NewMessage(input_message.Chat.ID, TXT_TOTAL)
	c.bot.Send(msg)
	if total_work, err = c.get_total_work(input_message); err != nil {
		return
	}

	c.bot.Send(tgapi.NewMessage(input_message.Chat.ID, TXT_VID))
	if vid, err = c.get_vid(input_message); err != nil {
		return
	}

	c.bot.Send(tgapi.NewMessage(input_message.Chat.ID, TXT_WORKER1))
	if workers1, gfot1, err = c.get_workers(input_message); err != nil {
		return
	}

	c.bot.Send(tgapi.NewMessage(input_message.Chat.ID, TXT_WORKER2))
	if workers2, gfot2, err = c.get_workers(input_message); err != nil {
		return
	}

	c.product_service.Calculate([]any{total_work, vid, workers1, gfot1, workers2, gfot2})
}

func (c *Commander) get_total_work(mes *tgapi.Message) (int, error) {
	var res int
	var err error
	//for {

	res, err = strconv.Atoi(strings.TrimSpace(mes.Text))
	if err != nil {
		c.bot.Send(tgapi.NewMessage(mes.Chat.ID, "Не могу разобрать число. Введете снова или не будем продолжать?"))
	}
	//}
	return res, err
}

func (c *Commander) get_vid(mes *tgapi.Message) (int, error) {
	res := 0
	return res, nil
}

func (c *Commander) get_workers(mes *tgapi.Message) (int, float64, error) {
	w := 0
	fot := 0.0
	return w, fot, nil
}

func init() {
	registered_commands["calc"] = (*Commander).calc
}
