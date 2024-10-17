package Spotify

import (
	"MusicApp/internal/domain"
	"net/http"
)

func (s ServiceSpotify) GetSpotifyTrackMetadataByLink(link string) (*domain.Song, error) {
	id, err := GetID(link)
	if err != nil {
		return nil, err
	}
	endpoint, err := s.MakeEndpointSpotifyTrackById(id)
	if err != nil {
		return nil, err
	}
	if err = s.RequestToken(); err != nil {
		return nil, err
	}

	Body, err := s.createAndExecuteRequest(http.MethodGet, endpoint, NilAuthToken, nil)
	if err != nil {
		return nil, err
	}

	return s.decodeRespTrackById(Body)

}
func (s ServiceSpotify) GetSpotifyTrackByMetadata(data domain.MetaData) (*domain.Song, error) {
	endpoint, err := s.MakeEndpointSpotifyTrackByMetadata(data.Title, data.Artist)
	if err != nil {
		return nil, err
	}
	if err = s.RequestToken(); err != nil {
		return nil, err
	}
	Body, err := s.createAndExecuteRequest(http.MethodGet, endpoint, NilAuthToken, nil)
	if err != nil {
		return nil, err
	}

	return s.decodeSpotifyRespTrackByMetadata(Body)

}

func (s ServiceSpotify) GetSpotifyPlaylistDataByLink(link string) (*domain.Playlist, error) {
	id, err := GetID(link)
	if err != nil {
		return nil, err
	}
	if err = s.RequestToken(); err != nil {
		return nil, err
	}
	endpoint, err := s.MakeEndpointSpotifyPlaylistById(id)
	if err != nil {
		return nil, err
	}
	Body, err := s.createAndExecuteRequest(http.MethodGet, endpoint, NilAuthToken, nil)
	if err != nil {
		return nil, err
	}

	return s.decodeRespPlaylistId(Body)

}

func (s ServiceSpotify) CreateSpotifyPlaylist(Title, AuthToken, SpotifyUserId string) (string, error) {
	endpoint, body, err := s.MakeEndpointCreateSpotifyPlaylist(Title, SpotifyUserId)
	if err != nil {
		return "", err
	}
	Body, err := s.createAndExecuteRequest(http.MethodPost, endpoint, AuthToken, body)
	if err != nil {
		return "", err
	}

	return s.decodeRespCreateSpotifyPlaylist(Body)
}
func (s ServiceSpotify) FillSpotifyPlaylist(YouTubePlaylist *domain.Playlist, SpotifyPlaylistId, AuthToken string) (*domain.Playlist, error) {

	return s.FillPlaylist(YouTubePlaylist, SpotifyPlaylistId, AuthToken)
}
