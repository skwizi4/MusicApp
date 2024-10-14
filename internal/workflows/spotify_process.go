package workflows

import (
	"MusicApp/internal/domain"
	"MusicApp/internal/errors"
	"fmt"
	tg "gopkg.in/tucnak/telebot.v2"
)

func (w WorkFlows) GetSpotifySong(msg *tg.Message) error {
	process := w.ProcessingSpotifySongByYoutubeMediaLink.GetOrCreate(msg.Chat.ID)

	switch process.Step {
	case domain.ProcessSpotifySongByYouTubeMediaLinkStart:
		w.SendMsg(msg, "Send link of song that you wanna find")
		if err := w.ProcessingSpotifySongByYoutubeMediaLink.UpdateStep(domain.ProcessSpotifySongByYouTubeMediaLinkEnd, msg.Chat.ID); err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)
			return err
		}
	case domain.ProcessSpotifySongByYouTubeMediaLinkEnd:
		if msg.Text == "/exit" {
			w.SendMsg(msg, "Process stopped")
			w.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)
			return nil
		}
		track, err := w.YouTubeHandler.GetYoutubeMediaByLink(msg.Text)
		if err != nil {
			if err.Error() == errors.ErrInvalidParamsSpotify {
				w.SendMsg(msg, errors.ErrInvalidParamsSpotify)
				w.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)
				return err
			}
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)
			return err
		}
		song, err := w.SpotifyHandler.GetSpotifySongByMetaData(&domain.MetaData{Title: track.Title, Artist: track.Artist})
		if err != nil {
			if err.Error() == errors.ErrInvalidParamsYoutube {
				w.SendMsg(msg, errors.ErrInvalidParamsYoutube)
				w.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)
				return err
			}
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)
			return err
		}
		SongPrint := fmt.Sprintf("Song Title: %s \n Song Artist: %s , \n Song link: %s", song.Artist, song.Title, song.Link)
		w.SendMsg(msg, SongPrint)
		w.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)

	}
	return nil
}
func (w WorkFlows) DeleteProcessingSpotifySongByYoutubeMediaLink(msg *tg.Message) {
	if err := w.ProcessingSpotifySongByYoutubeMediaLink.Delete(msg.Chat.ID); err != nil {
		w.logger.ErrorFrmt("Error deleting process:", err)
	}
}

/* ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ */
/* ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ */
/* ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ */

// CreateAndFillSpotifyPlaylist - todo refactor
func (w WorkFlows) CreateAndFillSpotifyPlaylist(msg *tg.Message) error {
	process := w.ProcessingCreateAndFillSpotifyPlaylists.GetOrCreate(msg.Chat.ID)
	switch process.Step {
	case domain.ProcessCreateAndFillSpotifyPlaylistStart:
		w.SendMsg(msg, "Send link of youtube playlist that you wanna transfer")
		if err := w.ProcessingCreateAndFillSpotifyPlaylists.UpdateStep(domain.ProcessCreateAndFillSpotifyPlaylistEnd, msg.Chat.ID); err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)
			return err
		}
	case domain.ProcessCreateAndFillSpotifyPlaylistEnd:
		playlist, err := w.YouTubeHandler.GetYoutubePlaylistByLink(msg.Text)
		if err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)
			return err
		}
		if err = w.ProcessingCreateAndFillSpotifyPlaylists.AddTitle(playlist.Title, msg.Chat.ID); err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)
			return err
		}
		if err = w.ProcessingCreateAndFillSpotifyPlaylists.AddSongs(playlist.Songs, msg.Chat.ID); err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingSpotifySongByYoutubeMediaLink(msg)
			return err
		}
		//Send Auth2.0 link
		//Create playlist
		//Fill playlist
		//Return playlist
	}
	return nil
}
