package product1

import (
	"fmt"
	"log"
	"strings"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/handlers"
	"github.com/mrmioxin/gak_telegram_bot/internal/storages/sessions"
)

const (
	WRONG_AGAIN   string = `Опять ошибка. `
	TXT_LAST5YEAR string = `Были ли страховые случаи за последние 5 лет?`
	WRONG_CALC    string = `Произошла ошибка при расчете. Проверьте введённые данные и порпобуйте повторить расчет сначала.`

	YES string = "Да"
	NO  string = "Нет"
)

type HandlerCalc struct {
	ins *Insurance
	handlers.Handler
}

type reguest struct {
	ok_text    string
	wrong_text string
	worker     func(c *HandlerCalc, mes *tgapi.Message) error
}

var (
	requestsListCalc  = make([]reguest, 0)
	yesNoKeyboardCalc = tgapi.NewInlineKeyboardMarkup(
		tgapi.NewInlineKeyboardRow(
			tgapi.NewInlineKeyboardButtonData(YES, "c yes"),
			tgapi.NewInlineKeyboardButtonData(NO, "c no"),
		))
)

func NewHandlerCalc(ins *Insurance, bot *tgapi.BotAPI, ses *sessions.Session, update tgapi.Update) *HandlerCalc {
	return &HandlerCalc{ins, handlers.Handler{Bot: bot, Ses: ses, Update: update}}
}

func (h *HandlerCalc) Execute() {
	err := requestsListCalc[h.Ses.IdxRequest].worker(h, h.Update.Message)
	var m tgapi.Message
	if err != nil {
		if err, ok := err.(ErrorBinIinNotFound); ok {
			log.Printf("HandlerCalc: idx=%v. %v\n", h.Ses.IdxRequest, err)
			mes := tgapi.NewMessage(h.Update.Message.Chat.ID, TXT_VID)
			m, _ = h.Bot.Send(mes)
			h.Ses.IdxRequest++
			h.Ses.LastMessageID = m.MessageID
			//h.Sessions.UpdateSession(h.Update.Message.Chat.ID, ses)
			return
		}
		log.Printf("error HandlerCalc: Idx=%v %v", h.Ses.IdxRequest, err)
		if h.Ses.LastRequestIsError {
			m, _ = h.Bot.Send(tgapi.NewEditMessageText(h.Update.Message.Chat.ID, h.Ses.LastMessageID, WRONG_AGAIN+requestsListCalc[h.Ses.IdxRequest].wrong_text))
		} else {
			m, _ = h.Bot.Send(tgapi.NewMessage(h.Update.Message.Chat.ID, requestsListCalc[h.Ses.IdxRequest].wrong_text))
		}
		h.Ses.LastRequestIsError = true
	} else {
		//log.Println("HandlerCalc: idx=", h.Ses.IdxRequest)
		mes := tgapi.NewMessage(h.Update.Message.Chat.ID, requestsListCalc[h.Ses.IdxRequest].ok_text)
		mes.ParseMode = "Markdown"
		switch h.Ses.IdxRequest {
		case 0: //skip VID it gotten from the internet by the BIN/IIN
			h.Ses.IdxRequest++
		case 1:
			editText := fmt.Sprintf("Основной вид экономической деятельности: %s - %s\n",
				h.ins.Vid, h.ins.VidDescr)
			m, _ = h.Bot.Send(tgapi.NewEditMessageText(h.Update.Message.Chat.ID, h.Ses.LastMessageID, editText))
		case 3: //need send the inline button Yes|No
			mes.ReplyMarkup = yesNoKeyboardCalc
		default:
		}
		m, _ = h.Bot.Send(mes)
		h.Ses.IdxRequest++
	}
	h.Ses.LastMessageID = m.MessageID
	h.Ses.LastRequestIsError = false

	//h.Sessions.UpdateSession(h.Update.Message.Chat.ID, ses)
}

func (h *HandlerCalc) ExecuteCallback() {
	callbackData := strings.Split(h.Update.CallbackQuery.Data, " ")
	editText := ""
	switch h.Ses.ActionName {
	case "calc":
		log.Println("HandlerCallback: start calc:", callbackData[1])
		if callbackData[1] == "yes" {
			editText = TXT_LAST5YEAR + " *" + YES + "*\n\n"
		} else {
			editText = TXT_LAST5YEAR + " *" + NO + "*\n\n"
		}
		mes := tgapi.NewEditMessageText(h.Update.CallbackQuery.Message.Chat.ID, h.Update.CallbackQuery.Message.MessageID, editText)
		mes.ParseMode = "Markdown"
		h.Bot.Send(mes)
		//h.Get_yes_no(callbackData[1])

		if txt, err := h.Get_yes_no(callbackData[1]); err != nil {
			log.Println("error HandlerCallback: err calc:", WRONG_CALC, err)
			editText = txt + WRONG_CALC
		} else {
			editText = txt
		}
		mes1 := tgapi.NewMessage(h.Update.CallbackQuery.Message.Chat.ID, editText)
		mes1.ParseMode = "Markdown"
		h.Bot.Send(mes1)

		h.Ses.ResetSession()
	default: // clear button
		h.Bot.Send(tgapi.NewEditMessageText(h.Update.CallbackQuery.Message.Chat.ID, h.Update.CallbackQuery.Message.MessageID,
			h.Update.CallbackQuery.Message.Text))
	}

}