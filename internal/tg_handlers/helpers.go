package tg_handlers

import (
	"MusicApp/internal/domain"
	"fmt"
	tg "gopkg.in/tucnak/telebot.v2"
	"log"
)

func (h Handler) GetMetadata(msg *tg.Message) error {
	process := h.processingFindSongByMetadata.GetOrCreate(msg.Chat.ID)

	switch process.Step {
	case domain.ProcessSpotifySongByMetadataStart:
		if _, err := h.bot.Send(msg.Sender, "Send Title of song that you wanna find"); err != nil {
			if err = h.processingFindSongByMetadata.Delete(msg.Chat.ID); err != nil {
				h.errChannel.HandleError(err)
				return err
			}
			h.errChannel.HandleError(err)

			return err
		}
		if err := h.processingFindSongByMetadata.UpdateStep(domain.ProcessSpotifySongByMetadataTitle, msg.Chat.ID); err != nil {
			h.errChannel.HandleError(err)
			return err
		}

	case domain.ProcessSpotifySongByMetadataTitle:

		if err := h.processingFindSongByMetadata.AddTitle(msg.Chat.ID, msg.Text); err != nil {
			h.errChannel.HandleError(err)
		}
		if _, err := h.bot.Send(msg.Sender, "Send Artist of song that you wanna find"); err != nil {
			if err = h.processingFindSongByMetadata.Delete(msg.Chat.ID); err != nil {
				h.errChannel.HandleError(err)
			}
			h.errChannel.HandleError(err)
		}
		if err := h.processingFindSongByMetadata.UpdateStep(domain.ProcessSpotifySongByMetadataEnd, msg.Chat.ID); err != nil {
			h.errChannel.HandleError(err)
			return err
		}

	}
	return nil
}

func (h Handler) HelpOut(msg *tg.Message) {
	formatString := fmt.Sprintf("			Commands: \n" +
		"	/FindSong - search songs in youtube and spotify by track metadata ( track title and artist) \n" +
		"	/SpotifySong - search song in spotify by youtube link of this track \n" +
		"	/YoutubeSong - search track in youtube by spotify link of this track \n " +
		"	/SpotifyPlaylist - search and fill spotify playlist by link of the playlist from youtube \n" +
		"	/YoutubePlaylist - search and fill youtube playlist by link of the playlist from spotify \n")
	if _, err := h.bot.Send(msg.Sender, formatString); err != nil {
		log.Fatal(err)
	}

}
