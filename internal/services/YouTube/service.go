package YouTube

import (
	"MusicApp/internal/config"
	"MusicApp/internal/domain"
	"errors"
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

// GetYoutubeMediaByID Tested
func (y ServiceYouTube) GetYoutubeMediaByID(link string) (*domain.Song, error) {
	isTrack, id, err := ParseYouTubeIDFromURL(link)
	if isTrack == "playlist" {
		return &domain.Song{}, errors.New("invalid link, its playlist link")
	}
	if err != nil {
		return &domain.Song{}, err
	}
	endpoint, err := y.CreateEndpointYoutubeMedia(id)
	if err != nil {
		return &domain.Song{}, err
	}
	req, err := y.CreateRequest(http.MethodGet, endpoint)
	if err != nil {
		return &domain.Song{}, err
	}
	resp, err := y.DoRequest(req)
	if err != nil {
		return &domain.Song{}, err
	}
	song, err := DecodeRespMediaById(resp)
	if err != nil {
		return &domain.Song{}, err
	}
	song.Link = link
	return song, nil
}

// GetYoutubePlaylistByID todo write tests + fix bugs
func (y ServiceYouTube) GetYoutubePlaylistByID(link string) (*domain.Playlist, error) {
	isPlaylist, id, err := ParseYouTubeIDFromURL(link)
	if isPlaylist == "track" {
		return &domain.Playlist{}, errors.New("invalid link, its track link")
	}
	if err != nil {
		return &domain.Playlist{}, err
	}
	endpoint, err := y.CreateEndpointYoutubePlaylist(id)
	if err != nil {
		return &domain.Playlist{}, err
	}
	req, err := y.CreateRequest(http.MethodGet, endpoint)
	if err != nil {
		return &domain.Playlist{}, err
	}
	resp, err := y.DoRequest(req)
	if err != nil {
		return &domain.Playlist{}, err
	}
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return &domain.Playlist{}, err
	//}
	//fmt.Println(string(body))
	playlist, err := DecodeRespPlaylistById(resp)
	if err != nil {
		return &domain.Playlist{}, err
	}
	return playlist, nil
}

//todo Complete finding media by metadata for playlist/ future fitch + write tests

func (y ServiceYouTube) GetYoutubeMediaByMetadata(data domain.MetaData) (*domain.Song, error) {
	return &domain.Song{}, nil
}

// todo Complete creating playlist and returning playlist structure

func (y ServiceYouTube) FillYoutubePlaylist(list domain.Playlist) (*domain.Playlist, error) {
	return &domain.Playlist{
		ExternalUrl: "",
	}, nil
}
