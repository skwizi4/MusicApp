package tg_handlers

import (
	"MusicApp/internal/domain"
	"MusicApp/internal/handlers"
	"github.com/skwizi4/lib/ErrChan"
	"github.com/skwizi4/lib/logs"
	tg "gopkg.in/tucnak/telebot.v2"
)

type Handler struct {
	bot                          *tg.Bot
	youtubeHandler               handlers.Youtube
	spotifyHandler               handlers.Spotify
	errChannel                   *ErrChan.ErrorChannel
	processingFindSongByMetadata *domain.ProcessingFindSongByMetadata
	logger                       logs.GoLogger
}

func New(bot *tg.Bot, youtubeHandler handlers.Youtube, spotifyHandler handlers.Spotify, processingFinSongByMetadata *domain.ProcessingFindSongByMetadata, errChan *ErrChan.ErrorChannel, logger logs.GoLogger) Handler {
	return Handler{
		bot:                          bot,
		youtubeHandler:               youtubeHandler,
		spotifyHandler:               spotifyHandler,
		errChannel:                   errChan,
		logger:                       logger,
		processingFindSongByMetadata: processingFinSongByMetadata,
	}
}

func (h Handler) YoutubeSong(msg *tg.Message) {
	err := h.youtubeHandler.GetMediaBySpotifyLink(msg)
	if err != nil {
		h.errChannel.HandleError(err)
	}

}
func (h Handler) SpotifySong(msg *tg.Message) {
	err := h.spotifyHandler.GetSpotifySongByYoutubeLink(msg)
	if err != nil {
		h.errChannel.HandleError(err)
	}

}
func (h Handler) Help(msg *tg.Message) {
	h.HelpOut(msg)
}

func (h Handler) FindSong(msg *tg.Message) {
	err := h.GetSongsByMetadata(msg)
	if err != nil {
		h.errChannel.HandleError(err)
		return
	}

}

func (h Handler) SpotifyPlaylist(msg *tg.Message) {

}
func (h Handler) YouTubePlaylist(msg *tg.Message) {

}
