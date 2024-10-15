package workflows

import (
	"MusicApp/internal/config"
	"MusicApp/internal/domain"
	"MusicApp/internal/handlers"
	"MusicApp/internal/repo/MongoDB"
	errs "errors"
	"github.com/skwizi4/lib/ErrChan"
	logger "github.com/skwizi4/lib/logs"
	"go.mongodb.org/mongo-driver/mongo"
	tg "gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

type WorkFlows struct {
	bot                                     *tg.Bot
	logger                                  logger.GoLogger
	ErrChannel                              *ErrChan.ErrorChannel
	ProcessingYoutubeMediaBySpotifySongLink *domain.ProcessingYoutubeMediaBySpotifySongLink
	ProcessingSpotifySongByYoutubeMediaLink *domain.ProcessingSpotifySongByYoutubeMediaLink
	ProcessingFindSongByMetadata            *domain.ProcessingFindSongByMetadata
	ProcessingCreateAndFillYoutubePlaylists *domain.ProcessingCreateAndFillYoutubePlaylists
	ProcessingCreateAndFillSpotifyPlaylists *domain.ProcessingCreateAndFillSpotifyPlaylists
	YouTubeHandler                          handlers.Youtube
	SpotifyHandler                          handlers.Spotify
	cfg                                     config.Config
	mongo                                   *MongoDB.MongoDB
}

func New(bot *tg.Bot, ProcessingYoutubeMediaBySpotifySongLink *domain.ProcessingYoutubeMediaBySpotifySongLink,
	ProcessingSpotifySongByYoutubeMediaLink *domain.ProcessingSpotifySongByYoutubeMediaLink,
	ProcessingFindSongByMetadata *domain.ProcessingFindSongByMetadata,
	ProcessingCreateAndFillYoutubePlaylists *domain.ProcessingCreateAndFillYoutubePlaylists,
	ProcessingCreateAndFillSpotifyPlaylists *domain.ProcessingCreateAndFillSpotifyPlaylists,
	logger logger.GoLogger, errChannel *ErrChan.ErrorChannel, youtubeHandler handlers.Youtube, spotifyHandler handlers.Spotify, cfg config.Config, mongo *MongoDB.MongoDB,
) *WorkFlows {
	if ProcessingYoutubeMediaBySpotifySongLink == nil || ProcessingSpotifySongByYoutubeMediaLink == nil || ProcessingFindSongByMetadata == nil || ProcessingCreateAndFillYoutubePlaylists == nil || ProcessingCreateAndFillSpotifyPlaylists == nil {
		log.Fatal("err: some process are nil")
	}
	return &WorkFlows{
		bot:                                     bot,
		ProcessingYoutubeMediaBySpotifySongLink: ProcessingYoutubeMediaBySpotifySongLink,
		ProcessingSpotifySongByYoutubeMediaLink: ProcessingSpotifySongByYoutubeMediaLink,
		ProcessingFindSongByMetadata:            ProcessingFindSongByMetadata,
		ProcessingCreateAndFillYoutubePlaylists: ProcessingCreateAndFillYoutubePlaylists,
		ProcessingCreateAndFillSpotifyPlaylists: ProcessingCreateAndFillSpotifyPlaylists,
		logger:                                  logger,
		ErrChannel:                              errChannel,
		YouTubeHandler:                          youtubeHandler,
		SpotifyHandler:                          spotifyHandler,
		cfg:                                     cfg,
		mongo:                                   mongo,
	}
}

func (w WorkFlows) SendMsg(msg *tg.Message, outText string) {
	if _, err := w.bot.Send(msg.Sender, outText); err != nil {
		log.Fatalf("Critical error: Telegram bot failed to send message: %v", err)
	}

}
func (w WorkFlows) CheckForToken(msg *tg.Message, UserProcess string) (string, error) {
	timeout := time.After(90 * time.Second)
	ticker := time.NewTicker(2500 * time.Millisecond)

	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return "", errs.New("authorization timeout")
		case <-ticker.C:
			data, err := w.mongo.Get(UserProcess)
			if err != nil {
				if errs.Is(err, mongo.ErrNoDocuments) {
					continue
				} else {
					return "", err
				}
			}

			if err = w.mongo.Delete(UserProcess); err != nil {
				return "", err

			}
			w.SendMsg(msg, "Authorization completed! Processing creating playlist...")
			return data.Token, nil
		}
	}
}
