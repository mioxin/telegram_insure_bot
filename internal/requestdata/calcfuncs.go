package requestdata

import (
	"strconv"
	"strings"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
)

func init() {

}

func (r *Requester) get_total_work(mes *tgapi.Message) {
	var res int
	var err error
	for {
		res, err = strconv.Atoi(strings.TrimSpace(mes.Text))
		if err == nil {
			break
		}
		r.bot.Send(tgapi.NewMessage(mes.Chat.ID, "Не могу разобрать число. Введете снова или не будем продолжать?"))
	}

	r.Product_service.(*services.Insurance).Total_work = res
	return
}

func (r *Requester) get_vid(mes *tgapi.Message) {
	res := 0
	r.Product_service.(*services.Insurance).Vid = res

	return
}

func (r *Requester) get_workers(mes *tgapi.Message) {
	w := 0
	fot := 0.0
	if r.Product_service.(*services.Insurance).Workers1 < 0 && r.Product_service.(*services.Insurance).Gfot1 < 0 {
		r.Product_service.(*services.Insurance).Workers1 = w
		r.Product_service.(*services.Insurance).Gfot1 = fot
	} else {
		r.Product_service.(*services.Insurance).Workers2 = w
		r.Product_service.(*services.Insurance).Gfot2 = fot
	}
	return
}
