package facade

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
	"testing"
)

func TestInterestingImage(t *testing.T) {
}

func average(t *testing.T) {
	imageFile := "/ssd/my-humpy/files/ec/ec8f6b2f-39b7-4f38-8863-c06bdff31df2/content"
	file, err := os.Open(imageFile)
	//file, err := os.Open("/ssd/my-humpy/files/97/9762df70-5b7f-4e09-96c2-eab93c19e052/content")
	//file, err := os.Open("/ssd/my-humpy/files/06/0672eec9-90e7-4e0c-9aea-d787a78827d4/content")
	//file, err := os.Open("/ssd/my-humpy/files/93/93173b24-9f88-42fc-9b3f-ed6bd063a6b5/thumbnail")

	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalln(err)
	}

	imgSize := img.Bounds().Size()

	var redSum float64
	var greenSum float64
	var blueSum float64

	for x := 0; x < imgSize.X; x++ {
		for y := 0; y < imgSize.Y; y++ {
			pixel := img.At(x, y)
			col := color.RGBAModel.Convert(pixel).(color.RGBA)

			redSum += float64(col.R)
			greenSum += float64(col.G)
			blueSum += float64(col.B)
		}
	}

	// Calculate the mean as 24-bit value
	imgArea := float64(imgSize.X * imgSize.Y)
	redAverage := redSum / imgArea
	greenAverage := greenSum / imgArea
	blueAverage := blueSum / imgArea
	mean := (redAverage * 256 * 256) + (greenAverage * 256) + blueAverage

	// Now work out the standard deviation
	redSum = 0
	greenSum = 0
	blueSum = 0
	var stddevTotal float64
	for x := 0; x < imgSize.X; x++ {
		for y := 0; y < imgSize.Y; y++ {
			pixel := img.At(x, y)
			col := color.RGBAModel.Convert(pixel).(color.RGBA)
			sum := float64(col.R)*float64(256*256) + float64(col.G)*float64(256) + float64(col.B)
			stddevTotal += (sum - mean) * (sum - mean)
		}
	}
	stddev := math.Sqrt(stddevTotal / imgArea)
	fmt.Printf(
		"Average colour: rgb(%.0f, %.0f, %.0f)\nStd dev: %f",
		redAverage,
		greenAverage,
		blueAverage,
		stddev,
	)
}
