package domain

import (
	"errors"
)

//ProcessSpotifySongs

func (p *ProcessSpotifySongs) GetOrCreate(chatID int64) ProcessSpotifySong {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return (*p)[idx]
	}
	return ProcessSpotifySong{
		chatID: chatID,
		step:   ProcessSpotifySongStart,
	}
}
func (p *ProcessSpotifySongs) AddSongID(songID string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].songId = songID
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessSpotifySongs) AddSong(song Song, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].song = song
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessSpotifySongs) UpdateStep(step string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].step = step
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessSpotifySongs) Delete(chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		*p = append((*p)[:idx], (*p)[idx+1:]...)
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

//ProcessSpotifyPlaylist
