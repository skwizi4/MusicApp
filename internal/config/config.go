package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type (
	Config struct {
		BotCfg     BotConfig     `json:"botCfg"`
		SpotifyCfg SpotifyConfig `json:"spotifyCfg"`
		YoutubeCfg YoutubeConfig `json:"youtubeCfg"`
		MongoDbCfg MongoDb       `json:"mongoDbCfg"`
	}
	BotConfig struct {
		Token string `json:"token"`
	}
	SpotifyConfig struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}
	YoutubeConfig struct {
		Key string `json:"key"`
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
