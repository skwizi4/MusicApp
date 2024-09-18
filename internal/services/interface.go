package services

import "MusicApp/internal/domain"

type (
	SpotifyService interface {
		TrackById(link string) (domain.Song, error)
		PlaylistById(link string) (domain.Playlist, error)
		TrackByName(trackName, artist string) (domain.Song, error)
	}
	YouTubeService interface {
		GetYoutubeSongByID(songID string) (*domain.Song, error)
		GetYoutubePlaylist(playlistID string) (*domain.Playlist, error)
		FindMusicByMetadata(data domain.MetaData) (*domain.Song, error)
	}
)
