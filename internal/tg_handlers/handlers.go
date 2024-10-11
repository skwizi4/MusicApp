package tg_handlers

import (
	tg "gopkg.in/tucnak/telebot.v2"
)

// YoutubeSong - completed
func (h Handler) YoutubeSong(msg *tg.Message) {
	h.errChannel.HandleError(h.GetYoutubeSong(msg))
}

// SpotifySong -  completed
func (h Handler) SpotifySong(msg *tg.Message) {
	h.errChannel.HandleError(h.GetSpotifySong(msg))
}

// Help - completed
func (h Handler) Help(msg *tg.Message) {
	h.HelpOut(msg)
}

// FindSong - completed
func (h Handler) FindSong(msg *tg.Message) {
	h.errChannel.HandleError(h.GetSongsByMetadata(msg))
}

// FillYoutubePlaylist  todo write handler
func (h Handler) FillYoutubePlaylist(msg *tg.Message) {
	h.errChannel.HandleError(h.CreateAndFillYoutubePlaylist(msg))
}

// FillSpotifyPlaylist todo - write handler
func (h Handler) FillSpotifyPlaylist(msg *tg.Message) {

}
