package Spotify

import (
	"MusicApp/internal/config"
	"MusicApp/internal/domain"
	"errors"
	logger "github.com/skwizi4/lib/logs"
	"net/http"
	"net/url"
)

func NewSpotifyService(cfg config.Config) ServiceSpotify {
	return ServiceSpotify{
		BaseUrl:      BaseUrl,
		ClientId:     cfg.SpotifyCfg.ClientID,
		ClientSecret: cfg.SpotifyCfg.ClientSecret,
		Logger:       logger.InitLogger(),
	}
}

func (s ServiceSpotify) GetSpotifyTrackById(link string) (domain.Song, error) {
	isTrack, id, err := ParseSpotifyIDFromURL(link)
	if isTrack == "playlist" {
		return domain.Song{}, errors.New("invalid link, its link of playlist")
	}
	if err != nil {
		return domain.Song{}, err
	}
	if err = s.RequestToken(); err != nil {
		return domain.Song{}, err
	}
	req, err := s.CreateRequest(http.MethodGet, "/v1/tracks/"+id)
	if err != nil {
		return domain.Song{}, err
	}
	resp, err := s.doRequest(req)
	if err != nil {
		return domain.Song{}, err
	}

	Track, err := s.decodeRespTrackById(resp)

	if err != nil {
		return domain.Song{}, err
	}

	return Track, nil

}
func (s ServiceSpotify) GetSpotifyPlaylistById(link string) (domain.Playlist, error) {
	isPlaylist, id, err := ParseSpotifyIDFromURL(link)
	if isPlaylist == "track" {
		return domain.Playlist{}, errors.New("invalid link, its link of track")
	}
	if err != nil {
		return domain.Playlist{}, err
	}
	req, err := s.CreateRequest(http.MethodGet, "/v1/playlists/"+id)
	if err != nil {
		return domain.Playlist{}, err
	}
	resp, err := s.doRequest(req)
	if err != nil {
		return domain.Playlist{}, err
	}
	Playlist, err := s.decodeRespPlaylistId(resp)
	if err != nil {
		return domain.Playlist{}, err
	}
	return Playlist, nil

}
func (s ServiceSpotify) GetSpotifyTrackByName(data domain.MetaData) (domain.Song, error) {
	req, err := s.CreateRequest(http.MethodGet, "/v1/search?q=track:"+url.QueryEscape(data.Title)+"+artist:"+url.QueryEscape(data.Artist)+"&type=track&limit=10&offset=5")
	if err != nil {
		return domain.Song{}, err
	}
	resp, err := s.doRequest(req)
	if err != nil {
		return domain.Song{}, err
	}
	Song, err := s.decodeRespTrackByName(resp)
	if err != nil {
		return domain.Song{}, err
	}

	return Song, nil

}
