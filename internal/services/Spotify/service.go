package Spotify

import (
	"MusicApp/internal/domain"
	"bytes"
	"fmt"
	"io"
	"log"
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
	endpoint, err := s.MakeEndpointSpotifyTrackByMetadata(data.Title, data.Artist)
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
	endpoint, body, err := s.MakeEndpointCreateSpotifyPlaylist(Title, SpotifyUserId)
	if err != nil {
		return "", err
	}
	resp, err := s.createAndExecuteCreateSpotifyPlaylistRequset(http.MethodPost, endpoint, AuthToken, body)
	if err != nil {
		return "", err
	}

	id, err := s.decodeRespCreateSpotifyPlaylist(resp)
	if err != nil {
		return "", err
	}

	return id, nil
}
func (s ServiceSpotify) FillSpotifyPlaylist(YouTubePlaylist *domain.Playlist, SpotifyPlaylistId, AuthToken string) (*domain.Playlist, error) {
	SpotifyPlaylist := &domain.Playlist{}
	for i, track := range YouTubePlaylist.Songs {
		Song, err := s.GetSpotifyTrackByMetadata(domain.MetaData{Title: track.Title, Artist: track.Artist})
		if err != nil {
			return nil, err
		}
		req, _ := http.NewRequest("POST", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", SpotifyPlaylistId), bytes.NewBuffer([]byte(fmt.Sprintf(`{"uris": ["spotify:track:%s"]}`, Song.Id))))
		req.Header.Add("Authorization", "Bearer "+AuthToken)
		req.Header.Add("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(body))
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
	SpotifyPlaylist.ExternalUrl = fmt.Sprintf("https://open.spotify.com/playlist/%s", SpotifyPlaylistId)
	return SpotifyPlaylist, nil
}
