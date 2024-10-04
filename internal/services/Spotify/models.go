package Spotify

import (
	logger "github.com/skwizi4/lib/logs"
)

const BaseUrl = "https://api.spotify.com"

type ServiceSpotify struct {
	BaseUrl              string
	Logger               logger.GoLogger
	ClientId             string
	ClientSecret         string
	Token                string
	SpotifyResponseToken struct {
		AccessToken string `json:"access_token"`
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
				DurationMs  int `json:"duration_ms"`
				ExternalURL struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
			} `json:"items"`
			Total int `json:"total"`
		} `json:"tracks"`
	}
)
