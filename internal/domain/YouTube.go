package domain

import "errors"

//ProcessingSpotifySongByYoutubeMediaId

func (p *ProcessingSpotifySongByYoutubeMediaId) GetOrCreate(chatID int64) ProcessYouTubeSong {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return (*p)[idx]
	}
	NewProcess := ProcessYouTubeSong{
		ChatID: chatID,
		Step:   ProcessSpotifySongByYouTubeMediaStart,
	}
	*p = append(*p, NewProcess)
	return NewProcess
}

func (p *ProcessingSpotifySongByYoutubeMediaId) IfExist(chatID int64) bool {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return true
	}
	return false
}

func (p *ProcessingSpotifySongByYoutubeMediaId) AddSong(song Song, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Song = song
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

func (p *ProcessingSpotifySongByYoutubeMediaId) UpdateStep(step string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Step = step
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

func (p *ProcessingSpotifySongByYoutubeMediaId) Delete(chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		*p = append((*p)[:idx], (*p)[idx+1:]...)
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

//ProcessingYoutubePlaylists

func (p *ProcessingFillYoutubePlaylists) GetOrCreate(chatID int64) ProcessFillYoutubePlaylist {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return (*p)[idx]
	}
	NewProcess := ProcessFillYoutubePlaylist{
		ChatID: chatID,
		Step:   ProcessFillYouTubePlaylistStart,
	}
	*p = append(*p, NewProcess)
	return NewProcess
}

func (p *ProcessingFillYoutubePlaylists) AddSongs(playlist Playlist, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Songs = playlist
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

func (p *ProcessingFillYoutubePlaylists) UpdateStep(step string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Step = step
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

func (p *ProcessingFillYoutubePlaylists) AddTitle(title string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Songs.Title = title
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

func (p *ProcessingFillYoutubePlaylists) AddAuthToken(AuthToken string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].AuthToken = AuthToken
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessingFillYoutubePlaylists) IfExist(chatID int64) bool {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return true
	}

	return false
}
func (p *ProcessingFillYoutubePlaylists) Delete(chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		*p = append((*p)[:idx], (*p)[idx+1:]...)
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
