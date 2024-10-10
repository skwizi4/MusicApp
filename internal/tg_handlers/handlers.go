package tg_handlers

import (
	"MusicApp/internal/domain"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	tg "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"time"
)

// YoutubeSong - completed
func (h Handler) YoutubeSong(msg *tg.Message) {
	if err := h.GetYoutubeSong(msg); err != nil {
		h.errChannel.HandleError(err)
	}
}

// SpotifySong -  completed
func (h Handler) SpotifySong(msg *tg.Message) {
	if err := h.GetSpotifySong(msg); err != nil {
		h.errChannel.HandleError(err)
	}
}

// Help - completed
func (h Handler) Help(msg *tg.Message) {
	h.HelpOut(msg)
}

// FindSong - completed
func (h Handler) FindSong(msg *tg.Message) {
	if err := h.GetSongsByMetadata(msg); err != nil {
		h.errChannel.HandleError(err)
	}
}

// FillYoutubePlaylist  todo write handler
func (h Handler) FillYoutubePlaylist(msg *tg.Message) {
	process := h.processingFillYoutubePlaylists.GetOrCreate(msg.Chat.ID)
	switch process.Step {
	case domain.ProcessFillYouTubePlaylistStart:
		h.SendMsg(msg, "send a link to the playlist from spotify that you want to transfer to youtube")
		if err := h.processingFillYoutubePlaylists.UpdateStep(domain.ProcessFillYouTubePlaylistSendAuthLink, msg.Chat.ID); err != nil {
			h.errChannel.HandleError(err)
			h.SendMsg(msg, "Error, try again")
			h.DeleteProcessingFindSongByMetadata(msg)
			return
		}
	case domain.ProcessFillYouTubePlaylistSendAuthLink:
		playlist, err := h.spotifyHandler.GetSpotifyPlaylistByLink(msg.Text)
		if err != nil {
			h.errChannel.HandleError(err)
			h.SendMsg(msg, "Error, try again")
			h.DeleteProcessingFindSongByMetadata(msg)
			return
		}
		if err = h.processingFillYoutubePlaylists.AddTitle(playlist.Title, msg.Chat.ID); err != nil {
			h.errChannel.HandleError(err)
			h.SendMsg(msg, "Error, try again")
			h.DeleteProcessingFindSongByMetadata(msg)
			return
		}
		if err = h.processingFillYoutubePlaylists.AddSongs(*playlist, msg.Chat.ID); err != nil {
			h.errChannel.HandleError(err)
			h.SendMsg(msg, "Error, try again")
			h.DeleteProcessingFindSongByMetadata(msg)
			return
		}
		h.SendMsg(msg, fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=https://www.googleapis.com/auth/youtube&access_type=offline&state=%s&prompt=consent",
			h.cfg.YoutubeCfg.ClientID, h.cfg.YoutubeCfg.RedirectUrl, strconv.FormatInt(msg.Sender.ID, 10)))
		h.CheckForToken(msg)
		process = h.processingFillYoutubePlaylists.GetOrCreate(msg.Chat.ID)
		youtubePlaylist, err := h.youtubeHandler.CreateAndFillYoutubePlaylist(process.Songs, process.AuthToken)
		if err != nil {
			h.errChannel.HandleError(err)
			h.SendMsg(msg, "Error, try later")
			h.DeleteProcessingFindSongByMetadata(msg)
			return
		}
		h.SendMsg(msg, fmt.Sprintf("Playlist title: %s \n Playlist link: %s", youtubePlaylist.Title, youtubePlaylist.ExternalUrl))

	}

}

func (h Handler) CheckForToken(msg *tg.Message) {
	stats := h.mongo.Health()

	if stats["message"] != "It's healthy" {
		log.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	} else {
		fmt.Println(stats["message"])
	}
	telegramID := strconv.FormatInt(msg.Sender.ID, 10)
	timeout := time.After(90 * time.Second)
	ticker := time.NewTicker(2500 * time.Millisecond)

	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			h.SendMsg(msg, "Authorization timeout. Please try again.")
			h.DeleteProcessingFindSongByMetadata(msg)
			return
		case <-ticker.C:
			fmt.Println(telegramID)
			data, err := h.mongo.Get(telegramID)
			if err != nil {
				// Игнорируем ошибку, если запись отсутствует
				if errors.Is(err, mongo.ErrNoDocuments) {
					fmt.Println(err.Error())
					continue // Продолжаем проверку
				} else {
					// Обрабатываем другие ошибки
					h.errChannel.HandleError(err)
					h.SendMsg(msg, "Error, try later.")
					h.DeleteProcessingFindSongByMetadata(msg)
					return
				}
			}
			if err = h.processingFillYoutubePlaylists.AddAuthToken(data.Token, msg.Chat.ID); err != nil {
				h.errChannel.HandleError(err)
				h.DeleteProcessingFindSongByMetadata(msg)
				return
			}
			if err = h.mongo.Delete(telegramID); err != nil {
				h.errChannel.HandleError(err)
				h.SendMsg(msg, "Error, try later")
				h.DeleteProcessingFindSongByMetadata(msg)

			}
			h.SendMsg(msg, "Authorization completed! Processing creating playlist...")
			return
		}
	}
}

// FillSpotifyPlaylist todo - write handler
func (h Handler) FillSpotifyPlaylist(msg *tg.Message) {

}
