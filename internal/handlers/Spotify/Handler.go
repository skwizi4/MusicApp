package Spotify

import (
	"MusicApp/internal/config"
	"MusicApp/internal/domain"
	"MusicApp/internal/services"
	"MusicApp/internal/services/Spotify"
	"fmt"
	"github.com/skwizi4/lib/ErrChan"
	tg "gopkg.in/tucnak/telebot.v2"
)

type Handler struct {
	bot                    *tg.Bot
	processingSpotifySongs *domain.ProcessingSpotifySongs
	errChannel             *ErrChan.ErrorChannel
	spotifyService         services.SpotifyService
	cfg                    config.Config
}

func New(bot *tg.Bot, processingSpotifySongs *domain.ProcessingSpotifySongs, errChan *ErrChan.ErrorChannel, cfg config.Config,
) Handler {
	return Handler{
		bot:                    bot,
		spotifyService:         Spotify.NewSpotifyService(cfg),
		processingSpotifySongs: processingSpotifySongs,
		errChannel:             errChan,
		cfg:                    cfg,
	}
}

func (h Handler) GetSongByYoutubeLink(msg *tg.Message) (*domain.Song, error) {
	process := h.processingSpotifySongs.GetOrCreate(msg.Chat.ID)

	switch process.Step {
	case domain.ProcessSpotifySongStart:
		if _, err := h.bot.Send(msg.Sender, "Send link of song that you wanna find"); err != nil {
			if err = h.processingSpotifySongs.Delete(msg.Chat.ID); err != nil {
				h.errChannel.HandleError(err)
				return nil, err
			}
			h.errChannel.HandleError(err)

			return nil, err
		}
		if err := h.processingSpotifySongs.UpdateStep(domain.ProcessSpotifySongEnd, msg.Chat.ID); err != nil {
			h.errChannel.HandleError(err)
			return nil, err
		}
	case domain.ProcessSpotifySongEnd:
		track, err := h.spotifyService.GetSpotifyTrackById(msg.Text)
		if err != nil {
			h.errChannel.HandleError(err)
			return nil, err
		}

		fmt.Println(track)
		//todo - fix error - tg: unsupported what argument (track)
		if _, err = h.bot.Send(msg.Sender, track); err != nil {
			h.errChannel.HandleError(err)
			return nil, err
		}

	}
	return nil, nil

}
func (h Handler) GetSongByMetaData(metadata domain.MetaData) (*domain.Song, error) {
	return nil, nil
}
func (h Handler) GetPlaylistByYoutubeLink(youtubeLink string) (*domain.Playlist, error) {
	return nil, nil
}

func (h Handler) GetPlaylistByMetaData(metadata domain.MetaData) (*domain.Playlist, error) {
	return nil, nil
}
