package facade

import (
	"humpy/src/data"
	"os"
)

type Thumbnailer interface {
	GetThumbnail(f *os.File) (data.Thumbnail, error)
}
