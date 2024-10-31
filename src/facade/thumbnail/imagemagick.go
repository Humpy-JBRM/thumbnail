package facade

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"thumbnailer/src/data"
	"thumbnailer/src/util"

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

func (t *imagemagickThumbnailer) GetThumbnail(u *url.URL) (data.Thumbnail, error) {
	// Step 1: Create a temporary dir
	dir, err := os.MkdirTemp("", "*-magick")
	if err != nil {
		return nil, fmt.Errorf("imagemagick.GetThumbnail(%s): Could not create temp dir: %s", u.Path, err.Error())
	}
	defer os.RemoveAll(dir)

	// Step 3: Convert the image to a thumbnail with ImageMagick
	cmd := "convert"
	thumbnailFile := filepath.Join(dir, filepath.Base(u.Path)+"-thumb.png")
	args := []string{
		"-thumbnail",
		"x300",
		u.Path,
		thumbnailFile,
	}
	err = util.RunCommand(dir, cmd, 10, map[string]string{}, args...)
	if err != nil {
		return nil, fmt.Errorf("imagemagick.GetThumbnail(%s): Could not generate thumbnail: %s", u.Path, err.Error())
	}

	// Step 5: Get the thumbnail bytes
	thumbnailBytes, err := os.ReadFile(thumbnailFile)
	if err != nil {
		return nil, fmt.Errorf("imagemagick.GetThumbnail(%s): Could not get thumbnail bytes: %s", u.Path, err.Error())
	}

	mt := mimetype.Detect(thumbnailBytes)
	if mt == nil {
		return nil, fmt.Errorf("imagemagick.GetThumbnail(%s): Could not get thumbnail mime type", u.Path)
	}

	// Step 4: Generate the thumbnail object and return
	return &data.ThumbnailImpl{
		MimeType: mt.String(),
		Size:     int64(len(thumbnailBytes)),
		Content:  thumbnailBytes,
	}, nil
}
