package requestdata

import (
	"log"
	"strconv"
	"strings"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
)

const WRONG_INPUT string = "Похоже вы ввели не число. Введете снова."

func init() {
	requests_list = append(requests_list, reguest{"get_total_work", (*Requester).get_total_work})
	requests_list = append(requests_list, reguest{"get_vid", (*Requester).get_vid})
	requests_list = append(requests_list, reguest{"get_workers", (*Requester).get_workers1})
	requests_list = append(requests_list, reguest{"get_workers", (*Requester).get_workers2})
}

func (r *Requester) get_total_work(mes *tgapi.Message) {
	r.Product_service.(*services.Insurance).Total_work = r.str2int(mes)
}

func (r *Requester) get_vid(mes *tgapi.Message) {
	r.Product_service.(*services.Insurance).Vid = r.str2int(mes)
}

func (r *Requester) get_workers1(mes *tgapi.Message) {
	r.Product_service.(*services.Insurance).Workers1, r.Product_service.(*services.Insurance).Gfot1 = r.str2pair(mes)
}

func (r *Requester) get_workers2(mes *tgapi.Message) {
	r.Product_service.(*services.Insurance).Workers2, r.Product_service.(*services.Insurance).Gfot2 = r.str2pair(mes)
}

func (r *Requester) str2int(mes *tgapi.Message) int {
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

func (r *Requester) str2pair(mes *tgapi.Message) (int, float64) {
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
			log.Println("requester: error in get_workers ", err)
			r.bot.Send(tgapi.NewMessage(mes.Chat.ID, WRONG_INPUT))
			continue
		}
		fot, err = strconv.ParseFloat(str_arr[1], 64)
		if err == nil {
			log.Println("requester: error in get_workers ", err)
			r.bot.Send(tgapi.NewMessage(mes.Chat.ID, WRONG_INPUT))
			continue
		}
		break
	}
	return w, fot
}
