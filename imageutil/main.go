package imageutil

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"strings"

	"github.com/disintegration/imaging"
)

type dimensions struct {
	x int
	y int
}

func (d *dimensions) scale(maxSizeLength int) {
	if d.x == d.y {
		d.x = maxSizeLength
		d.y = maxSizeLength
		return
	}

	isLandscape := d.x > d.y
	if isLandscape {
		d.y = (maxSizeLength * d.y) / d.x
		d.x = maxSizeLength
	} else {
		d.x = (maxSizeLength * d.x) / d.y
		d.y = maxSizeLength
	}
}

func ResizeImage(base64 []byte, maxSideLength int) []byte {
	reader := bytes.NewReader(base64)
	img, format, err := image.Decode(reader)
	if err != nil {
		log.Println("failed to decode base64 string", err, format)
		return []byte{}
	}

	dimensions := dimensions{x: img.Bounds().Size().X, y: img.Bounds().Size().Y}
	dimensions.scale(maxSideLength)
	resizedImg := imaging.Resize(img, dimensions.x, dimensions.y, imaging.CatmullRom)

	var resizedBytes bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&resizedBytes, resizedImg, nil)
	case "png":
		err = png.Encode(&resizedBytes, resizedImg)
	case "gif":
		err = gif.Encode(&resizedBytes, resizedImg, nil)
	default:
		err = errors.New("unrecognised image format")
	}
	if err != nil {
		log.Println("failed to encode image", err)
		return []byte{}
	}

	return resizedBytes.Bytes()
}

func resizeImageFromString(base64Str string, maxSideLength int) []byte {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Str))
	img, format, err := image.Decode(reader)
	if err != nil {
		log.Println("failed to decode base64 string", err, format)
		return []byte{}
	}

	dimensions := dimensions{x: img.Bounds().Size().X, y: img.Bounds().Size().Y}
	dimensions.scale(maxSideLength)
	resizedImg := imaging.Resize(img, dimensions.x, dimensions.y, imaging.CatmullRom)

	var resizedBytes bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&resizedBytes, resizedImg, nil)
	case "png":
		err = png.Encode(&resizedBytes, resizedImg)
	case "gif":
		err = gif.Encode(&resizedBytes, resizedImg, nil)
	default:
		err = errors.New("unrecognised image format")
	}
	if err != nil {
		log.Println("failed to encode image", err)
		return []byte{}
	}

	return resizedBytes.Bytes()
}
