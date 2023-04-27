package commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
)

const (
	TXT_VID     string = "Введите Основной вид экономической деятельности (одно 5-ти значное число)."
	TXT_WORKER1 string = `Количество работников с ежемесячым окладом более 10 минимальных заработных плат (далее МЗП) и ГФОТ по этим работникам в тенге через пробел.

Например: у вас 10 работников и ГФОТ 3 000 000. Нужно ввести "10 3000000".`
	TXT_WORKER2 string = "Количество работников с ежемесячым окладом менее или равного 10 МЗП и ГФОТ по этим работникам в тенге через пробел (аналогично ранее введенным данным)."
	TXT_TOTAL   string = `Результат расчета.
Общее количество работников: %v
Основной вид экономической деятельности: %v
Работники с ежемесячым окладом >10 МЗП: %v
ГФОТ работников с окладом >10 МЗП: %v
Работники с ежемесячым окладом <=10 МЗП: %v
ГФОТ работников с окладом <=10 МЗП: %v
----------------
%v
`
	WRONG_TXT string = `Произошла ошибка при расчете. Проверьте введённые данные.`
)

func init() {
	requests_list = append(requests_list, reguest{TXT_VID, (*Commander).get_total_work})
	requests_list = append(requests_list, reguest{TXT_WORKER1, (*Commander).get_vid})
	requests_list = append(requests_list, reguest{TXT_WORKER2, (*Commander).get_workers1})
	requests_list = append(requests_list, reguest{"", (*Commander).get_workers2})
}

func (r *Commander) get_total_work(mes *tgapi.Message) {
	r.Product_service.(*services.Insurance).Total_work = r.str2int(mes)
}

func (r *Commander) get_vid(mes *tgapi.Message) {
	r.Product_service.(*services.Insurance).Vid = r.str2int(mes)
}

func (r *Commander) get_workers1(mes *tgapi.Message) {
	r.Product_service.(*services.Insurance).Workers1, r.Product_service.(*services.Insurance).Gfot1 = r.str2pair(mes)
}

func (r *Commander) get_workers2(mes *tgapi.Message) {
	defer r.resetCommander()

	r.Product_service.(*services.Insurance).Workers2, r.Product_service.(*services.Insurance).Gfot2 = r.str2pair(mes)
	sum, err := (*r).Product_service.Calculate()
	if err != nil {
		msg := tgapi.NewMessage(mes.Chat.ID, WRONG_TXT)
		r.bot.Send(msg)
	}

	str := fmt.Sprintf(TXT_TOTAL, (*r).Product_service.(*services.Insurance).Total_work,
		(*r).Product_service.(*services.Insurance).Vid,
		(*r).Product_service.(*services.Insurance).Workers1,
		(*r).Product_service.(*services.Insurance).Gfot1,
		(*r).Product_service.(*services.Insurance).Workers2,
		(*r).Product_service.(*services.Insurance).Gfot2,
		sum)
	msg := tgapi.NewMessage(mes.Chat.ID, str)
	r.bot.Send(msg)

}

func (r *Commander) str2int(mes *tgapi.Message) int {
	var res int
	var err error
	for {
		res, err = strconv.Atoi(strings.TrimSpace(mes.Text))
		if err == nil {
			return res
		}
		r.bot.Send(tgapi.NewMessage(mes.Chat.ID, WRONG_INPUT))
	}
}

func (r *Commander) str2pair(mes *tgapi.Message) (int, float64) {
	w := 0
	fot := 0.0
	var err error
	for {
		str_arr := strings.Split(strings.TrimSpace(mes.Text), " ")
		if len(str_arr) != 2 {
			r.bot.Send(tgapi.NewMessage(mes.Chat.ID, "Введите ровно два числа разделенные пробелом."))
			continue
		}
		w, err = strconv.Atoi(str_arr[0])
		if err != nil {
			log.Println("Commander: error in get_workers ", err)
			r.bot.Send(tgapi.NewMessage(mes.Chat.ID, WRONG_INPUT))
			continue
		}
		fot, err = strconv.ParseFloat(str_arr[1], 64)
		if err != nil {
			log.Println("Commander: error in get_workers ", err)
			r.bot.Send(tgapi.NewMessage(mes.Chat.ID, WRONG_INPUT))
			continue
		}
		break
	}
	return w, fot
}
