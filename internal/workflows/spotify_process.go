package workflows

import (
	"MusicApp/internal/domain"
	"MusicApp/internal/errors"
	"fmt"
	tg "gopkg.in/tucnak/telebot.v2"
	"net/url"
	"strconv"
)

func (w WorkFlows) GetSpotifySong(msg *tg.Message) error {
	process := w.ProcessingSpotifySongByYoutubeMediaLink.GetOrCreate(msg.Chat.ID)

	switch process.Step {
	case domain.ProcessSpotifySongByYouTubeMediaLinkStart:
		w.SendMsg(msg, "Send link of youtube media that you wanna find")
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
			w.DeleteProcessingCreateAndFillSpotifyPlaylist(msg)
			return err
		}
	case domain.ProcessCreateAndFillSpotifyPlaylistEnd:
		if msg.Text == "/exit" {
			w.SendMsg(msg, "Process stopped")
			w.DeleteProcessingCreateAndFillSpotifyPlaylist(msg)
			return nil
		}
		YouTubePlaylist, err := w.YouTubeHandler.GetYoutubePlaylistByLink(msg.Text)
		if err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingCreateAndFillSpotifyPlaylist(msg)
			return err
		}
		if err = w.ProcessingCreateAndFillSpotifyPlaylists.AddTitle(YouTubePlaylist.Title, msg.Chat.ID); err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingCreateAndFillSpotifyPlaylist(msg)
			return err
		}
		if err = w.ProcessingCreateAndFillSpotifyPlaylists.AddSongs(YouTubePlaylist.Songs, msg.Chat.ID); err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingCreateAndFillSpotifyPlaylist(msg)
			return err
		}
		TelegramId := strconv.FormatInt(msg.Sender.ID, 10)

		params := url.Values{}
		params.Add("response_type", "code")
		params.Add("client_id", w.cfg.SpotifyCfg.ClientID)
		params.Add("scope", "playlist-modify-public")
		params.Add("redirect_uri", "http://localhost:8080/auth/spotify/callback")
		params.Add("state", TelegramId)

		w.SendMsg(msg, fmt.Sprintf("%s?%s", SpotifyOauthUrl, params.Encode()))
		UserProcess := fmt.Sprintf("SpotifyProcess%s", TelegramId)
		token, err := w.CheckForToken(msg, UserProcess)
		if err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingCreateAndFillSpotifyPlaylist(msg)
			return err
		}
		process = w.ProcessingCreateAndFillSpotifyPlaylists.GetOrCreate(msg.Chat.ID)

		playlistId, err := w.SpotifyHandler.CreateSpotifyPlaylist(process.Playlist.Title, token)
		if err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingCreateAndFillSpotifyPlaylist(msg)
			return err
		}
		SpotifyPlaylist, err := w.SpotifyHandler.FillSpotifyPlaylist(YouTubePlaylist, playlistId, token)
		if err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingCreateAndFillSpotifyPlaylist(msg)
			return err
		}
		w.SendMsg(msg, fmt.Sprintf("Playlist title: %s, \n PLaylist Owner:%s, \n PLaylist description: %s, \n Playlist link: %s",
			SpotifyPlaylist.Title, SpotifyPlaylist.Owner, SpotifyPlaylist.Description, SpotifyPlaylist.ExternalUrl))
		w.DeleteProcessingCreateAndFillSpotifyPlaylist(msg)

	}
	return nil
}

func (w WorkFlows) DeleteProcessingCreateAndFillSpotifyPlaylist(msg *tg.Message) {
	if err := w.ProcessingCreateAndFillSpotifyPlaylists.Delete(msg.Chat.ID); err != nil {
		w.logger.ErrorFrmt("Error deleting process:", err)
	}
}
