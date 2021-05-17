package directoryscaper

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
	"github.com/google/uuid"
)

// todo: there is no wav, need to fork taglib and write yours
var acceptableFileTypes []string = []string{".mp3", ".flac", ".m4a" /*".wav",*/, ".ogg", ".alac", ".m4p"}

// MusicLibrary is a map of a MusicFile's uuid to itself
type MusicLibrary map[string]MusicFile

// Struct to hold information about music files.
type MusicFile struct {
	Album           string
	AlbumArtist     string
	Artist          string
	Composer        string
	DiscNumber      int
	TotalDiscs      int
	FileType        string
	Format          string
	Genre           string
	Lyrics          string
	PictureData     []byte
	PictureMIMEType string
	Title           string
	TrackNumber     int
	TotalTracks     int
	Year            int
	Path            string
	Key             string
}

// PictureData and PictureMIMEType do not get populated
// this makes the server use too much memory and the library huge
// to avoid this we can just retrieve the picture before showing it as now playing or something
func buildMusicFile(meta tag.Metadata, path string, getPics bool) MusicFile {
	discNum, totalDiscs := meta.Disc()
	trackNum, totalTracks := meta.Track()

	musicFile := MusicFile{
		Album:       meta.Album(),
		AlbumArtist: meta.AlbumArtist(),
		Artist:      meta.Artist(),
		Composer:    meta.Composer(),
		DiscNumber:  discNum,
		TotalDiscs:  totalDiscs,
		FileType:    string(meta.FileType()),
		Format:      string(meta.Format()),
		Genre:       meta.Genre(),
		Title:       meta.Title(),
		TrackNumber: trackNum,
		TotalTracks: totalTracks,
		Year:        meta.Year(),
		Path:        path,
		Key:         uuid.New().String(),
	}

	picture := meta.Picture()
	if getPics && picture != nil {
		musicFile.PictureData = picture.Data
		musicFile.PictureMIMEType = picture.MIMEType
	}

	return musicFile
}

// GetAllInDirectory gets all music files in the specified base directory
// note: this is extremely slow on wsl
func GetAllInDirectory(baseDir string, getImages bool) (MusicLibrary, error) {
	if getImages {
		return getAllInDirectoryWithImages(baseDir)
	}
	return getAllInDirectoryNoImages(baseDir)
}

func getAllInDirectoryWithImages(baseDir string) (MusicLibrary, error) {
	log.Println("---- start of getAllInDirectoryWithImages")
	log.Println(baseDir)
	filePaths, err := getAllMusicFilePaths(baseDir)
	if err != nil {
		log.Println("error reading files in directory ", err)
	}

	musicLibrary := make(MusicLibrary)

	for _, path := range filePaths {
		meta := readFileMetaData(path)
		if meta != nil {
			mf := buildMusicFile(meta, path, true)
			musicLibrary[mf.Key] = mf
		} else {
			// we don't have any meta data :c
			mf := MusicFile{
				Path: path,
				Key:  uuid.New().String(),
			}
			musicLibrary[mf.Key] = mf
		}

	}

	log.Println("---- end of getAllInDirectoryWithImages")
	return musicLibrary, nil
}

func getAllMusicFilePaths(baseDir string) ([]string, error) {
	files := []string{}

	err := filepath.WalkDir(baseDir, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !file.IsDir() {
			for _, v := range acceptableFileTypes {
				if filepath.Ext(path) == v {
					files = append(files, path)
					break
				}
			}
		}
		return nil
	})
	return files, err
}

func readFileMetaData(path string) tag.Metadata {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("failed to read file %v: %v", path, err)
		return nil
	}
	defer file.Close()
	meta, err := tag.ReadFrom(file)
	if err != nil {
		log.Printf("error reading meta data: %v, %v", err, path)
		return nil
	}
	return meta
}

func getAllInDirectoryNoImages(baseDir string) (MusicLibrary, error) {
	log.Println("---- start of getAllInDirectoryNoPhotos")
	log.Println(baseDir)
	filePaths, err := getAllMusicFilePaths(baseDir)
	if err != nil {
		log.Println("error reading files in directory ", err)
	}

	musicLibrary := make(MusicLibrary)

	for _, path := range filePaths {
		meta := readFileMetaData(path)
		if meta != nil {
			mf := buildMusicFile(meta, path, false)
			musicLibrary[mf.Key] = mf
		} else {
			// we don't have any meta data :c
			mf := MusicFile{
				Path: path,
				Key:  uuid.New().String(),
			}
			musicLibrary[mf.Key] = mf
		}
	}

	log.Println("---- end of getAllInDirectoryNoPhotos")
	return musicLibrary, nil
}
