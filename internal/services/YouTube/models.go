package YouTube

import (
	"MusicApp/internal/config"
	logger "github.com/skwizi4/lib/logs"
	"net/http"
)

const (
	BaseApiUrl               = "https://www.googleapis.com/youtube/v3/"
	youtubeTrackDomen        = "https://www.youtube.com/watch?v="
	scope                    = "https://www.googleapis.com/auth/youtube"
	creatingPlaylistEndpoint = "playlists?part=snippet,status"
	fillingPlaylistEndpoint  = "playlistItems?part=snippet"
	NilAuthToken             = ""
)

type ServiceYouTube struct {
	Key       string
	BaseUrl   string
	logger    logger.GoLogger
	Token     string
	ClientID  string
	ServerUrl string
	Scope     string
	Client    *http.Client
}

func NewYouTubeService(cfg *config.Config) ServiceYouTube {
	return ServiceYouTube{
		BaseUrl:   BaseApiUrl,
		logger:    logger.InitLogger(),
		Key:       cfg.YoutubeCfg.Key,
		ClientID:  cfg.YoutubeCfg.ClientID,
		ServerUrl: cfg.YoutubeCfg.RedirectUrl,
		Scope:     scope,
		Client:    &http.Client{},
	}
}

type (
	youtubeMediaById struct {
		Items []struct {
			VideoId string `json:"id"`
			Snippet struct {
				ChanelName string `json:"channelTitle"`
				Title      string `json:"title"`
			} `json:"snippet"`
		} `json:"items"`
	}
	youtubePlaylistParamsById struct {
		Items []struct {
			Snippet struct {
				Title        string `json:"title"`
				ChannelTitle string `json:"ChannelTitle"`
			} `json:"snippet"`
		} `json:"items"`
	}
	youtubeResponsePlaylistMediaById struct {
		Items []struct {
			Snippet struct {
				Title        string `json:"title"`
				ChannelTitle string `json:"videoOwnerChannelTitle"`
				ResourceId   struct {
					VideoId string `json:"videoId"`
				} `json:"resourceId"`
			} `json:"snippet"`
		} `json:"items"`
		NextPageToken string `json:"nextPageToken"`
	}
	youtubeMediaByMetadata struct {
		Items []struct {
			Id struct {
				VideoId string `json:"videoId"`
			} `json:"id"`
			Snippet struct {
				ChanelName string `json:"channelTitle"`
				Title      string `json:"title"`
			} `json:"snippet"`
		} `json:"items"`
	}
	youtubePlaylistIdResp struct {
		ID string `json:"id"`
	}
)
