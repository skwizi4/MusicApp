package YouTube

import (
	"MusicApp/internal/config"
	"MusicApp/internal/domain"
	"MusicApp/internal/services"
	"MusicApp/internal/services/Spotify"
	YouTubeService "MusicApp/internal/services/YouTube"
	"fmt"
	"github.com/skwizi4/lib/ErrChan"
	"github.com/skwizi4/lib/logs"
	tg "gopkg.in/tucnak/telebot.v2"
	"log"
)

type Handler struct {
	bot                                   *tg.Bot
	ProcessingSpotifySongByYoutubeMediaId *domain.ProcessingSpotifySongByYoutubeMediaId
	errChannel                            *ErrChan.ErrorChannel
	spotifyService                        services.SpotifyService
	youtubeService                        services.YouTubeService
	cfg                                   config.Config
	logger                                logs.GoLogger
}

func New(bot *tg.Bot, cfg config.Config, processingYoutubeSongsById *domain.ProcessingSpotifySongByYoutubeMediaId, errChan *ErrChan.ErrorChannel, logger logs.GoLogger) Handler {
	return Handler{
		bot:                                   bot,
		spotifyService:                        Spotify.NewSpotifyService(cfg),
		youtubeService:                        YouTubeService.NewYouTubeService(cfg),
		ProcessingSpotifySongByYoutubeMediaId: processingYoutubeSongsById,
		errChannel:                            errChan,
		cfg:                                   cfg,
		logger:                                logger,
	}
}

func (h Handler) GetSpotifySongByYoutubeLink(msg *tg.Message) error {
	process := h.ProcessingSpotifySongByYoutubeMediaId.GetOrCreate(msg.Chat.ID)
	switch process.Step {
	case domain.ProcessSpotifySongByYouTubeMediaStart:
		if _, err := h.bot.Send(msg.Sender, "Send youtube link of media that you wanna find in Spotify"); err != nil {
			log.Fatal(err)

		}
		if err := h.ProcessingSpotifySongByYoutubeMediaId.UpdateStep(domain.ProcessSpotifySongByYouTubeMediaEnd, msg.Chat.ID); err != nil {
			h.errChannel.HandleError(err)
			return err
		}
	case domain.ProcessSpotifySongByYouTubeMediaEnd:
		track, err := h.youtubeService.GetYoutubeMediaByID(msg.Text)
		if err != nil {
			h.errChannel.HandleError(err)
			if _, err = h.bot.Send(msg.Sender, "Error"); err != nil {
				log.Fatal(err)

			}
			if err = h.ProcessingSpotifySongByYoutubeMediaId.Delete(msg.Chat.ID); err != nil {
				log.Fatal(err)

			}

			return err
		}
		song, err := h.spotifyService.GetSpotifyTrackByMetadata(domain.MetaData{Title: track.Title, Artist: track.Artist})
		if err != nil {
			h.errChannel.HandleError(err)
			if _, err = h.bot.Send(msg.Sender, "Error"); err != nil {
				log.Fatal(err)

			}
			if err = h.ProcessingSpotifySongByYoutubeMediaId.Delete(msg.Chat.ID); err != nil {
				log.Fatal(err)
			}
		}
		SongPrint := fmt.Sprintf("Spotify Song Title: %s \n Spotify Song Artist: %s , \n Spotify Song link: %s", song.Artist, song.Title, song.Link)
		if _, err = h.bot.Send(msg.Sender, SongPrint); err != nil {
			log.Fatal(err)
		}
		if err = h.ProcessingSpotifySongByYoutubeMediaId.Delete(msg.Chat.ID); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func (h Handler) GetSongByMetaData(metadata *domain.MetaData) (*domain.Song, error) {
	track, err := h.spotifyService.GetSpotifyTrackByMetadata(*metadata)
	if err != nil {
		h.errChannel.HandleError(err)
		return nil, err
	}

	return track, nil
}

func (h Handler) GetPlaylistBySpotifyLink(spotifyLink string) (*domain.Playlist, error) {
	return nil, nil
}

func (h Handler) GetPlaylistByMetaData(metadata domain.MetaData) (*domain.Playlist, error) {
	return nil, nil
}
