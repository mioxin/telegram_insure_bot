package handlers

import (
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mioxin/telegram_insure_bot/internal/storages/sessions"
)

type Handler struct {
	Bot    *tgapi.BotAPI
	Ses    *sessions.Session
	Update tgapi.Update
	User   string
}
