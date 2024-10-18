package workflows

import (
	"MusicApp/internal/domain"
	"MusicApp/internal/errors"
	"fmt"
	tg "gopkg.in/tucnak/telebot.v2"
	"strconv"
)

// GetYoutubeSong - completed
func (w WorkFlows) GetYoutubeSong(msg *tg.Message) error {
	process := w.ProcessingYoutubeMediaBySpotifySongLink.GetOrCreate(msg.Chat.ID)

	switch process.Step {
	case domain.ProcessYoutubeMediaBySpotifySongLinkStart:
		w.SendMsg(msg, "Send link of spotify song that you wanna find")
		if err := w.ProcessingYoutubeMediaBySpotifySongLink.UpdateStep(domain.ProcessYoutubeMediaBySpotifySongLinkEnd, msg.Chat.ID); err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingYoutubeMediaBySpotifySongID(msg)
			return err
		}
	case domain.ProcessYoutubeMediaBySpotifySongLinkEnd:
		if msg.Text == "/exit" {
			w.SendMsg(msg, "Process stopped")
			w.DeleteProcessingYoutubeMediaBySpotifySongID(msg)
			return nil
		}
		track, err := w.SpotifyHandler.GetSpotifySongByLink(msg.Text)
		if err != nil {
			if err.Error() == errors.ErrInvalidParamsSpotify {
				w.SendMsg(msg, errors.ErrInvalidParamsSpotify)
				w.DeleteProcessingYoutubeMediaBySpotifySongID(msg)
				return err
			}
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingYoutubeMediaBySpotifySongID(msg)
			return err
		}
		song, err := w.YouTubeHandler.GetYoutubeMediaByMetaData(&domain.MetaData{Title: track.Title, Artist: track.Artist})
		if err != nil {
			if err.Error() == errors.ErrInvalidParamsYoutube {
				w.SendMsg(msg, errors.ErrInvalidParamsYoutube)
				w.DeleteProcessingYoutubeMediaBySpotifySongID(msg)
				return err
			}

			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingYoutubeMediaBySpotifySongID(msg)
			return err
		}
		SongPrint := fmt.Sprintf("Song Title: %s \n Song Artist: %s , \n Song link: %s", song.Artist, song.Title, song.Link)
		w.SendMsg(msg, SongPrint)
		w.DeleteProcessingYoutubeMediaBySpotifySongID(msg)

	}
	return nil
}
func (w WorkFlows) DeleteProcessingYoutubeMediaBySpotifySongID(msg *tg.Message) {
	if err := w.ProcessingYoutubeMediaBySpotifySongLink.Delete(msg.Chat.ID); err != nil {
		w.logger.ErrorFrmt("Error deleting process:", err)
	}
}

/* ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ */
/* ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ */
/* ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ */

// CreateAndFillYoutubePlaylistHelper - Completed

func (w WorkFlows) CreateAndFillYoutubePlaylist(msg *tg.Message) error {
	process := w.ProcessingCreateAndFillYoutubePlaylists.GetOrCreate(msg.Chat.ID)
	switch process.Step {
	case domain.ProcessFillYouTubePlaylistStart:
		w.SendMsg(msg, "send a link to the playlist from spotify that you want to transfer to youtube")
		if err := w.ProcessingCreateAndFillYoutubePlaylists.UpdateStep(domain.ProcessFillYouTubePlaylistSendAuthLink, msg.Chat.ID); err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingCreateAndFillYoutubePlaylists(msg)
			return err
		}
	case domain.ProcessFillYouTubePlaylistSendAuthLink:
		if msg.Text == "/exit" {
			w.SendMsg(msg, "Process stopped")
			w.DeleteProcessingYoutubeMediaBySpotifySongID(msg)
			return nil
		}
		SpotifyPlaylist, err := w.SpotifyHandler.GetSpotifyPlaylistByLink(msg.Text)
		if err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingCreateAndFillYoutubePlaylists(msg)
			return err
		}
		if err = w.ProcessingCreateAndFillYoutubePlaylists.AddTitle(SpotifyPlaylist.Title, msg.Chat.ID); err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingCreateAndFillYoutubePlaylists(msg)
			return err
		}
		if err = w.ProcessingCreateAndFillYoutubePlaylists.AddSongs(SpotifyPlaylist, msg.Chat.ID); err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingCreateAndFillYoutubePlaylists(msg)
			return err
		}
		TelegramId := strconv.FormatInt(msg.Sender.ID, 10)
		w.SendMsg(msg, fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=https://www.googleapis.com/auth/youtube&access_type=offline&state=%s&prompt=consent",
			w.cfg.YoutubeCfg.ClientID, w.cfg.YoutubeCfg.RedirectUrl, TelegramId))
		UserProcess := fmt.Sprintf(YoutubeProcess, TelegramId)
		token, err := w.CheckForToken(msg, UserProcess)
		if err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingCreateAndFillYoutubePlaylists(msg)
			return err
		}
		process = w.ProcessingCreateAndFillYoutubePlaylists.GetOrCreate(msg.Chat.ID)
		playlistId, err := w.YouTubeHandler.CreateYoutubePlaylist(process.Playlist.Title, token)
		if err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingCreateAndFillYoutubePlaylists(msg)
			return err
		}
		youtubePlaylist, err := w.YouTubeHandler.FillYouTubePlaylist(SpotifyPlaylist, playlistId, token)
		if err != nil {
			w.SendMsg(msg, errors.ErrTryAgain)
			w.DeleteProcessingCreateAndFillYoutubePlaylists(msg)
			return err
		}
		w.SendMsg(msg, fmt.Sprintf("Playlist title: %s \n Playlist link: %s", youtubePlaylist.Title, youtubePlaylist.ExternalUrl))
		w.DeleteProcessingCreateAndFillYoutubePlaylists(msg)

	}
	return nil
}

func (w WorkFlows) DeleteProcessingCreateAndFillYoutubePlaylists(msg *tg.Message) {
	if err := w.ProcessingCreateAndFillYoutubePlaylists.Delete(msg.Chat.ID); err != nil {
		w.logger.ErrorFrmt("Error deleting process:", err)
	}
}
