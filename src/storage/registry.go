package storage

import (
	"fmt"

	"blob/src/config"
)

// DriverFactory creates a StorageDriver from a config.
type DriverFactory func(cfg *config.Config) (StorageDriver, error)

// factories registers available storage drivers by name.
var factories = map[string]DriverFactory{
	"filesystem": func(cfg *config.Config) (StorageDriver, error) {
		return NewFilesystemDriver(cfg.ResolveStoragePath()), nil
	},
}

// Register adds a new storage driver to the registry. It is a no-op if the
// name is already registered.
func Register(name string, factory DriverFactory) {
	if _, ok := factories[name]; !ok {
		factories[name] = factory
	}
}

// New creates a StorageDriver by name, defaulting to "filesystem".
func New(name string, cfg *config.Config) (StorageDriver, error) {
	if name == "" {
		name = "filesystem"
	}
	factory, ok := factories[name]
	if !ok {
		return nil, fmt.Errorf("unknown storage driver %q", name)
	}
	return factory(cfg)
}
