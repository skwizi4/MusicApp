package YouTube

import logger "github.com/skwizi4/lib/logs"

const BaseUrl = "https://www.googleapis.com/youtube/v3/"

type ServiceYouTube struct {
	Key     string
	BaseUrl string
	logger  logger.GoLogger
}

type (
	youtubeMediaById struct {
		Items []struct {
			Snippet struct {
				ChanelName string `json:"channelTitle"`
				Title      string `json:"title"`
			} `json:"snippet"`
		} `json:"items"`
	}
	youtubePlaylistById struct {
		Items []struct {
			Snippet struct {
				Title        string `json:"title"`
				ChannelTitle string `json:"videoOwnerChannelTitle"`
			} `json:"snippet"`
		} `json:"items"`
	}
	youtubeMediaByMetadata struct{}
)
