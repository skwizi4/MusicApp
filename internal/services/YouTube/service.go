package YouTube

import (
	"MusicApp/internal/config"
	"MusicApp/internal/domain"
	"fmt"
	logger "github.com/skwizi4/lib/logs"
	"net/http"
)

//todo Refactor

func NewYouTubeService(cfg config.Config) ServiceYouTube {
	return ServiceYouTube{
		BaseUrl:   BaseUrl,
		logger:    logger.InitLogger(),
		Key:       cfg.YoutubeCfg.Key,
		ClientID:  cfg.YoutubeCfg.ClientID,
		ServerUrl: ServerUrl,
		Scope:     scope,
	}
}

// GetYoutubeMediaByID Tested - Tested(OK)
func (y ServiceYouTube) GetYoutubeMediaByID(link string) (*domain.Song, error) {

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

// GetYoutubePlaylistByID - Tested(OK)
func (y ServiceYouTube) GetYoutubePlaylistByID(link string) (*domain.Playlist, error) {
	var playlist = &domain.Playlist{}
	id, err := GetID(link)
	if err != nil {
		return playlist, err
	}
	// Fill playlist params
	endpoint, err := y.CreateEndpointYoutubePlaylistParams(id)
	if err != nil {
		return playlist, err
	}
	resp, err := y.createAndExecuteRequest(http.MethodGet, endpoint)
	if err != nil {
		return playlist, err
	}

	err = FillPlaylistParams(resp, playlist)
	if err != nil {
		return playlist, err
	}
	fmt.Println(playlist)

	// fill Playlist media
	var NextPageToken string
	for {
		endpoint, err = y.CreateEndpointYoutubePlaylistSongs(id, NextPageToken)
		if err != nil {
			return playlist, err
		}
		resp, err = y.createAndExecuteRequest(http.MethodGet, endpoint)
		if err != nil {
			return playlist, err
		}
		playlist, err = FillPlaylistSongs(resp, playlist)
		if err != nil {
			return playlist, err
		}
		if playlist.NextPageToken == "" {
			return playlist, nil
		}
		NextPageToken = playlist.NextPageToken
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
func (y ServiceYouTube) CreateYoutubePlaylist(SpotifyPlaylist domain.Playlist, token string) (*domain.Playlist, error) {
	id, err := y.CreatePlaylist(token, SpotifyPlaylist.Title)
	if err != nil {
		return nil, err
	}
	YoutubePlaylist, err := y.FillYoutubePlaylist(token, id, SpotifyPlaylist.Songs)
	if err != nil {
		return nil, err
	}
	YoutubePlaylist.Owner = SpotifyPlaylist.Owner
	YoutubePlaylist.Title = SpotifyPlaylist.Title
	YoutubePlaylist.Description = "Playlist Created by tg-bot"
	YoutubePlaylist.ExternalUrl = fmt.Sprintf("https://youtube.com/playlist?list=%s", id)

	return YoutubePlaylist, nil
}
