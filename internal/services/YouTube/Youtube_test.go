package YouTube

import (
	"MusicApp/internal/app"
	"fmt"
	"testing"
)

func TestGetMediaById(t *testing.T) {
	a := app.New("test")
	a.InitValidator()
	a.PopulateConfig()
	youtubeService := New(a.Config)
	song, err := youtubeService.GetYoutubeMediaByID("https://www.youtube.com/watch?v=hTWKbfoikeg")
	if err != nil {
		t.Error(err)
	}
	youtubeService.logger.InfoFrmt("NameOfSong: ", song.Title)
	youtubeService.logger.InfoFrmt("ArtistName: ", song.Artist)
	youtubeService.logger.InfoFrmt("AlbumName: ", song.Album)
	youtubeService.logger.InfoFrmt("Link ", song.Link)

}
func TestGetPlaylistById(t *testing.T) {
	a := app.New("test")
	a.InitValidator()
	a.PopulateConfig()
	youtubeService := New(a.Config)
	playlist, err := youtubeService.GetYoutubePlaylistByID("https://youtube.com/playlist?list=PLbTTxxr-hMmzSGTsO5mdYLrvKY-RZFanp&si=sB4Nwh55aIM1fcwW")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(playlist)
}
