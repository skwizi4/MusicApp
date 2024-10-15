package services

import "MusicApp/internal/domain"

type (
	SpotifyService interface {
		GetSpotifyTrackMetadataByLink(link string) (*domain.Song, error) // ok
		GetSpotifyPlaylistDataByLink(link string) (*domain.Playlist, error)
		GetSpotifyTrackByMetadata(data domain.MetaData) (*domain.Song, error)
		CreateSpotifyPlaylist(Title, AuthToken, SpotifyUserId string) (string, error)
		FillSpotifyPlaylist(YouTubePlaylist *domain.Playlist, SpotifyPlaylistId, AuthToken string) (*domain.Playlist, error)
	}
	YouTubeService interface {
		GetYoutubeMediaByLink(link string) (*domain.Song, error)
		GetYoutubePlaylistDataByLink(link string) (*domain.Playlist, error)
		GetYoutubeMediaByMetadata(data domain.MetaData) (*domain.Song, error) // ok
		CreateYoutubePlaylist(Title string, token string) (string, error)
		FillYoutubePlaylist(SpotifyPlaylist *domain.Playlist, YoutubePlaylistId, token string) (*domain.Playlist, error)
	}
)
