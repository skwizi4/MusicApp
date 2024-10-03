package tg_handlers

import (
	"MusicApp/internal/domain"
	"MusicApp/internal/handlers"
	"fmt"
	"github.com/skwizi4/lib/ErrChan"
	tg "gopkg.in/tucnak/telebot.v2"
)

type Handler struct {
	bot                         *tg.Bot
	spotifyHandler              handlers.Spotify
	youtubeHandler              handlers.YouTube
	errChannel                  *ErrChan.ErrorChannel
	processingFinSongByMetadata domain.ProcessingFindSongByMetadata
}

func New(bot *tg.Bot, spotifyHandler handlers.Spotify, youtubeHandler handlers.YouTube, errChan *ErrChan.ErrorChannel) Handler {
	return Handler{
		bot:            bot,
		spotifyHandler: spotifyHandler,
		youtubeHandler: youtubeHandler,
		errChannel:     errChan,
	}
}

// todo - Добавить в основной тгешный хендлер хендлера спотифая и ютуба, дописать запрос к ютубу после получения данных
func (h Handler) SpotifySong(msg *tg.Message) {
	err := h.spotifyHandler.GetSongByYoutubeLink(msg)
	if err != nil {
		h.errChannel.HandleError(err)
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
	process := h.processingFinSongByMetadata.GetOrCreate(msg.Chat.ID)
	switch process.Step {
	case domain.ProcessSpotifySongByMetadataStart:
		err := h.GetMetadata(msg)
		if err != nil {
			h.errChannel.HandleError(err)
			return
		}
	case domain.ProcessSpotifySongByMetadataTitle:
		err := h.GetMetadata(msg)
		if err != nil {
			h.errChannel.HandleError(err)
			return
		}
	case domain.ProcessSpotifySongByMetadataEnd:
		spotifySong, err := h.spotifyHandler.GetSongByMetaData(&domain.MetaData{Title: process.Song.Title, Artist: process.Song.Artist})
		if err != nil {

			h.errChannel.HandleError(err)
			return
		}
		youtubeSong, err := h.youtubeHandler.GetSongByMetaData(&domain.MetaData{Title: process.Song.Title, Artist: process.Song.Artist})
		if err != nil {
			h.errChannel.HandleError(err)
			return
		}
		SongPrint := fmt.Sprintf("Spotify song title: %s \n Spotify song artist: %s , \n Spotify song link: %s \n \n Youtube song title: %s  \n Youtube song artist: %s \n Youtube song link: %s",
			spotifySong.Artist, spotifySong.Title, spotifySong.Link, youtubeSong.Title, youtubeSong.Artist, youtubeSong.Link)
		if _, err := h.bot.Send(msg.Sender, SongPrint); err != nil {
			h.errChannel.HandleError(err)
			return
		}
	}

}
