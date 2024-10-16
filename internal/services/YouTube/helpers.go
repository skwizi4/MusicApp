package YouTube

import (
	"MusicApp/internal/domain"
	errors "MusicApp/internal/errors"
	"bytes"
	"encoding/json"
	"fmt"
	errs "github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

//todo Refactor

func GetID(url string) (string, error) {
	re := regexp.MustCompile(`^https://www.youtube.com/watch\?v=([^&]+)`)
	matches := re.FindStringSubmatch(url)

	if len(matches) < 2 {
		re = regexp.MustCompile(`^https://youtube.com/playlist\?list=([^&]+)`)
		matches = re.FindStringSubmatch(url)
		if len(matches) < 2 {
			return "", errs.Errorf("Youtube url is invalid")
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
	resp, err := y.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errs.Errorf("Error reading response body: %v", err)
		}
		y.logger.ErrorFrmt("body: %s ", string(body))
		return nil, errs.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return &resp.Body, nil
}

func (y ServiceYouTube) CreateEndpointYoutubeMediaById(id string) (string, error) {
	if y.Key == "" {
		return "", errs.Errorf("Spotify key is empty")
	}
	return fmt.Sprintf("videos?id=%s&key=%s&part=snippet,status", id, y.Key), nil
}
func (y ServiceYouTube) CreateEndpointYoutubeMediaByMetadata(data domain.MetaData) (string, error) {
	if y.Key == "" {
		return "", errs.Errorf("Spotify key is empty")
	}
	return fmt.Sprintf("search?part=snippet&maxResults=1&q=%s+%s&key=%s", url.QueryEscape(data.Title), url.QueryEscape(data.Artist), y.Key), nil
}
func (y ServiceYouTube) CreateEndpointYoutubePlaylistParams(id string) (string, error) {
	if y.Key == "" {
		return "", errs.Errorf("Spotify key is empty")
	}
	return fmt.Sprintf("playlists?id=%s&key=%s&part=snippet,status", id, y.Key), nil
}
func (y ServiceYouTube) CreateEndpointYoutubePlaylistSongs(id, nextPageToken string) (string, error) {
	if y.Key == "" {
		return "", errs.Errorf("Spotify key is empty")
	}
	if nextPageToken == "" {
		return fmt.Sprintf("playlistItems?playlistId=%s&key=%s&part=snippet,status&maxResults=50", id, y.Key), nil
	}
	return fmt.Sprintf("playlistItems?playlistId=%s&key=%s&part=snippet,status&maxResults=50&pageToken=%s", id, y.Key, nextPageToken), nil
}

func DecodeRespMediaById(body *io.ReadCloser) (*domain.Song, error) {
	var Media youtubeMediaById
	if err := json.NewDecoder(*body).Decode(&Media); err != nil {
		return nil, errs.Errorf("can't decode response: %v", err)
	}
	if len(Media.Items) == 0 || Media.Items[0].Snippet.Title == "" || Media.Items[0].Snippet.ChanelName == "" || Media.Items[0].VideoId == "" {
		return nil, errs.Errorf(errors.ErrInvalidParamsYoutube)
	}
	return &domain.Song{
		Title:  Media.Items[0].Snippet.Title,
		Artist: Media.Items[0].Snippet.ChanelName,
		Link:   youtubeTrackDomen + Media.Items[0].VideoId,
		Id:     Media.Items[0].VideoId,
	}, nil
}
func DecodeRespMediaByMetadata(body *io.ReadCloser) (*domain.Song, error) {
	var Media youtubeMediaByMetadata
	if err := json.NewDecoder(*body).Decode(&Media); err != nil {
		return nil, errs.Errorf("can't decode response: %v", err)
	}
	if len(Media.Items) == 0 || Media.Items[0].Snippet.Title == "" || Media.Items[0].Snippet.ChanelName == "" || Media.Items[0].Id.VideoId == "" {
		return nil, errs.Errorf(errors.ErrInvalidParamsYoutube)
	}

	return &domain.Song{
		Title:  Media.Items[0].Snippet.Title,
		Artist: Media.Items[0].Snippet.ChanelName,
		Link:   youtubeTrackDomen + Media.Items[0].Id.VideoId,
		Id:     Media.Items[0].Id.VideoId,
	}, nil
}

func FillPlaylistParams(body *io.ReadCloser) (*domain.Playlist, error) {
	var Playlist youtubePlaylistParamsById
	if err := json.NewDecoder(*body).Decode(&Playlist); err != nil {
		return nil, errs.Errorf("can't decode response: %v", err)
	}
	if len(Playlist.Items) == 0 || Playlist.Items[0].Snippet.Title == "" || Playlist.Items[0].Snippet.ChannelTitle == "" {
		return nil, errs.Errorf(errors.ErrInvalidParamsYoutube)
	}
	return &domain.Playlist{
		Title: Playlist.Items[0].Snippet.Title,
		Owner: Playlist.Items[0].Snippet.ChannelTitle}, nil

}

func FillPlaylistSongs(body *io.ReadCloser, playlist *domain.Playlist) error {
	var Playlist youtubeResponsePlaylistMediaById
	if err := json.NewDecoder(*body).Decode(&Playlist); err != nil {
		return errs.Errorf("can't decode response: %v", err)
	}
	for _, media := range Playlist.Items {
		playlist.Songs = append(playlist.Songs, domain.Song{
			Title:  media.Snippet.Title,
			Artist: media.Snippet.ChannelTitle,
			Id:     media.Snippet.ResourceId.VideoId,
		})
	}
	playlist.NextPageToken = Playlist.NextPageToken
	return nil
}
func (y ServiceYouTube) CreatePlaylist(token, playlistTitle string) (string, error) {
	var PlaylistId youtubePlaylistIdResp
	payload := []byte(fmt.Sprintf(`{
		"snippet": {
			"title": "%s",
			"description": "Playlist created by tg_bot MusicApp"
		},
		"status": {
			"privacyStatus": "public"
		}
	}`, playlistTitle))
	req, err := http.NewRequest("POST", "https://www.googleapis.com/youtube/v3/playlists?part=snippet,status", bytes.NewBuffer(payload))
	if err != nil {
		return "", errs.Errorf("can't create request to create playlist: %v", err)

	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := y.Client.Do(req)
	if err != nil {
		return "", errs.Errorf("can't execute request to create playlist: %v", err)
	}
	if err = json.NewDecoder(resp.Body).Decode(&PlaylistId); err != nil {
		return "", errs.Errorf("can't decode response: %v", err)
	}
	return PlaylistId.ID, nil
}

func (y ServiceYouTube) WriteInPlaylist(token, playlistId string, SpotifyPlaylist *domain.Playlist) (*domain.Playlist, error) {
	songs := make([]domain.Song, len(SpotifyPlaylist.Songs))
	for i, track := range SpotifyPlaylist.Songs {
		song, err := y.GetYoutubeMediaByMetadata(domain.MetaData{Title: track.Title, Artist: track.Artist})
		if err != nil {
			return nil, errs.Errorf("Error getting song: %v", err)
		}
		payload := []byte(fmt.Sprintf(`{
			"snippet": {
				"playlistId": "%s",
				"resourceId": {
					"kind": "youtube#video",
					"videoId": "%s"
				}
			}
		}`, playlistId, song.Id))
		req, err := http.NewRequest("POST", "https://www.googleapis.com/youtube/v3/playlistItems?part=snippet", bytes.NewBuffer(payload))
		if err != nil {
			return nil, errs.Errorf("can't create request to create playlist: %v", err)

		}
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Content-Type", "application/json")
		time.Sleep(750 * time.Millisecond)
		resp, err := y.Client.Do(req)
		if resp.StatusCode != 200 {
			if resp.StatusCode == 409 {
				resp, err = http.DefaultClient.Do(req)

			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, errs.Errorf("can't read response: %v", err)
			}
			y.logger.ErrorFrmt(string(body))
			return nil, errs.Errorf(resp.Status)
		}

		if err != nil {
			return nil, errs.Errorf("can't execute request to create playlist: %v", err)
		}
		songs = append(songs, *song)
		if i == len(SpotifyPlaylist.Songs)-1 {
			break
		}

	}
	var YoutubePlaylist = &domain.Playlist{
		Owner:       SpotifyPlaylist.Owner,
		Title:       SpotifyPlaylist.Title,
		Description: "Playlist created by tg_bot MusicApp",
		ExternalUrl: fmt.Sprintf("https://youtube.com/playlist?list=%s", playlistId),
		Songs:       songs,
	}

	return YoutubePlaylist, nil
}
