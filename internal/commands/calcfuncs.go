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
	TXT_BIN       string = "Ваш БИН или ИИН."
	TXT_VID       string = "Введите Основной вид экономической деятельности (одно 5-ти значное число)."
	TXT_TOTAL     string = `Введите общее количество работников с учетом работников филиалов (одно число).`
	TXT_GFOT      string = `Введите ГФОТ.`
	TXT_LAST5YEAR string = `Были ли страховые случаи за последние 5 лет?`
	TXT_FINISH    string = `*Результат расчета.*
----------------
%v
`
	// _Общее количество работников:_ *%v*
	// _Основной вид экономической деятельности:_ *%v*
	// _Работники с ежемесячым окладом >10 МЗП:_ *%v*
	// _ГФОТ работников с окладом >10 МЗП:_ *%v*
	// _Работники с ежемесячым окладом <=10 МЗП:_ *%v*
	// _ГФОТ работников с окладом <=10 МЗП:_ *%v*
	WRONG_CALC   string = `Произошла ошибка при расчете. Проверьте введённые данные и порпобуйте повторить расчет сначала.`
	WRONG_1DIGIT string = "Введите одно число."
	WRONG_5SIGN  string = "Введите одно 5-значное число."
	WRONG_BIN    string = "ИИН или БИН введен не корректно."
)

type ErrorBinIinNotFound struct {
	binIin string
}

func (e ErrorBinIinNotFound) Error() string {
	return fmt.Sprintf("error BIN/IIN %v not found on https://stat.gov.kz/", e.binIin)
}

func init() {
	requestsListCalc = append(requestsListCalc, reguest{TXT_TOTAL, WRONG_BIN, (*Commander).binIin})
	requestsListCalc = append(requestsListCalc, reguest{TXT_TOTAL, WRONG_5SIGN, (*Commander).oked})
	requestsListCalc = append(requestsListCalc, reguest{TXT_GFOT, WRONG_1DIGIT, (*Commander).totalWorker})
	requestsListCalc = append(requestsListCalc, reguest{TXT_LAST5YEAR, WRONG_1DIGIT, (*Commander).gfot})
	//requestsListCalc = append(requestsListCalc, reguest{"", "", (*Commander).finishCalculate})
}

func (r *Commander) binIin(updMes *tgapi.Message) error {
	if !okBinIin(updMes.Text) {
		return fmt.Errorf("error binIin: invalid bin/iin %v", updMes.Text)
	}
	if comp, err := services.NewCompany(strings.TrimSpace(updMes.Text)); err != nil {
		return ErrorBinIinNotFound{updMes.Text}
	} else {
		r.Product_service.(*services.Insurance).Vid = comp.OkedCode
		r.Product_service.(*services.Insurance).VidDescr = comp.OkedName
		log.Println("binIin: user:", updMes.Chat.UserName, comp)
	}
	return nil
}

func (r *Commander) oked(updMes *tgapi.Message) error {
	if len(updMes.Text) != 5 {
		return fmt.Errorf("error oked: expexted 5 sign only %v", updMes.Text)
	}

	vid := strings.TrimSpace(updMes.Text)
	if descr, err := r.TypeOfBuseness.Get(vid); err != nil {
		return err
	} else {
		r.Product_service.(*services.Insurance).VidDescr = descr
	}
	r.Product_service.(*services.Insurance).Vid = vid
	return nil
}

func (r *Commander) totalWorker(updMes *tgapi.Message) error {
	var err error
	r.Product_service.(*services.Insurance).Total_work, err = r.str2int(updMes)
	if err != nil {
		return err
	}
	return nil
}

func (r *Commander) gfot(updMes *tgapi.Message) error {
	var err error
	r.Product_service.(*services.Insurance).Gfot, err = r.str2float(updMes)
	if err != nil {
		return err
	}
	return nil
}

func (r *Commander) Get_yes_no(callbackData string) (string, error) {
	var err error
	log.Println("Get_yes_no: start", callbackData, err)
	if callbackData == "yes" {
		r.Product_service.(*services.Insurance).EventInLast5Year = true
	}

	sum, err := (*r).Product_service.Calculate()
	if err != nil {
		return "", err
	}

	str := fmt.Sprintf(TXT_FINISH, sum)
	return str, nil
}

func (r *Commander) str2int(updMes *tgapi.Message) (int, error) {
	res, err := strconv.Atoi(strings.TrimSpace(updMes.Text))
	return res, err
}
func (r *Commander) str2float(updMes *tgapi.Message) (float64, error) {
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
