package Spotify

import (
	"MusicApp/internal/domain"
	"gopkg.in/tucnak/telebot.v2"
)

type Handler struct {
	bot telebot.Bot
}

func New() Handler {
	return Handler{}
}

func (h Handler) GetSongByYoutubeLink(youtubeLink string) (*domain.Song, error) {
	return nil, nil

}
func (h Handler) GetSongByMetaData(metadata domain.MetaData) (*domain.Song, error) {
	return nil, nil
}
func (h Handler) GetPlaylistByYoutubeLink(youtubeLink string) (*domain.Playlist, error) {
	return nil, nil
}

func (h Handler) GetPlaylistByMetaData(metadata domain.MetaData) (*domain.Playlist, error) {
	return nil, nil
}
