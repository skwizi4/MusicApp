package Spotify

import (
	"MusicApp/internal/config"
	"MusicApp/internal/domain"
	"MusicApp/internal/services"
	"MusicApp/internal/services/Spotify"
	YouTubeService "MusicApp/internal/services/YouTube"
	"errors"
	"fmt"
	"github.com/skwizi4/lib/ErrChan"
	"github.com/skwizi4/lib/logs"
	tg "gopkg.in/tucnak/telebot.v2"
	"log"
)

type Handler struct {
	bot                                   *tg.Bot
	ProcessingYoutubeMediaBySpotifySongID *domain.ProcessingYoutubeMediaBySpotifySongID
	errChannel                            *ErrChan.ErrorChannel
	spotifyService                        services.SpotifyService
	youtubeService                        services.YouTubeService
	cfg                                   config.Config
	logger                                logs.GoLogger
}

func New(bot *tg.Bot, processingSpotifySongs *domain.ProcessingYoutubeMediaBySpotifySongID,
	errChan *ErrChan.ErrorChannel, cfg config.Config, logger logs.GoLogger,
) Handler {
	return Handler{
		bot:                                   bot,
		spotifyService:                        Spotify.NewSpotifyService(cfg),
		youtubeService:                        YouTubeService.NewYouTubeService(cfg),
		ProcessingYoutubeMediaBySpotifySongID: processingSpotifySongs,
		errChannel:                            errChan,
		cfg:                                   cfg,
		logger:                                logger,
	}
}

func (h Handler) GetMediaBySpotifyLink(msg *tg.Message) error {
	process := h.ProcessingYoutubeMediaBySpotifySongID.GetOrCreate(msg.Chat.ID)

	switch process.Step {
	case domain.ProcessSpotifySongByIdStart:
		if _, err := h.bot.Send(msg.Sender, "Send link of song that you wanna find"); err != nil {
			log.Fatal("err in tg.bot :", err)
		}
		if err := h.ProcessingYoutubeMediaBySpotifySongID.UpdateStep(domain.ProcessSpotifySongByIdEnd, msg.Chat.ID); err != nil {
			log.Fatal("err in domain :", err)
		}
	case domain.ProcessSpotifySongByIdEnd:
		track, err := h.spotifyService.GetSpotifyTrackById(msg.Text)
		if err != nil {
			h.errChannel.HandleError(err)
			if _, err = h.bot.Send(msg.Sender, "Error"); err != nil {
				log.Fatal(err, errors.New("error in tg.bot"))
			}
			if err = h.ProcessingYoutubeMediaBySpotifySongID.Delete(msg.Chat.ID); err != nil {
				log.Fatal(err, errors.New("error in domain"))
			}

			return err
		}
		song, err := h.youtubeService.GetYoutubeMediaByMetadata(domain.MetaData{Title: track.Title, Artist: track.Artist})
		if err != nil {
			h.errChannel.HandleError(err)
			if _, err = h.bot.Send(msg.Sender, "Error"); err != nil {
				log.Fatal(err, errors.New("error in tg.bot"))

			}
			if err = h.ProcessingYoutubeMediaBySpotifySongID.Delete(msg.Chat.ID); err != nil {
				log.Fatal(err, errors.New("error in domain"))
			}
			return err
		}
		SongPrint := fmt.Sprintf("Song Title: %s \n Song Artist: %s , \n Song link: %s", song.Artist, song.Title, song.Link)
		if _, err = h.bot.Send(msg.Sender, SongPrint); err != nil {
			log.Fatal(err)
		}
		if err = h.ProcessingYoutubeMediaBySpotifySongID.Delete(msg.Chat.ID); err != nil {
			log.Fatal(err)
		}

	}
	return nil

}
func (h Handler) GetSongByMetaData(metadata *domain.MetaData) (*domain.Song, error) {

	track, err := h.youtubeService.GetYoutubeMediaByMetadata(domain.MetaData{Title: metadata.Title, Artist: metadata.Artist})
	if err != nil {
		h.errChannel.HandleError(err)
		return nil, err
	}

	return track, nil
}
func (h Handler) GetPlaylistByYoutubeLink(youtubeLink string) (*domain.Playlist, error) {
	return nil, nil
}

func (h Handler) GetPlaylistByMetaData(metadata domain.MetaData) (*domain.Playlist, error) {
	return nil, nil
}
