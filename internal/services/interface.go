package services

import "MusicApp/internal/domain"

type (
	SpotifyService interface {
		GetSpotifyTrackById(link string) (*domain.Song, error)
		GetSpotifyPlaylistById(link string) (*domain.Playlist, error)
		GetSpotifyTrackByMetadata(data domain.MetaData) (*domain.Song, error)
		FillSpotifyPlaylist(playlist domain.Playlist) (*domain.Playlist, error)
	}
	YouTubeService interface {
		GetYoutubeMediaByID(link string) (*domain.Song, error)
		GetYoutubePlaylistByID(link string) (*domain.Playlist, error)
		GetYoutubeMediaByMetadata(data domain.MetaData) (*domain.Song, error)
		FillYoutubePlaylist(playlist domain.Playlist) (*domain.Playlist, error)
	}
)
