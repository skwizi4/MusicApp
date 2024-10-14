package tg_handlers

import (
	tg "gopkg.in/tucnak/telebot.v2"
)

// GetYoutubeSong - completed
func (h Handler) GetYoutubeSong(msg *tg.Message) {
	h.errChannel.HandleError(h.GetYoutubeSongHelper(msg))
}

// GetSpotifySong -  completed
func (h Handler) GetSpotifySong(msg *tg.Message) {
	h.errChannel.HandleError(h.GetSpotifySongHelper(msg))
}

// Help - completed
func (h Handler) Help(msg *tg.Message) {
	h.HelpOut(msg)
}

// FindSong - completed
func (h Handler) FindSong(msg *tg.Message) {
	h.errChannel.HandleError(h.GetSongsByMetadataHelper(msg))
}

// CreateFillYoutubePlaylist  todo write handler
func (h Handler) CreateFillYoutubePlaylist(msg *tg.Message) {
	h.errChannel.HandleError(h.CreateAndFillYoutubePlaylistHelper(msg))
}

// CreateFillSpotifyPlaylist todo - write handler
func (h Handler) CreateFillSpotifyPlaylist(msg *tg.Message) {
	h.errChannel.HandleError(h.CreateAndFillSpotifyPlaylistHelper(msg))
}
