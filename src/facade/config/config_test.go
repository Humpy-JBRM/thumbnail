package facade

import (
	"testing"
)

func TestConfigSingleton(t *testing.T) {
	if Config() != Config() {
		t.Errorf("Config() singleton does not return the same pointer")
	}
}

func TestConfigFactorySingleton(t *testing.T) {
	if ConfigFactory() != ConfigFactory() {
		t.Errorf("ConfigFactory() singleton does not return the same pointer")
	}
}

func TestConfigFactoryReturnsNaiiveByDefault(t *testing.T) {
	cm, err := ConfigFactory().New()
	if err != nil {
		t.Error(err)
	}
	if cm == nil {
		t.Errorf("Expected a ConfigManager, but got nil")
	}
	if cm.Type() != CONFIG_NAIIVE {
		t.Errorf("Expected ConfigManager type '%s', but got '%s'", CONFIG_NAIIVE, cm.Type())
	}
}
