package commands

import (
	"log"
	"strings"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

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
