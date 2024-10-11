package domain

import (
	"errors"
	"reflect"
)

// todo - Заполнить структуры и написать к ним crud - методы

const (
	ProcessYoutubeMediaBySpotifySongLinkStart = "ProcessYoutubeMediaBySpotifySongLinkStart"
	ProcessYoutubeMediaBySpotifySongLinkEnd   = "ProcessYoutubeMediaBySpotifySongLinkEnd"

	ProcessSpotifySongByYouTubeMediaLinkStart = "ProcessSpotifySongByYouTubeMediaLinkStart"
	ProcessSpotifySongByYouTubeMediaLinkEnd   = "ProcessSpotifySongByYouTubeMediaLinkEnd"

	ProcessSongByMetadataStart  = "ProcessSongByMetadataStart"
	ProcessSongByMetadataTitle  = "ProcessSongByMetadataTitle"
	ProcessSongByMetadataArtist = "ProcessSongByMetadataArtist"
	ProcessSongByMetadataEnd    = "ProcessSongByMetadataEnd"

	ProcessSpotifyPlaylistStart = "ProcessSpotifyPlaylistStart"
	ProcessSpotifyPlaylistEnd   = "ProcessSpotifyPlaylistEnd"

	ProcessFillYouTubePlaylistStart        = "ProcessFillYouTubePlaylistStart"
	ProcessFillYouTubePlaylistSendAuthLink = "ProcessFillYouTubePlaylistSendAuthLink"
	ProcessFillYouTubePlaylistEnd          = "ProcessFillYouTubePlaylistEnd"

	ErrChatIDNotFound = "ChatID not found"
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

// Song - Structure of song(Youtube), media (Spotify)
type Song struct {
	Title  string
	Artist string
	Album  string
	Genre  string
	Link   string
	Id     string
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

// Youtube domain

// ProcessSong - Structure of handler "spotifyHandler";
// + process for many users and some methods
// todo - check songId
type ProcessSong struct {
	SongId        string
	Song          Song
	ChatID        int64
	Step          string
	IsGetMetadata bool
}
type ProcessingYoutubeMediaBySpotifySongLink []ProcessSong

type ProcessCreateAndFillSpotifyPlaylist struct {
	ChatID   int64
	Step     string
	Title    string
	Playlist Playlist
}
type ProcessingCreateAndFillSpotifyPlaylists []ProcessCreateAndFillSpotifyPlaylist

// Spotify domain

type ProcessYouTubeSong struct {
	ChatID int64
	Song   Song
	Step   string
}
type ProcessingSpotifySongByYoutubeMediaLink []ProcessYouTubeSong

type ProcessCreateAndFillYoutubePlaylist struct {
	Playlist  Playlist
	ChatID    int64
	Step      string
	AuthToken string
}
type ProcessingCreateAndFillYoutubePlaylists []ProcessCreateAndFillYoutubePlaylist

type ProcessFindSong struct {
	chatID   int64
	metadata MetaData
	step     string
}

// ProcessingFindSongs both handlers
/* ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
/* ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
/* ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
/* ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

type ProcessingFindSongs []ProcessFindSong

type ProcessingFindSongByMetadata []ProcessSong

//ProcessSpotifyByMetadata

func (p *ProcessingFindSongByMetadata) GetOrCreate(chatID int64) ProcessSong {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return (*p)[idx]
	}
	NewProcess := ProcessSong{
		ChatID: chatID,
		Step:   ProcessSongByMetadataStart,
	}
	*p = append(*p, NewProcess)
	return NewProcess
}
func (p *ProcessingFindSongByMetadata) IfExist(chatID int64) bool {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return true
	}
	return false

}

func (p *ProcessingFindSongByMetadata) UpdateStep(step string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Step = step
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessingFindSongByMetadata) Delete(chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		*p = append((*p)[:idx], (*p)[idx+1:]...)
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessingFindSongByMetadata) AddTitle(chatID int64, title string) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Song.Title = title
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessingFindSongByMetadata) AddArtist(chatID int64, artist string) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Song.Artist = artist
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessingFindSongByMetadata) ChangeIsGetMetadata(chatID int64, value bool) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].IsGetMetadata = value
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
