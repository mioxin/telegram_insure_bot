package product1

import (
	"fmt"
	"log"
	"strings"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/handlers"
	"github.com/mrmioxin/gak_telegram_bot/internal/storages/sessions"
	"github.com/mrmioxin/gak_telegram_bot/resources"
)

type HandlerCalc struct {
	ins *Insurance
	handlers.Handler
	vid            string
	vidDescr       string
	totalWork      int
	gfotValue      float64
	eventLast5Year bool
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
			tgapi.NewInlineKeyboardButtonData(resources.YES, "c yes"),
			tgapi.NewInlineKeyboardButtonData(resources.NO, "c no"),
		))
)

func NewHandlerCalc(ins *Insurance, bot *tgapi.BotAPI, ses *sessions.Session, update tgapi.Update) *HandlerCalc {
	var user string
	//user := update.Message.Chat.UserName
	log.Printf("NewHandlerCalc: update.Message=%***REMOVED***v\nupdate.CallbackQuery=%***REMOVED***v\n", update.Message, update.CallbackQuery)
	if update.Message != nil {
		user = update.Message.Chat.UserName
	}
	if update.CallbackQuery != nil {
		user = update.CallbackQuery.Message.Chat.UserName
	}

	return &HandlerCalc{ins: ins, Handler: handlers.Handler{Bot: bot, Ses: ses, Update: update, User: user}}
}

func (h *HandlerCalc) Execute() {
	// defer log.Println("[%s]session after Execute", h.Ses)
	idx := -1
	if h.Ses != nil {
		idx = h.Ses.IdxRequest
	}
	if idx < 0 || idx >= len(requestsListCalc) {
		log.Printf("[%s] HandlerCalc: invalid request index %d", h.User, idx)
		if h.Ses != nil {
			h.Ses.ResetSession()
		}
		return
	}
	err := requestsListCalc[h.Ses.IdxRequest].worker(h, h.Update.Message)
	var m tgapi.Message
	if err != nil {
		if binErr, ok := err.(ErrorBinIinNotFound); ok {
			log.Printf("[%s] HandlerCalc: idx=%v. %v\n", h.User, h.Ses.IdxRequest, binErr)
			mes := tgapi.NewMessage(h.Update.Message.Chat.ID, resources.TXT_VID)
			m, err = h.Bot.Send(mes)
			if err != nil {
				log.Printf("[%s] HandlerCalc: send retry prompt failed: %v", h.User, err)
				return
			}
			h.Ses.IdxRequest++
			h.Ses.LastMessageID = m.MessageID
			//h.Sessions.UpdateSession(h.Update.Message.Chat.ID, ses)
			return
		}
		log.Printf("[%s] error HandlerCalc: Idx=%v. %v\n", h.User, h.Ses.IdxRequest, err)
		if h.Ses.LastRequestIsError {
			m, err = h.Bot.Send(tgapi.NewEditMessageText(h.Update.Message.Chat.ID, h.Ses.LastMessageID, resources.WRONG_AGAIN+" "+requestsListCalc[h.Ses.IdxRequest].wrong_text))
		} else {
			m, err = h.Bot.Send(tgapi.NewMessage(h.Update.Message.Chat.ID, requestsListCalc[h.Ses.IdxRequest].wrong_text))
		}
		if err != nil {
			log.Printf("[%s] HandlerCalc: send validation error failed: %v", h.User, err)
			return
		}
		h.Ses.LastRequestIsError = true
	} else {
		// log.Println("[%s]HandlerCalc: idx=", h.Ses.IdxRequest)
		mes := tgapi.NewMessage(h.Update.Message.Chat.ID, requestsListCalc[h.Ses.IdxRequest].ok_text)
		mes.ParseMode = "Markdown"
		switch h.Ses.IdxRequest {
		case 0: //skip VID it gotten from the internet by the BIN/IIN
			h.Ses.IdxRequest++
		case 1:
			editText := fmt.Sprintf("Основной вид экономической деятельности: %s - %s\n",
				h.vid, h.vidDescr)
			m, err = h.Bot.Send(tgapi.NewEditMessageText(h.Update.Message.Chat.ID, h.Ses.LastMessageID, editText))
			if err != nil {
				log.Printf("[%s] HandlerCalc: edit business type failed: %v", h.User, err)
				return
			}
		case 3: //need send the inline button Yes|No
			mes.ReplyMarkup = yesNoKeyboardCalc
		default:
		}
		m, err = h.Bot.Send(mes)
		if err != nil {
			log.Printf("[%s] HandlerCalc: send prompt failed: %v", h.User, err)
			return
		}
		if h.Ses.IdxRequest < len(requestsListCalc)-1 {
			h.Ses.IdxRequest++
		}
	}
	h.Ses.LastMessageID = m.MessageID
	h.Ses.LastRequestIsError = false
}

func (h *HandlerCalc) ExecuteCallback() {
	log.Printf("[%s] HandlerCallback: start calc update: %***REMOVED***v\n", h.User, h.Update)
	callbackData := strings.Fields(h.Update.CallbackQuery.Data)
	if len(callbackData) != 2 || callbackData[0] != "c" || (callbackData[1] != "yes" && callbackData[1] != "no") {
		log.Printf("[%s] HandlerCallback: invalid callback data %q", h.User, h.Update.CallbackQuery.Data)
		return
	}
	editText := ""
	log.Printf("[%s] HandlerCallback: start calc callback Data: %v\n", h.User, callbackData[1])
	if callbackData[1] == "yes" {
		editText = resources.TXT_LAST5YEAR + " *" + resources.YES + "*\n\n"
	} else {
		editText = resources.TXT_LAST5YEAR + " *" + resources.NO + "*\n\n"
	}
	mes := tgapi.NewEditMessageText(h.Update.CallbackQuery.Message.Chat.ID, h.Update.CallbackQuery.Message.MessageID, editText)
	mes.ParseMode = "Markdown"
	if _, err := h.Bot.Send(mes); err != nil {
		log.Printf("[%s] HandlerCallback: edit callback failed: %v", h.User, err)
	}
	//h.Get_yes_no(callbackData[1])

	if txt, err := h.Get_yes_no(callbackData[1]); err != nil {
		log.Printf("[%s] error HandlerCallback: %v, %v\n", h.User, resources.WRONG_CALC, err)
		editText = txt + resources.WRONG_CALC
	} else {
		editText = txt
	}
	mes1 := tgapi.NewMessage(h.Update.CallbackQuery.Message.Chat.ID, editText)
	mes1.ParseMode = "Markdown"
	if _, err := h.Bot.Send(mes1); err != nil {
		log.Printf("[%s] HandlerCallback: send result failed: %v", h.User, err)
		return
	}

	h.Ses.ResetSession()

}
