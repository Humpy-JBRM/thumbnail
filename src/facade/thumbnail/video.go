package facade

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"thumbnailer/src/data"
	"thumbnailer/src/util"

	"github.com/gabriel-vasile/mimetype"
)

type videoThumbnailer struct {
	name string
}

func NewVideoThumbnailer() Thumbnailer {
	return &videoThumbnailer{
		name: "video",
	}
}

func (t *videoThumbnailer) GetThumbnail(u *url.URL) (data.Thumbnail, error) {
	// Step 1: Create a temporary dir
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, fmt.Errorf("video.GetThumbnail(%s): Could not create temp dir: %s", u.Path, err.Error())
	}
	defer os.RemoveAll(dir)

	// Step 3: Execute ffmpeg
	cmd := "ffmpeg"
	outputFile := filepath.Join(dir, filepath.Base(u.Path)+"-thumb.gif")
	args := []string{
		"-t",
		"30",
		"-i",
		u.Path,
		"-vf",
		"fps=5",
		"-s",
		"360x240",
		outputFile,
	}
	log.Printf("EXECUTE: %s %s", cmd, strings.Join(args, " "))
	err = util.RunCommand(dir, cmd, 10, map[string]string{}, args...)
	if err != nil {
		return nil, fmt.Errorf("video.GetThumbnail(%s): Could not generate thumbnail: %s", u.Path, err.Error())
	}
	thumbnailBytes, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, fmt.Errorf("video.GetThumbnail(%s): Could not get thumbnail bytes: %s", u.Path, err.Error())
	}

	mt := mimetype.Detect(thumbnailBytes)
	if mt == nil {
		return nil, fmt.Errorf("video.GetThumbnail(%s): Could not get thumbnail mime type", u.Path)
	}

	// Step 4: Generate the thumbnail object and return
	return &data.ThumbnailImpl{
		MimeType: mt.String(),
		Size:     int64(len(thumbnailBytes)),
		Content:  thumbnailBytes,
	}, nil
}
