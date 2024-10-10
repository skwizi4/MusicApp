package domain

import (
	"errors"
)

//processingSpotifySongByYoutubeMediaLink

func (p *ProcessingSpotifySongByYoutubeMediaLink) GetOrCreate(chatID int64) ProcessYouTubeSong {

	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return (*p)[idx]
	}
	NewProcess := ProcessYouTubeSong{
		ChatID: chatID,
		Step:   ProcessSpotifySongByYouTubeMediaLinkStart,
	}
	*p = append(*p, NewProcess)
	return NewProcess
}

func (p *ProcessingSpotifySongByYoutubeMediaLink) IfExist(chatID int64) bool {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return true
	}
	return false
}

func (p *ProcessingSpotifySongByYoutubeMediaLink) AddSong(song Song, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Song = song
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

func (p *ProcessingSpotifySongByYoutubeMediaLink) UpdateStep(step string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Step = step
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

func (p *ProcessingSpotifySongByYoutubeMediaLink) Delete(chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		*p = append((*p)[:idx], (*p)[idx+1:]...)
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

//ProcessingYoutubePlaylists

//ProcessSpotifyPlaylist

func (p *ProcessingSpotifyPlaylists) GetOrCreate(chatID int64) ProcessSpotifyPlaylist {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return (*p)[idx]
	}
	NewProcess := ProcessSpotifyPlaylist{
		chatID: chatID,
		step:   ProcessSpotifyPlaylistStart,
	}
	*p = append(*p, NewProcess)
	return NewProcess

}
func (p *ProcessingSpotifyPlaylists) AddSongs(playlist Playlist, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].playlist = playlist
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessingSpotifyPlaylists) UpdateStep(step string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].step = step
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessingSpotifyPlaylists) AddTitle(title string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].title = title
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessingSpotifyPlaylists) Delete(chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		*p = append((*p)[:idx], (*p)[idx+1:]...)
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
