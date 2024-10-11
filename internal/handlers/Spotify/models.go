package Spotify

import (
	"MusicApp/internal/config"
	"MusicApp/internal/services"
	"MusicApp/internal/services/Spotify"
	"github.com/skwizi4/lib/ErrChan"
	"github.com/skwizi4/lib/logs"
)

type Handler struct {
	errChannel     *ErrChan.ErrorChannel
	spotifyService services.SpotifyService
	cfg            *config.Config
	logger         logs.GoLogger
}

func New(errChan *ErrChan.ErrorChannel, cfg *config.Config, logger logs.GoLogger) Handler {
	return Handler{
		spotifyService: Spotify.NewSpotifyService(cfg),
		errChannel:     errChan,
		cfg:            cfg,
		logger:         logger,
	}
}
