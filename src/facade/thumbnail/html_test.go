package facade

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestThumbnailWebPage(t *testing.T) {
	u, err := url.Parse("http://news.bbc.co.uk")
	require.Nil(t, err)
	require.NotNil(t, u)

	thumbnailer, err := GetUrlThumbnailer(u)
	require.Nil(t, err)
	require.NotNil(t, thumbnailer)

	thumbnail, err := thumbnailer.GetThumbnail(u)
	require.Nil(t, err)
	require.NotNil(t, thumbnail)
	require.Nil(t, thumbnail.GetError())
}
