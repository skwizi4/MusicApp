package Spotify

import (
	"MusicApp/internal/config"
	"MusicApp/internal/domain"
	"MusicApp/internal/services"
	"MusicApp/internal/services/Spotify"
	YouTubeService "MusicApp/internal/services/YouTube"
	"fmt"
	"github.com/skwizi4/lib/ErrChan"
	tg "gopkg.in/tucnak/telebot.v2"
)

type Handler struct {
	bot                        *tg.Bot
	processingSpotifySongsById *domain.ProcessingSpotifySongsByID
	errChannel                 *ErrChan.ErrorChannel
	spotifyService             services.SpotifyService
	youtubeService             services.YouTubeService
	cfg                        config.Config
}

func New(bot *tg.Bot, processingSpotifySongs *domain.ProcessingSpotifySongsByID,
	errChan *ErrChan.ErrorChannel, cfg config.Config,
) Handler {
	return Handler{
		bot:                        bot,
		spotifyService:             Spotify.NewSpotifyService(cfg),
		youtubeService:             YouTubeService.NewYouTubeService(cfg),
		processingSpotifySongsById: processingSpotifySongs,
		errChannel:                 errChan,
		cfg:                        cfg,
	}
}

func (h Handler) GetSongByYoutubeLink(msg *tg.Message) error {
	process := h.processingSpotifySongsById.GetOrCreate(msg.Chat.ID)

	switch process.Step {
	case domain.ProcessSpotifySongByIdStart:
		if _, err := h.bot.Send(msg.Sender, "Send link of song that you wanna find"); err != nil {
			if err = h.processingSpotifySongsById.Delete(msg.Chat.ID); err != nil {
				h.errChannel.HandleError(err)
				return err
			}
			h.errChannel.HandleError(err)

			return err
		}
		if err := h.processingSpotifySongsById.UpdateStep(domain.ProcessSpotifySongByIdEnd, msg.Chat.ID); err != nil {
			h.errChannel.HandleError(err)
			return err
		}
	case domain.ProcessSpotifySongByIdEnd:
		track, err := h.spotifyService.GetSpotifyTrackById(msg.Text)
		if err != nil {
			h.errChannel.HandleError(err)
			if _, err = h.bot.Send(msg.Sender, "Error"); err != nil {
				return err

			}
			return err
		}
		song, err := h.youtubeService.GetYoutubeMediaByMetadata(domain.MetaData{Title: track.Title, Artist: track.Artist})
		if err != nil {
			if _, err = h.bot.Send(msg.Sender, "Error"); err != nil {
				return err

			}
			return err
		}
		SongPrint := fmt.Sprintf("Song Title: %s \n Song Artist: %s , \n Song link: %s", song.Artist, song.Title, song.Link)
		if _, err = h.bot.Send(msg.Sender, SongPrint); err != nil {
			h.errChannel.HandleError(err)
			return err
		}
		if err = h.processingSpotifySongsById.Delete(msg.Chat.ID); err != nil {
			h.errChannel.HandleError(err)
			return err
		}

	}
	return nil

}
func (h Handler) GetSongByMetaData(metadata *domain.MetaData) (*domain.Song, error) {

	//process := h.processingSpotifySongsByMetadata.GetOrCreate(msg.Chat.ID)
	//
	//switch process.Step {
	//case domain.ProcessSpotifySongByMetadataStart:
	//	if _, err := h.bot.Send(msg.Sender, "Send Title of song that you wanna find"); err != nil {
	//		if err = h.processingSpotifySongsByMetadata.Delete(msg.Chat.ID); err != nil {
	//			h.errChannel.HandleError(err)
	//			return err
	//		}
	//		h.errChannel.HandleError(err)
	//
	//		return err
	//	}
	//	if err := h.processingSpotifySongsByMetadata.UpdateStep(domain.ProcessSpotifySongByMetadataTitle, msg.Chat.ID); err != nil {
	//		h.errChannel.HandleError(err)
	//		return err
	//	}
	//case domain.ProcessSpotifySongByMetadataTitle:
	//
	//	if err := h.processingSpotifySongsByMetadata.AddTitle(msg.Chat.ID, msg.Text); err != nil {
	//		h.errChannel.HandleError(err)
	//	}
	//	if _, err := h.bot.Send(msg.Sender, "Send Artist of song that you wanna find"); err != nil {
	//		if err = h.processingSpotifySongsByMetadata.Delete(msg.Chat.ID); err != nil {
	//			h.errChannel.HandleError(err)
	//		}
	//		h.errChannel.HandleError(err)
	//	}
	//	if err := h.processingSpotifySongsByMetadata.UpdateStep(domain.ProcessSpotifySongByMetadataEnd, msg.Chat.ID); err != nil {
	//		h.errChannel.HandleError(err)
	//		return err
	//	}
	//case domain.ProcessSpotifySongByMetadataEnd:
	//	if err := h.processingSpotifySongsByMetadata.AddArtist(msg.Chat.ID, msg.Text); err != nil {
	//		h.errChannel.HandleError(err)
	//		return err
	//	}
	//	user := h.processingSpotifySongsByMetadata.GetOrCreate(msg.Chat.ID)

	track, err := h.spotifyService.GetSpotifyTrackByMetadata(domain.MetaData{Title: metadata.Title, Artist: metadata.Artist})
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
