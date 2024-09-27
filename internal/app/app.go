package app

import (
	"MusicApp/internal/config"
	"MusicApp/internal/domain"
	"MusicApp/internal/handlers"
	"MusicApp/internal/handlers/Spotify"
	"MusicApp/internal/handlers/YouTube"
	"MusicApp/internal/repo/MongoDB"
	"MusicApp/internal/tg_handlers"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/skwizi4/lib/ErrChan"
	logger "github.com/skwizi4/lib/logs"
	tg "gopkg.in/tucnak/telebot.v2"
	"time"
)

// App - structure of app
type App struct {
	appName                string
	bot                    *tg.Bot
	Config                 config.Config
	errChan                *ErrChan.ErrorChannel
	logger                 logger.GoLogger
	validator              *validator.Validate
	mongo                  *MongoDB.MongoDB
	spotifyHandler         handlers.Spotify
	youtubeHandler         handlers.YouTube
	telegramHandler        tg_handlers.Handler
	processingSpotifySongs domain.ProcessingSpotifySongs
}

// New - return new variation of application
func New(AppName string) App {
	return App{
		appName: AppName,
		logger:  logger.InitLogger(),
	}
}

// Run - run application
func (a *App) Run(ctx context.Context) {
	a.InitLogger()
	a.InitErrHandler(ctx)
	a.InitValidator()
	a.PopulateConfig()
	a.InitMongo()
	a.InitBot()
	a.InitHandlers()
	a.ListenTgBot()
}

// InitLogger - inits logger for application
func (a *App) InitLogger() {
	a.logger = logger.InitLogger()
	a.logger.InfoFrmt("InitLogger-Successfully")
}

// InitErrHandler - Inits error Handler channel
func (a *App) InitErrHandler(ctx context.Context) {
	a.errChan = ErrChan.InitErrChan(10, a.logger)
	go func() {
		for {
			select {
			case <-ctx.Done():
				a.errChan.Close()
				return
			}
		}
	}()
	a.errChan.Start()
	a.logger.InfoFrmt("InitErrorHandler-Successfully")
}

// InitValidator Инициализирут  валидатор
func (a *App) InitValidator() {
	a.validator = validator.New()
	a.logger.InfoFrmt("initValidator-Successfully")
}

// PopulateConfig Проверяет конфиг
func (a *App) PopulateConfig() {
	cfg, err := config.ParseConfig("/home/skwizi_4/code/MusicApp/config.json")
	if err != nil {
		a.logger.ErrorFrmt("error in parsing config: %s", err)
	}

	err = cfg.ValidateConfig(a.validator)
	if err != nil {
		a.logger.ErrorFrmt("error in config validation: %s", err)
	}
	a.Config = *cfg
	a.logger.InfoFrmt("InitConfig-Successfully")
}

// InitMongo Инициализируем монго
func (a *App) InitMongo() {
	var err error
	a.mongo, err = MongoDB.InitMongo(a.Config.MongoDbCfg.Uri, a.Config.MongoDbCfg.DataBaseName, a.Config.MongoDbCfg.CollectionName)
	if err != nil {
		a.logger.ErrorFrmt("error in initializing MongoDB  client: %s", err)
	}
	a.logger.InfoFrmt("InitMongo-Successfully")
}

// InitBot Инициализируем бота
func (a *App) InitBot() {
	botSettings := tg.Settings{
		Token:  a.Config.BotCfg.Token,
		Poller: &tg.LongPoller{Timeout: 1 * time.Second},
	}
	var err error
	if a.bot, err = tg.NewBot(botSettings); err != nil {
		a.logger.ErrorFrmt("Error Is occurred in InitTgBot, error: %s", err)
	}
	a.logger.InfoFrmt("InitTgBot-Successfully")
}

// InitHandlers - инициализирует хендлера
func (a *App) InitHandlers() {
	a.spotifyHandler = Spotify.New(a.bot, &a.processingSpotifySongs, a.errChan, a.Config)
	a.youtubeHandler = YouTube.New()
	a.telegramHandler = tg_handlers.New(a.bot, a.spotifyHandler, a.youtubeHandler, a.errChan)
}

// ListenTgBot - todo - отредактировать хендлера под задачи
func (a *App) ListenTgBot() {
	go a.bot.Handle("/SpotifySong", func(msg *tg.Message) {
		go a.telegramHandler.SpotifySong(msg)
	})
	go a.bot.Handle("SpotifyPlaylist", func(msg *tg.Message) {
		go a.telegramHandler.SpotifyPlaylist(msg)
	})
	go a.bot.Handle("/Help", func(msg *tg.Message) {
		go a.telegramHandler.Help(msg)
	})
	go a.bot.Handle("/YoutubeSong", func(msg *tg.Message) {
		go a.telegramHandler.YouTubeSong(msg)
	})
	go a.bot.Handle("/YoutubePlaylist", func(msg *tg.Message) {
		go a.telegramHandler.YouTubePlaylist(msg)
	})
	go a.bot.Handle("/FindSong", func(msg *tg.Message) {
		go a.telegramHandler.FindSong(msg)
	})
	a.bot.Handle(tg.OnText, func(msg *tg.Message) {
		switch {
		case a.processingSpotifySongs.IfExist(msg.Chat.ID):
			go a.telegramHandler.SpotifySong(msg)
		}
	})
	a.bot.Start()
	defer a.bot.Stop()

}
