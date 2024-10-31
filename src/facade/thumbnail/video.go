package facade

import (
	"fmt"
	"log"
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

func (t *videoThumbnailer) GetThumbnail(f *os.File) (data.Thumbnail, error) {
	// Step 1: Create a temporary dir
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, fmt.Errorf("video.GetThumbnail(%s): Could not create temp dir: %s", f.Name(), err.Error())
	}
	defer os.RemoveAll(dir)

	// Step 3: Execute ffmpeg
	cmd := "ffmpeg"
	outputFile := filepath.Join(dir, filepath.Base(f.Name())+"-thumb.gif")
	args := []string{
		"-t",
		"30",
		"-i",
		f.Name(),
		"-vf",
		"fps=5",
		"-s",
		"360x240",
		outputFile,
	}
	log.Printf("EXECUTE: %s %s", cmd, strings.Join(args, " "))
	err = util.RunCommand(dir, cmd, 10, map[string]string{}, args...)
	if err != nil {
		return nil, fmt.Errorf("video.GetThumbnail(%s): Could not generate thumbnail: %s", f.Name(), err.Error())
	}
	thumbnailBytes, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, fmt.Errorf("video.GetThumbnail(%s): Could not get thumbnail bytes: %s", f.Name(), err.Error())
	}

	mt := mimetype.Detect(thumbnailBytes)
	if mt == nil {
		return nil, fmt.Errorf("video.GetThumbnail(%s): Could not get thumbnail mime type", f.Name())
	}

	// Step 4: Generate the thumbnail object and return
	return &data.ThumbnailImpl{
		MimeType: mt.String(),
		Size:     int64(len(thumbnailBytes)),
		Content:  thumbnailBytes,
	}, nil
}
