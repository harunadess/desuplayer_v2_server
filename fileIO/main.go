package fileio

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var workingDirectory string

type chunk struct {
	bufsize int
	offset  int64
}

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Println("failed to get working directory ", err)
	}
	workingDirectory = wd
}

// AbsPath returns the resolved path based on the current working
// directory the program is running from
func AbsPath(path string) string {
	return workingDirectory + path
}

// ScrapeDirectory recursively finds all files in the base directory
// that have one of the fileTypes specified. It returns a slice of strings
// containing the paths found and nil, or an empty slice and error if there is one.
func ScrapeDirectory(basePath string, fileTypes []string) ([]string, error) {
	files := []string{}

	err := filepath.WalkDir(basePath, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !file.IsDir() {
			for _, v := range fileTypes {
				if strings.EqualFold(filepath.Ext(path), v) {
					files = append(files, path)
					break
				}
			}
		}
		return nil
	})
	return files, err
}

// WriteToJSON writes the specified interface to a JSON file
func WriteToJSON(a interface{}, path string) error {
	JSON, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}

	_, err = os.Lstat(path)
	if err == nil {
		log.Printf("removing file %v\n", path)
		RemoveFile(path)
	}

	err = os.WriteFile(path, JSON, fs.FileMode(os.O_WRONLY))
	if err != nil {
		log.Println("failed to write file ", err)
		return err
	}
	return nil
}

// Remove file removes the specified file from the system
func RemoveFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		log.Println("failed to remove file ", err)
		return err
	}
	return nil
}

// // ReadMusicLibraryFromJSON reads the json file at the specified path and converts it
// // into a directoryscraper.MusicLibrary struct
// func ReadMusicLibraryFromJSON(path string) (directoryscaper.MusicLibrary, error) {
// 	file, err := ReadSingleFile(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return makeLibrary(file)
// }

// func makeLibrary(file []byte) (directoryscaper.MusicLibrary, error) {
// 	library := make(directoryscaper.MusicLibrary)
// 	err := json.Unmarshal(file, &library)
// 	if err != nil {
// 		log.Println("failed to unmarshal json ", err)
// 		return nil, err
// 	}
// 	return library, nil
// }

// ReadSingleFile reads a single file specified by path
// and returns the bytes read
func ReadSingleFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Println("error opening file ", path, err)
		return []byte{}, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("error getting file details ", err)
		return []byte{}, err
	}

	bufferSize := 1024 * 1024 * 1024
	fileSize := int(fileInfo.Size())
	bs := make([]byte, fileSize)

	concurrency := fileSize / bufferSize
	chunkSizes := make([]chunk, concurrency)

	for i := 0; i < concurrency; i++ {
		chunkSizes[i].bufsize = bufferSize
		chunkSizes[i].offset = int64(bufferSize * i)
	}

	if remainder := fileSize % bufferSize; remainder != 0 {
		c := chunk{bufsize: remainder, offset: int64(concurrency * bufferSize)}
		concurrency++
		chunkSizes = append(chunkSizes, c)
	}

	var wg sync.WaitGroup
	wg.Add(concurrency)
	log.Println("num go routines: ", concurrency)

	for i := 0; i < concurrency; i++ {
		go func(chunkSizes []chunk, i int) {
			defer wg.Done()

			chunk := chunkSizes[i]
			buffer := make([]byte, chunk.bufsize)
			bytesRead, err := file.ReadAt(buffer, chunk.offset)
			if err != nil {
				log.Println(err)
				return
			}
			copy(bs[chunk.offset:], buffer[:bytesRead])
		}(chunkSizes, i)
	}

	wg.Wait()

	// old stuff that can't read over ~1GB of data
	// n, err := file.Read(bs)

	// if err != nil {
	// 	log.Println("error reading file ", err)
	// 	return []byte{}, err
	// }
	// log.Printf("read %v bytes of %v\n", n, path)
	return bs, nil
}
