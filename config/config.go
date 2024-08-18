package config

import (
	"io"
	"log/slog"
	"os"
	"path"

	"github.com/bryopsida/go-grpc-server-template/interfaces"
	"github.com/spf13/viper"
)

const (
	database_path_key        = "database.path"
	server_port_key          = "server.port"
	server_address_key       = "server.address"
	server_tls_enabled_key   = "server.tls.enabled"
	server_tls_cert_key      = "server.tls.cert"
	server_tls_cert_path_key = "server.tls.cert_path"
	server_tls_key_key       = "server.tls.key"
	server_tls_key_path_key  = "server.tls.key_path"
	server_tls_ca_key        = "server.tls.ca"
	server_tls_ca_path_key   = "server.tls.ca_path"
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
	c.viper.SetDefault(database_path_key, path.Join("data", "db"))
	c.viper.SetDefault(server_port_key, 50051)
	c.viper.SetDefault(server_address_key, "localhost")
	c.viper.SetDefault(server_tls_enabled_key, false)
	c.viper.SetDefault(server_tls_cert_key, "")
	c.viper.SetDefault(server_tls_cert_path_key, "")
	c.viper.SetDefault(server_tls_key_key, "")
	c.viper.SetDefault(server_tls_key_path_key, "")
	c.viper.SetDefault(server_tls_ca_key, "")
	c.viper.SetDefault(server_tls_ca_path_key, "")
}

func (c *viperConfig) initialize() {
	c.viper.SetConfigName("config")
	c.viper.SetConfigType("yaml")
	c.viper.AddConfigPath(".")
	c.viper.AutomaticEnv()
}

// GetDatabasePath returns the database path
func (c *viperConfig) GetDatabasePath() string {
	return c.viper.GetString(database_path_key)
}

func (c *viperConfig) GetServerPort() uint16 {
	return uint16(c.viper.GetInt(server_port_key))
}

func (c *viperConfig) GetServerAddress() string {
	return c.viper.GetString(server_address_key)
}

func (c *viperConfig) ifNilTryPath(primaryKey string, pathKey string) string {
	if c.viper.GetString(primaryKey) == "" {
		path := c.viper.GetString(pathKey)
		if path != "" {
			// Open the file
			file, err := os.Open(path)
			if err != nil {
				slog.Warn("Failed to open file from path", slog.String("path", path), slog.Any("error", err))
				return ""
			}
			defer file.Close()

			// Read the file contents
			content, err := io.ReadAll(file)
			if err != nil {
				slog.Warn("Failed to read file from path", slog.String("path", path), slog.Any("error", err))
				return ""
			}
			return string(content)
		}
		return ""
	} else {
		return c.viper.GetString(primaryKey)
	}
}

func (c *viperConfig) GetServerCert() string {
	return c.ifNilTryPath(server_tls_cert_key, server_tls_cert_path_key)
}

func (c *viperConfig) GetServerKey() string {
	return c.ifNilTryPath(server_tls_key_key, server_tls_key_path_key)
}

func (c *viperConfig) GetServerCA() string {
	return c.ifNilTryPath(server_tls_ca_key, server_tls_ca_path_key)
}

func (c *viperConfig) IsTlsEnabled() bool {
	return c.viper.GetBool(server_tls_enabled_key)
}
