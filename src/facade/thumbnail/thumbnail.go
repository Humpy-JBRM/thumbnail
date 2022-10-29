package facade

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func GetThumbnailer(f *os.File) (Thumbnailer, error) {
	mt, err := mimetype.DetectFile(f.Name())
	if err != nil {
		err := fmt.Errorf("GetThumbnailer(%s): Could not get mime type", f.Name())
		return NewUnknownThumbnailer(err), err
	}
	if mt == nil {
		err := fmt.Errorf("GetThumbnailer(%s): Could not get mime type", f.Name())
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

	log.Printf("GetThumbnailer(%s): mime type '%s' not supported", f.Name(), mt.String())
	return NewNullThumbnailer(), nil
}
