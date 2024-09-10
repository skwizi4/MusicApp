package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type (
	Config struct {
		BotToken     BotConfig     `json:"bot_token"`
		SpotifyToken SpotifyConfig `json:"spotify_token"`
		YoutubeToken YoutubeConfig `json:"youtube_token"`
		MongoDb      MongoDb       `json:"mongodb"`
	}
	BotConfig struct {
		Token string `json:"token"`
	}
	SpotifyConfig struct {
	}
	YoutubeConfig struct {
	}
	MongoDb struct {
		Uri            string `json:"Uri"`
		DataBaseName   string `json:"DataBaseName"`
		CollectionName string `json:"collection_name"`
	}
)

func ParseConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}
	var Cfg Config
	err = json.Unmarshal(file, &Cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}
	return &Cfg, nil
}
