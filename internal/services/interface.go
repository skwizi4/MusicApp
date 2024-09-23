package services

import "MusicApp/internal/domain"

type (
	SpotifyService interface {
		GetSpotifyTrackById(link string) (domain.Song, error)
		GetSpotifyPlaylistById(link string) (domain.Playlist, error)
		GetSpotifyTrackByName(data domain.MetaData) (domain.Song, error)
	}
	YouTubeService interface {
		GetYoutubeMediaByID(link string) (*domain.Song, error)
		GetYoutubePlaylistByID(link string) (*domain.Playlist, error)
		GetMediaByMetadata(data domain.MetaData) (*domain.Song, error)
	}
)
