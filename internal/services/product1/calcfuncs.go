package product1

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
	"github.com/mrmioxin/gak_telegram_bot/resources"
)

type ErrorBinIinNotFound struct {
	binIin string
}

func (e ErrorBinIinNotFound) Error() string {
	return fmt.Sprintf("error BIN/IIN %v not found on https://stat.gov.kz/", e.binIin)
}

func init() {
	requestsListCalc = append(requestsListCalc, reguest{resources.TXT_TOTAL, resources.WRONG_BIN, (*HandlerCalc).binIin})
	requestsListCalc = append(requestsListCalc, reguest{resources.TXT_TOTAL, resources.WRONG_5SIGN, (*HandlerCalc).oked})
	requestsListCalc = append(requestsListCalc, reguest{resources.TXT_GFOT, resources.WRONG_1DIGIT, (*HandlerCalc).totalWorker})
	requestsListCalc = append(requestsListCalc, reguest{resources.TXT_LAST5YEAR, resources.WRONG_1DIGIT, (*HandlerCalc).gfot})
	//requestsListCalc = append(requestsListCalc, reguest{"", "", (*HandlerCalc).finishCalculate})
}

func (r *HandlerCalc) binIin(updMes *tgapi.Message) error {
	if !okBinIin(updMes.Text) {
		return fmt.Errorf("error binIin: invalid bin/iin %v", updMes.Text)
	}
	if comp, err := services.NewCompany(strings.TrimSpace(updMes.Text)); err != nil {
		return ErrorBinIinNotFound{updMes.Text}
	} else {
		r.ins.Vid = comp.OkedCode
		r.ins.VidDescr = comp.OkedName
		log.Println("binIin: user:", updMes.Chat.UserName, comp)
	}
	return nil
}

func (r *HandlerCalc) oked(updMes *tgapi.Message) error {
	if len(updMes.Text) != 5 {
		return fmt.Errorf("error oked: expexted 5 sign only %v", updMes.Text)
	}

	vid := strings.TrimSpace(updMes.Text)
	if descr, err := r.ins.TypeOfBuseness.Get(vid); err != nil {
		return err
	} else {
		r.ins.VidDescr = descr
	}
	r.ins.Vid = vid
	return nil
}

func (r *HandlerCalc) totalWorker(updMes *tgapi.Message) error {
	var err error
	r.ins.Total_work, err = r.str2int(updMes)
	if err != nil {
		return err
	}
	return nil
}

func (r *HandlerCalc) gfot(updMes *tgapi.Message) error {
	var err error
	r.ins.Gfot, err = r.str2float(updMes)
	if err != nil {
		return err
	}
	return nil
}

func (r *HandlerCalc) Get_yes_no(callbackData string) (string, error) {
	var err error
	log.Println("Get_yes_no: start", callbackData, err)
	if callbackData == "yes" {
		r.ins.EventInLast5Year = true
	}

	sum, err := r.Calculate()
	if err != nil {
		return "", err
	}

	str := fmt.Sprintf(resources.TXT_FINISH, sum)
	return str, nil
}

func (r *HandlerCalc) Calculate() (string, error) {
	sum := 70000.0
	bonus := 23000.0

	//TODO calculate sum
	log.Println("Calculate 5 sec ...")
	time.Sleep(1 * time.Second)
	var str string = fmt.Sprintf("Сумма страховки: *%.2f тенге*\nВаша скидка: *%.2f тенге*", sum, bonus)
	return str, nil
}

func (r *HandlerCalc) str2int(updMes *tgapi.Message) (int, error) {
	res, err := strconv.Atoi(strings.TrimSpace(updMes.Text))
	return res, err
}
func (r *HandlerCalc) str2float(updMes *tgapi.Message) (float64, error) {
	res, err := strconv.ParseFloat(strings.TrimSpace(updMes.Text), 64)
	return res, err
}

func okBinIin(biniin string) bool {
	var ok bool
	arrBin := make([]int, 0)
	if len(biniin) != 12 {
		return false
	}
	for _, ch := range biniin {
		x, err := strconv.Atoi(string(ch))
		if err != nil {
			return false
		}
		arrBin = append(arrBin, x)
	}
	sum := 0
	for i := 1; i < 12; i++ {
		sum += i * arrBin[i-1]
	}
	if sum%11 == 10 {
		sum = 0
		k := 0
		for i := 1; i < 12; i++ {
			k = i + 2
			if k > 11 {
				k = k % 11
			}
			sum += k * arrBin[i-1]
		}
	}
	if sum%11 == arrBin[11] {
		ok = true
	}
	return ok
}
