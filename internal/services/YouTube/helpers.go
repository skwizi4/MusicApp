package YouTube

import (
	"MusicApp/internal/domain"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
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
func (y ServiceYouTube) createAndExecuteRequest(method, endpoint string) (*http.Response, error) {
	Url := y.BaseUrl + endpoint
	req, err := http.NewRequest(method, Url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != 200 {
		return resp, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return resp, nil
}

func (y ServiceYouTube) CreateEndpointYoutubeMediaById(id string) (string, error) {
	if y.Key == "" {
		return "", errors.New("YouTube key is empty")
	}
	endpoint := fmt.Sprintf("videos?id=%s&key=%s&part=snippet,status", id, y.Key)
	return endpoint, nil
}
func (y ServiceYouTube) CreateEndpointYoutubePlaylistSongs(id string) (string, error) {
	if y.Key == "" {
		return "", errors.New("YouTube key is empty")
	}
	endpoint := fmt.Sprintf("playlistItems?playlistId=%s&key=%s&part=snippet,status", id, y.Key)
	return endpoint, nil
}
func (y ServiceYouTube) CreateEndpointYoutubePlaylistParams(id string) (string, error) {
	if y.Key == "" {
		return "", errors.New("YouTube key is empty")
	}
	endpoint := fmt.Sprintf("playlists?id=%s&key=%s&part=snippet,status", id, y.Key)
	return endpoint, nil
}
func (y ServiceYouTube) CreateEndpointYoutubeMediaByMetadata(data domain.MetaData) (string, error) {
	if y.Key == "" {
		return "", errors.New("YouTube key is empty")
	}
	endpoint := fmt.Sprintf("search?part=snippet&maxResults=1&q=%s+%s&key=%s", url.QueryEscape(data.Title), url.QueryEscape(data.Artist), y.Key)
	return endpoint, nil
}

func DecodeRespMediaById(resp *http.Response) (*domain.Song, error) {
	var Media youtubeMediaById
	if err := json.NewDecoder(resp.Body).Decode(&Media); err != nil {
		return nil, errors.New("can't decode response")
	}

	return &domain.Song{
		Title:  Media.Items[0].Snippet.Title,
		Artist: Media.Items[0].Snippet.ChanelName,
		Link:   Media.Items[0].VideoId,
	}, nil
}
func DecodeRespMediaByMetadata(resp *http.Response) (*domain.Song, error) {
	var Media youtubeMediaByMetadata
	if err := json.NewDecoder(resp.Body).Decode(&Media); err != nil {
		return nil, errors.New("can't decode response")
	}
	return &domain.Song{
		Title:  Media.Items[0].Snippet.Title,
		Artist: Media.Items[0].Snippet.ChanelName,
		Link:   Media.Items[0].Id.VideoId,
	}, nil
}

func FillPlaylistParams(resp *http.Response, playlist *domain.Playlist) error {
	var Playlist youtubePlaylistById
	if err := json.NewDecoder(resp.Body).Decode(&Playlist); err != nil {
		return errors.New("can't decode response")
	}

	playlist.Title = Playlist.Items[0].Snippet.Title
	playlist.Owner = Playlist.Items[0].Snippet.ChannelTitle

	return nil
}

func FillPlaylist(resp *http.Response, playlist *domain.Playlist) (*domain.Playlist, error) {
	var Playlist youtubeResponsePlaylist
	if err := json.NewDecoder(resp.Body).Decode(&Playlist); err != nil {
		return &domain.Playlist{}, errors.New("can't decode response")
	}
	for _, media := range Playlist.Items {
		playlist.Songs = append(playlist.Songs, domain.Song{
			Title:  media.Snippet.Title,
			Artist: media.Snippet.ChannelTitle,
		})
	}
	return playlist, nil
}
func (y ServiceYouTube) CreatePlaylist(token, playlistTitle string) (string, error) {
	payload := []byte(fmt.Sprintf(`{
		"snippet": {
			"title": "%s",
			"description": ""
		},
		"status": {
			"privacyStatus": "public"
		}
	}`, playlistTitle))
	req, err := http.NewRequest("POST", "https://www.googleapis.com/youtube/v3/playlists?part=snippet,status", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	var PlaylistId youtubePlaylistIdResp
	if err = json.NewDecoder(resp.Body).Decode(&PlaylistId); err != nil {
		return "", errors.New("can't decode response")
	}
	return PlaylistId.ID, nil
}

func (y ServiceYouTube) FillYoutubePlaylist(token, playlistId string, tracks []domain.Song) (*domain.Playlist, error) {
	var YoutubePlaylist = &domain.Playlist{}
	for i, track := range tracks {
		song, err := y.GetYoutubeMediaByMetadata(domain.MetaData{Title: track.Title, Artist: track.Artist})
		if err != nil {
			return nil, err
		}
		payload := []byte(fmt.Sprintf(`{
			"snippet": {
				"playlistId": "%s",
				"resourceId": {
					"kind": "youtube#video",
					"videoId": "%s"
				}
			}
		}`, playlistId, song.Link))
		req, _ := http.NewRequest("POST", "https://www.googleapis.com/youtube/v3/playlistItems?part=snippet", bytes.NewBuffer(payload))
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logrus.Print(err)
		}
		if resp.StatusCode != 200 {
			return nil, errors.New(resp.Status)
		}
		YoutubePlaylist.Songs = append(YoutubePlaylist.Songs, *song)
		if i == len(tracks) {
			break
		}

	}
	return YoutubePlaylist, nil
}
