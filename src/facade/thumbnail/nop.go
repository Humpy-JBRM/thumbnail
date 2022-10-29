package facade

import (
	"humpy/src/data"
	"os"
)

type nullThumbnailer struct {
	name string
	err  error
}

func NewNullThumbnailer() Thumbnailer {
	return &nullThumbnailer{
		name: "null",
	}
}

func (t *nullThumbnailer) GetThumbnail(f *os.File) (data.Thumbnail, error) {
	return &data.ThumbnailImpl{
		Content: make([]byte, 0),
	}, nil
}
