package library

import (
	"testing"
)

func TestBuildLibrary(t *testing.T) {
	const basePath = "/mnt/d/Jorta/Music/iTunes/iTunes Media/Music/Suisei Hoshimachi"
	err := BuildLibrary(basePath)
	if err != nil {
		t.Error(err)
	}
}
