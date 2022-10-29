package api

import (
	"fmt"
	thumbnailer "humpy/src/facade/thumbnail"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// PUT /api/thumbnail
//
// This is a multipart-form upload with a single field called "file" (type "file")
func ThumbnailFile(c *gin.Context) {
	// Get the file being uploaded
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("ThumbnailFile(): %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	mpf, err := file.Open()
	if err != nil {
		log.Printf("ThumbnailFile(): %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Save it to a temporary directory
	tempFile, err := ioutil.TempFile("", "*-"+file.Filename)
	if err != nil {
		log.Printf("ThumbnailFile(): %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Make sure the file gets cleaned up
	defer os.Remove(tempFile.Name())
	io.Copy(tempFile, mpf)
	mpf.Close()
	tempFile.Close()

	// Get the thumbnailer that deals with this type of file
	thumber, err := thumbnailer.GetThumbnailer(tempFile)
	if err != nil {
		log.Printf("ThumbnailFile(): %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Generate the thumbnail
	thumbnail, err := thumber.GetThumbnail(tempFile)
	if err != nil {
		log.Printf("ThumbnailFile(): %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// return the thumbnail
	c.Status(http.StatusOK)
	c.Writer.Header().Add("Content-disposition", "inline; filename="+file.Filename+"-thumb.png")
	c.Writer.Header().Add("Content-type", "image/png")
	c.Writer.Header().Add("Content-length", fmt.Sprint(len(thumbnail.GetContent())))
	c.Writer.Write(thumbnail.GetContent())
}
