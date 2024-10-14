package Youtube

import (
	"MusicApp/internal/domain"
)

func (h Handler) GetYoutubeMediaByLink(youtubeLink string) (*domain.Song, error) {

	return h.youtubeService.GetYoutubeMediaByLink(youtubeLink)

}
func (h Handler) GetYoutubeMediaByMetaData(metadata *domain.MetaData) (*domain.Song, error) {
	return h.youtubeService.GetYoutubeMediaByMetadata(domain.MetaData{Title: metadata.Title, Artist: metadata.Artist})
}
func (h Handler) GetYoutubePlaylistByLink(youtubeLink string) (*domain.Playlist, error) {
	return h.youtubeService.GetYoutubePlaylistDataByLink(youtubeLink)
}

func (h Handler) CreateAndFillYoutubePlaylist(Title string, AuthToken string) (string, error) {
	return h.youtubeService.CreateYoutubePlaylist(Title, AuthToken)
}
func (h Handler) FillYouTubePlaylist(playlist domain.Playlist, YoutubePlaylistId, AuthToken string) (*domain.Playlist, error) {
	return h.youtubeService.FillYoutubePlaylist(playlist, YoutubePlaylistId, AuthToken)
}
