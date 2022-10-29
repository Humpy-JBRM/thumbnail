package facade

// A facade which encapsulates a SQL config.
//
// The naiive implementation is an in-memory CONFIG
import (
	"fmt"
	"log"
	"strings"
	"sync"
)

var configManagerOnce sync.Once
var configManager ConfigManager

func initialise() {
	cmm, err := ConfigFactory().New()
	if err != nil {
		if strings.Index(err.Error(), "ERROR|") == 0 {
			log.Fatal(err)
		} else {
			log.Fatalf("ERROR|facade/config.Config()|Could not create instance|%s", err.Error())
		}
	}
	configManager = cmm
}

func Config() ConfigManager {
	configManagerOnce.Do(initialise)
	if configManager == nil {
		initialise()
	}
	return configManager
}

type ConfigManagerType string

const (
	CONFIG_UNKNOWN ConfigManagerType = "unknown"
	CONFIG_NAIIVE  ConfigManagerType = "naiive"
)

type configManagerFactoryImpl struct {
	cmType ConfigManagerType
}

var configManagerFactoryOnce sync.Once
var configManagerFactory ConfigManagerFactory

func ConfigFactory() ConfigManagerFactory {
	cmType := CONFIG_NAIIVE

	configManagerFactoryOnce.Do(func() {
		configManagerFactory = &configManagerFactoryImpl{
			cmType: cmType,
		}
	})
	return configManagerFactory
}

func (cmf *configManagerFactoryImpl) SetType(cmType ConfigManagerType) ConfigManagerFactory {
	cmf.cmType = cmType
	return cmf
}

func (cmf *configManagerFactoryImpl) New() (ConfigManager, error) {
	switch cmf.cmType {
	case CONFIG_NAIIVE, "":
		return SetupViperClient()

	default:
		return nil, fmt.Errorf("ERROR|facade/config|No implementation for config manager type '%s'", cmf.cmType)
	}
}
