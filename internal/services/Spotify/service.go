package Spotify

import (
	"MusicApp/internal/config"
	"MusicApp/internal/domain"
	logger "github.com/skwizi4/lib/logs"
	"net/http"
)

// todo - Refactor

func NewSpotifyService(cfg *config.Config) ServiceSpotify {
	return ServiceSpotify{
		BaseUrl:      BaseUrl,
		ClientId:     cfg.SpotifyCfg.ClientID,
		ClientSecret: cfg.SpotifyCfg.ClientSecret,
		Logger:       logger.InitLogger(),
	}
}

func (s ServiceSpotify) GetSpotifyTrackMetadataByLink(link string) (*domain.Song, error) {
	id, err := GetID(link)

	if err != nil {
		return nil, err
	}
	endpoint, err := s.CreateEndpointSpotifyTrackById(id)
	if err != nil {
		return nil, err
	}
	if err = s.RequestToken(); err != nil {
		return nil, err
	}

	resp, err := s.createAndExecuteRequest(http.MethodGet, endpoint)
	if err != nil {
		return nil, err
	}

	Track, err := s.decodeRespTrackById(resp)
	if err != nil {
		return nil, err
	}

	return Track, nil

}

// GetSpotifyPlaylistById todo - Check bugs
func (s ServiceSpotify) GetSpotifyPlaylistDataByLink(link string) (*domain.Playlist, error) {
	id, err := GetID(link)

	if err != nil {
		return nil, err
	}
	err = s.RequestToken()
	if err != nil {
		return nil, err
	}
	endpoint, err := s.CreateEndpointSpotifyPlaylistById(id)
	if err != nil {
		return nil, err
	}
	resp, err := s.createAndExecuteRequest(http.MethodGet, endpoint)
	if err != nil {
		return nil, err
	}
	Playlist, err := s.decodeRespPlaylistId(resp)
	if err != nil {
		return nil, err
	}
	return Playlist, nil

}

// todo GetSpotifyTrackByMetadata  -refactor
func (s ServiceSpotify) GetSpotifyTrackByMetadata(data domain.MetaData) (*domain.Song, error) {
	endpoint, err := s.CreateEndpointSpotifyTrackByMetadata(data.Title, data.Artist)
	if err != nil {
		return nil, err
	}
	if err = s.RequestToken(); err != nil {
		return nil, err
	}
	resp, err := s.createAndExecuteRequest(http.MethodGet, endpoint)
	if err != nil {
		return nil, err
	}

	Song, err := s.decodeRespTrackByName(resp)
	if err != nil {
		return nil, err
	}
	return Song, nil

}
func (s ServiceSpotify) CreateAndFillSpotifyPlaylist(playlist domain.Playlist) (*domain.Playlist, error) {
	return nil, nil
}
