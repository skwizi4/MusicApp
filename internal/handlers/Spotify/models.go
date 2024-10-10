package Spotify

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
	bot                                   *tg.Bot
	ProcessingSpotifySongByYoutubeMediaId *domain.ProcessingSpotifySongByYoutubeMediaLink
	errChannel                            *ErrChan.ErrorChannel
	spotifyService                        services.SpotifyService
	youtubeService                        services.YouTubeService
	cfg                                   *config.Config
	logger                                logs.GoLogger
}

func New(bot *tg.Bot, cfg *config.Config, processingYoutubeSongsById *domain.ProcessingSpotifySongByYoutubeMediaLink, errChan *ErrChan.ErrorChannel, logger logs.GoLogger) Handler {
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
