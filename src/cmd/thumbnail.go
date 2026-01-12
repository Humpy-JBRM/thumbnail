package cmd

import (
	"log"
	"thumbnailer/src/api"
	"thumbnailer/src/data"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var thumbnailCmd = &cobra.Command{
	Use:   "thumbnail",
	Short: "Runs the Humpy thumbnail server",
	Run:   RunThumbnail,
}

var listenAddress string

func init() {
	thumbnailCmd.PersistentFlags().StringVarP(&listenAddress, "listen", "l", "0.0.0.0:8001", "Thumbnail service listen address (default is 0.0.0.0:8001)")
}

func RunThumbnail(cmd *cobra.Command, args []string) {
	// Validate args and flags
	if listenAddress == "" {
		listenAddress = viper.GetString(data.CONF_LISTEN)
	}
	if listenAddress == "" {
		log.Fatal("ERROR|RunThumbnail()|Must have a value for -listen or config.listen")
	}

	log.Println("Running Humpy thumbnail on " + listenAddress)

	// Create the router
	ginRouter, err := api.NewRouter()
	if err != nil {
		log.Printf("ERROR|RunThumbnail()|Could not start server: %s", err.Error())
	}

	// Set the routes for this service
	ginRouter.POST("/api/thumbnail", api.HandleThumbnail)

	// If config.server.enabled is true, set up the static file server
	// so that we:
	//
	//	- save thumbnails to disk when they are generated
	//	- serve thumbnails from disk when they already exist
	if viper.GetBool(data.CONF_SERVER_ENABLED) {
		staticPath := viper.GetString(data.CONF_SERVER_PATH)
		if staticPath == "" {
			log.Fatal("ERROR|RunThumbnail()|config.server.enabled is true but config.server.path is not set")
		}
		log.Printf("INFO|RunThumbnail()|Serving thumbnails from %s|", staticPath)
		ginRouter.GET("/thumbnail/:id", api.ServeThumbnail)
	}

	// Start the service
	err = ginRouter.Run(listenAddress)
	if err != nil {
		log.Fatalf("ERROR|RunThumbnail()|Could not start server: %s", err.Error())
	}
}
