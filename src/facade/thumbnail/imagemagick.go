package facade

import (
	"fmt"
	"humpy/src/data"
	"humpy/src/util"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

type imagemagickThumbnailer struct {
	name string
}

func NewImagemagickThumbnailer() Thumbnailer {
	return &imagemagickThumbnailer{
		name: "imagemagick",
	}
}

func (t *imagemagickThumbnailer) GetThumbnail(f *os.File) (data.Thumbnail, error) {
	// Step 1: Create a temporary dir
	dir, err := ioutil.TempDir("", "*-magick")
	if err != nil {
		return nil, fmt.Errorf("imagemagick.GetThumbnail(%s): Could not create temp dir: %s", f.Name(), err.Error())
	}
	defer os.RemoveAll(dir)

	// Step 3: Convert the image to a thumbnail with ImageMagick
	cmd := "convert"
	thumbnailFile := filepath.Join(dir, filepath.Base(f.Name())+"-thumb.png")
	args := []string{
		"-thumbnail",
		"x300",
		f.Name(),
		thumbnailFile,
	}
	err = util.RunCommand(dir, cmd, 10, map[string]string{}, args...)
	if err != nil {
		return nil, fmt.Errorf("imagemagick.GetThumbnail(%s): Could not generate thumbnail: %s", f.Name(), err.Error())
	}

	// Step 5: Get the thumbnail bytes
	thumbnailBytes, err := ioutil.ReadFile(thumbnailFile)
	if err != nil {
		return nil, fmt.Errorf("imagemagick.GetThumbnail(%s): Could not get thumbnail bytes: %s", f.Name(), err.Error())
	}

	mt := mimetype.Detect(thumbnailBytes)
	if mt == nil {
		return nil, fmt.Errorf("imagemagick.GetThumbnail(%s): Could not get thumbnail mime type", f.Name())
	}

	// Step 4: Generate the thumbnail object and return
	return &data.ThumbnailImpl{
		MimeType: mt.String(),
		Size:     int64(len(thumbnailBytes)),
		Content:  thumbnailBytes,
	}, nil
}
