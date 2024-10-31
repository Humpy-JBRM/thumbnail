package facade

import (
	"bytes"
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

type officeThumbnailer struct {
	name string
}

func NewOfficeThumbnailer() Thumbnailer {
	return &officeThumbnailer{
		name: "office",
	}
}

func (t *officeThumbnailer) GetThumbnail(u *url.URL) (data.Thumbnail, error) {
	// Step 1: Create a temporary dir
	dir, err := os.MkdirTemp("", "*-office")
	if err != nil {
		log.Printf("office.GetThumbnail(%s): Could not create temp dir: %s", u.Path, err.Error())
		return nil, fmt.Errorf("office.GetThumbnail(%s): Could not create temp dir: %s", u.Path, err.Error())
	}
	defer os.RemoveAll(dir)

	// Step 3: Execute libreoffice to convert the file to an image
	cmd := "lowriter"
	args := []string{
		"-env:UserInstallation=file://" + dir,
		"--invisible",
		"--headless",
		"--convert-to",
		"png",
		"--outdir",
		dir,
		u.Path,
	}
	err = util.RunCommand(dir, cmd, 10, map[string]string{}, args...)
	if err != nil {
		log.Printf("office.GetThumbnail(%s): Could not generate thumbnail: %s", u.Path, err.Error())
		return nil, fmt.Errorf("office.GetThumbnail(%s): Could not generate thumbnail: %s", u.Path, err.Error())
	}

	sourceFile := filepath.Join(dir, filepath.Base(u.Path)+"[0]")
	if _, err := os.Stat(sourceFile); err != nil {
		// Perhaps there was only 1 page in the doc
		sourceFile = filepath.Join(dir, strings.ReplaceAll(filepath.Base(u.Path), filepath.Ext(u.Path), ".png"))
	}
	// Step 4: Convert this image to a thumbnail with ImageMagick
	// convert -thumbnail x300 /tmp/ScottLogic.png /tmp/thumb.png
	cmd = "convert"
	thumbnailFile := filepath.Join(dir, filepath.Base(u.Path)+"-thumb.png")
	args = []string{
		"-thumbnail",
		"x300",
		sourceFile,
		thumbnailFile,
	}
	err = util.RunCommand(dir, cmd, 10, map[string]string{}, args...)
	if err != nil {
		fmt.Printf("office.GetThumbnail(%s): Could not generate thumbnail: %s", u.Path, err.Error())
		return nil, fmt.Errorf("office.GetThumbnail(%s): Could not generate thumbnail: %s", u.Path, err.Error())
	}

	// Step 5: Get the thumbnail bytes
	thumbnailBytes, err := os.ReadFile(thumbnailFile)
	if err != nil {
		fmt.Printf("office.GetThumbnail(%s): Could not get thumbnail bytes: %s", u.Path, err.Error())
		return nil, fmt.Errorf("office.GetThumbnail(%s): Could not get thumbnail bytes: %s", u.Path, err.Error())
	}

	// Step 6: Set the width and height of the thumbnail
	w, h, _ := getImageDimensions(bytes.NewReader(thumbnailBytes))

	mt := mimetype.Detect(thumbnailBytes)
	if mt == nil {
		fmt.Printf("office.GetThumbnail(%s): Could not get thumbnail mime type", u.Path)
		return nil, fmt.Errorf("office.GetThumbnail(%s): Could not get thumbnail mime type", u.Path)
	}

	// Step 7: Generate the thumbnail object and return
	return &data.ThumbnailImpl{
		MimeType: mt.String(),
		Size:     int64(len(thumbnailBytes)),
		Content:  thumbnailBytes,
		Width:    int64(w),
		Height:   int64(h),
	}, nil
}
