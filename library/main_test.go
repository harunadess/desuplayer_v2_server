package library

import (
	"fmt"
	"testing"
)

func TestBuildLibrary(t *testing.T) {
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

func skTestGetAllAlbums(t *testing.T) {
	albums := GetAllAlbums()
	fmt.Println(len(albums))
}
