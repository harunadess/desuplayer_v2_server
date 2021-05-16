package fileio

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"

	"github.com/jordanjohnston/desuplayer_v2/directoryscaper"
)

var workingDirectory string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Println("failed to get working directory ", err)
	}
	workingDirectory = wd
}

func absPath(fp string) string {
	return workingDirectory + fp
}

// WriteToJSON writes the specified interface to a JSON file
func WriteToJSON(a interface{}, path string) bool {
	JSON, err := json.Marshal(a)
	if err != nil {
		log.Println("failed to marshal json ", err)
		return false
	}

	_, err = os.Lstat(path)
	if err == nil {
		RemoveFile(path)
	}

	err = os.WriteFile(absPath(path), JSON, fs.FileMode(os.O_WRONLY))
	if err != nil {
		log.Println("failed to write file ", err)
		return false
	}
	return true
}

// Remove file removes the specified file from the system
func RemoveFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		log.Println("failed to remove file ", err)
		return err
	}
	return nil
}

// ReadMusicLibraryFromJSON reads the json file at the specified path and converts it
// into a directoryscraper.MusicLibrary struct
func ReadMusicLibraryFromJSON(path string) (directoryscaper.MusicLibrary, error) {
	file, err := ReadSingleFile(absPath(path))
	if err != nil {
		return nil, err
	}
	return makeLibrary(file)
}

func makeLibrary(file []byte) (directoryscaper.MusicLibrary, error) {
	library := make(directoryscaper.MusicLibrary)
	err := json.Unmarshal(file, &library)
	if err != nil {
		log.Println("failed to unmarshal json ", err)
		return nil, err
	}
	return library, nil
}

// ReadSingleFile reads a single file specified by path
// and returns the bytes read
func ReadSingleFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Println("error opening file ", path, err)
		return []byte{}, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("error getting file details ", err)
		return []byte{}, err
	}

	bs := make([]byte, fileInfo.Size())
	n, err := file.Read(bs)
	if err != nil {
		log.Println("error reading file ", err)
		return []byte{}, err
	}
	log.Printf("read %v bytes\n", n)

	return bs, nil
}
