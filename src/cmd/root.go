package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "humpy",
	Short: "Humpy is a way to mirror and declutter all your digital stuff",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

var cfgFile string

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "config.yml", "config file (default is config.yml)")

	rootCmd.AddCommand(thumbnailCmd)

	// Microservices
}

func initConfig() {
	if configFile, isSet := os.LookupEnv("CONFIG_FILE"); isSet {
		cfgFile = configFile
	}
	if cfgFile != "" {
		err := LoadConfig(cfgFile)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func LoadConfig(cfgFile string) error {
	log.Printf("INFO|cmd.initConfig()|Reading config from %s|", cfgFile)
	cfDir := filepath.Dir(cfgFile)
	if cfDir == "" {
		cfDir = "."
	}
	viper.AddConfigPath(cfDir)
	viper.SetConfigType(filepath.Ext(cfgFile)[1:])
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("ERROR|cmd.initConfig()|Could not read file %s|%s", viper.ConfigFileUsed(), err.Error())
	}

	for _, key := range viper.AllKeys() {
		log.Printf("%s: %v\n", key, viper.Get(key))
	}

	return nil
}
