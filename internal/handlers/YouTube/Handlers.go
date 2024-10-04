package YouTube

import (
	"MusicApp/internal/config"
	"MusicApp/internal/domain"
	"MusicApp/internal/services"
	"MusicApp/internal/services/Spotify"
	YouTubeService "MusicApp/internal/services/YouTube"
	"github.com/skwizi4/lib/ErrChan"
	"github.com/skwizi4/lib/logs"
	tg "gopkg.in/tucnak/telebot.v2"
)

type Handler struct {
	bot                        *tg.Bot
	processingYoutubeSongsById *domain.ProcessingYoutubeMediaById
	errChannel                 *ErrChan.ErrorChannel
	spotifyService             services.SpotifyService
	youtubeService             services.YouTubeService
	cfg                        config.Config
	logger                     logs.GoLogger
}

func New(bot *tg.Bot, cfg config.Config, processingYoutubeSongsById *domain.ProcessingYoutubeMediaById, errChan *ErrChan.ErrorChannel, logger logs.GoLogger) Handler {
	return Handler{
		bot:                        bot,
		spotifyService:             Spotify.NewSpotifyService(cfg),
		youtubeService:             YouTubeService.NewYouTubeService(cfg),
		processingYoutubeSongsById: processingYoutubeSongsById,
		errChannel:                 errChan,
		cfg:                        cfg,
		logger:                     logger,
	}
}

func (h Handler) GetSongBySpotifyLink(msg *tg.Message) error {
	return nil
}

func (h Handler) GetSongByMetaData(metadata *domain.MetaData) (*domain.Song, error) {
	track, err := h.youtubeService.GetYoutubeMediaByMetadata(domain.MetaData{Title: metadata.Title, Artist: metadata.Artist})
	if err != nil {
		h.errChannel.HandleError(err)
		return nil, err
	}

	return track, nil
}

func (h Handler) GetPlaylistBySpotifyLink(spotifyLink string) (*domain.Playlist, error) {
	return nil, nil
}

func (h Handler) GetPlaylistByMetaData(metadata domain.MetaData) (*domain.Playlist, error) {
	return nil, nil
}
