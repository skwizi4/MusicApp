package Spotify

import (
	"MusicApp/internal/config"
	"MusicApp/internal/domain"
	"MusicApp/internal/errors"
	"MusicApp/internal/services"
	"MusicApp/internal/services/Spotify"
	YouTubeService "MusicApp/internal/services/YouTube"
	"fmt"
	"github.com/skwizi4/lib/ErrChan"
	"github.com/skwizi4/lib/logs"
	tg "gopkg.in/tucnak/telebot.v2"
)

type Handler struct {
	bot                                   *tg.Bot
	ProcessingSpotifySongByYoutubeMediaId *domain.ProcessingSpotifySongByYoutubeMediaId
	errChannel                            *ErrChan.ErrorChannel
	spotifyService                        services.SpotifyService
	youtubeService                        services.YouTubeService
	cfg                                   *config.Config
	logger                                logs.GoLogger
}

func New(bot *tg.Bot, cfg *config.Config, processingYoutubeSongsById *domain.ProcessingSpotifySongByYoutubeMediaId, errChan *ErrChan.ErrorChannel, logger logs.GoLogger) Handler {
	return Handler{
		bot:                                   bot,
		spotifyService:                        Spotify.NewSpotifyService(cfg),
		youtubeService:                        YouTubeService.NewYouTubeService(cfg),
		ProcessingSpotifySongByYoutubeMediaId: processingYoutubeSongsById,
		errChannel:                            errChan,
		cfg:                                   cfg,
		logger:                                logger,
	}
}

func (h Handler) GetSpotifySongByYoutubeLink(msg *tg.Message) error {
	process := h.ProcessingSpotifySongByYoutubeMediaId.GetOrCreate(msg.Chat.ID)
	switch process.Step {
	case domain.ProcessSpotifySongByYouTubeMediaStart:

		h.SendMsg(msg, "Send youtube link of media that you wanna find in Youtube")

		if err := h.ProcessingSpotifySongByYoutubeMediaId.UpdateStep(domain.ProcessSpotifySongByYouTubeMediaEnd, msg.Chat.ID); err != nil {
			h.SendMsg(msg, "Error, try again")
			h.DeleteProcess(msg)
			return err
		}
	case domain.ProcessSpotifySongByYouTubeMediaEnd:
		if msg.Text == "/exit" {
			h.SendMsg(msg, "Process stopped")
			h.DeleteProcess(msg)
			return nil
		}
		track, err := h.youtubeService.GetYoutubeMediaByID(msg.Text)
		if err != nil {
			if err.Error() == errors.ErrInvalidParamsYoutube {
				h.SendMsg(msg, errors.ErrInvalidParamsYoutube)
				h.DeleteProcess(msg)
				return err
			} else {
				h.SendMsg(msg, "Error, try again")
				h.DeleteProcess(msg)
				return err
			}

		}
		song, err := h.spotifyService.GetSpotifyTrackByMetadata(domain.MetaData{Title: track.Title, Artist: track.Artist})
		if err != nil {
			if err.Error() == errors.ErrInvalidParamsSpotify {
				h.SendMsg(msg, errors.ErrInvalidParamsSpotify)
				h.DeleteProcess(msg)
				return err
			} else {
				h.SendMsg(msg, "Error, try again")
				h.DeleteProcess(msg)
				return err
			}

		}
		SongPrint := fmt.Sprintf("Youtube Song Title: %s \n Youtube Song Artist: %s , \n Youtube Song link: %s", song.Title, song.Artist, song.Link)
		h.SendMsg(msg, SongPrint)
		h.DeleteProcess(msg)
	}
	return nil
}

func (h Handler) GetSpotifySongByMetaData(metadata *domain.MetaData) (*domain.Song, error) {
	track, err := h.spotifyService.GetSpotifyTrackByMetadata(*metadata)
	if err != nil {
		h.errChannel.HandleError(err)
		return nil, err
	}

	return track, nil
}

func (h Handler) GetSpotifyPlaylistByYoutubeLink(spotifyLink string) (*domain.Playlist, error) {
	return nil, nil
}

func (h Handler) FillSpotifyPlaylist(metadata domain.MetaData) (*domain.Playlist, error) {
	return nil, nil
}
