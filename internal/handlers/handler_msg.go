package handlers

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/services"
	"github.com/mrmioxin/gak_telegram_bot/internal/storages/sessions"
)

type HandlerMessage struct {
	Handler
	Products map[string]services.IService
}

func NewHandlerMessage(bot *tgapi.BotAPI, ses *sessions.Session, prods map[string]services.IService, update tgapi.Update) *HandlerMessage {
	return &HandlerMessage{Handler{bot, ses, update}, prods}
}

func (h *HandlerMessage) Execute() {
	log.Printf("HandlerMess Execute start: %#v", h.Ses)
	updText, userName := "", ""
	if h.Ses == nil || h.Ses.ActionName == "" {
		//clear button if callback
		if h.Update.CallbackQuery != nil {
			log.Printf("HandlerMess Execute clear callback %v", h.Update)
			updText = "CallBack " + h.Update.CallbackQuery.Data
			userName = h.Update.CallbackQuery.Message.Chat.UserName
			h.Bot.Send(tgapi.NewEditMessageText(h.Update.CallbackQuery.Message.Chat.ID, h.Update.CallbackQuery.Message.MessageID, ""))
			//h.Update.CallbackQuery.Message.Text))
		} else {
			updText = "Message " + h.Update.Message.Text
			userName = h.Update.Message.Chat.UserName
		}
		log.Printf("error HandlerMessage: not found ActionName in session of %v (user %v)",
			updText, userName)
		return //fmt.Errorf("error HandlerMain: not found ActionName in session of Message \"%v\" (user %v)", h.update.Message.Text, h.update.Message.Chat.UserName)
	}

	h.Products[h.Ses.ActionName].Execute(h.Bot, h.Ses, h.Update)
}
