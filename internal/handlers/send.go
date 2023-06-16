package handlers

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const SEND string = `Вложите сканы документов содержащих следующие данные:
- общее количество работников с учетом работников филиалов (одно число);
- основной вид экономической деятельности;
- количество работников с ежемесячым окладом более 10 минимальных заработных плат (далее МЗП);
- ГФОТ по работникам с ежемесячым окладом более 10 МЗП;
- количество работников с ежемесячым окладом менее или равного 10 МЗП;
- ГФОТ по работникам с ежемесячым окладом менее или равного 10 МЗП.

На основании этих документов менеджер расчитает вам страховку с учетом скидки или кешбэка.`

func (h *HandlerCommands) send(input_message *tgapi.Message) (string, int) {
	log.Printf("send: [%s] %s", input_message.From.UserName, input_message.Text)

	msg := tgapi.NewMessage(input_message.Chat.ID, SEND)
	m, _ := h.bot.Send(msg)
	return "send", m.MessageID
}

func init() {
	registered_commands["send"] = RegisteredCommand{Description: "Отправка информации для расчета страховки менеджером.", Worker: (*HandlerCommands).send, ShowInHelp: true}
}
