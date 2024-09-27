package tg_handlers

import (
	"MusicApp/internal/handlers"
	"github.com/skwizi4/lib/ErrChan"
	tg "gopkg.in/tucnak/telebot.v2"
)

type Handler struct {
	bot            *tg.Bot
	spotifyHandler handlers.Spotify
	youtubeHandler handlers.YouTube
	errChan        *ErrChan.ErrorChannel
}

func New(bot *tg.Bot, spotifyHandler handlers.Spotify, youtubeHandler handlers.YouTube, errChan *ErrChan.ErrorChannel) Handler {
	return Handler{
		bot:            bot,
		spotifyHandler: spotifyHandler,
		youtubeHandler: youtubeHandler,
		errChan:        errChan,
	}
}

// todo - Добавить в основной тгешный хендлер хендлера спотифая и ютуба, дописать запрос к ютубу после получения данных
func (h Handler) SpotifySong(msg *tg.Message) {
	err := h.spotifyHandler.GetSongByYoutubeLink(msg)
	if err != nil {
		h.errChan.HandleError(err)
	}

}
func (h Handler) Help(msg *tg.Message) {

}
func (h Handler) YouTubeSong(msg *tg.Message) {

}
func (h Handler) SpotifyPlaylist(msg *tg.Message) {

}
func (h Handler) YouTubePlaylist(msg *tg.Message) {

}
func (h Handler) FindSong(msg *tg.Message) {

}
