package facade

import (
	"humpy/src/data"
	"os"
)

type unknownThumbnailer struct {
	name string
	err  error
}

func NewUnknownThumbnailer(err error) Thumbnailer {
	return &unknownThumbnailer{
		name: "unknown",
		err:  err,
	}
}

func (t *unknownThumbnailer) GetThumbnail(file *os.File) (data.Thumbnail, error) {
	panic("TODO(john): Implement unknown thumbnailer")
}
