package facade

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"thumbnailer/src/data"
	"thumbnailer/src/util"

	"github.com/gabriel-vasile/mimetype"
)

type pdfThumbnailer struct {
	name string
}

func NewPdfThumbnailer() Thumbnailer {
	return &pdfThumbnailer{
		name: "pdf",
	}
}

func (t *pdfThumbnailer) GetThumbnail(u *url.URL) (data.Thumbnail, error) {
	// Step 1: Create a temporary dir
	dir, err := os.MkdirTemp("", "*-pdf")
	if err != nil {
		return nil, fmt.Errorf("pdf.GetThumbnail(%s): Could not create temp dir: %s", u.Path, err.Error())
	}
	defer os.RemoveAll(dir)

	// Step 3: Execute pdf2image
	cmd := "pdf2image"
	args := []string{
		"-o",
		dir,
		"--pages",
		"1",
		u.Path,
	}
	outputFile := filepath.Join(dir, strings.ReplaceAll(filepath.Base(u.Path), filepath.Ext(u.Path), ".png"))
	err = util.RunCommand(dir, cmd, 10, map[string]string{}, args...)
	if err != nil {
		return nil, fmt.Errorf("pdf.GetThumbnail(%s): Could not generate thumbnail: %s", u.Path, err.Error())
	}

	// Step 4: Convert this image to a thumbnail with ImageMagick
	// convert -thumbnail x300 /tmp/ScottLogic.png /tmp/thumb.png
	cmd = "convert"
	thumbnailFile := filepath.Join(dir, filepath.Base(u.Path)+"-thumb.png")
	args = []string{
		"-thumbnail",
		"x300",
		outputFile,
		thumbnailFile,
	}
	err = util.RunCommand(dir, cmd, 10, map[string]string{}, args...)
	if err != nil {
		fmt.Printf("pdf.GetThumbnail(%s): Could not generate thumbnail: %s", u.Path, err.Error())
		return nil, fmt.Errorf("pdf.GetThumbnail(%s): Could not generate thumbnail: %s", u.Path, err.Error())
	}

	thumbnailBytes, err := os.ReadFile(thumbnailFile)
	if err != nil {
		return nil, fmt.Errorf("pdf.GetThumbnail(%s): Could not get thumbnail bytes: %s", u.Path, err.Error())
	}

	// Step 5: get the width and height of the thumbnail
	w, h, _ := getImageDimensions(bytes.NewReader(thumbnailBytes))

	// Step 6: Get the mime type of the thumbnail
	mt := mimetype.Detect(thumbnailBytes)
	if mt == nil {
		return nil, fmt.Errorf("pdf.GetThumbnail(%s): Could not get thumbnail mime type", u.Path)
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
