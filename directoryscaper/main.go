package directoryscaper

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
)

var acceptableFileTypes [10]string = [10]string{".mp3", ".flac", ".m4a", ".wav", ".ogg", ".aac", ".alac", ".m4p", ".raw", ".wma"}

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
}

// todo: this takes like 5 minutes to run and get all of the files
// and that's only with like 7100 files
// it's specifically due to the tag library
// it is probably worth storing this as json or something

// GetAllInDirectory gets all music files in the specified base directory
func GetAllInDirectory(baseDir string) ([]MusicFile, error) {
	files := []MusicFile{}
	err := filepath.WalkDir(baseDir, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !file.IsDir() {
			for _, v := range acceptableFileTypes {
				if strings.HasSuffix(file.Name(), v) {
					if err != nil {
						return err
					}
					meta := readFileMetaData(path)

					discNum, totalDiscs := meta.Disc()
					trackNum, totalTracks := meta.Track()
					picture := meta.Picture()

					musicFile := MusicFile{
						Album:           meta.Album(),
						AlbumArtist:     meta.AlbumArtist(),
						Artist:          meta.Artist(),
						Composer:        meta.Composer(),
						DiscNumber:      discNum,
						TotalDiscs:      totalDiscs,
						FileType:        string(meta.FileType()),
						Format:          string(meta.Format()),
						Genre:           meta.Genre(),
						PictureData:     picture.Data,
						PictureMIMEType: picture.MIMEType,
						Title:           meta.Title(),
						TrackNumber:     trackNum,
						TotalTracks:     totalTracks,
						Year:            meta.Year(),
						Path:            path,
					}
					files = append(files, musicFile)
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
	meta, err := tag.ReadFrom(file)
	if err != nil {
		log.Printf("error reading meta data: %v", err)
		return nil
	}
	return meta
}
