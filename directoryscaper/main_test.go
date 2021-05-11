package directoryscaper

import (
	"fmt"
	"testing"
)

func TestGetAllInDirectory(t *testing.T) {
	const baseDir string = "/mnt/d/Users/Jorta/Music"
	files, err := GetAllInDirectory(baseDir)
	if err != nil {
		t.Error(err)
		return
	}
	if len(files) < 1 {
		t.Error("Expected to find some files, found", len(files))
	}

	fmt.Printf("len=%v\n", len(files))
}

func TestGetAllInDirectoryInline(t *testing.T) {
	const baseDir string = "/mnt/d/Users/Jorta/Music"
	files, err := GetAllInDirectoryInline(baseDir)
	if err != nil {
		t.Error(err)
		return
	}
	if len(files) < 1 {
		t.Error("Expected to find some files, found", len(files))
	}

	fmt.Printf("len=%v\n", len(files))
}

func BenchmarkGetAllInDirectory(b *testing.B) {
	const baseDir string = "/mnt/d/Users/Jorta/Music"
	for i := 0; i < b.N; i++ {
		files, err := GetAllInDirectory(baseDir)
		if err != nil {
			b.Error(err)
			return
		}
		if len(files) < 1 {
			b.Error("Expected to find some files, found", len(files))
		}

		fmt.Printf("len=%v\n", len(files))
	}
}
