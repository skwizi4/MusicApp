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

/*-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
//Creating and Executing Requests
/*-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
func (y ServiceYouTube) createAndExecuteRequest(method, endpoint, AuthToken string, Body io.Reader) (*io.ReadCloser, error) {
	if method == "" || endpoint == "" {
		return nil, errs.Errorf("method or endpoint are nil")
	}

	Url := y.BaseUrl + endpoint

	switch {
	case method == http.MethodGet:
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
	case method == http.MethodPost && endpoint == fillingPlaylistEndpoint:
		fmt.Println("here1")
		if Body == nil || AuthToken == "" {
			return nil, errs.Errorf("body or AuthToken is nil")
		}
		req, err := http.NewRequest(method, Url, Body)
		if err != nil {
			return nil, errs.Errorf("can't create request to create playlist: %v", err)

		}
		req.Header.Add("Authorization", "Bearer "+AuthToken)
		req.Header.Add("Content-Type", "application/json")
		resp, err := y.Client.Do(req)
		if err != nil {
			return nil, errs.Errorf("can't execute request to create playlist: %v", err)
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, errs.Errorf("can't read response: %v", err)
			}
			y.logger.ErrorFrmt(string(body))
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		return &resp.Body, nil
	case method == http.MethodPost && endpoint == creatingPlaylistEndpoint:
		if Body == nil || AuthToken == "" {
			return nil, errs.Errorf("body or AuthToken is nil")
		}

		req, err := http.NewRequest(method, Url, Body)
		if err != nil {
			return nil, errs.Errorf("can't create request to create playlist: %v", err)

		}
		req.Header.Add("Authorization", "Bearer "+AuthToken)
		req.Header.Add("Content-Type", "application/json")
		resp, err := y.Client.Do(req)
		if err != nil {
			return nil, errs.Errorf("can't execute request to create playlist: %v", err)
		}
		time.Sleep(750 * time.Millisecond)
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

		return &resp.Body, nil
	}

	return nil, errs.Errorf("method or endpoint invalid")
}

/*-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
/*Creating Endpoints
/*-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

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

func (y ServiceYouTube) MakeEndpointCreateYoutubePlaylist(title string) (string, io.Reader, error) {
	if title == "" {
		return "", nil, errs.Errorf(" title is empty")
	}
	payload := []byte(fmt.Sprintf(`{
	"snippet": {
		"title": "%s",
		"description": "Playlist created by tg_bot MusicApp"
	},
	"status": {
		"privacyStatus": "public"
	}
}`, title))
	return creatingPlaylistEndpoint, bytes.NewBuffer(payload), nil
}
func (y ServiceYouTube) MakeEndpointFillingYoutubePlaylist(playlistId, songId string) (string, io.Reader, error) {
	payload := []byte(fmt.Sprintf(`{
			"snippet": {
				"playlistId": "%s",
				"resourceId": {
					"kind": "youtube#video",
					"videoId": "%s"
				}
			}
		}`, playlistId, songId))
	return fillingPlaylistEndpoint, bytes.NewBuffer(payload), nil
}

/*-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
/*Decoders
/*-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

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
func (y ServiceYouTube) DecodeRespCreatePlaylist(body *io.ReadCloser) (string, error) {
	var PlaylistId youtubePlaylistIdResp
	if err := json.NewDecoder(*body).Decode(&PlaylistId); err != nil {
		return "", errs.Errorf("can't decode response: %v", err)
	}
	if PlaylistId.ID == "" {
		return "", errs.Errorf("nil playlist id")
	}
	fmt.Println(PlaylistId.ID)
	return PlaylistId.ID, nil
}

/*-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
/*FillingPlaylists
/*-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

func (y ServiceYouTube) FillPlaylistParams(body *io.ReadCloser) (*domain.Playlist, error) {
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

func (y ServiceYouTube) FillPlaylistSongs(body *io.ReadCloser, playlist *domain.Playlist) error {
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

func (y ServiceYouTube) FillPlaylist(token, youtubePlaylistId string, SpotifyPlaylist *domain.Playlist) (*domain.Playlist, error) {
	songs := make([]domain.Song, len(SpotifyPlaylist.Songs))

	for i, spotifyTrack := range SpotifyPlaylist.Songs {
		youtubeMedia, err := y.GetYoutubeMediaByMetadata(domain.MetaData{Title: spotifyTrack.Title, Artist: spotifyTrack.Artist})
		if err != nil {
			return nil, errs.Errorf("Error getting song: %v", err)
		}
		endpoint, body, err := y.MakeEndpointFillingYoutubePlaylist(youtubePlaylistId, youtubeMedia.Id)
		if err != nil {
			return nil, err
		}
		if _, err = y.createAndExecuteRequest(http.MethodPost, endpoint, token, body); err != nil {
			return nil, err
		}
		songs = append(songs, *youtubeMedia)
		if i == len(SpotifyPlaylist.Songs)-1 {
			break
		}

	}
	if len(songs) == 0 {
		return nil, errs.Errorf("playlist is nil")
	}
	YoutubePlaylist := &domain.Playlist{
		Owner:       SpotifyPlaylist.Owner,
		Title:       SpotifyPlaylist.Title,
		Description: "Playlist created by tg_bot MusicApp",
		ExternalUrl: fmt.Sprintf("https://youtube.com/playlist?list=%s", youtubePlaylistId),
		Songs:       songs,
	}

	return YoutubePlaylist, nil
}
