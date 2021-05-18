package library

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/dhowden/tag"
	"github.com/google/uuid"
	"github.com/jordanjohnston/desuplayer_v2/fileio"
	"github.com/jordanjohnston/desuplayer_v2/tags"
)

// MusicLibrary is a structure holding the different types of mappings
// to individual music files
type MusicLibrary struct {
	Artists   ArtistMap
	Albums    AlbumMap
	Genres    GenreMap
	Playlists PlaylistMap
	Songs     SongMap
}

// ArtistMap is a map of an Artist (string) to a list of MusicFile uuids
type ArtistMap map[string][]string

// Album is a struct that represents an album
// Contains Artist and Title
type Album struct {
	Artist string
	Title  string
}

// AlbumMap is a map of an Album to a list of MusicFile uuids
type AlbumMap map[string][]string

// GenreMap is a map of a Genre (string) to a list of MusicFile uuids
type GenreMap map[string][]string

// PlaylistMap is a map of a Playlist to a list of MusicFile uuids
type PlaylistMap map[string][]string

// SongMap is a map of a MusicFile uuid to a BasicMusicFile
type SongMap map[string]BasicMusicFile

// BasicMusicFile is a basic representation of a music file
// and some of the more important attributes
type BasicMusicFile struct {
	Album       string
	Artist      string
	AlbumArtist string
	DiscNumber  int
	FileType    string
	Genre       string
	Title       string
	TrackNumber int
	Path        string
	Key         string
}

// MusicFile is a more complete representation of a music file
type MusicFile struct {
	BasicMusicFile
	Composer        string
	FileType        string
	Format          string
	Lyrics          string
	PictureData     []byte
	PictureMIMEType string
	TotalTracks     int
	Year            int
}

// todo: there is no wav, need to fork taglib and write yours

var library *MusicLibrary

// UnloadLibrary clears the current library
func UnloadLibrary() {
	library = nil
}

// BuildLibrary creates the library using the basePath
// passed in.
func BuildLibrary(basePath string) error {
	log.Println("scraping directory:", basePath)
	paths, err := fileio.ScrapeDirectory(basePath, tags.AcceptableFileTypes)
	if err != nil {
		return errors.New("Error scraping directory: " + err.Error())
	}
	log.Printf("got %v files\n", len(paths))

	library = &MusicLibrary{
		Artists:   make(ArtistMap),
		Albums:    make(AlbumMap),
		Genres:    make(GenreMap),
		Playlists: make(PlaylistMap),
		Songs:     make(SongMap),
	}
	err = fillLibrary(paths)
	if err != nil {
		return errors.New("Error filling library: " + err.Error())
	}

	return nil
}

func fillLibrary(paths []string) error {
	for _, path := range paths {
		metaData, err := tags.ReadTags(path)
		if err != nil {
			log.Printf("failed to read file %v\n", path)
			addUnknownToLibrary(path)
			continue
		}
		addBasicToLibrary(metaData, path)
	}

	err := SaveLibrary()
	if err != nil {
		log.Println("failed to save library: ", err)
	}

	return nil
}

func addUnknownToLibrary(path string) {
	basicFile := BasicMusicFile{
		Album:       "Unknown",
		Artist:      "Unknown",
		AlbumArtist: "Unknown",
		DiscNumber:  -1,
		FileType:    "Unknown", //todo: get filetype from ext
		Genre:       "Unknown",
		Title:       "Unknown",
		TrackNumber: -1,
		Path:        path,
		Key:         uuid.NewString(),
	}
	album := (basicFile.Artist + "_" + basicFile.Album)
	library.Artists[basicFile.Artist] = append(library.Artists[basicFile.Artist], basicFile.Key)
	library.Albums[album] = append(library.Albums[album], basicFile.Key)
	library.Artists[basicFile.Artist] = append(library.Artists[basicFile.Artist], basicFile.Key)
	library.Genres[basicFile.Genre] = append(library.Genres[basicFile.Genre], basicFile.Key)
	library.Songs[basicFile.Key] = basicFile
}

func addBasicToLibrary(metaData tag.Metadata, path string) {
	disc, _ := metaData.Disc()
	track, _ := metaData.Track()
	basicFile := BasicMusicFile{
		Album:       metaData.Album(),
		Artist:      metaData.Artist(),
		AlbumArtist: metaData.AlbumArtist(),
		DiscNumber:  disc,
		FileType:    string(metaData.FileType()),
		Genre:       metaData.Genre(),
		Title:       metaData.Title(),
		TrackNumber: track,
		Path:        path,
		Key:         uuid.NewString(),
	}
	if basicFile.Album == "" {
		basicFile.Album = "Unknown"
	}
	if basicFile.Artist == "" {
		basicFile.Artist = "Unknown"
	}
	if basicFile.AlbumArtist == "" {
		basicFile.AlbumArtist = ""
	}
	if basicFile.Genre == "" {
		basicFile.Genre = "Unknown"
	}

	album := (basicFile.Artist + "_" + basicFile.Album)
	library.Artists[basicFile.Artist] = append(library.Artists[basicFile.Artist], basicFile.Key)
	library.Albums[album] = append(library.Albums[album], basicFile.Key)
	library.Artists[basicFile.Artist] = append(library.Artists[basicFile.Artist], basicFile.Key)
	library.Genres[basicFile.Genre] = append(library.Genres[basicFile.Genre], basicFile.Key)
	library.Songs[basicFile.Key] = basicFile
}

// SaveLibrary saves the library as JSON
func SaveLibrary() error {
	return fileio.WriteToJSON(*library, fileio.AbsPath("/library.json"))
}

// LoadLibrary loads the library from it's JSON file.
func LoadLibrary() {
	file, err := fileio.ReadSingleFile(fileio.AbsPath("/library.json"))
	if err != nil {
		log.Println("error reading library file ", err)
		return
	}
	library = &MusicLibrary{
		Artists:   make(ArtistMap),
		Albums:    make(AlbumMap),
		Genres:    make(GenreMap),
		Playlists: make(PlaylistMap),
		Songs:     make(SongMap),
	}

	err = json.Unmarshal(file, library)
	if err != nil {
		log.Println("error unmarshalling json ", err)
	}
}

// AsJson returns the library as JSON
func AsJson() ([]byte, error) {
	return json.Marshal(*library)
}

func GetSong(key string) (BasicMusicFile, bool) {
	song, ok := library.Songs[key]
	return song, ok
}

func GetArtist(key string) ([]BasicMusicFile, bool) {
	artist, ok := library.Artists[key]
	if !ok {
		return []BasicMusicFile{}, ok
	}
	songs := make([]BasicMusicFile, len(artist))
	for i, v := range artist {
		songs[i] = library.Songs[v]
	}
	return songs, ok
}

func GetAlbum(key string) ([]BasicMusicFile, bool) {
	album, ok := library.Albums[key]
	if !ok {
		return []BasicMusicFile{}, ok
	}
	songs := make([]BasicMusicFile, len(album))
	for i, v := range album {
		songs[i] = library.Songs[v]
	}
	return songs, ok
}

func GetGenre(key string) ([]BasicMusicFile, bool) {
	genre, ok := library.Genres[key]
	if !ok {
		return []BasicMusicFile{}, ok
	}
	songs := make([]BasicMusicFile, len(genre))
	for i, v := range genre {
		songs[i] = library.Songs[v]
	}
	return songs, ok
}

func GetAllArtists() ([]string, bool) {
	artists := make([]string, len(library.Artists))
	i := 0
	for k := range library.Artists {
		artists[i] = k
		i++
	}
	return artists, true
}

func GetAllAlbums() ([]string, bool) {
	albums := make([]string, len(library.Albums))
	i := 0
	for k := range library.Albums {
		albums[i] = k
		i++
	}
	return albums, true
}

func GetAllGenres() ([]string, bool) {
	genres := make([]string, len(library.Genres))
	i := 0
	for k := range library.Genres {
		genres[i] = k
		i++
	}
	return genres, true
}
