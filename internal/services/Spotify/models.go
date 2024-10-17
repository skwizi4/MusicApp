package Spotify

import (
	"MusicApp/internal/config"
	logger "github.com/skwizi4/lib/logs"
	"net/http"
)

const (
	TokenEndpoint = "https://accounts.spotify.com/api/token"
	BaseApiUrl    = "https://api.spotify.com"
	NilAuthToken  = ""
)

type ServiceSpotify struct {
	BaseUrl      string
	Logger       logger.GoLogger
	ClientId     string
	ClientSecret string
	Token        string
	Client       *http.Client
}

func NewSpotifyService(cfg *config.Config) ServiceSpotify {
	return ServiceSpotify{
		BaseUrl:      BaseApiUrl,
		ClientId:     cfg.SpotifyCfg.ClientID,
		ClientSecret: cfg.SpotifyCfg.ClientSecret,
		Logger:       logger.InitLogger(),
		Client:       &http.Client{},
	}
}

// spotifyTrackById, spotifyPlaylistById, spotifySongByMetadata Структуры для работы с api
type (
	spotifyTrackById struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
		Album struct {
			Name string `json:"name"`
		} `json:"album"`
	}
	spotifyPlaylistById struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Tracks      struct {
			Items []struct {
				Track struct {
					ID      string `json:"id"`
					Name    string `json:"name"`
					Artists []struct {
						Name string `json:"name"`
					} `json:"artists"`
					Album struct {
						Name string `json:"name"`
					} `json:"album"`
					DurationMs  int `json:"duration_ms"`
					ExternalURL struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
				} `json:"track"`
			} `json:"items"`
			Total int `json:"total"`
		} `json:"tracks"`
		Owner struct {
			DisplayName string `json:"display_name"`
			ID          string `json:"id"`
		} `json:"owner"`
		ExternalURL struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
	}
	spotifySongByMetadata struct {
		Tracks struct {
			Items []struct {
				ID      string `json:"id"`
				Name    string `json:"name"`
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
				Album struct {
					Name string `json:"name"`
				} `json:"album"`
				ExternalURL struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
			} `json:"items"`
		} `json:"tracks"`
	}
	PlaylistCreateRequest struct {
		Name string `json:"name"`
	}
	spotifyResponseToken struct {
		AccessToken string `json:"access_token"`
	}
	PlaylistIdResponse struct {
		Id string `json:"id"`
	}
)
