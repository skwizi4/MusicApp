package tg_handlers

import (
	"MusicApp/internal/config"
	"MusicApp/internal/domain"
	"MusicApp/internal/handlers"
	"MusicApp/internal/repo/MongoDB"
	"MusicApp/internal/workflows"
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
	ProcessingCreateAndFillYoutubePlaylists *domain.ProcessingCreateAndFillYoutubePlaylists
	ProcessingYoutubeMediaBySpotifySongLink *domain.ProcessingYoutubeMediaBySpotifySongLink
	ProcessingSpotifySongByYoutubeMediaLink *domain.ProcessingSpotifySongByYoutubeMediaLink
	logger                                  logs.GoLogger
	cfg                                     *config.Config
	mongo                                   *MongoDB.MongoDB
	workFlows                               *workflows.WorkFlows
}

func New(bot *tg.Bot, youtubeHandler handlers.Youtube, spotifyHandler handlers.Spotify, ProcessingFinSongByMetadata *domain.ProcessingFindSongByMetadata,
	ProcessingFillYoutubePlaylists *domain.ProcessingCreateAndFillYoutubePlaylists, ProcessingYoutubeMediaBySpotifySongID *domain.ProcessingYoutubeMediaBySpotifySongLink,
	ProcessingSpotifySongByYoutubeMediaLink *domain.ProcessingSpotifySongByYoutubeMediaLink,
	cfg *config.Config, mongo *MongoDB.MongoDB, errChan *ErrChan.ErrorChannel, logger logs.GoLogger, workFlows *workflows.WorkFlows) Handler {
	return Handler{
		bot:                                     bot,
		youtubeHandler:                          youtubeHandler,
		spotifyHandler:                          spotifyHandler,
		errChannel:                              errChan,
		logger:                                  logger,
		processingFindSongByMetadata:            ProcessingFinSongByMetadata,
		ProcessingCreateAndFillYoutubePlaylists: ProcessingFillYoutubePlaylists,
		ProcessingYoutubeMediaBySpotifySongLink: ProcessingYoutubeMediaBySpotifySongID,
		ProcessingSpotifySongByYoutubeMediaLink: ProcessingSpotifySongByYoutubeMediaLink,
		cfg:                                     cfg,
		mongo:                                   mongo,
		workFlows:                               workFlows,
	}
}
