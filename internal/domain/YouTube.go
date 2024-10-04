package domain

import "errors"

//ProcessingYoutubeMediaById

func (p *ProcessingYoutubeMediaById) GetOrCreate(chatID int64) ProcessYouTubeSong {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return (*p)[idx]
	}
	NewProcess := ProcessYouTubeSong{
		chatID: chatID,
		step:   ProcessYouTubeSongStart,
	}
	*p = append(*p, NewProcess)
	return NewProcess
}

func (p *ProcessingYoutubeMediaById) AddSongID(songID string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].songId = songID
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

func (p *ProcessingYoutubeMediaById) AddSong(song Song, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].song = song
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

func (p *ProcessingYoutubeMediaById) UpdateStep(step string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].step = step
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

func (p *ProcessingYoutubeMediaById) Delete(chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		*p = append((*p)[:idx], (*p)[idx+1:]...)
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

//ProcessingYoutubePlaylists

func (p *ProcessingYoutubePlaylists) GetOrCreate(chatID int64) ProcessYoutubePlaylist {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return (*p)[idx]
	}
	NewProcess := ProcessYoutubePlaylist{
		chatID: chatID,
		step:   ProcessYouTubePlaylistStart,
	}
	*p = append(*p, NewProcess)
	return NewProcess
}

func (p *ProcessingYoutubePlaylists) AddSongs(playlist Playlist, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].songs = append((*p)[idx].songs, playlist)
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

func (p *ProcessingYoutubePlaylists) UpdateStep(step string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].step = step
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

func (p *ProcessingYoutubePlaylists) AddTitle(title string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].title = title
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

func (p *ProcessingYoutubePlaylists) Delete(chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		*p = append((*p)[:idx], (*p)[idx+1:]...)
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
