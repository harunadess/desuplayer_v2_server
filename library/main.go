package library

import (
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

func UnloadLibrary() {
	library = nil
}

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

func SaveLibrary() error {
	return fileio.WriteToJSON(*library, fileio.AbsPath("/library.json"))
}

// todo: this func
func LoadLibrary() {
	// data, err := fileio.ReadMusicLibraryFromJSON(fileio.AbsPath("/library.json"))
	// if err != nil {
	// 	return
	// }

	// library = data
}

// func SetLibrary(lib directoryscaper.MusicLibrary) {
// 	library = lib
// }

// func GetLibrary() directoryscaper.MusicLibrary {
// 	return library
// }

// // todo: refactor errors
// // GetFromLibrary gets a single track from the library
// func GetFromLibrary(key string) ([]byte, error) {
// 	if library == nil {
// 		return []byte{}, errors.New("library not initialised (no library file found)")
// 	}

// 	file, ok := library[key]
// 	if !ok {
// 		log.Println("file not found in library (invalid key)")
// 		return []byte{}, errors.New("file not found in library")
// 	}

// 	fileContents, err := fileio.ReadSingleFile(file.Path)
// 	if err != nil {
// 		log.Println("error reading file (io error)", err)
// 		return []byte{}, errors.New("error reading file")
// 	}
// 	return fileContents, nil
// }
