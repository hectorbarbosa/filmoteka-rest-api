package envvar

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

//go:generate counterfeiter -o envvartesting/provider.gen.go . Provider

// Provider ...
type Provider interface {
	Get(key string) (string, error)
}

// Configuration ...
type Configuration struct{}

// Load read the env filename and load it into ENV for this process.
func Load(filename string) error {
	if err := godotenv.Load(filename); err != nil {
		return fmt.Errorf("loading env var file: %w", err)
	}

	return nil
}

// New ...
func New() *Configuration {
	return &Configuration{}
}

// Get returns the value from environment variable `<key>`. When an environment variable `<key>_SECURE` exists
// the provider is used for getting the value.
func (c *Configuration) Get(key string) (string, error) {
	res := os.Getenv(key)

	return res, nil
}
