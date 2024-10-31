package facade

import (
	"net/url"
	"thumbnailer/src/data"
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

func (t *unknownThumbnailer) GetThumbnail(u *url.URL) (data.Thumbnail, error) {
	panic("TODO(john): Implement unknown thumbnailer")
}
