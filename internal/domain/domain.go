package domain

import (
	"reflect"
)

// todo - Заполнить структуры и написать к ним crud - методы

const (
	ProcessSpotifySongByIdStart = "ProcessSpotifySongByIdStart"
	ProcessSpotifySongByIdEnd   = "ProcessSpotifySongByIdEnd"

	ProcessSpotifySongByMetadataStart  = "ProcessSpotifySongByMetadataStart"
	ProcessSpotifySongByMetadataTitle  = "ProcessSpotifySongByMetadataTitle"
	ProcessSpotifySongByMetadataArtist = "ProcessSpotifySongByMetadataArtist"
	ProcessSpotifySongByMetadataEnd    = "ProcessSpotifySongByMetadataEnd"

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
	Title  string
	Artist string
	Album  string
	Genre  string
	Link   string
}

type MetaData struct {
	Title  string
	Artist string
}
type Playlist struct {
	Songs         []Song
	Title         string
	Owner         string
	Description   string
	ExternalUrl   string
	NextPageToken string
}

// ProcessSpotifySong - Structure of handler "spotifyHandler";
// + process for many users and some methods
type ProcessSpotifySong struct {
	SongId string
	Song   Song
	ChatID int64
	Step   string
}
type ProcessingSpotifySongsByID []ProcessSpotifySong

type ProcessSpotifyPlaylist struct {
	chatID   int64
	step     string
	title    string
	playlist Playlist
}
type ProcessingSpotifyPlaylists []ProcessSpotifyPlaylist

type ProcessYouTubeSong struct {
	chatID int64
	songId string
	song   Song
	step   string
}
type ProcessingYoutubeSongs []ProcessYouTubeSong

type ProcessYoutubePlaylist struct {
	title  string
	songs  []Playlist
	chatID int64
	step   string
}
type ProcessingYoutubePlaylists []ProcessYoutubePlaylist

type ProcessFindSong struct {
	chatID   int64
	metadata MetaData
	step     string
}
type ProcessingFindSongs []ProcessFindSong

//todo - refactor

type ProcessingFindSongByMetadata []ProcessSpotifySong
