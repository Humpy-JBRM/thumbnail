package facade

import (
	"net/url"
	"thumbnailer/src/data"
)

type Thumbnailer interface {
	GetThumbnail(u *url.URL) (data.Thumbnail, error)
}
