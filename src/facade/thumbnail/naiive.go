package facade

import (
	"net/url"
	"thumbnailer/src/data"
)

type naiiveThumbnailer struct {
	name string
}

func NewNaiiveThumbnailer(err error) Thumbnailer {
	return &naiiveThumbnailer{
		name: "naiive",
	}
}

func (t *naiiveThumbnailer) GetThumbnail(u *url.URL) (data.Thumbnail, error) {
	panic("TODO(john): Implement naive thumbnailer")
}
