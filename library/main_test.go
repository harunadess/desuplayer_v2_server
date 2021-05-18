package library

import (
	"testing"
)

func TestBuildLibrary(t *testing.T) {
	const basePath = "/mnt/d/Users/Jorta/Music"
	err := BuildLibrary(basePath)
	if err != nil {
		t.Error(err)
	}
}
