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

func (y ServiceYouTube) GetYoutubeMediaByID(link string) (*domain.Song, error) {
	isTrack, id, err := ParseYouTubeIDFromURL(link)
	if isTrack == "playlist" {
		return &domain.Song{}, errors.New("invalid link, its playlist link")
	}
	if err != nil {
		return &domain.Song{}, err
	}
	endpoint, err := y.CreateEndpoint(id)
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
	var MediaByIdStruct youtubeMediaById
	err = DecodedBody(resp, &MediaByIdStruct)
	if err != nil {
		return &domain.Song{}, err
	}

	return &domain.Song{
		Title:  MediaByIdStruct.Items[0].Snippet.Title,
		Artist: MediaByIdStruct.Items[0].Snippet.ChanelName,
	}, nil
}
func (y ServiceYouTube) GetYoutubePlaylistByID(link string) (*domain.Playlist, error) {
	isPlaylist, id, err := ParseYouTubeIDFromURL(link)
	if isPlaylist == "playlist" {
		return &domain.Playlist{}, errors.New("invalid link, its track link")
	}
	if err != nil {
		return &domain.Playlist{}, err
	}
	endpoint, err := y.CreateEndpoint(id)
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
	var MediaByIdStruct youtubePlaylistById
	err = DecodedBody(resp, &MediaByIdStruct)
	if err != nil {
		return &domain.Playlist{}, err
	}
	playlist := &domain.Playlist{
		Title: "", /*todo Create structure for youtube playlist */
	}
	//todo create filling playlist
	return playlist, nil
}

//todo Complete finding media by metadata for playlist/ future fitch

func (y ServiceYouTube) GetMediaByMetadata(data domain.MetaData) (*domain.Song, error) {
	return &domain.Song{}, nil
}

// todo Complete creating playlist and returning playlist structure

func (y ServiceYouTube) FillYoutubePlaylist(list domain.Playlist) (*domain.Playlist, error) {
	return &domain.Playlist{
		ExternalUrl: "",
	}, nil
}
