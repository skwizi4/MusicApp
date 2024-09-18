package Spotify

import (
	"MusicApp/internal/domain"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (s ServiceSpotify) CreateApiKey() {
	//todo Create method that give s.ApiKey needed value
	s.ApiKey = ""
}
func (s ServiceSpotify) CreateRequest(method, endpoint string) (*http.Request, error) {
	url := s.BaseUrl + endpoint

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+s.ApiKey)
	return req, nil
}
func (s ServiceSpotify) doRequest(req *http.Request) (*http.Response, error) {
	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return &http.Response{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return resp, nil
}
func (s ServiceSpotify) decodeRespTrackById(response *http.Response) (domain.Song, error) {
	var track spotifyTrackById
	if err := json.NewDecoder(response.Body).Decode(&track); err != nil {
		return domain.Song{}, err
	}

	return domain.Song{
		Title:  track.Name,
		Artist: track.Artists[0].Name,
		Album:  track.Album.Name,
	}, nil
}
func (s ServiceSpotify) decodeRespPlaylistId(response *http.Response) (domain.Playlist, error) {
	var playlist spotifyPlaylistById
	if err := json.NewDecoder(response.Body).Decode(&playlist); err != nil {
		return domain.Playlist{}, err
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
	return p, nil
}
func (s ServiceSpotify) decodeRespTrackByName(response *http.Response) (domain.Song, error) {
	var track spotifySongByName
	if err := json.NewDecoder(response.Body).Decode(&track); err != nil {
		return domain.Song{}, err
	}

	return domain.Song{
		Title:  track.Tracks.Items[0].Name,
		Artist: track.Tracks.Items[0].Artists[0].Name,
		Album:  track.Tracks.Items[0].Album.Name,
	}, nil
}
func ParseTrackIDFromURL(link string) (string, error) {
	if strings.Contains(link, "https://open.spotify.com/track/") {
		id := strings.Split(link, "https://open.spotify.com/track/")[1]
		return strings.Split(id, "?")[0], nil
	}
	return "", errors.New("can't parse ID")
}
func ParsePlaylistIDFromURL(link string) (string, error) {
	if strings.Contains(link, "https://open.spotify.com/playlist/") {
		id := strings.Split(link, "https://open.spotify.com/playlist/")[1]
		return strings.Split(id, "?")[0], nil
	}
	return "", errors.New("can't parse ID")
}
