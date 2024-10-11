package workflows

import (
	"MusicApp/internal/domain"
	"MusicApp/internal/errors"
	"fmt"
	tg "gopkg.in/tucnak/telebot.v2"
	"log"
)

func (w WorkFlows) GetSongsByMetadata(msg *tg.Message) error {
	process := w.ProcessingFindSongByMetadata.GetOrCreate(msg.Chat.ID)

	switch process.Step {
	case domain.ProcessSongByMetadataStart:
		w.SendMsg(msg, "Send Title of song that you wanna find")

		if err := w.ProcessingFindSongByMetadata.UpdateStep(domain.ProcessSongByMetadataTitle, msg.Chat.ID); err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingFindSongByMetadata(msg)
			return err
		}

	case domain.ProcessSongByMetadataTitle:

		if err := w.ProcessingFindSongByMetadata.AddTitle(msg.Chat.ID, msg.Text); err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingFindSongByMetadata(msg)
			return err
		}
		w.SendMsg(msg, "Send Artist of song that you wanna find")

		if err := w.ProcessingFindSongByMetadata.UpdateStep(domain.ProcessSongByMetadataArtist, msg.Chat.ID); err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingFindSongByMetadata(msg)
			return err
		}
	case domain.ProcessSongByMetadataArtist:
		if err := w.ProcessingFindSongByMetadata.AddArtist(msg.Chat.ID, msg.Text); err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingFindSongByMetadata(msg)
			return err
		}
		w.SendMsg(msg, "wait a second ....")

		if err := w.ProcessingFindSongByMetadata.ChangeIsGetMetadata(msg.Chat.ID, true); err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingFindSongByMetadata(msg)
			return err
		}
		err := w.Search(msg)
		if err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
		}

	}
	return nil
}
func (w WorkFlows) DeleteProcessingFindSongByMetadata(msg *tg.Message) {
	if err := w.ProcessingFindSongByMetadata.Delete(msg.Chat.ID); err != nil {
		log.Println("Error deleting process:", err)
	}
}
func (w WorkFlows) Search(msg *tg.Message) error {
	process := w.ProcessingFindSongByMetadata.GetOrCreate(msg.Chat.ID)
	spotifySong, err := w.SpotifyHandler.GetSpotifySongByMetaData(&domain.MetaData{Title: process.Song.Title, Artist: process.Song.Artist})
	if err != nil {
		if err.Error() == errors.ErrInvalidParamsSpotify {
			w.SendMsg(msg, fmt.Sprintf("Error: check input parameters. We cant find song with title: %s, and artist: %s ",
				process.Song.Title, process.Song.Artist))
			w.DeleteProcessingFindSongByMetadata(msg)
			return err
		}
		w.SendMsg(msg, errors.ErrTryAgain)
		w.DeleteProcessingFindSongByMetadata(msg)
		return err
	}
	youtubeSong, err := w.YouTubeHandler.GetYoutubeMediaByMetaData(&domain.MetaData{Title: process.Song.Title, Artist: process.Song.Artist})
	if err != nil {
		if err.Error() == errors.ErrInvalidParamsYoutube {
			w.SendMsg(msg, fmt.Sprintf("Error: check input parameters. We cant find song with title: %s, and artist: %s ",
				process.Song.Title, process.Song.Artist))
			w.DeleteProcessingFindSongByMetadata(msg)
			return err

		}
		w.SendMsg(msg, errors.ErrTryAgain)
		w.DeleteProcessingFindSongByMetadata(msg)
		return err
	}

	SongPrint := fmt.Sprintf("Spotify song title: %s; \n Spotify song artist: %s; \n Spotify song link: %s; \n \n Youtube song title: %s;  \n Youtube song artist: %s; \n Youtube song link: %s.",
		spotifySong.Title, spotifySong.Artist, spotifySong.Link, youtubeSong.Title, youtubeSong.Artist, youtubeSong.Link)
	w.SendMsg(msg, SongPrint)
	w.DeleteProcessingFindSongByMetadata(msg)
	return nil
}
