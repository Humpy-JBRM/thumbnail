package facade

import (
	"fmt"
	"image"
	"io"
)

func getImageDimensions(imageBytes io.Reader) (width int, height int, err error) {
	img, _, err := image.DecodeConfig(imageBytes)
	if err != nil {
		return 0, 0, fmt.Errorf("GetImageDimensions(): %s", err.Error())
	}

	return img.Width, img.Height, nil
}
