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
	youtubeService := NewYouTubeService(a.Config)
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
	youtubeService := NewYouTubeService(a.Config)
	playlist, err := youtubeService.GetYoutubePlaylistByID("https://youtube.com/playlist?list=PLinB3MHWKsOORBP1pMnDLh_arJrmq-ieS&si=eSAVwwTjUvoZOZGu")
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Playlist title:", playlist.Title, "\n Playlist Owner: ", playlist.Owner, "\n num of media ", len(playlist.Songs))
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
	youtubeService := NewYouTubeService(a.Config)
	song, err := youtubeService.GetYoutubeMediaByMetadata(domain.MetaData{Title: "close eyes ", Artist: "dvrst"})
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Song Title:", song.Title, "\n Artist:", song.Artist, "\n SongId:", song.Link)
}

// TestServiceYouTube_FillYoutubePlaylist - OK
func TestServiceYouTube_FillYoutubePlaylist(t *testing.T) {
	a := app.New("test")
	a.InitValidator()
	a.PopulateConfig()
	youtubeService := NewYouTubeService(a.Config)
	acessToken := "4/0AVG7fiT4m7tJnQE3hjmvBBXqDzr1ZXHyndO8dI1TIJ4ErJiHY3DEgifhZVazRfLG4deoqQ"
	songs := make([]domain.Song, 0)
	songs = append(songs, domain.Song{Title: "Sonne", Artist: "Rammstein"})
	playlist, err := youtubeService.CreateYoutubePlaylist(domain.Playlist{Title: "Example", Songs: songs}, acessToken)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(playlist)

}
