package handlers

import (
	"MusicApp/internal/domain"
)

// todo - Написать структуры: MetaData, Song, Playlist ( directory - domain)

type (
	Youtube interface {
		GetYoutubeMediaByLink(youtubeLink string) (*domain.Song, error)
		GetYoutubeMediaByMetaData(metadata *domain.MetaData) (*domain.Song, error)
		GetYoutubePlaylistByLink(youtubeLink string) (*domain.Playlist, error)
		CreateAndFillYoutubePlaylist(playlist domain.Playlist, AuthToken string) (*domain.Playlist, error)
	}

	Spotify interface {
		GetSpotifySongByLink(spotifyLink string) (*domain.Song, error)
		GetSpotifySongByMetaData(metadata *domain.MetaData) (*domain.Song, error)
		GetSpotifyPlaylistByLink(spotifyLink string) (*domain.Playlist, error)
		FillSpotifyPlaylist(playlist domain.Playlist) (*domain.Playlist, error)
	}
)
