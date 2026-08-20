package services

import (
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mioxin/telegram_insure_bot/internal/storages/sessions"
)

type IService interface {
	//Calculate() (string, error)
	Execute(bot *tgapi.BotAPI, ses *sessions.Session, update tgapi.Update)
}
