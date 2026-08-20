package handlers

import (
	"fmt"
	"io"
	"log"
	"os"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mioxin/telegram_insure_bot/resources"
)

func (h *HandlerCommands) about(input_message *tgapi.Message) (string, int) {
	var msg tgapi.Chattable
	var err error
	var closer io.Closer
	log.Printf("about: [%s] %s", input_message.From.UserName, input_message.Text)
	if msg, closer, err = h.videoMessage(input_message); err != nil {
		log.Printf("error in About command: %v.\n The text message \"about\" will be send\n", err)
		msg = tgapi.NewMessage(input_message.Chat.ID, fmt.Sprintf(resources.ABOUT, resources.COMPANY))
	}
	if closer != nil {
		defer closer.Close()
	}

	m, err := h.Bot.Send(msg)
	if err != nil {
		log.Printf("error in About command: %v.\n %v\n", m, err)

	}

	if err == nil && m.Video != nil && !msg.(tgapi.VideoConfig).UseExisting {
		log.Printf("About command: The video file id %v saved to files_id.json.\nmsg.(tgapi.VideoConfig).UseExisting: %v", m.Video.FileID, msg.(tgapi.VideoConfig).UseExisting)
		h.FilesId.SetFileId(resources.VIDEO_ABOUT, "bot", m.Video.FileID)
	}
	return "", m.MessageID
}

func (h *HandlerCommands) videoMessage(input_message *tgapi.Message) (tgapi.Chattable, io.Closer, error) {
	var vmsg tgapi.VideoConfig
	if file_id, err := h.FilesId.GetFileId(resources.VIDEO_ABOUT, "bot"); err != nil {
		log.Printf("About command: %v. The video \"about\" will be send by file\n", err)
		if file, err := os.OpenFile(resources.VIDEO_ABOUT, os.O_RDONLY, 0755); err != nil {
			return nil, nil, fmt.Errorf("error About videoMessage: %v", err)
		} else {
			vmsg = tgapi.NewVideoUpload(input_message.Chat.ID,
				tgapi.FileReader{Name: "video about company", Reader: file, Size: -1})
			vmsg.Caption = fmt.Sprintf(resources.ABOUT, resources.COMPANY)
			return vmsg, file, nil
		}
	} else {
		vmsg = tgapi.NewVideoShare(input_message.Chat.ID, file_id)
		log.Printf("About command: %v. The video \"about\" will be send by id\n", err)
	}
	vmsg.Caption = fmt.Sprintf(resources.ABOUT, resources.COMPANY)

	return vmsg, nil, nil
}

func init() {
	registered_commands["about"] = RegisteredCommand{Description: "Коротко о боте.", Worker: (*HandlerCommands).about, ShowInHelp: true}
}
