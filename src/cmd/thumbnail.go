package cmd

import (
	"log"
	"thumbnailer/src/api"

	"github.com/spf13/cobra"
)

var thumbnailCmd = &cobra.Command{
	Use:   "thumbnail",
	Short: "Runs the Humpy thumbnail server",
	Run:   RunThumbnail,
}

var thumbnailAddress string

func init() {
	thumbnailCmd.PersistentFlags().StringVarP(&thumbnailAddress, "listen", "l", "0.0.0.0:8000", "Thumbnail service listen address (default is 0.0.0.0:8000)")
}

func RunThumbnail(cmd *cobra.Command, args []string) {
	// Validate args and flags
	if thumbnailAddress == "" {
		log.Fatal("RunThumbnail(): Must have a value for -listen")
	}

	log.Println("Running Humpy thumbnail on " + thumbnailAddress)

	// Create the router
	ginRouter, err := api.NewRouter()
	if err != nil {
		log.Printf("RunThumbnail(): Could not start server: %s", err.Error())
	}

	// Set the routes for this service
	ginRouter.POST("/api/thumbnail", api.ThumbnailFile)

	// Start the service
	err = ginRouter.Run(thumbnailAddress)
	if err != nil {
		log.Fatalf("RunThumbnail(): Could not start server: %s", err.Error())
	}
}
