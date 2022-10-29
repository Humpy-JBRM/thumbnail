package facade

import (
	"humpy/src/data"
	"os"
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
