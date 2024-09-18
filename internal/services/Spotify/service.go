package Spotify

import (
	"MusicApp/internal/domain"
	logger "github.com/skwizi4/lib/logs"
	"net/http"
	"net/url"
)

func NewSpotifyService(token string) ServiceSpotify {
	return ServiceSpotify{
		BaseUrl: BaseUrl,
		ApiKey:  token,
		Logger:  logger.InitLogger(),
	}
}
func (s ServiceSpotify) TrackById(link string) (domain.Song, error) {
	id, err := ParseTrackIDFromURL(link)
	if err != nil {
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
func (s ServiceSpotify) PlaylistById(link string) (domain.Playlist, error) {
	id, err := ParsePlaylistIDFromURL(link)
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
func (s ServiceSpotify) TrackByName(trackName, artist string) (domain.Song, error) {
	req, err := s.CreateRequest(http.MethodGet, "/v1/search?q=track:"+url.QueryEscape(trackName)+"+artist:"+url.QueryEscape(artist)+"&type=track&limit=10&offset=5")
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
