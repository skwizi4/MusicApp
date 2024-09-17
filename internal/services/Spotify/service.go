package Spotify

import (
	"MusicApp/internal/domain"
	logger "github.com/skwizi4/lib/logs"
	"net/http"
)

func NewSpotifyService(BaseUrl, token string) ServiceSpotify {
	return ServiceSpotify{
		BaseUrl: BaseUrl,
		ApiKey:  token,
		Logger:  logger.InitLogger(),
	}
}
func (s ServiceSpotify) Track(link string) (domain.Song, error) {
	id, err := ParseIDFromURL(link)
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
	Track, err := s.DecodeRespTrack(resp)
	if err != nil {
		return domain.Song{}, err
	}
	return Track, nil

}
func (s ServiceSpotify) GetResponse() {}
func (s ServiceSpotify) PostRequest() {

}
func (s ServiceSpotify) PostResponse() {}
