package fileIO

import (
	"os"
	"testing"

	"github.com/jordanjohnston/desuplayer_v2/directoryscaper"
)

type MusicFileRaw struct {
	Meta map[string]interface{}
	Path string
}

func TestWriteToJSON(t *testing.T) {
	const baseDir = "/mnt/d/Users/Jorta/Music/Alestorm"
	files, err := directoryscaper.GetAllInDirectory(baseDir)
	if err != nil {
		t.Error(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Error("failed to get working directory ", err)
	}

	fp := wd + "/library.json"

	success := WriteToJSON(files, fp)
	if !success {
		t.Error("failed to write file")
	}
}
