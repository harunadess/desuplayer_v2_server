package fileUtil

import (
	"log"
	"os"
)

func ReadSingleFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Println("error opening file ", path, err)
		return []byte{}, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("error getting file details ", err)
		return []byte{}, err
	}

	bs := make([]byte, fileInfo.Size())
	n, err := file.Read(bs)
	if err != nil {
		log.Println("error reading file ", err)
		return []byte{}, err
	}
	log.Printf("read %v bytes\n", n)

	return bs, nil
}
