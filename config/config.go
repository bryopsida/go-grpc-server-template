package config

import (
	"path"

	"github.com/bryopsida/go-grpc-server-template/interfaces"
	"github.com/spf13/viper"
)

type viperConfig struct {
	viper *viper.Viper
}

// NewViperConfig creates a new viperConfig instance
func NewViperConfig() interfaces.IConfig {
	config := viperConfig{viper: viper.New()}
	config.setDefaults()
	config.initialize()
	return &config
}

func (c *viperConfig) setDefaults() {
	c.viper.SetDefault("database.path", path.Join("data", "db"))
}

func (c *viperConfig) initialize() {
	c.viper.SetConfigName("config")
	c.viper.SetConfigType("yaml")
	c.viper.AddConfigPath(".")
	c.viper.AutomaticEnv()
}

// GetDatabasePath returns the database path
func (c *viperConfig) GetDatabasePath() string {
	return c.viper.GetString("database.path")
}
