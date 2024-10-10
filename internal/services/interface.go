package services

import "MusicApp/internal/domain"

type (
	SpotifyService interface {
		GetSpotifyTrackMetadataByLink(link string) (*domain.Song, error) // ok
		GetSpotifyPlaylistDataByLink(link string) (*domain.Playlist, error)
		GetSpotifyTrackByMetadata(data domain.MetaData) (*domain.Song, error)
		CreateAndFillSpotifyPlaylist(playlist domain.Playlist) (*domain.Playlist, error)
	}
	YouTubeService interface {
		GetYoutubeMediaByLink(link string) (*domain.Song, error)
		GetYoutubePlaylistDataByLink(link string) (*domain.Playlist, error)
		GetYoutubeMediaByMetadata(data domain.MetaData) (*domain.Song, error) // ok
		CreateAndFillYoutubePlaylist(playlist domain.Playlist, token string) (*domain.Playlist, error)
	}
)
