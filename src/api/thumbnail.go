package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"thumbnailer/src/data"
	thumbnailer "thumbnailer/src/facade/thumbnail"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// POST /api/thumbnail
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

	err := fmt.Errorf("ThumbnailFile(): Need 'file' or 'url' value")
	c.Writer.Header().Add("X-Error", err.Error())
	c.AbortWithError(http.StatusBadRequest, err)
}

func ThumbnailUrl(c *gin.Context) {
	u, err := url.Parse(c.Request.PostFormValue("url"))
	if err != nil {
		log.Printf("ERROR|ThumbnailUrl()|%s", err.Error())
		c.Writer.Header().Add("X-Error", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	thumber, err := thumbnailer.GetUrlThumbnailer(u)
	if err != nil {
		log.Printf("ERROR|ThumbnailUrl()|%s", err.Error())
		c.Writer.Header().Add("X-Error", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Generate the thumbnail
	thumbnail, err := thumber.GetThumbnail(u)
	if err != nil {
		log.Printf("ERROR|ThumbnailFile()|%s", err.Error())
		c.Writer.Header().Add("X-Error", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// If we are saving thumbnails (config.server.enabled is true), save it to disk
	var id string
	if viper.GetBool(data.CONF_SERVER_ENABLED) {
		id = c.Request.PostFormValue("id")
		if id == "" {
			log.Printf("ERROR|ThumbnailFile()|No 'id' value provided and '%s' is true; thumbnail not saved to disk", data.CONF_SERVER_ENABLED)
			c.Writer.Header().Add("X-Error", fmt.Sprintf("No 'id' value provided and '%s' is true; thumbnail not saved to disk", data.CONF_SERVER_ENABLED))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		err := saveThumbnail(id, thumbnail)
		if err != nil {
			log.Printf("ERROR|ThumbnailFile()|%s", err.Error())
			c.Writer.Header().Add("X-Error", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	// return the thumbnail
	c.Status(http.StatusOK)
	c.Writer.Header().Add("Content-disposition", "inline; filename="+u.Host+"-thumb"+thumbnail.GetExtension())
	c.Writer.Header().Add("Content-type", "image/png")
	c.Writer.Header().Add("Content-length", fmt.Sprint(len(thumbnail.GetContent())))
	if id != "" {
		c.Writer.Header().Add("X-Thumbnail-Id", id)
	}
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

	// If we are saving thumbnails (config.server.enabled is true), save it to disk
	var id string
	if viper.GetBool(data.CONF_SERVER_ENABLED) {
		id = c.Request.PostFormValue("id")
		if id == "" {
			log.Printf("ERROR|ThumbnailFile()|No 'id' value provided and '%s' is true; thumbnail not saved to disk", data.CONF_SERVER_ENABLED)
			c.Writer.Header().Add("X-Error", fmt.Sprintf("No 'id' value provided and '%s' is true; thumbnail not saved to disk", data.CONF_SERVER_ENABLED))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		err := saveThumbnail(id, thumbnail)
		if err != nil {
			log.Printf("ERROR|ThumbnailFile()|%s", err.Error())
			c.Writer.Header().Add("X-Error", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	// return the thumbnail
	c.Status(http.StatusOK)
	c.Writer.Header().Add("Content-disposition", "inline; filename="+file.Filename+"-thumb"+thumbnail.GetExtension())
	c.Writer.Header().Add("Content-type", "image/png")
	c.Writer.Header().Add("Content-length", fmt.Sprint(len(thumbnail.GetContent())))
	if id != "" {
		c.Writer.Header().Add("X-Thumbnail-Id", id)
	}
	c.Writer.Write(thumbnail.GetContent())
}

// GET /thumbnail/:id
func ServeThumbnail(c *gin.Context) {
	id := c.Param("id")
	thumbnailFile, err := getStaticPath(id)
	if err != nil {
		log.Printf("ERROR|ServeThumbnail()|%s", err.Error())
		c.Writer.Header().Add("X-Error", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.File(thumbnailFile)
}

func getStaticPath(id string) (string, error) {
	staticRoot := viper.GetString(data.CONF_SERVER_PATH)
	if staticRoot == "" {
		return "", fmt.Errorf("ERROR|getStaticPath()|config.server.enabled is true but config.server.path is not set")
	}
	thumbnailFile := filepath.Join(staticRoot, id)
	return thumbnailFile, nil
}

func saveThumbnail(id string, thumbnail data.Thumbnail) error {
	thumbnailFile, err := getStaticPath(id)
	if err != nil {
		return fmt.Errorf("ERROR|saveThumbnail()|%s", err.Error())
	}
	_ = os.MkdirAll(filepath.Dir(thumbnailFile), 0755)
	err = os.WriteFile(thumbnailFile, thumbnail.GetContent(), 0644)
	if err != nil {
		return fmt.Errorf("ERROR|saveThumbnail()|%s", err.Error())
	}
	return nil
}
