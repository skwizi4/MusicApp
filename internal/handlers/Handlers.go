package handlers

import (
	"MusicApp/internal/domain"
	tg "gopkg.in/tucnak/telebot.v2"
)

// todo - Написать структуры: MetaData, Song, Playlist ( directory - domain)

type (
	Youtube interface {
		GetMediaBySpotifyLink(msg *tg.Message) error
		GetSongByMetaData(metadata *domain.MetaData) (*domain.Song, error)
		GetPlaylistByYoutubeLink(youtubeLink string) (*domain.Playlist, error)
		GetPlaylistByMetaData(metadata domain.MetaData) (*domain.Playlist, error)
	}

	Spotify interface {
		GetSpotifySongByYoutubeLink(msg *tg.Message) error
		GetSongByMetaData(metadata *domain.MetaData) (*domain.Song, error)
		GetPlaylistBySpotifyLink(spotifyLink string) (*domain.Playlist, error)
		GetPlaylistByMetaData(metadata domain.MetaData) (*domain.Playlist, error)
	}
)
