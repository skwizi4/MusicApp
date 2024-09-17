package domain

import (
	"errors"
)

//ProcessingSpotifySongs

func (p *ProcessingSpotifySongs) GetOrCreate(chatID int64) ProcessSpotifySong {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return (*p)[idx]
	}
	NewProcess := ProcessSpotifySong{
		ChatID: chatID,
		Step:   ProcessSpotifySongStart,
	}
	*p = append(*p, NewProcess)
	return NewProcess
}
func (p *ProcessingSpotifySongs) AddSongID(songID string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].SongId = songID
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessingSpotifySongs) AddSong(song Song, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Song = song
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessingSpotifySongs) UpdateStep(step string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Step = step
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessingSpotifySongs) Delete(chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		*p = append((*p)[:idx], (*p)[idx+1:]...)
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

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
