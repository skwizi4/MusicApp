package tg_handlers

import (
	"MusicApp/internal/domain"
	"MusicApp/internal/errors"
	"fmt"
	tg "gopkg.in/tucnak/telebot.v2"
	"log"
)

func (h Handler) GetSongsByMetadata(msg *tg.Message) error {
	process := h.processingFindSongByMetadata.GetOrCreate(msg.Chat.ID)

	switch process.Step {
	case domain.ProcessSpotifySongByMetadataStart:
		h.SendMsg(msg, "Send Title of song that you wanna find")

		if err := h.processingFindSongByMetadata.UpdateStep(domain.ProcessSpotifySongByMetadataTitle, msg.Chat.ID); err != nil {
			h.SendMsg(msg, errors.ErrTryAgain)
			h.DeleteProcess(msg)
			return err
		}

	case domain.ProcessSpotifySongByMetadataTitle:

		if err := h.processingFindSongByMetadata.AddTitle(msg.Chat.ID, msg.Text); err != nil {
			h.SendMsg(msg, errors.ErrTryAgain)
			h.DeleteProcess(msg)
			return err
		}
		h.SendMsg(msg, "Send Artist of song that you wanna find")

		if err := h.processingFindSongByMetadata.UpdateStep(domain.ProcessSpotifySongByMetadataArtist, msg.Chat.ID); err != nil {
			h.SendMsg(msg, errors.ErrTryAgain)
			h.DeleteProcess(msg)
			return err
		}
	case domain.ProcessSpotifySongByMetadataArtist:
		if err := h.processingFindSongByMetadata.AddArtist(msg.Chat.ID, msg.Text); err != nil {
			h.SendMsg(msg, errors.ErrTryAgain)
			h.DeleteProcess(msg)
			return err
		}
		h.SendMsg(msg, "wait a second ....")

		if err := h.processingFindSongByMetadata.ChangeIsGetMetadata(msg.Chat.ID, true); err != nil {
			h.SendMsg(msg, errors.ErrTryAgain)
			h.DeleteProcess(msg)
			return err
		}
		err := h.Search(msg)
		if err != nil {
			h.SendMsg(msg, errors.ErrTryAgain)
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
func (h Handler) SendMsg(msg *tg.Message, outText string) {
	if _, err := h.bot.Send(msg.Sender, outText); err != nil {
		log.Fatalf("Critical error: Telegram bot failed to send message: %v", err)
	}
}
func (h Handler) DeleteProcess(msg *tg.Message) {
	if err := h.processingFindSongByMetadata.Delete(msg.Chat.ID); err != nil {
		log.Println("Error deleting process:", err)
	}
}

func (h Handler) Search(msg *tg.Message) error {
	process := h.processingFindSongByMetadata.GetOrCreate(msg.Chat.ID)
	spotifySong, err := h.spotifyHandler.GetSongByMetaData(&domain.MetaData{Title: process.Song.Title, Artist: process.Song.Artist})
	if err != nil {
		if err.Error() == errors.ErrInvalidParamsSpotify {
			h.SendMsg(msg, fmt.Sprintf("Error: check input parameters. We cant find song with title: %s, and artist: %s ",
				process.Song.Title, process.Song.Artist))
			h.DeleteProcess(msg)
			return err
		}
		h.SendMsg(msg, errors.ErrTryAgain)
		h.DeleteProcess(msg)
		return err
	}
	youtubeSong, err := h.youtubeHandler.GetSongByMetaData(&domain.MetaData{Title: process.Song.Title, Artist: process.Song.Artist})
	if err != nil {
		if err.Error() == errors.ErrInvalidParamsYoutube {
			h.SendMsg(msg, fmt.Sprintf("Error: check input parameters. We cant find song with title: %s, and artist: %s ",
				process.Song.Title, process.Song.Artist))
			h.DeleteProcess(msg)
			return err

		}
		h.SendMsg(msg, errors.ErrTryAgain)
		h.DeleteProcess(msg)
		return err
	}

	SongPrint := fmt.Sprintf("Spotify song title: %s; \n Spotify song artist: %s; \n Spotify song link: %s; \n \n Youtube song title: %s;  \n Youtube song artist: %s; \n Youtube song link: %s.",
		spotifySong.Artist, spotifySong.Title, spotifySong.Link, youtubeSong.Title, youtubeSong.Artist, youtubeSong.Link)
	h.SendMsg(msg, SongPrint)
	h.DeleteProcess(msg)
	return nil
}
