package Youtube

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
	ProcessingYoutubeMediaBySpotifySongID *domain.ProcessingYoutubeMediaBySpotifySongLink
	errChannel                            *ErrChan.ErrorChannel
	spotifyService                        services.SpotifyService
	youtubeService                        services.YouTubeService
	cfg                                   *config.Config
	logger                                logs.GoLogger
}

func New(bot *tg.Bot, processingSpotifySongs *domain.ProcessingYoutubeMediaBySpotifySongLink,
	errChan *ErrChan.ErrorChannel, cfg *config.Config, logger logs.GoLogger,
) Handler {
	return Handler{
		bot:                                   bot,
		spotifyService:                        Spotify.NewSpotifyService(cfg),
		youtubeService:                        YouTubeService.NewYouTubeService(cfg),
		ProcessingYoutubeMediaBySpotifySongID: processingSpotifySongs,
		errChannel:                            errChan,
		cfg:                                   cfg,
		logger:                                logger,
	}
}
