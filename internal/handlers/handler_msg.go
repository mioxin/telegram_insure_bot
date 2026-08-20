package handlers

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mioxin/telegram_insure_bot/internal/services"
	"github.com/mioxin/telegram_insure_bot/internal/storages/sessions"
	"github.com/mioxin/telegram_insure_bot/resources"
)

type HandlerMessage struct {
	Handler
	Products map[string]services.IService
}

func NewHandlerMessage(bot *tgapi.BotAPI, ses *sessions.Session, prods map[string]services.IService, update tgapi.Update) *HandlerMessage {
	user := ""
	if update.CallbackQuery != nil && update.CallbackQuery.Message != nil && update.CallbackQuery.Message.Chat != nil {
		user = update.CallbackQuery.Message.Chat.UserName
	} else if update.Message != nil && update.Message.Chat != nil {
		user = update.Message.Chat.UserName
	}

	return &HandlerMessage{Handler{bot, ses, update, user}, prods}
}

func (h *HandlerMessage) Execute() {
	log.Printf("[%s] HandlerMess Execute start: %#v", h.User, h.Ses)
	if h.Update.CallbackQuery != nil && (h.Update.CallbackQuery.Message == nil || h.Update.CallbackQuery.Message.Chat == nil) {
		log.Printf("[%s] HandlerMessage: incomplete callback", h.User)
		return
	}
	updText := ""
	if h.Ses == nil || h.Ses.ActionName == "" {
		//clear button if callback
		if h.Update.CallbackQuery != nil {
			log.Printf("[%s] HandlerMess Execute clear callback %v", h.User, h.Update)
			updText = "CallBack " + h.Update.CallbackQuery.Data
			if _, err := h.Bot.Send(tgapi.NewEditMessageText(h.Update.CallbackQuery.Message.Chat.ID, h.Update.CallbackQuery.Message.MessageID, "")); err != nil {
				log.Printf("[%s] HandlerMessage: clear callback failed: %v", h.User, err)
			}
			//h.Update.CallbackQuery.Message.Text))
		} else {
			updText = "Message " + h.Update.Message.Text
		}
		log.Printf("error [%s] HandlerMessage: not found ActionName in session of %v\n", h.User, updText)
		return //fmt.Errorf("error HandlerMain: not found ActionName in session of Message \"%v\" (user %v)", h.update.Message.Text, h.update.Message.Chat.UserName)
	}

	if h.Update.CallbackQuery != nil && h.Ses.ActionName == "get" {
		if _, ok := h.Ses.AccessCommand["adm"]; !ok {
			log.Printf("[%s] HandlerMessage: denied get callback", h.User)
			if _, err := h.Bot.Send(tgapi.NewMessage(h.Update.CallbackQuery.Message.Chat.ID, resources.WRONG_ACCESS)); err != nil {
				log.Printf("[%s] HandlerMessage: callback denial message failed: %v", h.User, err)
			}
			return
		}
	}

	product, ok := h.Products[h.Ses.ActionName]
	if !ok || product == nil {
		log.Printf("[%s] HandlerMessage: unknown product %q", h.User, h.Ses.ActionName)
		h.Ses.ActionName = ""
		return
	}
	product.Execute(h.Bot, h.Ses, h.Update)
}
