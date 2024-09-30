package YouTube

import (
	"MusicApp/internal/domain"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

//todo Refactor

func GetID(url string) (string, error) {
	re := regexp.MustCompile(`^https://www.youtube.com/watch\?v=([^&]+)`)
	matches := re.FindStringSubmatch(url)

	if len(matches) < 2 {
		re = regexp.MustCompile(`^https://youtube.com/playlist\?list=([^&]+)`)
		matches = re.FindStringSubmatch(url)
		if len(matches) < 2 {
			return "", errors.New("YouTube URL Not Found")
		}
		return matches[1], nil
	}

	return matches[1], nil
}

func (y ServiceYouTube) createAndExecuteRequest(method, endpoint string) (*io.ReadCloser, error) {
	Url := y.BaseUrl + endpoint
	req, err := http.NewRequest(method, Url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return &resp.Body, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return &resp.Body, nil
}

func (y ServiceYouTube) CreateEndpointYoutubeMediaById(id string) (endpoint string, err error) {
	if y.Key == "" {
		return "", errors.New("YouTube key is empty")
	}
	return fmt.Sprintf("videos?id=%s&key=%s&part=snippet,status", id, y.Key), nil
}
func (y ServiceYouTube) CreateEndpointYoutubePlaylistSongs(id, nextPageToken string) (endpoint string, err error) {
	if y.Key == "" {
		return "", errors.New("YouTube key is empty")
	}
	if nextPageToken == "" {
		return fmt.Sprintf("playlistItems?playlistId=%s&key=%s&part=snippet,status&maxResults=50", id, y.Key), nil
	}
	return fmt.Sprintf("playlistItems?playlistId=%s&key=%s&part=snippet,status&maxResults=50&pageToken=%s", id, y.Key, nextPageToken), nil
}
func (y ServiceYouTube) CreateEndpointYoutubePlaylistParams(id string) (endpoint string, err error) {
	if y.Key == "" {
		return "", errors.New("YouTube key is empty")
	}
	return fmt.Sprintf("playlists?id=%s&key=%s&part=snippet,status", id, y.Key), nil
}
func (y ServiceYouTube) CreateEndpointYoutubeMediaByMetadata(data domain.MetaData) (endpoint string, err error) {
	if y.Key == "" {
		return "", errors.New("YouTube key is empty")
	}
	return fmt.Sprintf("search?part=snippet&maxResults=1&q=%s+%s&key=%s", url.QueryEscape(data.Title), url.QueryEscape(data.Artist), y.Key), nil
}

func DecodeRespMediaById(body *io.ReadCloser) (*domain.Song, error) {
	var Media youtubeMediaById
	if err := json.NewDecoder(*body).Decode(&Media); err != nil {
		return nil, errors.New("can't decode response")
	}

	return &domain.Song{
		Title:  Media.Items[0].Snippet.Title,
		Artist: Media.Items[0].Snippet.ChanelName,
		Link:   youtubeTrackDomen + Media.Items[0].VideoId,
	}, nil
}
func DecodeRespMediaByMetadata(body *io.ReadCloser) (*domain.Song, error) {
	var Media youtubeMediaByMetadata
	if err := json.NewDecoder(*body).Decode(&Media); err != nil {
		return nil, errors.New("can't decode response")
	}
	return &domain.Song{
		Title:  Media.Items[0].Snippet.Title,
		Artist: Media.Items[0].Snippet.ChanelName,
		Link:   Media.Items[0].Id.VideoId,
	}, nil
}

func FillPlaylistParams(body *io.ReadCloser, playlist *domain.Playlist) error {
	var Playlist youtubePlaylistParamsById
	if err := json.NewDecoder(*body).Decode(&Playlist); err != nil {
		return errors.New("can't decode response")
	}

	playlist.Title = Playlist.Items[0].Snippet.Title
	playlist.Owner = Playlist.Items[0].Snippet.ChannelTitle

	return nil
}

func FillPlaylistSongs(body *io.ReadCloser, playlist *domain.Playlist) (*domain.Playlist, error) {
	var Playlist youtubeResponsePlaylistMediaById
	if err := json.NewDecoder(*body).Decode(&Playlist); err != nil {
		return &domain.Playlist{}, errors.New("can't decode response")
	}
	for _, media := range Playlist.Items {
		playlist.Songs = append(playlist.Songs, domain.Song{
			Title:  media.Snippet.Title,
			Artist: media.Snippet.ChannelTitle,
		})
	}
	playlist.NextPageToken = Playlist.NextPageToken
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
