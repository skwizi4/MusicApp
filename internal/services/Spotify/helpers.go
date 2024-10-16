package Spotify

import (
	"MusicApp/internal/domain"
	"MusicApp/internal/errors"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	errs "github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func GetID(url string) (string, error) {
	re := regexp.MustCompile(`^https://open\.spotify\.com/track/([^?]+)`)
	matches := re.FindStringSubmatch(url)

	if len(matches) < 2 {
		re = regexp.MustCompile(`^^https://open\.spotify\.com/playlist/([^?]+)`)
		matches = re.FindStringSubmatch(url)
		if len(matches) < 2 {
			return "", errs.Errorf("Spotify url is invalid")
		}
		return matches[1], nil
	}

	return matches[1], nil
}

func (s ServiceSpotify) createAndExecuteRequest(method, endpoint string) (*io.ReadCloser, error) {
	Url := s.BaseUrl + endpoint
	req, err := http.NewRequest(method, Url, nil)
	if err != nil {
		return nil, err
	}
	if s.Token == "" {
		return nil, errs.Errorf("token is empty")
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, errs.Errorf("Spotify request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errs.Errorf("cant decode response: %v", err)
		}
		s.Logger.ErrorFrmt(string(body))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return &resp.Body, nil
}

func (s ServiceSpotify) MakeEndpointSpotifyTrackById(id string) (string, error) {
	if id == "" {
		return "", errs.Errorf("id  is empty")
	}
	return "/v1/tracks/" + id, nil
}

func (s ServiceSpotify) MakeEndpointSpotifyTrackByMetadata(title, artist string) (string, error) {

	if title == "" || artist == "" {
		return "", errs.Errorf("title or artist is empty")
	}
	return "/v1/search?q=https%3A%2F%2Fapi.spotify.com%2Fv1%2Fsearch%3Fquery%3Dtrack%3A" + url.QueryEscape(title) + "artist%3A" +
		url.QueryEscape(artist) + "%26type%3Dtrack%26offset%3D5%26limit%3D10&type=track", nil
}
func (s ServiceSpotify) MakeEndpointSpotifyPlaylistById(id string) (string, error) {
	if id == "" {
		return "", errs.Errorf("id is empty")
	}
	return "/v1/playlists/" + id, nil
}
func (s ServiceSpotify) MakeEndpointCreateSpotifyPlaylist(title, SpotifyUserId string) (string, io.Reader, error) {
	if title == "" || SpotifyUserId == "" {
		return "", nil, errs.Errorf("error, title or SpotifyUserId is nil")
	}
	body, err := json.Marshal(PlaylistCreateRequest{Name: title})
	if err != nil {
		return "", nil, errs.Errorf("can't marshal response: %v", err)
	}

	//todo refactor (  get user id and use it )
	return fmt.Sprintf("/v1/users/" + "ya9tkj10lvw2uva3iy11sby75" + "/playlists"), bytes.NewBuffer(body), nil
}

func (s ServiceSpotify) createAndExecuteCreateSpotifyPlaylistRequset(method, endpoint, AuthToken string, body io.Reader) (*io.ReadCloser, error) {
	if method == "" || endpoint == "" || body == nil {
		return nil, errs.Errorf("arguments  are nil")
	}
	Url := s.BaseUrl + endpoint
	req, err := http.NewRequest(method, Url, body)
	if err != nil {
		return nil, errs.Errorf("failed create request : %v", err)
	}
	if AuthToken == "" {
		return nil, errs.New("AuthToken is empty")
	}

	req.Header.Set("Authorization", "Bearer "+AuthToken)
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, errs.Errorf("Spotify request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errs.Errorf("can't read response: %v", err)
		}
		s.Logger.ErrorFrmt(string(respBody))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return &resp.Body, nil
}

func (s ServiceSpotify) decodeRespTrackById(body *io.ReadCloser) (*domain.Song, error) {
	var track spotifyTrackById
	if err := json.NewDecoder(*body).Decode(&track); err != nil {
		return nil, errs.Errorf("can't decode response: %v", err)
	}
	if track.Name == "" || &track == nil || track.Artists[0].Name == "" {
		return nil, errs.Errorf(errors.ErrInvalidParamsSpotify)
	}
	return &domain.Song{
		Title:  track.Name,
		Artist: track.Artists[0].Name,
		Album:  track.Album.Name,
	}, nil
}
func (s ServiceSpotify) decodeRespPlaylistId(body *io.ReadCloser) (*domain.Playlist, error) {
	var spotifyPlaylistResponse spotifyPlaylistById
	if err := json.NewDecoder(*body).Decode(&spotifyPlaylistResponse); err != nil {
		return nil, errs.Errorf("can't decode response: %v", err)
	}
	playlist := &domain.Playlist{
		Title:       spotifyPlaylistResponse.Name,
		Description: spotifyPlaylistResponse.Description,
		Owner:       spotifyPlaylistResponse.Owner.DisplayName,
		ExternalUrl: spotifyPlaylistResponse.ExternalURL.Spotify,
		Songs:       make([]domain.Song, len(spotifyPlaylistResponse.Tracks.Items)),
	}
	for i, song := range spotifyPlaylistResponse.Tracks.Items {
		playlist.Songs[i] = domain.Song{
			Title:  song.Track.Name,
			Artist: song.Track.Artists[0].Name,
			Album:  song.Track.Album.Name,
		}
	}
	return playlist, nil
}
func (s ServiceSpotify) decodeSpotifyRespTrackByMetadata(body *io.ReadCloser) (*domain.Song, error) {

	var track spotifySongByMetadata
	if err := json.NewDecoder(*body).Decode(&track); err != nil {
		return nil, errs.Errorf("can't decode response: %v", err)
	}
	if len(track.Tracks.Items) == 0 || track.Tracks.Items[0].Name == "" || track.Tracks.Items[0].Artists[0].Name == "" || track.Tracks.Items[0].ExternalURL.Spotify == "" {
		return nil, errs.Errorf(errors.ErrInvalidParamsSpotify)
	}
	return &domain.Song{
		Title:  track.Tracks.Items[0].Name,
		Artist: track.Tracks.Items[0].Artists[0].Name,
		Album:  track.Tracks.Items[0].Album.Name,
		Link:   track.Tracks.Items[0].ExternalURL.Spotify,
		Id:     track.Tracks.Items[0].ID,
	}, nil
}
func (s ServiceSpotify) decodeRespCreateSpotifyPlaylist(body *io.ReadCloser) (string, error) {
	var playlistId PlaylistIdResponse
	if err := json.NewDecoder(*body).Decode(&playlistId); err != nil {
		return "", errs.Errorf("can't decode response: %v", err)
	}
	if playlistId.Id != "" {
		return playlistId.Id, nil
	}
	return "", errs.New("id not found in response")
}

func (s *ServiceSpotify) RequestToken() error {
	req, err := s.buildTokenRequest()
	if err != nil {
		return errs.Errorf("can't build request: %v", err)
	}
	resp, err := s.Client.Do(req)
	if err != nil {
		return errs.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errs.Errorf("can't read response: %v", err)
		}
		s.Logger.ErrorFrmt(string(body))
		return errs.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if err = s.extractToken(resp.Body); err != nil {
		return errs.Errorf("can't extract token: %v", err)
	}

	return nil
}
func (s *ServiceSpotify) buildTokenRequest() (*http.Request, error) {
	tokenData := url.Values{}
	tokenData.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(tokenData.Encode()))
	if err != nil {
		return nil, errs.Errorf("failed to create request: %w", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(s.ClientId + ":" + s.ClientSecret))
	req.Header.Add("Authorization", "Basic "+authHeader)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}
func (s *ServiceSpotify) extractToken(body io.ReadCloser) error {
	var tokenResponse spotifyResponseToken
	data, err := io.ReadAll(body)
	if err != nil {
		return errs.Errorf("failed to read response body: %w", err)
	}
	if err = json.Unmarshal(data, &tokenResponse); err != nil {
		return errs.Errorf("failed to unmarshal response: %w", err)
	}
	if tokenResponse.AccessToken != "" {
		s.Token = tokenResponse.AccessToken
	}
	return errs.Errorf("AccessToken not found in response")
}

func (s *ServiceSpotify) FillPlaylist(YouTubePlaylist *domain.Playlist, SpotifyPlaylistId, AuthToken string) (*domain.Playlist, error) {
	songs := make([]domain.Song, len(YouTubePlaylist.Songs))
	for i, track := range YouTubePlaylist.Songs {
		Song, err := s.GetSpotifyTrackByMetadata(domain.MetaData{Title: track.Title, Artist: track.Artist})
		if err != nil {
			return nil, errs.Errorf("can't get Song by metadata: %v", err)
		}
		req, _ := http.NewRequest("POST", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", SpotifyPlaylistId), bytes.NewBuffer([]byte(fmt.Sprintf(`{"uris": ["spotify:track:%s"]}`, Song.Id))))
		req.Header.Add("Authorization", "Bearer "+AuthToken)
		req.Header.Add("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, errs.Errorf("failed to execute request: %v", err)
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			s.Logger.ErrorFrmt(string(body))
			return nil, errs.Errorf(resp.Status)
		}

		songs = append(songs, *Song)
		if i == len(YouTubePlaylist.Songs) {
			break
		}
	}
	SpotifyPlaylist := &domain.Playlist{
		Title:       YouTubePlaylist.Title,
		Owner:       YouTubePlaylist.Owner,
		Description: "Playlist created by tg_bot MusicApp",
		ExternalUrl: fmt.Sprintf("https://open.spotify.com/playlist/%s", SpotifyPlaylistId),
		Songs:       songs,
	}

	return SpotifyPlaylist, nil
}
