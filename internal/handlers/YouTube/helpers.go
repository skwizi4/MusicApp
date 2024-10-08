package YouTube

import (
	tg "gopkg.in/tucnak/telebot.v2"
	"log"
)

func (h Handler) SendMsg(msg *tg.Message, outText string) {
	if _, err := h.bot.Send(msg.Sender, outText); err != nil {
		log.Fatalf("Critical error: Telegram bot failed to send message: %v", err)
	}
}
func (h Handler) DeleteProcess(msg *tg.Message) {
	if err := h.ProcessingSpotifySongByYoutubeMediaId.Delete(msg.Chat.ID); err != nil {
		log.Println("Error deleting process:", err)
	}
}
