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
	DATABASE_PATH_KEY        = "database.path"
	SERVER_PORT_KEY          = "server.port"
	SERVER_ADDRESS_KEY       = "server.address"
	SERVER_TLS_ENABLED_KEY   = "server.tls.enabled"
	SERVER_TLS_CERT_KEY      = "server.tls.cert"
	SERVER_TLS_CERT_PATH_KEY = "server.tls.cert_path"
	SERVER_TLS_KEY_KEY       = "server.tls.key"
	SERVER_TLS_KEY_PATH_KEY  = "server.tls.key_path"
	SERVER_TLS_CA_KEY        = "server.tls.ca"
	SERVER_TLS_CA_PATH_KEY   = "server.tls.ca_path"
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
	c.viper.SetDefault(DATABASE_PATH_KEY, path.Join("data", "db"))
	c.viper.SetDefault(SERVER_PORT_KEY, 50051)
	c.viper.SetDefault(SERVER_ADDRESS_KEY, "localhost")
	c.viper.SetDefault(SERVER_TLS_ENABLED_KEY, false)
	c.viper.SetDefault(SERVER_TLS_CERT_KEY, "")
	c.viper.SetDefault(SERVER_TLS_CERT_PATH_KEY, "")
	c.viper.SetDefault(SERVER_TLS_KEY_KEY, "")
	c.viper.SetDefault(SERVER_TLS_KEY_PATH_KEY, "")
	c.viper.SetDefault(SERVER_TLS_CA_KEY, "")
	c.viper.SetDefault(SERVER_TLS_CA_PATH_KEY, "")
}

func (c *viperConfig) initialize() {
	c.viper.SetConfigName("config")
	c.viper.SetConfigType("yaml")
	c.viper.AddConfigPath(".")
	c.viper.AutomaticEnv()
}

// GetDatabasePath returns the database path
func (c *viperConfig) GetDatabasePath() string {
	return c.viper.GetString(DATABASE_PATH_KEY)
}

func (c *viperConfig) GetServerPort() uint16 {
	return uint16(c.viper.GetInt(SERVER_PORT_KEY))
}

func (c *viperConfig) GetServerAddress() string {
	return c.viper.GetString(SERVER_ADDRESS_KEY)
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
	return c.ifNilTryPath(SERVER_TLS_CERT_KEY, SERVER_TLS_CERT_PATH_KEY)
}

func (c *viperConfig) GetServerKey() string {
	return c.ifNilTryPath(SERVER_TLS_KEY_KEY, SERVER_TLS_KEY_PATH_KEY)
}

func (c *viperConfig) GetServerCA() string {
	return c.ifNilTryPath(SERVER_TLS_CA_KEY, SERVER_TLS_CA_PATH_KEY)
}

func (c *viperConfig) IsTlsEnabled() bool {
	return c.viper.GetBool(SERVER_TLS_ENABLED_KEY)
}
