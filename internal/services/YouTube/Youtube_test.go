package YouTube

import (
	"MusicApp/internal/app"
	"MusicApp/internal/domain"
	"fmt"
	"testing"
)

// TestServiceYouTube_GetYoutubeMediaByID -  OK
func TestServiceYouTube_GetYoutubeMediaByID(t *testing.T) {
	a := app.New("test")
	a.InitValidator()
	a.PopulateConfig()
	youtubeService := New(a.Config)
	song, err := youtubeService.GetYoutubeMediaByID("https://www.youtube.com/watch?v=hTWKbfoikeg")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("NameOfSong: ", song.Title)
	fmt.Println("ArtistName: ", song.Artist)
	fmt.Println("Link ", song.Link)

}

// TestServiceYouTube_GetYoutubePlaylistByID - OK
func TestServiceYouTube_GetYoutubePlaylistByID(t *testing.T) {
	a := app.New("test")
	a.InitValidator()
	a.PopulateConfig()
	youtubeService := New(a.Config)
	playlist, err := youtubeService.GetYoutubePlaylistByID("https://youtube.com/playlist?list=PLbTTxxr-hMmzSGTsO5mdYLrvKY-RZFanp&si=sB4Nwh55aIM1fcwW")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Playlist title:", playlist.Title, "\n Playlist Owner: ", playlist.Owner)
	for _, song := range playlist.Songs {
		fmt.Println("Song Title:", song.Title)
		fmt.Println("Song Artist:", song.Artist)
	}
}

// TestServiceYouTube_GetYoutubeMediaByMetadata - OK
func TestServiceYouTube_GetYoutubeMediaByMetadata(t *testing.T) {
	a := app.New("test")
	a.InitValidator()
	a.PopulateConfig()
	youtubeService := New(a.Config)
	song, err := youtubeService.GetYoutubeMediaByMetadata(domain.MetaData{Title: "close eyes ", Artist: "dvrst"})
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Song Title:", song.Title, "\n Artist:", song.Artist, "\n SongId:", song.Link)
}
func TestServiceYouTube_FillYoutubePlaylist(t *testing.T) {}
