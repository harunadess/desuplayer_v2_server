package fileio

import (
	"fmt"
	"testing"

	"github.com/jordanjohnston/desuplayer_v2/directoryscaper"
)

func TestWriteToJSON(t *testing.T) {
	const baseDir = "/mnt/d/Jorta/Music"
	files, err := directoryscaper.GetAllInDirectory(baseDir, true)
	if err != nil {
		t.Error(err)
	}

	fp := "/library.json"

	fmt.Println("num files", len(files))

	success := WriteToJSON(files, AbsPath(fp))
	if !success {
		t.Error("failed to write file")
	}
}

func TestReadFromJSON(t *testing.T) {
	fp := "../library.json"
	data, err := ReadMusicLibraryFromJSON(fp)
	if err != nil {
		t.Error(err)
		return
	}

	if data == nil {
		t.Error("library was nil")
		return
	}

	for k, v := range data {
		fmt.Println(k, v)
		break
	}
}
