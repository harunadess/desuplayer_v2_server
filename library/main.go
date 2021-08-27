package library

import (
	"encoding/json"
	"errors"
	"log"
	"sort"
	"strings"

	"github.com/dhowden/tag"
	"github.com/jordanjohnston/desuplayer_v2/fileio"
	"github.com/jordanjohnston/desuplayer_v2/imageutil"
	"github.com/jordanjohnston/desuplayer_v2/tags"
)

// MusicLibrary holds information about the library
// It contains a map of keys (artistName_albumTitle) -> Album
// It also contains a list of sorted keys (aristName_albumTitle)
type MusicLibrary struct {
	Albums       map[string]Album
	SortedAlbums []string
	Playlists    map[string][]Song
}

type Song struct {
	Title       string
	Artist      string
	Discnumber  int
	Tracknumber int
	Filetype    string
	Path        string
}

type Album struct {
	Title       string
	Artist      string
	Genre       string
	Picturedata []byte
	Picturetype string
	Songs       map[string]Song
}

type Artist struct {
	Name   string
	Albums map[string]Album
}

const (
	newAdd = iota
	existing
)

const imageSize = 200

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
		Albums:    make(map[string]Album),
		Playlists: make(map[string][]Song),
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
			addUnknownToLibrary(path)
			continue
		}
		addBasicToLibrary(metaData, path)
	}

	createSortedArtistList()

	err := SaveLibrary()
	if err != nil {
		log.Println("failed to save library: ", err)
	}

	return nil
}

func getAlbum(key string, meta tag.Metadata) (Album, int) {
	album, ok := library.Albums[key]
	returnCode := existing
	if !ok {
		returnCode = newAdd
		album = Album{
			Picturedata: []byte{},
			Picturetype: "",
			Songs:       make(map[string]Song),
		}
		if meta == nil {
			album.Title = "Unknown"
			album.Genre = "Unknown"
		} else {
			picture := meta.Picture()
			album.Title = meta.Album()
			album.Artist = meta.AlbumArtist()
			if album.Artist == "" {
				album.Artist = meta.Artist()
			}
			album.Genre = meta.Genre()
			if picture != nil {
				album.Picturedata = imageutil.ResizeImage(picture.Data, imageSize)
				album.Picturetype = picture.MIMEType
			}
		}
	}
	return album, returnCode
}

func buildSong(path string, meta tag.Metadata) Song {
	song := Song{
		Path: path,
	}
	if meta == nil {
		pathSplit := strings.Split(path, ".")
		fType := pathSplit[len(pathSplit)-1]
		song.Discnumber = 0
		song.Filetype = strings.ToUpper(fType)
		song.Title = "Unknown"
		song.Tracknumber = 0
	} else {
		disc, _ := meta.Disc()
		track, _ := meta.Track()
		song.Discnumber = disc
		song.Filetype = string(meta.FileType())
		song.Title = meta.Title()
		song.Tracknumber = track
	}
	return song
}

func addUnknownToLibrary(path string) {
	albumKey := "unknown_unknown"
	album, r := getAlbum(albumKey, nil)
	if r == newAdd {
		library.Albums[albumKey] = album
	}
	song := buildSong(path, nil)
	album.Songs[song.Path] = song
}

func addBasicToLibrary(metaData tag.Metadata, path string) {
	albumKey := sanitizeName(metaData.Artist() + "_" + metaData.Album())
	album, r := getAlbum(albumKey, metaData)
	if r == newAdd {
		library.Albums[albumKey] = album
	}
	song := buildSong(path, metaData)
	album.Songs[song.Path] = song
}

func sanitizeName(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func createSortedArtistList() {
	assumedLength := float64(len(library.Albums)) / 1.1
	set := make(map[string]bool, int(assumedLength))
	for k := range library.Albums {
		inserting := removePrefixesFromNames(k)
		set[inserting] = false
	}
	library.SortedAlbums = make([]string, len(set))
	i := 0
	for k := range set {
		library.SortedAlbums[i] = k
		i++
	}

	sort.Slice(library.SortedAlbums, func(i, j int) bool {
		return library.SortedAlbums[i] < library.SortedAlbums[j]
	})
}

func removePrefixesFromNames(s string) string {
	prefixes := []string{"the ", "a "}
	modifiedS := s
	for _, p := range prefixes {
		if len(modifiedS) < len(p) {
			continue
		}

		slice := modifiedS[0:len(p)]
		if slice == p {
			modifiedS = modifiedS[len(p):]
		}
	}

	return modifiedS
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
	library = &MusicLibrary{}

	err = json.Unmarshal(file, library)
	if err != nil {
		log.Println("error unmarshalling json ", err)
	}
}

// AsJson returns the library as JSON
func AsJson() ([]byte, error) {
	return json.Marshal(*library)
}

// GetAllAlbums returns a list of albums, sorted by artist and album title
func GetAllAlbums() []Album {
	if library == nil {
		return make([]Album, 0)
	}
	albums := make([]Album, len(library.Albums))
	i := 0
	for _, v := range library.SortedAlbums {
		albums[i] = library.Albums[v]
		i++
	}

	return albums
}
