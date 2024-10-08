package domain

import (
	"errors"
)

//ProcessingYoutubeMediaBySpotifySongID

func (p *ProcessingYoutubeMediaBySpotifySongID) GetOrCreate(chatID int64) ProcessSpotifySong {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return (*p)[idx]
	}
	NewProcess := ProcessSpotifySong{
		ChatID: chatID,
		Step:   ProcessSpotifySongByIdStart,
	}
	*p = append(*p, NewProcess)
	return NewProcess
}
func (p *ProcessingYoutubeMediaBySpotifySongID) IfExist(chatID int64) bool {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return true
	}
	return false

}

func (p *ProcessingYoutubeMediaBySpotifySongID) UpdateStep(step string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Step = step
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessingYoutubeMediaBySpotifySongID) Delete(chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		*p = append((*p)[:idx], (*p)[idx+1:]...)
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

//ProcessSpotifyByMetadata

func (p *ProcessingFindSongByMetadata) GetOrCreate(chatID int64) ProcessSpotifySong {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return (*p)[idx]
	}
	NewProcess := ProcessSpotifySong{
		ChatID: chatID,
		Step:   ProcessSpotifySongByMetadataStart,
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
