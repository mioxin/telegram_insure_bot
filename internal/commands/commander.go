package commands

import (
	"fmt"
	"log"
	"strings"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
	"github.com/mrmioxin/gak_telegram_bot/internal/sessions"
)

const (
	WRONG_AGAIN  string = `Опять ошибка. `
	WRONG_INPUT  string = `Введите команду или выберите её из меню.`
	WRONG_ACCESS string = "Извините, пока доступ закрыт."
	YES          string = "Да"
	NO           string = "Нет"
)

type reguest struct {
	ok_text    string
	wrong_text string
	worker     func(c *Commander, mes *tgapi.Message) error
}
type RegisteredCommand struct {
	Description string
	Worker      func(c *Commander, mes *tgapi.Message) (string, int)
	ShowInHelp  bool
}

var (
	requestsListCalc    = make([]reguest, 0)
	registered_commands = map[string]RegisteredCommand{}
	yesNoKeyboardCalc   = tgapi.NewInlineKeyboardMarkup(
		tgapi.NewInlineKeyboardRow(
			tgapi.NewInlineKeyboardButtonData(YES, "c yes"),
			tgapi.NewInlineKeyboardButtonData(NO, "c no"),
		))
)

type ITypeOfBusiness interface {
	Get(vid string) (string, error)
}

type IService interface {
	Calculate() (string, error)
}

type IConfig interface {
	IsAccess(user string) bool
}
type ISessions interface {
	GetSession(id int64) (*sessions.Session, error)
	GetIdsByUser(user string) []int64
	UpdateSession(id int64, ses *sessions.Session) error
	AddSession(id int64, ses *sessions.Session)
}
type Commander struct {
	bot             *tgapi.BotAPI
	Config          IConfig
	Product_service IService
	Sessions        ISessions
	TypeOfBuseness  ITypeOfBusiness
}

func NewCommander(bot *tgapi.BotAPI, conf IConfig, serv IService, ses ISessions, tob ITypeOfBusiness) *Commander {
	return &Commander{bot, conf, serv, ses, tob}
}

func (cmder *Commander) HandlerMain(update tgapi.Update) error {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			log.Printf("recover panic in HandlerMain %v:", panicVal)
		}
	}()

	if update.CallbackQuery != nil {
		//processing callbacks...
		cmder.HandlerCallback(update)
		return nil
	}
	if update.Message.IsCommand() {
		cmder.HandlerCommand(update)
	} else {
		ses, err := cmder.Sessions.GetSession(update.Message.Chat.ID)
		if err != nil {
			log.Printf("error HandlerMain: not found session for Message \"%v\" (user %v)", update.Message.Text, update.Message.Chat.UserName)
			return err
		}
		switch ses.ActionName {
		case "calc":
			cmder.HandlerCalc(update, ses)
		default:
		}
	}
	return nil
}

func (cmder *Commander) HandlerCommand(update tgapi.Update) {
	// If we got a message
	if command, ok := registered_commands[update.Message.Command()]; ok {
		ses, err := cmder.Sessions.GetSession(update.Message.Chat.ID)
		if err != nil {
			ses = sessions.NewSession(update.Message.Chat.UserName)
			if cmder.Config.IsAccess(update.Message.Chat.UserName) {
				ses.AccessCommand["all"] = struct{}{}
			} else {
				ses.AccessCommand["about"] = struct{}{}
			}
			cmder.Sessions.AddSession(update.Message.Chat.ID, ses)
		}
		_, okAll := ses.AccessCommand["all"]
		_, okCmd := ses.AccessCommand[update.Message.Command()]
		if okAll || okCmd {
			ses.ActionName, ses.LastMessageID = command.Worker(cmder, update.Message)
			cmder.Sessions.UpdateSession(update.Message.Chat.ID, ses)
			//log.Println(ses)
		} else {
			msg := tgapi.NewMessage(update.Message.Chat.ID, WRONG_ACCESS)
			cmder.bot.Send(msg)
			log.Printf("HandlerCommand: deny acces fo @%s on /%s", update.Message.Chat.UserName, update.Message.Command())
		}
	}
}

func (cmder *Commander) HandlerCalc(update tgapi.Update, ses *sessions.Session) {
	err := requestsListCalc[ses.IdxRequest].worker(cmder, update.Message)
	var m tgapi.Message
	if err != nil {
		if err, ok := err.(ErrorBinIinNotFound); ok {
			log.Printf("HandlerCalc: idx=%v. %v\n", ses.IdxRequest, err)
			mes := tgapi.NewMessage(update.Message.Chat.ID, TXT_VID)
			m, _ = cmder.bot.Send(mes)
			ses.IdxRequest++
			ses.LastMessageID = m.MessageID
			cmder.Sessions.UpdateSession(update.Message.Chat.ID, ses)
			return
		}
		log.Printf("error HandlerCalc: Idx=%v %v", ses.IdxRequest, err)
		if ses.LastRequestIsError {
			m, _ = cmder.bot.Send(tgapi.NewEditMessageText(update.Message.Chat.ID, ses.LastMessageID, WRONG_AGAIN+requestsListCalc[ses.IdxRequest].wrong_text))
		} else {
			m, _ = cmder.bot.Send(tgapi.NewMessage(update.Message.Chat.ID, requestsListCalc[ses.IdxRequest].wrong_text))
		}
		ses.LastRequestIsError = true
	} else {
		//log.Println("HandlerCalc: idx=", ses.IdxRequest)
		mes := tgapi.NewMessage(update.Message.Chat.ID, requestsListCalc[ses.IdxRequest].ok_text)
		mes.ParseMode = "Markdown"
		switch ses.IdxRequest {
		case 0: //skip VID it gotten from the internet by the BIN/IIN
			ses.IdxRequest++
		case 1:
			editText := fmt.Sprintf("Основной вид экономической деятельности: %s - %s\n",
				cmder.Product_service.(*services.Insurance).Vid, cmder.Product_service.(*services.Insurance).VidDescr)
			m, _ = cmder.bot.Send(tgapi.NewEditMessageText(update.Message.Chat.ID, ses.LastMessageID, editText))
		case 3: //need send the inline button Yes|No
			mes.ReplyMarkup = yesNoKeyboardCalc
		default:
		}
		m, _ = cmder.bot.Send(mes)
		ses.IdxRequest++
	}
	ses.LastMessageID = m.MessageID
	ses.LastRequestIsError = false
	// if ses.IdxRequest >= len(requestsListCalc) {
	// 	ses.ResetSession()
	// }
	cmder.Sessions.UpdateSession(update.Message.Chat.ID, ses)
}

func (cmder *Commander) HandlerCallback(update tgapi.Update) {
	ses, err := cmder.Sessions.GetSession(update.CallbackQuery.Message.Chat.ID)
	if err != nil { // clear button
		cmder.bot.Send(tgapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID,
			update.CallbackQuery.Message.Text))
		return
	}
	callbackData := strings.Split(update.CallbackQuery.Data, " ")
	editText := ""
	switch ses.ActionName {
	case "calc":
		log.Println("HandlerCallback: start calc:", callbackData[1])
		if callbackData[1] == "yes" {
			editText = TXT_LAST5YEAR + " *" + YES + "*\n\n"
		} else {
			editText = TXT_LAST5YEAR + " *" + NO + "*\n\n"
		}
		mes := tgapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, editText)
		mes.ParseMode = "Markdown"
		cmder.bot.Send(mes)
		//cmder.Get_yes_no(callbackData[1])

		if txt, err := cmder.Get_yes_no(callbackData[1]); err != nil {
			log.Println("error HandlerCallback: err calc:", WRONG_CALC, err)
			editText = txt + WRONG_CALC
		} else {
			editText = txt
		}
		mes1 := tgapi.NewMessage(update.CallbackQuery.Message.Chat.ID, editText)
		mes1.ParseMode = "Markdown"
		cmder.bot.Send(mes1)

		ses.ResetSession()
		cmder.Sessions.UpdateSession(update.CallbackQuery.Message.Chat.ID, ses)
	default:
	}

}

// func (cmder *Commander) WatchConfig(isModify chan any) {
// 	for range isModify {

// 	}
// }
