package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	thumbnailer "thumbnailer/src/facade/thumbnail"

	"github.com/gin-gonic/gin"
)

// PUT /api/thumbnail
//
// This is a multipart-form upload with a single field called "file" (type "file")
func HandleThumbnail(c *gin.Context) {
	file, _ := c.FormFile("file")
	if file != nil {
		ThumbnailFile(c)
		return
	}

	urlValue := c.Request.PostFormValue("url")
	if urlValue != "" {
		ThumbnailUrl(c)
		return
	}

	c.AbortWithError(http.StatusBadRequest, fmt.Errorf("ThumbnailFile(): Need 'file' or 'url' value"))
}

func ThumbnailUrl(c *gin.Context) {
	u, err := url.Parse(c.Request.PostFormValue("url"))
	if err != nil {
		log.Printf("ThumbnailUrl(): %s", err.Error())
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("ThumbnailUrl(): %s", err.Error()))
		return
	}

	thumber, err := thumbnailer.GetUrlThumbnailer(u)
	if err != nil {
		log.Printf("ThumbnailUrl(): %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Generate the thumbnail
	thumbnail, err := thumber.GetThumbnail(u)
	if err != nil {
		log.Printf("ThumbnailFile(): %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// return the thumbnail
	c.Status(http.StatusOK)
	c.Writer.Header().Add("Content-disposition", "inline; filename="+u.Host+"-thumb.png")
	c.Writer.Header().Add("Content-type", "image/png")
	c.Writer.Header().Add("Content-length", fmt.Sprint(len(thumbnail.GetContent())))
	c.Writer.Write(thumbnail.GetContent())
}

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
	tempFile, err := os.CreateTemp("", "*-"+file.Filename)
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
	thumber, err := thumbnailer.GetFileThumbnailer(tempFile.Name())
	if err != nil {
		log.Printf("ThumbnailFile(): %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Generate the thumbnail
	fileUrl, err := url.Parse(fmt.Sprintf("file://%s", tempFile.Name()))
	if err != nil {
		log.Printf("ThumbnailFile(): %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	thumbnail, err := thumber.GetThumbnail(fileUrl)
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
