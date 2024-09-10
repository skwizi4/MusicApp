package handlers

// todo - Написать структуры: MetaData, Song, Playlist ( directory - domain)

type (
	Spotify interface {
		GetSongByYoutubeLink(youtubeLink string) (*domain.Song, error)
		GetSongByMetaData(metadata domain.MetaData) (*domain.Song, error)
		GetPlaylistByYoutubeLink(youtubeLink string) (*domain.Playlist, error)
		GetPlaylistByMetaData(metadata domain.MetaData) (*domain.Playlist, error)
	}

	YouTube interface {
		GetSongBySpotifyLink(spotifyLink string) (*domain.Song, error)
		GetSongByMetaData(metadata domain.MetaData) (*domain.Song, error)
		GetPlaylistBySpotifyLink(spotifyLink string) (*domain.Playlist, error)
		GetPlaylistByMetaData(metadata domain.MetaData) (*domain.Playlist, error)
	}
)
