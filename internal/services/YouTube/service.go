package YouTube

import (
	"MusicApp/internal/domain"
	"net/http"
)

// GetYoutubeMediaByLink  - Tested(OK)
func (y ServiceYouTube) GetYoutubeMediaByLink(link string) (*domain.Song, error) {
	id, err := GetID(link)
	if err != nil {
		return nil, err
	}
	endpoint, err := y.CreateEndpointYoutubeMediaById(id)
	if err != nil {
		return nil, err
	}
	resp, err := y.createAndExecuteRequest(http.MethodGet, endpoint, NilAuthToken, nil)
	if err != nil {
		return nil, err
	}

	return DecodeRespMediaById(resp)
}

// GetYoutubePlaylistDataByLink - Tested(OK)
func (y ServiceYouTube) GetYoutubePlaylistDataByLink(link string) (*domain.Playlist, error) {
	id, err := GetID(link)
	if err != nil {
		return nil, err
	}
	endpoint, err := y.CreateEndpointYoutubePlaylistParams(id)
	if err != nil {
		return nil, err
	}
	resp, err := y.createAndExecuteRequest(http.MethodGet, endpoint, NilAuthToken, nil)
	if err != nil {
		return nil, err
	}

	playlist, err := y.FillPlaylistParams(resp)
	if err != nil {
		return nil, err
	}

	for {
		endpoint, err = y.CreateEndpointYoutubePlaylistSongs(id, playlist.NextPageToken)
		if err != nil {
			return nil, err
		}

		resp, err = y.createAndExecuteRequest(http.MethodGet, endpoint, NilAuthToken, nil)
		if err != nil {
			return nil, err
		}
		if err = y.FillPlaylistSongs(resp, playlist); err != nil {
			return nil, err
		}
		if playlist.NextPageToken == "" {
			return playlist, nil
		}
	}

}

// GetYoutubeMediaByMetadata - Tested(OK)
func (y ServiceYouTube) GetYoutubeMediaByMetadata(data domain.MetaData) (*domain.Song, error) {
	endpoint, err := y.CreateEndpointYoutubeMediaByMetadata(data)
	if err != nil {
		return nil, err
	}
	resp, err := y.createAndExecuteRequest(http.MethodGet, endpoint, NilAuthToken, nil)
	if err != nil {
		return nil, err
	}
	return DecodeRespMediaByMetadata(resp)
}

// CreateYoutubePlaylist - fix bugs

func (y ServiceYouTube) CreateYoutubePlaylist(Title string, token string) (string, error) {
	endpoint, Body, err := y.MakeEndpointCreateYoutubePlaylist(Title)
	if err != nil {
		return "", err
	}
	body, err := y.createAndExecuteRequest(http.MethodPost, endpoint, token, Body)
	if err != nil {
		return "", err
	}

	return y.DecodeRespCreatePlaylist(body)

}

func (y ServiceYouTube) FillYoutubePlaylist(SpotifyPlaylist *domain.Playlist, YouTubePlaylistId, token string) (*domain.Playlist, error) {
	return y.FillPlaylist(token, YouTubePlaylistId, SpotifyPlaylist)
}
