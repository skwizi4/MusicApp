package domain

import "errors"

//ProcessingYoutubeMediaBySpotifySongLink

func (p *ProcessingYoutubeMediaBySpotifySongLink) GetOrCreate(chatID int64) ProcessSong {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return (*p)[idx]
	}
	NewProcess := ProcessSong{
		ChatID: chatID,
		Step:   ProcessYoutubeMediaBySpotifySongLinkStart,
	}
	*p = append(*p, NewProcess)
	return NewProcess
}
func (p *ProcessingYoutubeMediaBySpotifySongLink) IfExist(chatID int64) bool {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return true
	}
	return false

}

func (p *ProcessingYoutubeMediaBySpotifySongLink) UpdateStep(step string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Step = step
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessingYoutubeMediaBySpotifySongLink) Delete(chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		*p = append((*p)[:idx], (*p)[idx+1:]...)
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

func (p *ProcessingCreateAndFillYoutubePlaylists) GetOrCreate(chatID int64) ProcessCreateAndFillYoutubePlaylist {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return (*p)[idx]
	}
	NewProcess := ProcessCreateAndFillYoutubePlaylist{
		ChatID: chatID,
		Step:   ProcessFillYouTubePlaylistStart,
	}
	*p = append(*p, NewProcess)
	return NewProcess
}

func (p *ProcessingCreateAndFillYoutubePlaylists) AddSongs(playlist Playlist, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Songs = playlist
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

func (p *ProcessingCreateAndFillYoutubePlaylists) UpdateStep(step string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Step = step
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

func (p *ProcessingCreateAndFillYoutubePlaylists) AddTitle(title string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].Songs.Title = title
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}

func (p *ProcessingCreateAndFillYoutubePlaylists) AddAuthToken(AuthToken string, chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		(*p)[idx].AuthToken = AuthToken
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
func (p *ProcessingCreateAndFillYoutubePlaylists) IfExist(chatID int64) bool {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		return true
	}

	return false
}
func (p *ProcessingCreateAndFillYoutubePlaylists) Delete(chatID int64) error {
	if idx := FindUserIndex(*p, chatID); idx != -1 {
		*p = append((*p)[:idx], (*p)[idx+1:]...)
		return nil
	}
	return errors.New(ErrChatIDNotFound)
}
