package facade

import (
	"os"
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

func (t *naiiveThumbnailer) GetThumbnail(f *os.File) (data.Thumbnail, error) {
	panic("TODO(john): Implement naive thumbnailer")
}
