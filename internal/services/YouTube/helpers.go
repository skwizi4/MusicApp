package YouTube

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func ParseYouTubeIDFromURL(link string) (string, string, error) {
	if strings.Contains(link, "https://www.youtube.com/watch?v=") {
		id := strings.Split(link, "https://www.youtube.com/watch?v=")[1]
		return "track", strings.Split(id, "?")[0], nil
	}
	if strings.Contains(link, "https://youtube.com/playlist?list=") {
		id := strings.Split(link, "https://youtube.com/playlist?list=")[1]
		return "playlist", strings.Split(id, "?")[0], nil
	}
	return "", "", errors.New("can't parse ID, invalid URL format")
}
func (y ServiceYouTube) CreateRequest(method, endpoint string) (*http.Request, error) {
	Url := y.BaseUrl + endpoint
	req, err := http.NewRequest(method, Url, nil)
	if err != nil {
		return nil, err
	}
	// Func with token//
	return req, nil
}
func (y ServiceYouTube) DoRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return resp, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return resp, nil
}
func (y ServiceYouTube) CreateEndpoint(id string) (string, error) {
	if y.Key == "" {
		return "", errors.New("YouTube key is empty")
	}
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?id=%s&part=snippet,statistics&key=%s", id, y.Key)
	return url, nil
}
func DecodedBody(resp *http.Response, ResponseStruct interface{}) error {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &ResponseStruct)
	if err != nil {
		return err
	}

	return nil
}
