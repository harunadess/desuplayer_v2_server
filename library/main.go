package library

import (
	"errors"
	"log"
	"os"

	"github.com/jordanjohnston/desuplayer_v2/directoryscaper"
	"github.com/jordanjohnston/desuplayer_v2/fileio"
	"github.com/jordanjohnston/desuplayer_v2/fileutil"
)

var library directoryscaper.MusicLibrary

func LoadLibrary() {
	wd, err := os.Getwd()
	if err != nil {
		log.Println("failed to get working directory ", err)
		return
	}

	fp := wd + "/library.json"
	data, err := fileio.ReadFromJSON(fp)
	if err != nil {
		log.Println("Error reading library.json ", err)
		return
	}

	library = data
}

// todo: refactor errors
// GetFromLibrary gets a single track from the library
func GetFromLibrary(key string) ([]byte, error) {
	if library == nil {
		return []byte{}, errors.New("library not initialised (library was nil)")
	}

	file, ok := library[key]
	if !ok {
		log.Println("file not found in library (invalid key)")
		return []byte{}, errors.New("file not found in library (invalid key)")
	}

	fileContents, err := fileutil.ReadSingleFile(file.Path)
	if err != nil {
		log.Println("error reading file", err)
		return []byte{}, errors.New("error reading file (io error)")
	}
	return fileContents, nil
}
