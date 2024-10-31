package facade

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func GetUrlThumbnailer(u *url.URL) (Thumbnailer, error) {
	if u == nil {
		err := fmt.Errorf("GetUrlThumbnailer(%v): URL cannot be nil", u)
		return NewUnknownThumbnailer(err), err
	}

	if strings.EqualFold(u.Scheme, "file") {
		// This is a file url.
		// Delegate to the file thumbnailer
		return GetFileThumbnailer(u.Path)
	}

	return NewHtmlThumbnailer(), nil
}

func GetFileThumbnailer(path string) (Thumbnailer, error) {
	mt, err := mimetype.DetectFile(path)
	if err != nil {
		err := fmt.Errorf("GetFileThumbnailer(%s): Could not get mime type", path)
		return NewUnknownThumbnailer(err), err
	}
	if mt == nil {
		err := fmt.Errorf("GetFileThumbnailer(%s): Could not get mime type", path)
		return NewUnknownThumbnailer(err), err
	}

	if strings.Contains(mt.String(), "application/pdf") {
		// PDF document
		return NewPdfThumbnailer(), nil
	}

	if strings.Contains(mt.String(), "image/") && !strings.Contains(mt.String(), "image/x-") && !strings.Contains(mt.String(), "svg") {
		// Image
		return NewImagemagickThumbnailer(), nil
	}

	if strings.Contains(mt.String(), "video/") {
		// Image
		return NewVideoThumbnailer(), nil
	}

	if strings.Contains(mt.String(), "application/vnd.openxmlformats-officedocument") && !strings.Contains(mt.String(), "image/x-") {
		// Office document
		return NewOfficeThumbnailer(), nil
	}

	log.Printf("GetFileThumbnailer(%s): mime type '%s' not supported", path, mt.String())
	return NewNullThumbnailer(), nil
}
