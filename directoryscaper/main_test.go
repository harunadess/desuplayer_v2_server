package directoryscaper

import (
	"fmt"
	"testing"
)

func TestGetAllInDirectory(t *testing.T) {
	const baseDir = "/mnt/d/Users/Jorta/Music"
	files, err := GetAllInDirectory(baseDir)
	if err != nil {
		t.Error(err)
		return
	}
	if len(files) < 1 {
		t.Error("Expected to find some files, found", len(files))
	}

	first := files[0]
	fmt.Printf("len=%v, cap=%v\n", len(files), cap(files))
	fmt.Printf("first file %v\n", first)
}
