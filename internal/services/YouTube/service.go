package YouTube

import (
	"MusicApp/internal/domain"
	"errors"
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
	resp, err := y.createAndExecuteRequest(http.MethodGet, endpoint)
	if err != nil {
		return nil, err
	}
	song, err := DecodeRespMediaById(resp)
	if err != nil {
		return nil, err
	}

	return song, nil
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
	resp, err := y.createAndExecuteRequest(http.MethodGet, endpoint)
	if err != nil {
		return nil, err
	}

	playlist, err := FillPlaylistParams(resp)
	if err != nil {
		return nil, err
	}

	for {
		endpoint, err = y.CreateEndpointYoutubePlaylistSongs(id, playlist.NextPageToken)
		if err != nil {
			return nil, err
		}

		resp, err = y.createAndExecuteRequest(http.MethodGet, endpoint)
		if err != nil {
			return nil, err
		}
		if err = FillPlaylistSongs(resp, playlist); err != nil {
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

	resp, err := y.createAndExecuteRequest(http.MethodGet, endpoint)
	if err != nil {
		return nil, err
	}

	song, err := DecodeRespMediaByMetadata(resp)
	if err != nil {
		return nil, err
	}
	return song, nil
}

// CreateYoutubePlaylist - Tested(OK)
func (y ServiceYouTube) CreateYoutubePlaylist(Title string, token string) (string, error) {
	id, err := y.CreatePlaylist(token, Title)
	if err != nil {
		return "", err
	}
	return id, nil

}

func (y ServiceYouTube) FillYoutubePlaylist(SpotifyPlaylist *domain.Playlist, YouTubePlaylistId, token string) (*domain.Playlist, error) {
	YoutubePlaylist, err := y.WriteInPlaylist(token, YouTubePlaylistId, SpotifyPlaylist)
	if err != nil {
		return nil, err
	}
	if YoutubePlaylist == nil {
		return nil, errors.New("YouTubePlaylist is nil")
	}

	return YoutubePlaylist, nil
}
