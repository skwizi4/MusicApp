package services

import "MusicApp/internal/domain"

type (
	SpotifyService interface {
		GetSpotifySongByID(songID string) (*domain.Song, error)
		GetSpotifyPlaylist(playlistID string) (*domain.Playlist, error)
		FindMusicByMetadata(data domain.MetaData) (*domain.Song, error)
	}
	YouTubeService interface {
		GetYoutubeSongByID(songID string) (*domain.Song, error)
		GetYoutubePlaylist(playlistID string) (*domain.Playlist, error)
		FindMusicByMetadata(data domain.MetaData) (*domain.Song, error)
	}
)
