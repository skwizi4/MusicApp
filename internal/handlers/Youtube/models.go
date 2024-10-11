package Youtube

import (
	"MusicApp/internal/config"
	"MusicApp/internal/services"
	YouTubeService "MusicApp/internal/services/YouTube"
	"github.com/skwizi4/lib/ErrChan"
	"github.com/skwizi4/lib/logs"
)

type Handler struct {
	errChannel     *ErrChan.ErrorChannel
	youtubeService services.YouTubeService
	cfg            *config.Config
	logger         logs.GoLogger
}

func New(
	errChan *ErrChan.ErrorChannel, cfg *config.Config, logger logs.GoLogger,
) Handler {
	return Handler{
		youtubeService: YouTubeService.NewYouTubeService(cfg),
		errChannel:     errChan,
		cfg:            cfg,
		logger:         logger,
	}
}
