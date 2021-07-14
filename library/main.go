package library

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"strings"

	"github.com/dhowden/tag"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/jordanjohnston/desuplayer_v2/fileio"
	"github.com/jordanjohnston/desuplayer_v2/tags"
)

// As a library
// We want to store limited information
// i.e. we don't want to have to re-store like 10000000000 MB of data
// So, we should have a full library :tm: and a basic library
// basic library just contains the SongMap

// When we load in the full library, we then create a full library
// by running some funcs to fill out the rest of the things
// since that is basically a bunch of duplication otherwise

// As an artist, we want to have a list of albums/songs
// As an album, we want to have the title, artist, year, album art
// maybe this would get around the mass fuckery
// of not being able to do anything with images

// MusicLibrary holds an array of artists
// These artists have >= 1 Album
// These albums contain several songs
// These songs are each a file
// Playlists are a list of songs (keys)

// this ain't gonna work because they're slices
// they should be maps of some identifier to a list
// not just a list..
type MusicLibrary struct {
	Artists   map[string]Artist
	Playlists map[string][]Song
}

type Song struct {
	Title       string
	Discnumber  int
	Tracknumber int
	Filetype    string
	Path        string
	Key         string
}

type Album struct {
	Title       string
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

const imageSize = 300

// Hierarchy of Generality of Songs
// Genre -> Artist -> Album -> Song -> SongFile
// idk if we give a shit about genre tbh
// other than just meta of the actual album

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
		Artists:   make(map[string]Artist),
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

func getArtist(key string, meta tag.Metadata) (Artist, int) {
	artist, ok := library.Artists[key]
	returnCode := existing
	if !ok {
		returnCode = newAdd
		artist = Artist{Albums: make(map[string]Album)}
		if meta == nil {
			artist.Name = "Unknown"
		} else {
			artist.Name = meta.Artist()
		}
	}
	return artist, returnCode
}

func getAlbum(artist Artist, key string, meta tag.Metadata) (Album, int) {
	album, ok := artist.Albums[key]
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
			album.Genre = meta.Genre()
			if picture != nil {
				album.Picturedata = resizePicture(picture.Data)
				album.Picturetype = picture.MIMEType
			}
		}
	}
	return album, returnCode
}

// todo: break this out into a separate file / helper
// todo: handle non-square images
func resizePicture(pictureData []byte) []byte {
	reader := bytes.NewReader(pictureData)
	img, format, err := image.Decode(reader)
	if err != nil {
		log.Println("failed to decode base64 string", err, format)
		return []byte{}
	}
	resizedImg := imaging.Resize(img, imageSize, imageSize, imaging.CatmullRom)

	var resizedBytes bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&resizedBytes, resizedImg, nil)
	case "png":
		err = png.Encode(&resizedBytes, resizedImg)
	case "gif":
		err = gif.Encode(&resizedBytes, resizedImg, nil)
	default:
		err = errors.New("unrecognised image format")
	}
	if err != nil {
		log.Println("failed to encode image", err)
		return []byte{}
	}

	return resizedBytes.Bytes()
}

func buildSong(path string, meta tag.Metadata) Song {
	song := Song{
		Path: path,
		Key:  uuid.NewString(),
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
	artistKey := "unknown"
	albumKey := "unknown_unknown"
	artist, r := getArtist(artistKey, nil)
	if r == newAdd {
		library.Artists[artistKey] = artist
	}
	album, r := getAlbum(artist, albumKey, nil)
	if r == newAdd {
		artist.Albums[albumKey] = album
	}
	song := buildSong(path, nil)
	album.Songs[song.Key] = song
}

func addBasicToLibrary(metaData tag.Metadata, path string) {
	artistKey := sanitizeName(metaData.Artist())
	albumKey := sanitizeName(metaData.Artist() + "_" + metaData.Album())
	artist, r := getArtist(artistKey, metaData)
	if r == newAdd {
		library.Artists[artistKey] = artist
	}
	album, r := getAlbum(artist, albumKey, metaData)
	if r == newAdd {
		artist.Albums[albumKey] = album
	}
	song := buildSong(path, metaData)
	album.Songs[song.Key] = song
}

func sanitizeName(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
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

func GetAllArtists() map[string]Artist {
	if library == nil {
		return make(map[string]Artist)
	}
	return library.Artists
}
