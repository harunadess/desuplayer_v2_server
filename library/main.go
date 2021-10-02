package library

import (
	"encoding/json"
	"errors"
	"log"
	"path/filepath"
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
	AlbumTitle  string
	AlbumArtist string
	Artist      string
	Discnumber  int
	Tracknumber int
	Filetype    string
	Path        string
}

type Album struct {
	Title       string
	AlbumArtist string
	Artist      string
	Genre       string
	Picturedata []byte
	Picturetype string
	Songs       map[string]Song
	AlbumKey    string
}

type Artist struct {
	Name   string
	Albums map[string]Album
}

type SongMeta struct {
	Song
	Genre       string
	Picturedata []byte
	Picturetype string
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
			log.Printf("failed to read meta for: %v\n", path)
			// addUnknownToLibrary(path)
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
			AlbumKey:    key,
		}
		if meta == nil {
			album.Title = "Unknown"
			album.Genre = "Unknown"
		} else {
			picture := meta.Picture()
			album.Title = meta.Album()
			album.Artist = meta.AlbumArtist()
			album.AlbumArtist = meta.AlbumArtist()
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
		song.Artist = "Unknown"
		song.AlbumArtist = "Unknown"
		song.AlbumTitle = "Unknown"
		song.Tracknumber = 0
	} else {
		disc, _ := meta.Disc()
		track, _ := meta.Track()
		song.Discnumber = disc
		song.Filetype = string(meta.FileType())
		song.Title = meta.Title()
		song.Artist = meta.Artist()
		song.AlbumArtist = meta.AlbumArtist()
		song.AlbumTitle = meta.Album()
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
	artistName := metaData.AlbumArtist()
	if artistName == "" {
		artistName = metaData.Artist()
	}
	albumKey := sanitizeName(artistName + "_" + metaData.Album())
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
	library.SortedAlbums = make([]string, len(library.Albums))
	i := 0
	for k := range library.Albums {
		library.SortedAlbums[i] = removePrefixesFromNames(k)
		i++
	}
	sort.Strings(library.SortedAlbums)
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

func loadLibraryForTest(path string) {
	file, err := fileio.ReadSingleFile(path)
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

// GetSong returns the bytes of the song at the specified path, or an empty byte slice if not found
func GetSong(path string) ([]byte, bool) {
	found := false
	for _, v := range tags.AcceptableFileTypes {
		if found = strings.EqualFold(filepath.Ext(path), v); found {
			break
		}
	}

	if !found {
		log.Println("path was not of a valid type: ", path)
		return make([]byte, 0), false
	}

	bs, err := fileio.ReadSingleFile(path)
	if err != nil {
		log.Println("error reading file ", err)
		return bs, false
	}
	return bs, true
}

func GetSongMeta(path string, albumTitle string, albumArtist string) (SongMeta, bool) {
	found := false
	for _, v := range tags.AcceptableFileTypes {
		if found = strings.EqualFold(filepath.Ext(path), v); found {
			break
		}
	}

	if !found {
		log.Println("path was not of a valid type: ", path)
		return SongMeta{}, false
	}

	albumKey := sanitizeName(albumArtist + "_" + albumTitle)
	album, ok := library.Albums[albumKey]
	if !ok {
		log.Println("did not find album for album key: ", albumKey, albumArtist, albumTitle)
		return SongMeta{}, false
	}

	song, ok := album.Songs[path]
	if !ok {
		log.Println("did not find song for path: ", path)
		return SongMeta{}, false
	}

	meta := SongMeta{
		Song:        song,
		Genre:       album.Genre,
		Picturetype: album.Picturetype,
		Picturedata: album.Picturedata,
	}

	return meta, true
}
