package YouTube

import (
	"MusicApp/internal/domain"
	tg "gopkg.in/tucnak/telebot.v2"
)

type Handler struct {
	bot tg.Bot
}

func New() Handler {
	return Handler{}
}

func (h Handler) GetSongBySpotifyLink(msg *tg.Message) error {
	return nil
}

func (h Handler) GetSongByMetaData(metadata *domain.MetaData) (*domain.Song, error) {
	return nil, nil
}

func (h Handler) GetPlaylistBySpotifyLink(spotifyLink string) (*domain.Playlist, error) {
	return nil, nil
}

func (h Handler) GetPlaylistByMetaData(metadata domain.MetaData) (*domain.Playlist, error) {
	return nil, nil
}
