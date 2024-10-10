package tg_handlers

import (
	"MusicApp/internal/domain"
	"MusicApp/internal/errors"
	"fmt"
	tg "gopkg.in/tucnak/telebot.v2"
	"log"
)

// MultiPurposeMethods

func (h Handler) SendMsg(msg *tg.Message, outText string) {
	if _, err := h.bot.Send(msg.Sender, outText); err != nil {
		log.Fatalf("Critical error: Telegram bot failed to send message: %v", err)
	}
}

/* ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ */

// YoutubeSong helper - completed

func (h Handler) GetYoutubeSong(msg *tg.Message) error {
	process := h.processingYoutubeMediaBySpotifySongLink.GetOrCreate(msg.Chat.ID)

	switch process.Step {
	case domain.ProcessYoutubeMediaBySpotifySongLinkStart:
		h.SendMsg(msg, "Send link of song that you wanna find")
		if err := h.processingYoutubeMediaBySpotifySongLink.UpdateStep(domain.ProcessYoutubeMediaBySpotifySongLinkEnd, msg.Chat.ID); err != nil {
			h.SendMsg(msg, errors.ErrTryAgain)
			h.DeleteProcessingYoutubeMediaBySpotifySongID(msg)
			return err
		}
	case domain.ProcessYoutubeMediaBySpotifySongLinkEnd:
		if msg.Text == "/exit" {
			h.SendMsg(msg, "Process stopped")
			h.DeleteProcessingYoutubeMediaBySpotifySongID(msg)
			return nil
		}
		track, err := h.spotifyHandler.GetSpotifySongByLink(msg.Text)
		if err != nil {
			if err.Error() == errors.ErrInvalidParamsSpotify {
				h.SendMsg(msg, errors.ErrInvalidParamsSpotify)
				h.DeleteProcessingYoutubeMediaBySpotifySongID(msg)
				return err
			}
			h.SendMsg(msg, errors.ErrTryAgain)
			h.DeleteProcessingYoutubeMediaBySpotifySongID(msg)
			return err
		}
		song, err := h.youtubeHandler.GetYoutubeMediaByMetaData(&domain.MetaData{Title: track.Title, Artist: track.Artist})
		if err != nil {
			if err.Error() == errors.ErrInvalidParamsYoutube {
				h.SendMsg(msg, errors.ErrInvalidParamsYoutube)
				h.DeleteProcessingYoutubeMediaBySpotifySongID(msg)
				return err
			}

			h.SendMsg(msg, errors.ErrTryAgain)
			h.DeleteProcessingYoutubeMediaBySpotifySongID(msg)
			return err
		}
		SongPrint := fmt.Sprintf("Song Title: %s \n Song Artist: %s , \n Song link: %s", song.Artist, song.Title, song.Link)
		h.SendMsg(msg, SongPrint)
		h.DeleteProcessingYoutubeMediaBySpotifySongID(msg)

	}
	return nil
}
func (h Handler) DeleteProcessingYoutubeMediaBySpotifySongID(msg *tg.Message) {
	if err := h.processingYoutubeMediaBySpotifySongLink.Delete(msg.Chat.ID); err != nil {
		log.Println("Error deleting process:", err)
	}
}

/* ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ */

//Spotify helper - completed

func (h Handler) GetSpotifySong(msg *tg.Message) error {
	process := h.processingSpotifySongByYoutubeMediaLink.GetOrCreate(msg.Chat.ID)

	switch process.Step {
	case domain.ProcessSpotifySongByYouTubeMediaLinkStart:
		h.SendMsg(msg, "Send link of song that you wanna find")
		if err := h.processingSpotifySongByYoutubeMediaLink.UpdateStep(domain.ProcessSpotifySongByYouTubeMediaLinkEnd, msg.Chat.ID); err != nil {
			h.SendMsg(msg, errors.ErrTryAgain)
			h.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)
			return err
		}
	case domain.ProcessSpotifySongByYouTubeMediaLinkEnd:
		if msg.Text == "/exit" {
			h.SendMsg(msg, "Process stopped")
			h.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)
			return nil
		}
		track, err := h.youtubeHandler.GetYoutubeMediaByLink(msg.Text)
		if err != nil {
			if err.Error() == errors.ErrInvalidParamsSpotify {
				h.SendMsg(msg, errors.ErrInvalidParamsSpotify)
				h.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)
				return err
			}
			h.SendMsg(msg, errors.ErrTryAgain)
			h.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)
			return err
		}
		song, err := h.spotifyHandler.GetSpotifySongByMetaData(&domain.MetaData{Title: track.Title, Artist: track.Artist})
		if err != nil {
			if err.Error() == errors.ErrInvalidParamsYoutube {
				h.SendMsg(msg, errors.ErrInvalidParamsYoutube)
				h.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)
				return err
			}
			h.SendMsg(msg, errors.ErrTryAgain)
			h.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)
			return err
		}
		SongPrint := fmt.Sprintf("Song Title: %s \n Song Artist: %s , \n Song link: %s", song.Artist, song.Title, song.Link)
		h.SendMsg(msg, SongPrint)
		h.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)

	}
	return nil
}
func (h Handler) DeleteProcessingSpotifySongByYoutubeMediaLink(msg *tg.Message) {
	if err := h.processingSpotifySongByYoutubeMediaLink.Delete(msg.Chat.ID); err != nil {
		log.Println("Error deleting process:", err)
	}
}

/* ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ */

// HelpOut helper - completed
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

/* ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ */

// GetSongsByMetadata helper - completed

func (h Handler) GetSongsByMetadata(msg *tg.Message) error {
	process := h.processingFindSongByMetadata.GetOrCreate(msg.Chat.ID)

	switch process.Step {
	case domain.ProcessSongByMetadataStart:
		h.SendMsg(msg, "Send Title of song that you wanna find")

		if err := h.processingFindSongByMetadata.UpdateStep(domain.ProcessSongByMetadataTitle, msg.Chat.ID); err != nil {
			h.SendMsg(msg, errors.ErrTryAgain)
			h.DeleteProcessingFindSongByMetadata(msg)
			return err
		}

	case domain.ProcessSongByMetadataTitle:

		if err := h.processingFindSongByMetadata.AddTitle(msg.Chat.ID, msg.Text); err != nil {
			h.SendMsg(msg, errors.ErrTryAgain)
			h.DeleteProcessingFindSongByMetadata(msg)
			return err
		}
		h.SendMsg(msg, "Send Artist of song that you wanna find")

		if err := h.processingFindSongByMetadata.UpdateStep(domain.ProcessSongByMetadataArtist, msg.Chat.ID); err != nil {
			h.SendMsg(msg, errors.ErrTryAgain)
			h.DeleteProcessingFindSongByMetadata(msg)
			return err
		}
	case domain.ProcessSongByMetadataArtist:
		if err := h.processingFindSongByMetadata.AddArtist(msg.Chat.ID, msg.Text); err != nil {
			h.SendMsg(msg, errors.ErrTryAgain)
			h.DeleteProcessingFindSongByMetadata(msg)
			return err
		}
		h.SendMsg(msg, "wait a second ....")

		if err := h.processingFindSongByMetadata.ChangeIsGetMetadata(msg.Chat.ID, true); err != nil {
			h.SendMsg(msg, errors.ErrTryAgain)
			h.DeleteProcessingFindSongByMetadata(msg)
			return err
		}
		err := h.Search(msg)
		if err != nil {
			h.SendMsg(msg, errors.ErrTryAgain)
		}

	}
	return nil
}
func (h Handler) DeleteProcessingFindSongByMetadata(msg *tg.Message) {
	if err := h.processingFindSongByMetadata.Delete(msg.Chat.ID); err != nil {
		log.Println("Error deleting process:", err)
	}
}
func (h Handler) Search(msg *tg.Message) error {
	process := h.processingFindSongByMetadata.GetOrCreate(msg.Chat.ID)
	spotifySong, err := h.spotifyHandler.GetSpotifySongByMetaData(&domain.MetaData{Title: process.Song.Title, Artist: process.Song.Artist})
	if err != nil {
		if err.Error() == errors.ErrInvalidParamsSpotify {
			h.SendMsg(msg, fmt.Sprintf("Error: check input parameters. We cant find song with title: %s, and artist: %s ",
				process.Song.Title, process.Song.Artist))
			h.DeleteProcessingFindSongByMetadata(msg)
			return err
		}
		h.SendMsg(msg, errors.ErrTryAgain)
		h.DeleteProcessingFindSongByMetadata(msg)
		return err
	}
	youtubeSong, err := h.youtubeHandler.GetYoutubeMediaByMetaData(&domain.MetaData{Title: process.Song.Title, Artist: process.Song.Artist})
	if err != nil {
		if err.Error() == errors.ErrInvalidParamsYoutube {
			h.SendMsg(msg, fmt.Sprintf("Error: check input parameters. We cant find song with title: %s, and artist: %s ",
				process.Song.Title, process.Song.Artist))
			h.DeleteProcessingFindSongByMetadata(msg)
			return err

		}
		h.SendMsg(msg, errors.ErrTryAgain)
		h.DeleteProcessingFindSongByMetadata(msg)
		return err
	}

	SongPrint := fmt.Sprintf("Spotify song title: %s; \n Spotify song artist: %s; \n Spotify song link: %s; \n \n Youtube song title: %s;  \n Youtube song artist: %s; \n Youtube song link: %s.",
		spotifySong.Title, spotifySong.Artist, spotifySong.Link, youtubeSong.Title, youtubeSong.Artist, youtubeSong.Link)
	h.SendMsg(msg, SongPrint)
	h.DeleteProcessingFindSongByMetadata(msg)
	return nil
}

/* ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ */

// FillYoutubePlaylist helper - uncompleted
