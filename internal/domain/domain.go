package domain

import (
	"reflect"
)

// todo - Заполнить структуры и написать к ним crud - методы

const (
	ProcessSpotifySongStart = "ProcessSpotifySongStart"
	ProcessSpotifySongEnd   = "ProcessSpotifySongEnd"

	ProcessSpotifyPlaylistStart = "ProcessSpotifyPlaylistStart"
	ProcessSpotifyPlaylistEnd   = "ProcessSpotifyPlaylistEnd"

	ProcessYouTubeSongStart = "ProcessYouTubeSongStart"
	ProcessYouTubeSongEnd   = "ProcessYouTubeSongEnd"

	ProcessYouTubePlaylistStart = "ProcessYouTubePlaylistStart"
	ProcessYouTubePlaylistEnd   = "ProcessYouTubePlaylistEnd"

	ProcessFindSongStart = "ProcessFindSongStart"
	ProcessFindSongEnd   = "ProcessFindSongEnd"

	ErrChatIDNotFound = "chatID not found"
)

func FindUserIndex(value interface{}, chatID int64) int {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return -1
	}

	for i := 0; i < v.Len(); i++ {
		user := v.Index(i)
		chatIDField := user.FieldByName("ChatID")
		if !chatIDField.IsValid() || chatIDField.Kind() != reflect.Int64 {
			continue
		}
		if chatIDField.Int() == chatID {
			return i
		}
	}

	return -1
}

// Song - Structure of song(Spotify), media (YouTube)
type Song struct {
	title    string
	artist   string
	album    string
	genre    string
	duration string
}

type MetaData struct {
	title  string
	artist string
}
type Playlist struct {
	songs []Song
}

// ProcessSpotifySong - Structure of handler "spotifyHandler";
// + process for many users and some methods
type ProcessSpotifySong struct {
	songId string
	song   Song
	chatID int64
	step   string
}
type ProcessSpotifySongs []ProcessSpotifySong

type ProcessSpotifyPlaylist struct {
	title string
	songs []Playlist
}
type ProcessYouTubeSong struct {
	songId string
	song   Song
}
type ProcessYoutubePlaylist struct {
	title string
	songs []Playlist
}
type ProcessFindSong struct {
	metadata MetaData
}
