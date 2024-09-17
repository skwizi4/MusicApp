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
func (s ServiceSpotify) DecodeRespTrack(response *http.Response) (domain.Song, error) {
	var track spotifyTrack
	if err := json.NewDecoder(response.Body).Decode(&track); err != nil {
		return domain.Song{}, err
	}

	return domain.Song{
		Title:  track.Name,
		Artist: track.Artists[0].Name,
		Album:  track.Album.Name,
	}, nil
}
func ParseIDFromURL(link string) (string, error) {
	if strings.Contains(link, "https://open.spotify.com/track/") {
		id := strings.Split(link, "https://open.spotify.com/track/")[1]
		return strings.Split(id, "?")[0], nil
	}
	return "", errors.New("can't parse ID")
}
