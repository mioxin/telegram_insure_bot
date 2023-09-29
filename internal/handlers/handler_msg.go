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
	user := ""
	if update.CallbackQuery != nil {
		user = update.CallbackQuery.Message.Chat.UserName
	} else {
		user = update.Message.Chat.UserName
	}

	return &HandlerMessage{Handler{bot, ses, update, user}, prods}
}

func (h *HandlerMessage) Execute() {
	log.Printf("[%s] HandlerMess Execute start: %#v", h.User, h.Ses)
	updText := ""
	if h.Ses == nil || h.Ses.ActionName == "" {
		//clear button if callback
		if h.Update.CallbackQuery != nil {
			log.Printf("[%s] HandlerMess Execute clear callback %v", h.User, h.Update)
			updText = "CallBack " + h.Update.CallbackQuery.Data
			h.Bot.Send(tgapi.NewEditMessageText(h.Update.CallbackQuery.Message.Chat.ID, h.Update.CallbackQuery.Message.MessageID, ""))
			//h.Update.CallbackQuery.Message.Text))
		} else {
			updText = "Message " + h.Update.Message.Text
		}
		log.Printf("error [%s] HandlerMessage: not found ActionName in session of %v\n", h.User, updText)
		return //fmt.Errorf("error HandlerMain: not found ActionName in session of Message \"%v\" (user %v)", h.update.Message.Text, h.update.Message.Chat.UserName)
	}

	h.Products[h.Ses.ActionName].Execute(h.Bot, h.Ses, h.Update)
}
