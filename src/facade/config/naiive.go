package facade

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type naiiveConfigManagerImpl struct {
}

func SetupViperClient() (ConfigManager, error) {
	if os.Getenv("CONFIG_FILE") == "" {
		defaultConfig, _ := os.Stat("config.yml")
		if defaultConfig != nil {
			// Config file is not explicitly set, and there is a ./config.yml
			os.Setenv("CONFIG_FILE", "config.yml")
		}
	}

	if os.Getenv("CONFIG_FILE") != "" {
		config := filepath.Base(os.Getenv("CONFIG_FILE"))
		dir := filepath.Dir(os.Getenv("CONFIG_FILE"))
		if dir == "" {
			dir = "."
		}
		ext := filepath.Ext(config)
		if len(ext) < 2 {
			return &unknownConfigManager{}, fmt.Errorf("config/naiive.SetupViperClient(): Could not read '%s' ($CONFIG_FILE)|expected it's the extention (%s) of it's base (%s) to have length >1", os.Getenv("CONFIG_FILE"), ext, config)
		}
		viper.SetConfigName(config)
		viper.SetConfigType(ext[1:])
		viper.AddConfigPath(dir)
		err := viper.ReadInConfig()
		if err != nil {
			return &unknownConfigManager{}, fmt.Errorf("config/naiive.SetupViperClient(): Could not read '%s' ($CONFIG_FILE): %s", os.Getenv("CONFIG_FILE"), err.Error())
		}
		for _, k := range viper.AllKeys() {
			log.Printf("%s: %s=%v", os.Getenv("CONFIG_FILE"), k, viper.Get(k))
		}

	}

	return &naiiveConfigManagerImpl{}, nil
}

func (cm *naiiveConfigManagerImpl) Type() ConfigManagerType {
	return CONFIG_NAIIVE
}

func (cm *naiiveConfigManagerImpl) GetString(name string, defaultValue ...string) (string, error) {
	defaultVal := ""
	if len(defaultValue) > 0 {
		defaultVal = defaultValue[0]
	}
	val := viper.Get(name)
	if val == nil {
		return defaultVal, nil
	}

	return fmt.Sprint(val), nil
}

func (cm *naiiveConfigManagerImpl) GetStringArray(name string, defaultValue ...[]string) ([]string, error) {
	defaultVal := []string{}
	if len(defaultValue) > 0 {
		defaultVal = defaultValue[0]
	}
	val := viper.GetStringSlice(name)
	if val == nil {
		return defaultVal, nil
	}

	return val, nil
}

func (cm *naiiveConfigManagerImpl) GetInt(name string, defaultValue ...int) (int, error) {
	defaultVal := 0
	if len(defaultValue) > 0 {
		defaultVal = defaultValue[0]
	}
	val := viper.Get(name)
	if val == nil {
		return defaultVal, nil
	}

	return val.(int), nil
}

func (cm *naiiveConfigManagerImpl) GetBool(name string, defaultValue ...bool) (bool, error) {
	defaultVal := false
	if len(defaultValue) > 0 {
		defaultVal = defaultValue[0]
	}
	val := viper.Get(name)
	if val == nil {
		return defaultVal, nil
	}

	return val.(bool), nil
}

func (cm *naiiveConfigManagerImpl) GetSecret(name string, defaultValue ...string) (string, error) {
	return "", nil
}

func (cm *naiiveConfigManagerImpl) SetString(name string, value string) error {
	viper.Set(name, value)
	return nil
}

func (cm *naiiveConfigManagerImpl) SetInt(name string, value int) error {
	viper.Set(name, value)
	return nil
}

func (cm *naiiveConfigManagerImpl) SetBool(name string, value bool) error {
	viper.Set(name, value)
	return nil
}

func (cm *naiiveConfigManagerImpl) GetValue(name string) (interface{}, error) {
	return viper.Get(name), nil
}

func (cm *naiiveConfigManagerImpl) Reset() {
	viper.Reset()
}

func (cm *naiiveConfigManagerImpl) IsDevMode() bool {
	return false
}

func (cm *naiiveConfigManagerImpl) GetAllValues(root string) (map[string]string, error) {
	panic("TODO(john): Implement naive.GetAllValues()")
}
