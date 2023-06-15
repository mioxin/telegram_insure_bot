package commands

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/internal/sessions"
)

type RegisteredCommand struct {
	Description string
	Worker      func(c *Commander, mes *tgapi.Message) (string, int)
	ShowInHelp  bool
}

var registered_commands = map[string]RegisteredCommand{}

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
