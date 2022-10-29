package facade

import (
	"bytes"
	"fmt"
	"humpy/src/data"
	"humpy/src/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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

func (t *pdfThumbnailer) GetThumbnail(f *os.File) (data.Thumbnail, error) {
	// Step 1: Create a temporary dir
	dir, err := ioutil.TempDir("", "*-pdf")
	if err != nil {
		return nil, fmt.Errorf("pdf.GetThumbnail(%s): Could not create temp dir: %s", f.Name(), err.Error())
	}
	defer os.RemoveAll(dir)

	// Step 3: Execute pdf2image
	cmd := "pdf2image"
	args := []string{
		"-o",
		dir,
		"--pages",
		"1",
		f.Name(),
	}
	outputFile := filepath.Join(dir, strings.ReplaceAll(filepath.Base(f.Name()), filepath.Ext(f.Name()), ".png"))
	err = util.RunCommand(dir, cmd, 10, map[string]string{}, args...)
	if err != nil {
		return nil, fmt.Errorf("pdf.GetThumbnail(%s): Could not generate thumbnail: %s", f.Name(), err.Error())
	}

	// Step 4: Convert this image to a thumbnail with ImageMagick
	// convert -thumbnail x300 /tmp/ScottLogic.png /tmp/thumb.png
	cmd = "convert"
	thumbnailFile := filepath.Join(dir, filepath.Base(f.Name())+"-thumb.png")
	args = []string{
		"-thumbnail",
		"x300",
		outputFile,
		thumbnailFile,
	}
	err = util.RunCommand(dir, cmd, 10, map[string]string{}, args...)
	if err != nil {
		fmt.Printf("pdf.GetThumbnail(%s): Could not generate thumbnail: %s", f.Name(), err.Error())
		return nil, fmt.Errorf("pdf.GetThumbnail(%s): Could not generate thumbnail: %s", f.Name(), err.Error())
	}

	thumbnailBytes, err := ioutil.ReadFile(thumbnailFile)
	if err != nil {
		return nil, fmt.Errorf("pdf.GetThumbnail(%s): Could not get thumbnail bytes: %s", f.Name(), err.Error())
	}

	// Step 5: get the width and height of the thumbnail
	w, h, _ := getImageDimensions(bytes.NewReader(thumbnailBytes))

	// Step 6: Get the mime type of the thumbnail
	mt := mimetype.Detect(thumbnailBytes)
	if mt == nil {
		return nil, fmt.Errorf("pdf.GetThumbnail(%s): Could not get thumbnail mime type", f.Name())
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
