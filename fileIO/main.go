package fileIO

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
)

// WriteToJSON writes the specified interface to a JSON file
func WriteToJSON(a interface{}, fp string) bool {
	JSON, err := json.Marshal(a)
	if err != nil {
		log.Println("failed to marshal json ", err)
		return false
	}

	_, err = os.Lstat(fp)
	if err == nil {
		err = os.Remove(fp)
		if err != nil {
			log.Println("failed to remove file ", err)
		}
	}

	err = os.WriteFile(fp, JSON, fs.FileMode(os.O_WRONLY))
	if err != nil {
		log.Println("failed to write file ", err)
		return false
	}
	return true
}
