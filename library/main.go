package library

import (
	"errors"
	"log"

	"github.com/jordanjohnston/desuplayer_v2/directoryscaper"
	"github.com/jordanjohnston/desuplayer_v2/fileio"
)

var library directoryscaper.MusicLibrary

func LoadLibrary() {
	data, err := fileio.ReadMusicLibraryFromJSON("/library.json")
	if err != nil {
		return
	}

	library = data
}

// todo: refactor errors
// GetFromLibrary gets a single track from the library
func GetFromLibrary(key string) ([]byte, error) {
	if library == nil {
		return []byte{}, errors.New("library not initialised (no library file found)")
	}

	file, ok := library[key]
	if !ok {
		log.Println("file not found in library (invalid key)")
		return []byte{}, errors.New("file not found in library")
	}

	fileContents, err := fileio.ReadSingleFile(file.Path)
	if err != nil {
		log.Println("error reading file (io error)", err)
		return []byte{}, errors.New("error reading file")
	}
	return fileContents, nil
}
