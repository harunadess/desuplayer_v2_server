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

// WriteToJSON writes the specified interface to a JSON file
func WriteToJSON(a interface{}, fp string) bool {
	JSON, err := json.Marshal(a)
	if err != nil {
		log.Println("failed to marshal json ", err)
		return false
	}

	_, err = os.Lstat(fp)
	if err == nil {
		err = os.Remove(fp)
		if err != nil {
			log.Println("failed to remove file ", err)
		}
	}

	absFp := workingDirectory + fp
	err = os.WriteFile(absFp, JSON, fs.FileMode(os.O_WRONLY))
	if err != nil {
		log.Println("failed to write file ", err)
		return false
	}
	return true
}

func ReadFromJSON(fp string) (directoryscaper.MusicLibrary, error) {
	absFp := workingDirectory + fp
	file, err := os.ReadFile(absFp)
	if err != nil {
		log.Println("failed to read file ", err)
		return nil, err
	}

	JSONData := make(directoryscaper.MusicLibrary)
	err = json.Unmarshal(file, &JSONData)
	if err != nil {
		log.Println("failed to unmarshal json ", err)
		return nil, err
	}

	return JSONData, nil
}
