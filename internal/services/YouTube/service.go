package YouTube

import (
	"MusicApp/internal/config"
	"MusicApp/internal/domain"
	"errors"
	"fmt"
	logger "github.com/skwizi4/lib/logs"
	"net/http"
)

func New(cfg config.Config) ServiceYouTube {
	return ServiceYouTube{
		BaseUrl: BaseUrl,
		logger:  logger.InitLogger(),
		Key:     cfg.YoutubeCfg.Key,
	}
}

// GetYoutubeMediaByID Tested - OK
func (y ServiceYouTube) GetYoutubeMediaByID(link string) (*domain.Song, error) {
	song := &domain.Song{}
	isTrack, id, err := ParseYouTubeIDFromURL(link)
	if isTrack == "playlist" {
		return song, errors.New("invalid link, its playlist link")
	}
	if err != nil {
		return song, err
	}
	endpoint, err := y.CreateEndpointYoutubeMediaById(id)
	if err != nil {
		return song, err
	}
	resp, err := y.createAndExecuteRequest(http.MethodGet, endpoint)
	if err != nil {
		return song, err
	}
	song, err = DecodeRespMediaById(resp)
	if err != nil {
		return song, err
	}

	return song, nil
}

// GetYoutubePlaylistByID - OK
func (y ServiceYouTube) GetYoutubePlaylistByID(link string) (*domain.Playlist, error) {
	var playlist = &domain.Playlist{}
	isPlaylist, id, err := ParseYouTubeIDFromURL(link)
	if isPlaylist == "track" || err != nil {
		return playlist, errors.New("invalid link, its track link")
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

	// fill Playlist media
	endpoint, err = y.CreateEndpointYoutubePlaylistSongs(id)
	if err != nil {
		return playlist, err
	}
	resp, err = y.createAndExecuteRequest(http.MethodGet, endpoint)
	if err != nil {
		return playlist, err
	}
	playlist, err = FillPlaylist(resp, playlist)
	if err != nil {
		return playlist, err
	}
	return playlist, nil
}

// GetYoutubeMediaByMetadata - OK
func (y ServiceYouTube) GetYoutubeMediaByMetadata(data domain.MetaData) (*domain.Song, error) {
	song := &domain.Song{}
	endpoint, err := y.CreateEndpointYoutubeMediaByMetadata(data)
	if err != nil {
		return song, err
	}
	resp, err := y.createAndExecuteRequest(http.MethodGet, endpoint)
	if err != nil {
		return song, err
	}
	song, err = DecodeRespMediaByMetadata(resp)
	if err != nil {
		return song, err
	}
	return song, nil
}

// todo Test last method of service in prod

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
