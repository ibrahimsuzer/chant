package cmd

import (
	"errors"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/viper"
)

const (
	envPrefix = "CHANT"
)

// AppConfig configures app service.
type AppConfig struct {
	Verbose bool
}

// Validate validates the configuration.
func (a AppConfig) Validate() error {
	return validation.ValidateStruct(&a)
}

// configuration holds all externally provided configuration values.
type configuration struct {
	viper *viper.Viper
	// App configuration
	App AppConfig
}

// NewConfiguration creates a new external configuration holder.
func NewConfiguration(v *viper.Viper) *configuration {
	return &configuration{viper: v} //nolint:exhaustivestruct
}

// Configure sets default viper configuration.
func (c *configuration) Configure() {
	c.viper.SetConfigName("config")
	c.viper.AddConfigPath(".")

	// Environment variables
	c.viper.SetEnvPrefix(envPrefix)
	c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	c.viper.AllowEmptyEnv(true)
	c.viper.AutomaticEnv()
}

// Load fills the config from config.* files or environment variables.
func (c *configuration) Load() error {
	if err := c.viper.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if ok := errors.As(err, &viper.ConfigFileNotFoundError{}); !ok {
			return fmt.Errorf("failed to parse configuration: %v", err)
		}
	}

	if err := c.viper.Unmarshal(&c.App); err != nil {
		return fmt.Errorf("failed to unmarshall the configuration: %w", err)
	}
	return nil
}
