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

func skTestGetAllAlbums(t *testing.T) {
	albums := GetAllAlbums()
	fmt.Println(len(albums))
}

/*
	{
  "Title": "Daydream maiden brain concept",
  "Artist": "Various Artists",
  "Discnumber": 0,
  "Tracknumber": 3,
  "Filetype": "MP3",
  "Path": "D:\\Users\\Jorta\\Music\\Yuzuki Yukari\\月の詩 V - ツキノウタ\\03. デイドリィム乙女脳内構想.mp3"
}
*/

func TestGetSongMeta(t *testing.T) {
	const libraryFilePath = "/mnt/d/Users/Jorta/Documents/Coding/go/src/github.com/jordanjohnston/desuplayer_v2/library.json"
	const pathToSong = "D:\\Users\\Jorta\\Music\\Yuzuki Yukari\\月の詩 V - ツキノウタ\\03. デイドリィム乙女脳内構想.mp3"
	const songArtist = "Various Artists"
	const songAlbumTitle = "Tsuki no Uta V"
	loadLibraryForTest(libraryFilePath)
	meta, ok := GetSongMeta(pathToSong, songAlbumTitle, songArtist)
	if !ok {
		t.Fatal("Did not get meta - failed")
	}
	fmt.Println(meta)
}
