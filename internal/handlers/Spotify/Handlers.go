package Spotify

import (
	"MusicApp/internal/domain"
	"github.com/skwizi4/lib/ErrChan"
	tg "gopkg.in/tucnak/telebot.v2"
)

type Handler struct {
	bot                    tg.Bot
	processingSpotifySongs domain.ProcessingSpotifySongs
	errChannel             *ErrChan.ErrorChannel
}

func New() Handler {
	return Handler{}
}

func (h Handler) GetSongByYoutubeLink(msg tg.Message) (*domain.Song, error) {
	process := h.processingSpotifySongs.GetOrCreate(msg.Chat.ID)
	switch process.Step {
	case domain.ProcessSpotifySongStart:
		if _, err := h.bot.Send(msg.Sender, "Send link of song that you wanna find"); err != nil {
			if err = h.processingSpotifySongs.Delete(msg.Chat.ID); err != nil {
				h.errChannel.HandleError(err)
				return nil, err
			}
			h.errChannel.HandleError(err)

			return nil, err
		}
		if err := h.processingSpotifySongs.UpdateStep(domain.ProcessSpotifySongEnd, msg.Chat.ID); err != nil {
			h.errChannel.HandleError(err)
			return nil, err
		}
	case domain.ProcessSpotifySongEnd:

		if err := h.processingSpotifySongs.AddSongID( /* todo - Написать получение ID песни из ссылки*/ "", msg.Chat.ID); err != nil {
			h.errChannel.HandleError(err)
			return nil, err
		}
		if err := h.processingSpotifySongs.AddSong( /*todo - Написать запрос к API Spotify и получение песни*/ domain.Song{}, msg.Chat.ID); err != nil {
			h.errChannel.HandleError(err)
			return nil, err
		}

	}
	return nil, nil

}
func (h Handler) GetSongByMetaData(metadata domain.MetaData) (*domain.Song, error) {
	return nil, nil
}
func (h Handler) GetPlaylistByYoutubeLink(youtubeLink string) (*domain.Playlist, error) {
	return nil, nil
}

func (h Handler) GetPlaylistByMetaData(metadata domain.MetaData) (*domain.Playlist, error) {
	return nil, nil
}
