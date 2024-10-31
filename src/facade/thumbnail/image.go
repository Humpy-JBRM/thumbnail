package facade

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"thumbnailer/src/data"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	"github.com/gabriel-vasile/mimetype"
	"github.com/nfnt/resize"
)

type imageThumbnailer struct {
	name string
}

func NewImageThumbnailer() Thumbnailer {
	return &imageThumbnailer{
		name: "image",
	}
}

func (t *imageThumbnailer) GetThumbnail(f *os.File) (data.Thumbnail, error) {
	// Step 1: Decode the image
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("image.GetThumbnail(%s): Could not decode image: %s", f.Name(), err.Error())
	}
	im, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, fmt.Errorf("image.GetThumbnail(%s): Could not decode image: %s", f.Name(), err.Error())
	}

	newImage := im
	if im.Bounds().Size().X > 640 {
		newImage = resize.Resize(640, 0, im, resize.Lanczos3)
	} else if im.Bounds().Size().X > 480 {
		newImage = resize.Resize(0, 480, im, resize.Lanczos3)
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	err = png.Encode(buf, newImage)
	if err != nil {
		return nil, fmt.Errorf("image.GetThumbnail(%s): Could not encode image: %s", f.Name(), err.Error())
	}

	mt := mimetype.Detect(buf.Bytes())
	if mt == nil {
		return nil, fmt.Errorf("image.GetThumbnail(%s): Could not get thumbnail mime type", f.Name())
	}

	// Step 4: Generate the thumbnail object and return
	return &data.ThumbnailImpl{
		MimeType: mt.String(),
		Size:     int64(len(buf.Bytes())),
		Content:  buf.Bytes(),
	}, nil
}
