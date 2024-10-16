package tg_handlers

import (
	"fmt"
	tg "gopkg.in/tucnak/telebot.v2"
	"log"
)

func (h Handler) GetYoutubeSongHelper(msg *tg.Message) error {
	return h.workFlows.GetYoutubeSong(msg)
}

func (h Handler) GetSpotifySongHelper(msg *tg.Message) error {
	return h.workFlows.GetSpotifySong(msg)
}

func (h Handler) SendMsg(msg *tg.Message, outText string) {
	if _, err := h.bot.Send(msg.Sender, outText); err != nil {
		log.Fatalf("Critical error: Telegram bot failed to send message: %v", err)
	}
}

// HelpOut helper - completed
func (h Handler) HelpOut(msg *tg.Message) {
	formatString := fmt.Sprintf("			Commands: \n" +
		"	/FindSong - search songs in youtube and spotify by track metadata ( track title and artist) \n" +
		"	/SpotifySong - search song in spotify by youtube link of this track \n" +
		"	/YoutubeSong - search track in youtube by spotify link of this track \n " +
		"	/FillSpotifyPlaylist - search and fill spotify playlist by link of the playlist from youtube \n" +
		"	/FillYoutubePlaylist - search and fill youtube playlist by link of the playlist from spotify \n")
	if _, err := h.bot.Send(msg.Sender, formatString); err != nil {
		log.Fatal(err)
	}

}
func (h Handler) GetSongsByMetadataHelper(msg *tg.Message) error {
	return h.workFlows.GetSongsByMetadata(msg)
}
func (h Handler) CreateAndFillYoutubePlaylistHelper(msg *tg.Message) error {
	return h.workFlows.CreateAndFillYoutubePlaylist(msg)
}
func (h Handler) CreateAndFillSpotifyPlaylistHelper(msg *tg.Message) error {
	return h.workFlows.CreateAndFillSpotifyPlaylist(msg)
}
