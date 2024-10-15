package Spotify

import (
	"MusicApp/internal/domain"
)

func (h Handler) GetSpotifySongByLink(spotifyLink string) (*domain.Song, error) {
	return h.spotifyService.GetSpotifyTrackMetadataByLink(spotifyLink)
}

func (h Handler) GetSpotifySongByMetaData(metadata *domain.MetaData) (*domain.Song, error) {
	return h.spotifyService.GetSpotifyTrackByMetadata(*metadata)
}

func (h Handler) GetSpotifyPlaylistByLink(youtubeLink string) (*domain.Playlist, error) {
	return h.spotifyService.GetSpotifyPlaylistDataByLink(youtubeLink)
}

// CreateSpotifyPlaylist todo refactor
func (h Handler) CreateSpotifyPlaylist(Title, AuthToken, SpotifyUserId string) (string, error) {
	return h.spotifyService.CreateSpotifyPlaylist(Title, AuthToken, SpotifyUserId)

}
func (h Handler) FillSpotifyPlaylist(playlist *domain.Playlist, playlistId, AuthToken string) (*domain.Playlist, error) {
	return h.spotifyService.FillSpotifyPlaylist(playlist, playlistId, AuthToken)
}
