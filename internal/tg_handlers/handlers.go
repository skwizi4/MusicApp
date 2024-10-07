package tg_handlers

import (
	"MusicApp/internal/domain"
	"MusicApp/internal/handlers"
	"fmt"
	"github.com/skwizi4/lib/ErrChan"
	"github.com/skwizi4/lib/logs"
	tg "gopkg.in/tucnak/telebot.v2"
	"log"
)

const ErrProcessingFindSongByMetadata = "cant end processingFindSongByMetadata"
const ErrInvalidParamsSpotify = "wrong request (invalid params for Spotify)"
const ErrInvalidParamsYoutube = "wrong request (invalid params for Youtube)"

type Handler struct {
	bot                          *tg.Bot
	youtubeHandler               handlers.Youtube
	spotifyHandler               handlers.Spotify
	errChannel                   *ErrChan.ErrorChannel
	processingFindSongByMetadata *domain.ProcessingFindSongByMetadata
	logger                       logs.GoLogger
}

func New(bot *tg.Bot, youtubeHandler handlers.Youtube, spotifyHandler handlers.Spotify, processingFinSongByMetadata *domain.ProcessingFindSongByMetadata, errChan *ErrChan.ErrorChannel, logger logs.GoLogger) Handler {
	return Handler{
		bot:                          bot,
		youtubeHandler:               youtubeHandler,
		spotifyHandler:               spotifyHandler,
		errChannel:                   errChan,
		logger:                       logger,
		processingFindSongByMetadata: processingFinSongByMetadata,
	}
}

func (h Handler) YoutubeSong(msg *tg.Message) {
	err := h.youtubeHandler.GetMediaBySpotifyLink(msg)
	if err != nil {
		h.errChannel.HandleError(err)
	}

}
func (h Handler) SpotifySong(msg *tg.Message) {
	err := h.spotifyHandler.GetSpotifySongByYoutubeLink(msg)
	if err != nil {
		h.errChannel.HandleError(err)
	}

}
func (h Handler) Help(msg *tg.Message) {
	h.HelpOut(msg)
}
func (h Handler) FindSong(msg *tg.Message) {
	process := h.processingFindSongByMetadata.GetOrCreate(msg.Chat.ID)
	switch process.Step {
	case domain.ProcessSpotifySongByMetadataStart:
		err := h.GetMetadata(msg)
		if err != nil {
			h.errChannel.HandleError(err)
			return
		}
	case domain.ProcessSpotifySongByMetadataTitle:

		err := h.GetMetadata(msg)
		if err != nil {
			h.errChannel.HandleError(err)
			return
		}
	case domain.ProcessSpotifySongByMetadataArtist:
		err := h.GetMetadata(msg)
		if err != nil {
			h.errChannel.HandleError(err)
			return
		}
	case domain.ProcessSpotifySongByMetadataEnd:
		if err := h.processingFindSongByMetadata.AddArtist(msg.Chat.ID, msg.Text); err != nil {

			h.errChannel.HandleError(err)
			err = h.processingFindSongByMetadata.Delete(msg.Chat.ID)
			if err != nil {
				log.Fatal(ErrProcessingFindSongByMetadata)
				return
			}
		}
		if _, err := h.bot.Send(msg.Sender, "wait a second ...."); err != nil {
			if err = h.processingFindSongByMetadata.Delete(msg.Chat.ID); err != nil {
				h.errChannel.HandleError(err)
				log.Fatal(ErrProcessingFindSongByMetadata)
				return

			}

		}
		process = h.processingFindSongByMetadata.GetOrCreate(msg.Chat.ID)
		spotifySong, err := h.spotifyHandler.GetSongByMetaData(&domain.MetaData{Title: process.Song.Title, Artist: process.Song.Artist})
		if err != nil {
			if err.Error() == ErrInvalidParamsSpotify {
				if _, err = h.bot.Send(msg.Sender, fmt.Sprintf("Error: check input parameters. We cant find song with title: %s, and artist: %s ",
					process.Song.Title, process.Song.Artist)); err != nil {
					if err = h.processingFindSongByMetadata.Delete(msg.Chat.ID); err != nil {
						log.Fatal(ErrProcessingFindSongByMetadata)
					}
					h.errChannel.HandleError(err)
				}
				if err = h.processingFindSongByMetadata.Delete(msg.Chat.ID); err != nil {
					log.Fatal(ErrProcessingFindSongByMetadata)
				}
				h.errChannel.HandleError(err)
				return
			}

		}
		youtubeSong, err := h.youtubeHandler.GetSongByMetaData(&domain.MetaData{Title: process.Song.Title, Artist: process.Song.Artist})
		if err != nil {
			if err.Error() == ErrInvalidParamsYoutube {
				if _, err = h.bot.Send(msg.Sender, fmt.Sprintf("Error: check input parameters. We cant find song with title: %s, and artist: %s ",
					process.Song.Title, process.Song.Artist)); err != nil {
					h.errChannel.HandleError(err)
					if err = h.processingFindSongByMetadata.Delete(msg.Chat.ID); err != nil {
						log.Fatal(ErrProcessingFindSongByMetadata)
						return
					}
				}
				h.errChannel.HandleError(err)
				return
			}

		}
		SongPrint := fmt.Sprintf("Spotify song title: %s; \n Spotify song artist: %s; \n Spotify song link: %s; \n \n Youtube song title: %s;  \n Youtube song artist: %s; \n Youtube song link: %s.",
			spotifySong.Artist, spotifySong.Title, spotifySong.Link, youtubeSong.Title, youtubeSong.Artist, youtubeSong.Link)
		if _, err = h.bot.Send(msg.Sender, SongPrint); err != nil {
			err = h.processingFindSongByMetadata.Delete(msg.Chat.ID)
			if err != nil {
				log.Fatal(ErrProcessingFindSongByMetadata)
				return
			}
			h.errChannel.HandleError(err)
			return
		}
		err = h.processingFindSongByMetadata.Delete(msg.Chat.ID)
		if err != nil {
			log.Fatal(ErrProcessingFindSongByMetadata)
			return
		}
	}

}

func (h Handler) SpotifyPlaylist(msg *tg.Message) {

}
func (h Handler) YouTubePlaylist(msg *tg.Message) {

}
