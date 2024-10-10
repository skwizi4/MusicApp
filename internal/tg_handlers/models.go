package tg_handlers

import (
	"MusicApp/internal/config"
	"MusicApp/internal/domain"
	"MusicApp/internal/handlers"
	"MusicApp/internal/repo/MongoDB"
	"github.com/skwizi4/lib/ErrChan"
	"github.com/skwizi4/lib/logs"
	tg "gopkg.in/tucnak/telebot.v2"
)

type Handler struct {
	bot                                     *tg.Bot
	youtubeHandler                          handlers.Youtube
	spotifyHandler                          handlers.Spotify
	errChannel                              *ErrChan.ErrorChannel
	processingFindSongByMetadata            *domain.ProcessingFindSongByMetadata
	processingFillYoutubePlaylists          *domain.ProcessingCreateAndFillYoutubePlaylists
	processingYoutubeMediaBySpotifySongLink *domain.ProcessingYoutubeMediaBySpotifySongLink
	processingSpotifySongByYoutubeMediaLink *domain.ProcessingSpotifySongByYoutubeMediaLink
	logger                                  logs.GoLogger
	cfg                                     *config.Config
	mongo                                   *MongoDB.MongoDB
}

func New(bot *tg.Bot, youtubeHandler handlers.Youtube, spotifyHandler handlers.Spotify, ProcessingFinSongByMetadata *domain.ProcessingFindSongByMetadata,
	ProcessingFillYoutubePlaylists *domain.ProcessingCreateAndFillYoutubePlaylists, ProcessingYoutubeMediaBySpotifySongID *domain.ProcessingYoutubeMediaBySpotifySongLink,
	ProcessingSpotifySongByYoutubeMediaLink *domain.ProcessingSpotifySongByYoutubeMediaLink,
	cfg *config.Config, mongo *MongoDB.MongoDB, errChan *ErrChan.ErrorChannel, logger logs.GoLogger) Handler {
	return Handler{
		bot:                                     bot,
		youtubeHandler:                          youtubeHandler,
		spotifyHandler:                          spotifyHandler,
		errChannel:                              errChan,
		logger:                                  logger,
		processingFindSongByMetadata:            ProcessingFinSongByMetadata,
		processingFillYoutubePlaylists:          ProcessingFillYoutubePlaylists,
		processingYoutubeMediaBySpotifySongLink: ProcessingYoutubeMediaBySpotifySongID,
		processingSpotifySongByYoutubeMediaLink: ProcessingSpotifySongByYoutubeMediaLink,
		cfg:                                     cfg,
		mongo:                                   mongo,
	}
}
