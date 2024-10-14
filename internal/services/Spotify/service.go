package Spotify

import (
	"MusicApp/internal/config"
	"MusicApp/internal/domain"
	"bytes"
	"fmt"
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

	Song, err := s.decodeSpotifyRespTrackByMetadata(resp)
	if err != nil {
		return nil, err
	}
	return Song, nil

}
func (s ServiceSpotify) CreateSpotifyPlaylist(Title, AuthToken, SpotifyUserId string) (string, error) {

	endpoint, body, err := s.CreateEndpointCreateAndFillSpotifyPlaylist(Title, SpotifyUserId)
	if err != nil {
		return "", err
	}

	resp, err := s.createAndExecuteCreateSpotifyPlaylistRequset(http.MethodGet, endpoint, body)
	if err != nil {
		return "", err
	}

	id, err := s.decodeRespCreateSpotifyPlaylist(resp)
	if err != nil {
		return "", err
	}

	return id, nil
}
func (s ServiceSpotify) FillSpotifyPlaylist(YouTubePlaylist domain.Playlist, AuthToken, SpotifyPlaylistId string) (*domain.Playlist, error) {
	SpotifyPlaylist := &domain.Playlist{}
	for i, track := range YouTubePlaylist.Songs {
		Song, err := s.GetSpotifyTrackByMetadata(domain.MetaData{Title: track.Title, Artist: track.Artist})
		if err != nil {
			return nil, err
		}
		req, _ := http.NewRequest("POST", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", SpotifyPlaylistId), bytes.NewBuffer([]byte(fmt.Sprintf(`{"uris": ["spotify:track:%s"]}`, Song.Id))))
		req.Header.Add("Authorization", "Bearer "+AuthToken)
		req.Header.Add("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf(resp.Status)
		}
		SpotifyPlaylist.Songs = append(SpotifyPlaylist.Songs, *Song)
		if i == len(YouTubePlaylist.Songs) {
			break
		}
	}
	SpotifyPlaylist.Owner = YouTubePlaylist.Owner
	SpotifyPlaylist.Title = YouTubePlaylist.Title
	SpotifyPlaylist.Description = "Playlist created by tg_bot MusicApp"
	SpotifyPlaylist.ExternalUrl = fmt.Sprintf("https://youtube.com/playlist?list=%s", SpotifyPlaylistId)
	return SpotifyPlaylist, nil
}
