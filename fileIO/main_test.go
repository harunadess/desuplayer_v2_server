package fileio

import (
	"testing"
)

func TestAbsPath(t *testing.T) {
	const absPath = "/mnt/d/Users/Jorta/Documents/Coding/go/src/github.com/jordanjohnston/desuplayer_v2/fileio"
	const path = "/library.json"
	want := absPath + path

	if got := AbsPath(path); got != want {
		t.Errorf("Absolute path did not match: got %v, want %v\n", got, want)
	}
}

func TestWriteToJSON(t *testing.T) {
	const baseDir = "/mnt/d/Users/Jorta/Music"
	fileTypesA := []string{".mp3", ".flac"}
	fileTypesB := []string{".MP3", ".FLAC"}

	pathsA, err := ScrapeDirectory(baseDir, fileTypesA)
	if err != nil {
		t.Error(err)
	}

	pathsB, err := ScrapeDirectory(baseDir, fileTypesB)
	if err != nil {
		t.Error(err)
	}

	if len(pathsA) != len(pathsB) {
		t.Error("len(pathsA) != len(pathsB)")
	}
}
