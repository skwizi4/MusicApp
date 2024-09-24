package YouTube

import (
	"MusicApp/internal/domain"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func ParseYouTubeIDFromURL(link string) (string, string, error) {
	if strings.Contains(link, "https://www.youtube.com/watch?v=") {
		id := strings.Split(link, "https://www.youtube.com/watch?v=")[1]
		return "track", strings.Split(id, "?")[0], nil
	}
	if strings.Contains(link, "https://youtube.com/playlist?list=") {
		id := strings.Split(link, "https://youtube.com/playlist?list=")[1]
		return "playlist", strings.Split(id, "?")[0], nil
	}
	return "", "", errors.New("can't parse ID, invalid URL format")
}
func (y ServiceYouTube) CreateRequest(method, endpoint string) (*http.Request, error) {
	Url := y.BaseUrl + endpoint
	req, err := http.NewRequest(method, Url, nil)
	if err != nil {
		return nil, err
	}
	// Func with token//
	return req, nil
}

// DoRequest todo fix bugs
func (y ServiceYouTube) DoRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	fmt.Println(resp.Body)
	if resp.StatusCode != 200 {
		return resp, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return resp, nil
}
func (y ServiceYouTube) CreateEndpointYoutubeMedia(id string) (string, error) {
	if y.Key == "" {
		return "", errors.New("YouTube key is empty")
	}
	url := fmt.Sprintf("videos?id=%s&key=%s&part=snippet,status", id, y.Key)
	return url, nil
}
func (y ServiceYouTube) CreateEndpointYoutubePlaylist(id string) (string, error) {
	if y.Key == "" {
		return "", errors.New("YouTube key is empty")
	}
	url := fmt.Sprintf("playlists?id=%s&key=%s&part=snippet,status", id, y.Key)
	return url, nil
}

func DecodeRespMediaById(resp *http.Response) (*domain.Song, error) {
	var Media youtubeMediaById
	if err := json.NewDecoder(resp.Body).Decode(&Media); err != nil {
		return nil, errors.New("can't decode response")
	}
	return &domain.Song{
		Title:  Media.Items[0].Snippet.Title,
		Artist: Media.Items[0].Snippet.ChanelName,
	}, nil
}
func DecodeRespPlaylistById(resp *http.Response) (*domain.Playlist, error) {
	var Playlist youtubePlaylistById
	if err := json.NewDecoder(resp.Body).Decode(&Playlist); err != nil {
		return &domain.Playlist{}, errors.New("can't decode response")
	}
	return nil, nil
}
