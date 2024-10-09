package handlers

import (
	"MusicApp/internal/domain"
	tg "gopkg.in/tucnak/telebot.v2"
)

// todo - Написать структуры: MetaData, Song, Playlist ( directory - domain)

type (
	Youtube interface {
		GetMediaBySpotifyLink(msg *tg.Message) error
		GetYoutubeMediaByMetaData(metadata *domain.MetaData) (*domain.Song, error)
		GetYoutubePlaylistBySpotifyLink(youtubeLink string) (*domain.Playlist, error)
		FillYoutubePlaylist(playlist domain.Playlist, AuthToken string) (*domain.Playlist, error)
	}

	Spotify interface {
		GetSpotifySongByYoutubeLink(msg *tg.Message) error
		GetSpotifySongByMetaData(metadata *domain.MetaData) (*domain.Song, error)
		GetSpotifyPlaylistByYoutubeLink(spotifyLink string) (*domain.Playlist, error)
		FillSpotifyPlaylist(metadata domain.MetaData) (*domain.Playlist, error)
	}
)
