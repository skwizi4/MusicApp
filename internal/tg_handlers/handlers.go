package tg_handlers

import (
	tg "gopkg.in/tucnak/telebot.v2"
)

type Handler struct {
}

func New() Handler {
	return Handler{}
}
func (h Handler) SpotifySong(msg *tg.Message) {

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
