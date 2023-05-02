package commands

import (
	"fmt"
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
_Общее количество работников:_ *%v*
_Основной вид экономической деятельности:_ *%v*
_Работники с ежемесячым окладом >10 МЗП:_ *%v*
_ГФОТ работников с окладом >10 МЗП:_ *%v*
_Работники с ежемесячым окладом <=10 МЗП:_ *%v*
_ГФОТ работников с окладом <=10 МЗП:_ *%v*
----------------
%v
`
	WRONG_CALC   string = `Произошла ошибка при расчете. Проверьте введённые данные и порпобуйте повторить расчет сначала.`
	WRONG_1DIGIT string = "Введите одно число."
	WRONG_2DIGIT string = `Введите ровно 2 числа через пробел.`
	WRONG_5SIGN  string = "Введите одно 5-значное число."
)

func init() {
	requests_list = append(requests_list, reguest{TXT_VID, WRONG_1DIGIT, (*Commander).get_total_work})
	requests_list = append(requests_list, reguest{TXT_WORKER1, WRONG_5SIGN, (*Commander).get_vid})
	requests_list = append(requests_list, reguest{TXT_WORKER2, WRONG_2DIGIT, (*Commander).get_workers1})
	requests_list = append(requests_list, reguest{"", WRONG_2DIGIT, (*Commander).get_workers2})
}

func (r *Commander) get_total_work(mes *tgapi.Message) error {
	var err error
	r.Product_service.(*services.Insurance).Total_work, err = r.str2int(mes)
	if err != nil {
		return err
	}
	return nil
}

func (r *Commander) get_vid(mes *tgapi.Message) error {
	var err error
	r.Product_service.(*services.Insurance).Vid, err = r.str2int(mes)
	if err != nil {
		return err
	}
	if len(mes.Text) != 5 {
		return fmt.Errorf("error: expexted 5 sign only %v", mes.Text)
	}

	return nil
}

func (r *Commander) get_workers1(mes *tgapi.Message) error {
	var err error
	r.Product_service.(*services.Insurance).Workers1, r.Product_service.(*services.Insurance).Gfot1, err = r.str2pair(mes)
	if err != nil {
		return err
	}
	return nil
}

func (r *Commander) get_workers2(mes *tgapi.Message) error {
	var err error
	r.Product_service.(*services.Insurance).Workers2, r.Product_service.(*services.Insurance).Gfot2, err = r.str2pair(mes)
	if err != nil {
		return err
	}

	err = finishCalculate(r, mes)
	if err != nil {
		requests_list[r.Idx].wrong_text = WRONG_CALC
		return err
	}
	return nil
}

func finishCalculate(r *Commander, mes *tgapi.Message) error {
	sum, err := (*r).Product_service.Calculate()
	if err != nil {
		return err
	}

	str := fmt.Sprintf(TXT_FINISH, (*r).Product_service.(*services.Insurance).Total_work,
		(*r).Product_service.(*services.Insurance).Vid,
		(*r).Product_service.(*services.Insurance).Workers1,
		(*r).Product_service.(*services.Insurance).Gfot1,
		(*r).Product_service.(*services.Insurance).Workers2,
		(*r).Product_service.(*services.Insurance).Gfot2,
		sum)
	requests_list[r.Idx].ok_text = str
	// msg := tgapi.NewMessage(mes.Chat.ID, str)
	//msg.ReplyMarkup = tgapi.NewRemoveKeyboard(true)
	// msg.ParseMode = "Markdown"
	// r.bot.Send(msg)
	return nil
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
