package product1

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
	"github.com/mrmioxin/gak_telegram_bot/internal/storages/sessions"
	"github.com/mrmioxin/gak_telegram_bot/resources"
)

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
	tob, err := services.NewFileTypeOfBusns(resources.TYPE_OF_BUSNS_FILENAME)
	if err != nil {
		log.Fatal("<<<", resources.TYPE_OF_BUSNS_FILENAME, ">>> ", err)
	}

	return &Insurance{Name: name, TypeOfBuseness: tob}
}

func (ins *Insurance) Execute(bot *tgapi.BotAPI, ses *sessions.Session, update tgapi.Update) {
	ins.Handler = NewHandlerCalc(ins, bot, ses, update)
	if update.CallbackQuery != nil {
		ins.Handler.ExecuteCallback()
	} else {
		ins.Handler.Execute()
	}
}
