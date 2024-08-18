package interfaces

// IConfig is an interface for configuration
type IConfig interface {
	// GetDatabasePath returns the database path
	GetDatabasePath() string
	GetServerAddress() string
	GetServerPort() uint16
	GetServerCert() string
	GetServerKey() string
	GetServerCA() string
	IsTLSEnabled() bool
}
