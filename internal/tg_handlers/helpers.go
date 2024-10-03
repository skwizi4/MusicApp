package tg_handlers

import (
	"MusicApp/internal/domain"
	tg "gopkg.in/tucnak/telebot.v2"
)

func (h Handler) GetMetadata(msg *tg.Message) error {
	process := h.processingFinSongByMetadata.GetOrCreate(msg.Chat.ID)

	switch process.Step {
	case domain.ProcessSpotifySongByMetadataStart:
		if _, err := h.bot.Send(msg.Sender, "Send Title of song that you wanna find"); err != nil {
			if err = h.processingFinSongByMetadata.Delete(msg.Chat.ID); err != nil {
				h.errChannel.HandleError(err)
				return err
			}
			h.errChannel.HandleError(err)

			return err
		}
		if err := h.processingFinSongByMetadata.UpdateStep(domain.ProcessSpotifySongByMetadataTitle, msg.Chat.ID); err != nil {
			h.errChannel.HandleError(err)
			return err
		}
	case domain.ProcessSpotifySongByMetadataTitle:

		if err := h.processingFinSongByMetadata.AddTitle(msg.Chat.ID, msg.Text); err != nil {
			h.errChannel.HandleError(err)
		}
		if _, err := h.bot.Send(msg.Sender, "Send Artist of song that you wanna find"); err != nil {
			if err = h.processingFinSongByMetadata.Delete(msg.Chat.ID); err != nil {
				h.errChannel.HandleError(err)
			}
			h.errChannel.HandleError(err)
		}
		if err := h.processingFinSongByMetadata.UpdateStep(domain.ProcessSpotifySongByMetadataEnd, msg.Chat.ID); err != nil {
			h.errChannel.HandleError(err)
			return err
		}

	}
	return nil
}
