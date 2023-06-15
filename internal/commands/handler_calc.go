package commands

import (
	"fmt"
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
	"github.com/mrmioxin/gak_telegram_bot/internal/sessions"
)

type reguest struct {
	ok_text    string
	wrong_text string
	worker     func(c *Commander, mes *tgapi.Message) error
}

var (
	requestsListCalc  = make([]reguest, 0)
	yesNoKeyboardCalc = tgapi.NewInlineKeyboardMarkup(
		tgapi.NewInlineKeyboardRow(
			tgapi.NewInlineKeyboardButtonData(YES, "c yes"),
			tgapi.NewInlineKeyboardButtonData(NO, "c no"),
		))
)

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
