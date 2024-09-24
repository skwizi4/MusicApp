package Spotify

import (
	"MusicApp/internal/domain"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (s ServiceSpotify) CreateRequest(method, endpoint string) (*http.Request, error) {
	Url := s.BaseUrl + endpoint

	req, err := http.NewRequest(method, Url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+s.Token)
	return req, nil
}
func (s ServiceSpotify) doRequest(req *http.Request) (*http.Response, error) {
	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return &http.Response{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return resp, nil
}
func (s ServiceSpotify) decodeRespTrackById(response *http.Response) (*domain.Song, error) {
	var track spotifyTrackById
	if err := json.NewDecoder(response.Body).Decode(&track); err != nil {
		return nil, err
	}

	return &domain.Song{
		Title:  track.Name,
		Artist: track.Artists[0].Name,
		Album:  track.Album.Name,
	}, nil
}
func (s ServiceSpotify) decodeRespPlaylistId(response *http.Response) (*domain.Playlist, error) {
	var playlist spotifyPlaylistById
	if err := json.NewDecoder(response.Body).Decode(&playlist); err != nil {
		return nil, err
	}
	var p domain.Playlist
	p.Title = playlist.Name
	p.Description = playlist.Description
	p.Owner = playlist.Owner.DisplayName
	p.ExternalUrl = playlist.ExternalURL.Spotify
	for i, song := range playlist.Tracks.Items {
		p.Songs[i] = domain.Song{
			Title:  song.Track.Name,
			Artist: song.Track.Artists[0].Name,
			Album:  song.Track.Album.Name,
		}
	}
	return &p, nil
}
func (s ServiceSpotify) decodeRespTrackByName(response *http.Response) (*domain.Song, error) {
	var track spotifySongByName
	if err := json.NewDecoder(response.Body).Decode(&track); err != nil {
		return nil, err
	}

	return &domain.Song{
		Title:  track.Tracks.Items[0].Name,
		Artist: track.Tracks.Items[0].Artists[0].Name,
		Album:  track.Tracks.Items[0].Album.Name,
	}, nil
}
func ParseSpotifyIDFromURL(link string) (string, string, error) {
	if strings.Contains(link, "https://open.spotify.com/track/") {
		id := strings.Split(link, "https://open.spotify.com/track/")[1]
		return "track", strings.Split(id, "?")[0], nil
	}
	if strings.Contains(link, "https://open.spotify.com/playlist/") {
		id := strings.Split(link, "https://open.spotify.com/playlist/")[1]
		return "playlist", strings.Split(id, "?")[0], nil
	}
	return "", "", errors.New("can't parse ID, invalid URL format")
}

func (s *ServiceSpotify) RequestToken() error {
	tokenData := url.Values{}
	tokenData.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(tokenData.Encode()))
	if err != nil {
		s.Logger.ErrorFrmt("failed to create request: ", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(s.ClientId + ":" + s.ClientSecret))
	req.Header.Add("Authorization", "Basic "+authHeader)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.Logger.ErrorFrmt("failed to execute request: ", err)
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		s.Logger.ErrorFrmt("unexpected status code: ", resp.StatusCode)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.Logger.ErrorFrmt("failed to read response body: ", err)
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var tokenResponse map[string]interface{}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		s.Logger.ErrorFrmt("failed to unmarshal response: ", err)
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	token, ok := tokenResponse["access_token"].(string)
	if !ok {
		s.Logger.Error("access_token not found in response")
		return fmt.Errorf("access_token not found in response")
	}

	s.Token = token
	s.Logger.Info("Key successfully retrieved")

	return nil
}
