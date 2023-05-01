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
	TXT_TOTAL string = `Введите общее количество работников с учетом работников филиалов (одно число).`

	TXT_VID     string = "Введите Основной вид экономической деятельности (одно 5-ти значное число)."
	TXT_WORKER1 string = `Количество работников с ежемесячым окладом более 10 минимальных заработных плат (далее МЗП) и ГФОТ по этим работникам в тенге через пробел.

Например: у вас 10 работников и ГФОТ 3 000 000. Нужно ввести "10 3000000".`
	TXT_WORKER2 string = "Количество работников с ежемесячым окладом менее или равного 10 МЗП и ГФОТ по этим работникам в тенге через пробел (аналогично ранее введенным данным)."
	TXT_FINISH  string = `*Результат расчета.*
_Общее количество работников:_ %v
_Основной вид экономической деятельности:_ %v
_Работники с ежемесячым окладом >10 МЗП:_ %v
_ГФОТ работников с окладом >10 МЗП:_ %v
_Работники с ежемесячым окладом <=10 МЗП:_ %v
_ГФОТ работников с окладом <=10 МЗП:_ %v
----------------
%v
`
	WRONG_CALC   string = `Произошла ошибка при расчете. Проверьте введённые данные.`
	WRONG_1DIGIT string = "Введите одно число."
	WRONG_2DIGIT string = `Введите ровно 2 числа через пробел.`
	WRONG_5SIGN  string = "Введите одно 5-значное число."
)

func init() {
	requests_list = append(requests_list, reguest{"", (*Commander).get_total_work})
	requests_list = append(requests_list, reguest{TXT_VID, (*Commander).get_vid})
	requests_list = append(requests_list, reguest{TXT_WORKER1, (*Commander).get_workers1})
	requests_list = append(requests_list, reguest{TXT_WORKER2, (*Commander).get_workers2})
}

func (r *Commander) get_total_work(mes *tgapi.Message) {
	// ses, err := r.getSession(mes.Chat.ID)
	// if err != nil {
	// 	log.Printf("error commander.Run: while calc operation %v\n get update: %v", err, mes)
	// } else {
	res, err := r.str2int(mes)
	if err == nil {
		r.Product_service.(*services.Insurance).Total_work = res
		r.Idx++
		r.ErrorInput = false
		return
	}
	r.bot.Send(tgapi.NewMessage(mes.Chat.ID, WRONG_1DIGIT))
	// }
	// ses.
	r.ErrorInput = true
}

func (r *Commander) get_vid(mes *tgapi.Message) {
	// ses, err := r.getSession(mes.Chat.ID)
	// if err != nil {
	// 	log.Printf("error commander.Run: while calc operation %v\n get update: %v", err, mes)
	// } else {
	res, err := r.str2int(mes)
	if err == nil && len(strconv.Itoa(res)) == 5 {
		r.Product_service.(*services.Insurance).Vid = res
		r.Idx++
		r.ErrorInput = false
		return
	}
	r.bot.Send(tgapi.NewMessage(mes.Chat.ID, WRONG_5SIGN))
	// }
	// ses
	r.ErrorInput = true
}

func (r *Commander) get_workers1(mes *tgapi.Message) {
	var err error
	// ses, err := r.getSession(mes.Chat.ID)
	// if err != nil {
	// 	log.Printf("error commander.Run: while calc operation %v\n get update: %v", err, mes)
	// } else {
	r.Product_service.(*services.Insurance).Workers1, r.Product_service.(*services.Insurance).Gfot1, err = r.str2pair(mes)
	if err != nil {
		log.Println("Commander: error in get_workers ", err)
		r.bot.Send(tgapi.NewMessage(mes.Chat.ID, WRONG_2DIGIT))
		r.ErrorInput = true

		return
	}
	r.Idx++
	r.ErrorInput = false
	// }
}

func (r *Commander) get_workers2(mes *tgapi.Message) {
	var err error
	// ses, err := r.getSession(mes.Chat.ID)
	// if err != nil {
	// 	log.Printf("error commander.Run: while calc operation %v\n get update: %v", err, mes)
	// } else {
	r.Product_service.(*services.Insurance).Workers2, r.Product_service.(*services.Insurance).Gfot2, err = r.str2pair(mes)
	if err != nil {
		log.Println("Commander: error in get_workers ", err)
		r.bot.Send(tgapi.NewMessage(mes.Chat.ID, WRONG_2DIGIT))
		r.ErrorInput = true
		return
	}

	finishCalculate(r, mes)
	r.ResetSession()
	r.ErrorInput = false
	// }
}

func finishCalculate(r *Commander, mes *tgapi.Message) {
	sum, err := (*r).Product_service.Calculate()
	if err != nil {
		msg := tgapi.NewMessage(mes.Chat.ID, WRONG_CALC)
		r.bot.Send(msg)
	}

	str := fmt.Sprintf(TXT_FINISH, (*r).Product_service.(*services.Insurance).Total_work,
		(*r).Product_service.(*services.Insurance).Vid,
		(*r).Product_service.(*services.Insurance).Workers1,
		(*r).Product_service.(*services.Insurance).Gfot1,
		(*r).Product_service.(*services.Insurance).Workers2,
		(*r).Product_service.(*services.Insurance).Gfot2,
		sum)
	msg := tgapi.NewMessage(mes.Chat.ID, str)
	//msg.ReplyMarkup = tgapi.NewRemoveKeyboard(true)
	msg.ParseMode = "Markdown"
	r.bot.Send(msg)

}

func (r *Commander) str2int(mes *tgapi.Message) (int, error) {
	res, err := strconv.Atoi(strings.TrimSpace(mes.Text))
	return res, err
}

func (r *Commander) str2pair(mes *tgapi.Message) (int, float64, error) {
	w := 0
	fot := 0.0
	var err error
	str_arr := strings.Split(strings.TrimSpace(mes.Text), " ")
	if len(str_arr) != 2 {
		return w, fot, fmt.Errorf("error in str2pair: lenght != 2")
	}
	w, err = strconv.Atoi(str_arr[0])
	if err != nil {
		return w, fot, err
	}
	fot, err = strconv.ParseFloat(str_arr[1], 64)
	if err != nil {
		return w, fot, err
	}

	return w, fot, err
}
