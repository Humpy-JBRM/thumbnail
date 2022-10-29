package facade

// ConfigManager encapsulates the behaviour that we want from a
// configuration manager component, abstracting away the
// implementation.
//
// The naiive implementation runs entirely in memory, and has no
// external dependencies.
//
// A production implementation would be a facade for vault / k8s secrets.

type ConfigApi interface {
	GetString(string, ...string) (string, error)
	GetStringArray(string, ...[]string) ([]string, error)
	GetInt(string, ...int) (int, error)
	GetBool(string, ...bool) (bool, error)
	GetSecret(string, ...string) (string, error)
	GetValue(string) (interface{}, error)
	SetString(string, string) error
	SetInt(string, int) error
	SetBool(string, bool) error

	// Helper function which says whether or not dev mode is enabled
	IsDevMode() bool

	GetAllValues(roto string) (map[string]string, error)
}

type ConfigManager interface {
	ConfigApi
	Type() ConfigManagerType
	Reset()
}

// ConfigManagerFactory produces ConfigManager instances based on
// the criteria we set.
//
// This is how we can flip implementations without refactoring
// any code
type ConfigManagerFactory interface {
	SetType(cmType ConfigManagerType) ConfigManagerFactory
	New() (ConfigManager, error)
}
