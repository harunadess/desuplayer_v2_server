package library

import (
	"fmt"
	"testing"
)

func skTestBuildLibrary(t *testing.T) {
	const basePath = "/mnt/d/Jorta/Music/iTunes/iTunes Media/Music"
	err := BuildLibrary(basePath)
	if err != nil {
		t.Error(err)
	}
}

func skTestCreatedSortedArtists(t *testing.T) {
	LoadLibrary()
	createSortedArtistList()
	err := SaveLibrary()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAllAlbums(t *testing.T) {
	LoadLibrary()
	albums := GetAllAlbums()
	fmt.Println(len(albums))

	for _, v := range albums {
		if v.Title == "" && v.Artist == "" {
			t.Fatal("property was empty or null", v.Title, v.Artist)
		}
	}
}
