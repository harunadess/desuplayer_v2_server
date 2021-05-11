package fileio

import (
	"fmt"
	"os"
	"testing"

	"github.com/jordanjohnston/desuplayer_v2/directoryscaper"
)

func TestWriteToJSON(t *testing.T) {
	const baseDir = "/mnt/d/Users/Jorta/Music"
	files, err := directoryscaper.GetAllInDirectory(baseDir)
	if err != nil {
		t.Error(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Error("failed to get working directory ", err)
	}

	fp := wd + "/library.json"

	fmt.Println("num files", len(files))

	success := WriteToJSON(files, fp)
	if !success {
		t.Error("failed to write file")
	}
}

func TestReadFromJSON(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error("failed to get working directory ", err)
	}

	fp := wd + "/library.json"
	data, err := ReadFromJSON(fp)
	if err != nil {
		t.Error(err)
		return
	}

	if data == nil {
		t.Error("library was nil")
		return
	}

	f := data["000d05f0-cf8a-4267-abfc-55d50e75f90b"]
	fmt.Println(f.Album)
	fmt.Println(f.AlbumArtist)
	fmt.Println(f.Artist)
	fmt.Println(f.Composer)
	fmt.Println(f.DiscNumber)
	fmt.Println(f.FileType)
	fmt.Println(f.Format)
	fmt.Println(f.Genre)
}
