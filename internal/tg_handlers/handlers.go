package tg_handlers

import (
	tg "gopkg.in/tucnak/telebot.v2"
)

type Handler struct {
}

func New() Handler {
	return Handler{}
}

// todo - Добавить в основной тгешный хендлер хендлера спотика и ютуба, дописать запрос к ютубу после получения данных
func (h Handler) SpotifySong(msg *tg.Message) {
	// SpotifyTrack, err := h.Spotify.GetSongByYoutubeLink()
	//if err != nil{
	//	h.ErrChan.HandleError(err)
	//}
	// YouTubeMedia, err := h.Youtube.GetSongByMetaData()
	//if err != nil{
	//	h.ErrChan.HandleError(err)
	//}
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
