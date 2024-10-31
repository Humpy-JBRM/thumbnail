package facade

import (
	"os"
	"thumbnailer/src/data"
)

type Thumbnailer interface {
	GetThumbnail(f *os.File) (data.Thumbnail, error)
}
