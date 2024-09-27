package YouTube

import logger "github.com/skwizi4/lib/logs"

const BaseUrl = "https://www.googleapis.com/youtube/v3/"
const ServerUrl = "http://localhost:8080/authToken"
const YoutubeTrackDomen = "https://www.youtube.com/watch?v="

type ServiceYouTube struct {
	Key       string
	BaseUrl   string
	logger    logger.GoLogger
	Token     string
	ClientID  string
	ServerUrl string
}

type (
	youtubeMediaById struct {
		Items []struct {
			VideoId string `json:"id"`
			Snippet struct {
				ChanelName string `json:"channelTitle"`
				Title      string `json:"title"`
			} `json:"snippet"`
		} `json:"items"`
	}
	youtubePlaylistParamsById struct {
		Items []struct {
			Snippet struct {
				Title        string `json:"title"`
				ChannelTitle string `json:"ChannelTitle"`
			} `json:"snippet"`
		} `json:"items"`
	}
	youtubeResponsePlaylistMediaById struct {
		Items []struct {
			Snippet struct {
				Title        string `json:"title"`
				ChannelTitle string `json:"videoOwnerChannelTitle"`
			} `json:"snippet"`
		} `json:"items"`
		NextPageToken string `json:"nextPageToken"`
	}
	youtubeMediaByMetadata struct {
		Items []struct {
			Id struct {
				VideoId string `json:"videoId"`
			} `json:"id"`
			Snippet struct {
				ChanelName string `json:"channelTitle"`
				Title      string `json:"title"`
			} `json:"snippet"`
		} `json:"items"`
	}
	youtubePlaylistIdResp struct {
		ID string `json:"id"`
	}
)
