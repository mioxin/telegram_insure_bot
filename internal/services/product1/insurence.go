package product1

import (
	"fmt"
	"log"
	"time"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
	"github.com/mrmioxin/gak_telegram_bot/internal/storages/sessions"
)

const TYPE_OF_BUSNS_FILENAME string = "vid.txt"

type ITypeOfBusiness interface {
	Get(vid string) (string, error)
}

type Insurance struct {
	Name, VidDescr, BinIin, Vid string
	Total_work                  int
	Gfot                        float64
	EventInLast5Year            bool
	TypeOfBuseness              ITypeOfBusiness
	Handler                     *HandlerCalc
}

func NewInsurence(name string) services.IService {
	tob, err := services.NewFileTypeOfBusns(TYPE_OF_BUSNS_FILENAME)
	if err != nil {
		log.Fatal("<<<", TYPE_OF_BUSNS_FILENAME, ">>> ", err)
	}

	return &Insurance{Name: name, TypeOfBuseness: tob}
}

func (ins *Insurance) Calculate() (string, error) {
	sum := 70000.0
	bonus := 23000.0

	//TODO calculate sum
	log.Println("Calculate 5 sec ...")
	time.Sleep(1 * time.Second)
	var str string = fmt.Sprintf("Сумма страховки: *%.2f тенге*\nВаша скидка: *%.2f тенге*", sum, bonus)
	return str, nil
}

func (ins *Insurance) Execute(bot *tgapi.BotAPI, ses *sessions.Session, update tgapi.Update) {
	ins.Handler = NewHandlerCalc(ins, bot, ses, update)
	if update.CallbackQuery != nil {
		ins.Handler.ExecuteCallback()
	} else {
		ins.Handler.Execute()
	}
}
