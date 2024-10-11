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

func (p *ProcessingCreateAndFillSpotifyPlaylists) GetOrCreate(chatID int64) ProcessCreateAndFillSpotifyPlaylist {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return (*p)[idx]
	}
	NewProcess := ProcessCreateAndFillSpotifyPlaylist{
		ChatID: chatID,
		Step:   ProcessSpotifyPlaylistStart,
	}
	*p = append(*p, NewProcess)
	return NewProcess

}
func (p *ProcessingCreateAndFillSpotifyPlaylists) AddSongs(playlist Playlist, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Playlist = playlist
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessingCreateAndFillSpotifyPlaylists) UpdateStep(step string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Step = step
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessingCreateAndFillSpotifyPlaylists) AddTitle(title string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Title = title
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessingCreateAndFillSpotifyPlaylists) Delete(chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		*p = append((*p)[:idx], (*p)[idx+1:]...)
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
