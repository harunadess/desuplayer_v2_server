package fileUtil

import (
	"testing"
)

func TestGetSingleFile(t *testing.T) {
	filePath := "/mnt/d/Users/Jorta/Music/Alestorm/Back Through Time (Limited Edition)/01 - Back Through Time.mp3"
	file, err := ReadSingleFile(filePath)
	if err != nil {
		t.Error(err)
	}

	if len(file) <= 0 {
		t.Error("non-error issue reading file")
	}
}
