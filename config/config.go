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
	databasePathkey      = "database.path"
	serverPortKey        = "server.port"
	serverAddressKey     = "server.address"
	serverTLSEnabledKey  = "server.tls.enabled"
	serverTLSCertKey     = "server.tls.cert"
	serverTLSCertPathKey = "server.tls.cert_path"
	serverTLSKeyKey      = "server.tls.key"
	serverTLSKeyPathKey  = "server.tls.key_path"
	serverTLSCaKey       = "server.tls.ca"
	serverTLSCaPathKey   = "server.tls.ca_path"
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
	c.viper.SetDefault(databasePathkey, path.Join("data", "db"))
	c.viper.SetDefault(serverPortKey, 50051)
	c.viper.SetDefault(serverAddressKey, "localhost")
	c.viper.SetDefault(serverTLSEnabledKey, false)
	c.viper.SetDefault(serverTLSCertKey, "")
	c.viper.SetDefault(serverTLSCertPathKey, "")
	c.viper.SetDefault(serverTLSKeyKey, "")
	c.viper.SetDefault(serverTLSKeyPathKey, "")
	c.viper.SetDefault(serverTLSCaKey, "")
	c.viper.SetDefault(serverTLSCaPathKey, "")
}

func (c *viperConfig) initialize() {
	c.viper.SetConfigName("config")
	c.viper.SetConfigType("yaml")
	c.viper.AddConfigPath(".")
	c.viper.AutomaticEnv()
}

// GetDatabasePath returns the database path
func (c *viperConfig) GetDatabasePath() string {
	return c.viper.GetString(databasePathkey)
}

func (c *viperConfig) GetServerPort() uint16 {
	return uint16(c.viper.GetInt(serverPortKey))
}

func (c *viperConfig) GetServerAddress() string {
	return c.viper.GetString(serverAddressKey)
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
	}
	return c.viper.GetString(primaryKey)
}

func (c *viperConfig) GetServerCert() string {
	return c.ifNilTryPath(serverTLSCertKey, serverTLSCertPathKey)
}

func (c *viperConfig) GetServerKey() string {
	return c.ifNilTryPath(serverTLSKeyKey, serverTLSKeyPathKey)
}

func (c *viperConfig) GetServerCA() string {
	return c.ifNilTryPath(serverTLSCaKey, serverTLSCaPathKey)
}

func (c *viperConfig) IsTLSEnabled() bool {
	return c.viper.GetBool(serverTLSEnabledKey)
}
