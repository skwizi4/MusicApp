package Spotify

import (
	logger "github.com/skwizi4/lib/logs"
)

type ServiceSpotify struct {
	ApiKey  string
	BaseUrl string
	Logger  logger.GoLogger
}

// SpotifyTrack Структура для ответа от Spotify API на запрос трека
type spotifyTrack struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	Album struct {
		Name string `json:"name"`
	} `json:"album"`
}
