// Package processor is a massively concurrent thumbnailler for generating
// thumbnails for directories of hundreds of thousands of files.
// this package is designed to generate JPEG thumbnails with a specified max
// width and height, and preserves the input aspect ratio.
//
// This package is primarily so fast because of libvips/bimg, and a fast os.Walk
// implementation by github.com/MichaelTJones/walk
package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	bimg "gopkg.in/h2non/bimg.v1"
)

var (
	width       = 300
	height      = 300
	quality     = 80
	destination = ""
)

// Process implements filepath.WalkFunc and will thumbnail any images it finds.
func Process(path string, info os.FileInfo, err error) error {
	ext := filepath.Ext(path)

	switch strings.ToLower(ext) {
	case ".jpeg", ".jpg", ".gif", ".png":
	default:
		return nil
	}

	buf, err := bimg.Read(path)
	if err != nil {
		return err
	}

	img := bimg.NewImage(buf)
	meta, err := img.Metadata()
	if err != nil {
		return err
	}

	w, h := ratio(meta.Size, width, height)

	b, err := img.Process(bimg.Options{
		Width:   w,
		Height:  h,
		Type:    bimg.JPEG,
		Quality: quality,
	})
	if err != nil {
		return err
	}

	f := filename(path, destination)

	if _, err = os.Stat(f); os.IsNotExist(err) {
		fmt.Printf("wrote: %s\n", f)

		return bimg.Write(f, b)
	}

	fmt.Printf("skip:  %s\n", f)
	return nil
}

// Expand is a function that will return an absolute path to any given path.
func Expand(path string) string {
	file, err := filepath.Abs(path)
	if err != nil {
		return path
	}

	return file
}

// SetHeight sets the max height of generated thumbnails.
func SetHeight(h int) {
	height = h
}

// SetWidth sets the max width of generated thumbnails
func SetWidth(w int) {
	width = w
}

// SetQuality sets the quality level (1-100) for JPEG output.
func SetQuality(q int) {
	quality = q
}

// SetDestination sets the target folder for thumbnails.
func SetDestination(dst string) {
	destination = Expand(dst)
}

func filename(path, target string) string {
	path = filepath.Clean(path)
	file := filepath.Base(path)

	return filepath.Join(target, ext(file))
}

func ext(path string) string {
	ext := filepath.Ext(path)

	return strings.Replace(path, ext, ".jpeg", 1)
}

// ratio grabs an
//
// Image data: (wi, hi) and define ri = wi / hi
// Target Size: (ws, hs) and define rt = wt / ht
//
//Scaled image dimensions:
//
// rs > ri ? (wi * hs/hi, hs) : (ws, hi * ws/wi)
func ratio(meta bimg.ImageSize, width, height int) (int, int) {
	ri := meta.Width / meta.Height
	rt := width / height

	if rt > ri {
		return meta.Width * height / meta.Height, height
	}

	return width, meta.Height * width / meta.Width
}
