package handlers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrmioxin/gak_telegram_bot/resources"
)

func (h *HandlerCommands) about(input_message *tgapi.Message) (string, int) {
	var msg tgapi.Chattable
	var err error
	log.Printf("about: [%s] %s", input_message.From.UserName, input_message.Text)
	if msg, err = h.videoMessage(input_message); err != nil {
		log.Printf("error in About command: %v.\n The text message \"about\" will be send\n", err)
		msg = tgapi.NewMessage(input_message.Chat.ID, resources.ABOUT)
	}

	m, err := h.Bot.Send(msg)
	if err != nil {
		log.Printf("error in About command: %v.\n %v\n", m, err)

	}

	if m.Video != nil && !msg.(tgapi.VideoConfig).UseExisting {
		log.Printf("About command: The video file id %v saved to files_id.json.\nmsg.(tgapi.VideoConfig).UseExisting: %v", m.Video.FileID, msg.(tgapi.VideoConfig).UseExisting)
		h.FilesId.SetFileId(resources.VIDEO_ABOUT, "", m.Video.FileID)
	}
	return "", m.MessageID
}

func (h *HandlerCommands) videoMessage(input_message *tgapi.Message) (tgapi.Chattable, error) {
	var vmsg tgapi.VideoConfig
	if file_id, err := h.FilesId.GetFileId(filepath.Join(resources.VIDEO_ABOUT)); err != nil {
		log.Printf("About command: %v. The video \"about\" will be send by file\n", err)
		if file, err := os.OpenFile(resources.VIDEO_ABOUT, os.O_RDONLY, 0755); err != nil {
			return nil, fmt.Errorf("error About videoMessage: %v", err)
		} else {
			vmsg = tgapi.NewVideoUpload(input_message.Chat.ID,
				tgapi.FileReader{Name: "video about company", Reader: file, Size: -1})
		}
	} else {
		vmsg = tgapi.NewVideoShare(input_message.Chat.ID, file_id)
		log.Printf("About command: %v. The video \"about\" will be send by id\n", err)
	}
	vmsg.Caption = resources.ABOUT

	return vmsg, nil
}

func init() {
	registered_commands["about"] = RegisteredCommand{Description: "Коротко о боте.", Worker: (*HandlerCommands).about, ShowInHelp: true}
}
