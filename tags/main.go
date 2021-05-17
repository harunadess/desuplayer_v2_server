package tags

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
)

var AcceptableFileTypes []string = []string{".mp3", ".flac", ".m4a" /*".wav",*/, ".ogg", ".alac", ".m4p"}

func ReadTags(path string) (tag.Metadata, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("failed to read file %v\n", path)
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Printf("failed to stat file %v\n", path)
		return nil, err
	}

	switch strings.ToLower(filepath.Ext(info.Name())) {
	case AcceptableFileTypes[0]:
		return readMp3(file)
	case AcceptableFileTypes[1]:
		fallthrough
	case AcceptableFileTypes[4]:
		return readFlac(file)
	case AcceptableFileTypes[2]:
		fallthrough
	case AcceptableFileTypes[5]:
		return readMp4(file)
	case AcceptableFileTypes[3]:
		return readOgg(file)
	default:
		return nil, errors.New("not an accepted type")
	}
}

func readMp3(file *os.File) (tag.Metadata, error) {
	tagMeta, err := tag.ReadID3v2Tags(file)
	if err != nil {
		tagMeta, err = tag.ReadID3v1Tags(file)
		if err != nil {
			return nil, errors.New("failed to read mp3: " + err.Error())
		}
	}
	return tagMeta, nil
}

func readMp4(file *os.File) (tag.Metadata, error) {
	tagMeta, err := tag.ReadAtoms(file)
	if err != nil {
		return nil, errors.New("failed to read mp4: " + err.Error())
	}
	return tagMeta, nil
}

func readFlac(file *os.File) (tag.Metadata, error) {
	tagMeta, err := tag.ReadFLACTags(file)
	if err != nil {
		return nil, errors.New("failed to read flac: " + err.Error())
	}
	return tagMeta, nil
}

func readOgg(file *os.File) (tag.Metadata, error) {
	tagMeta, err := tag.ReadOGGTags(file)
	if err != nil {
		return nil, errors.New("failed to read ogg: " + err.Error())
	}
	return tagMeta, nil
}
