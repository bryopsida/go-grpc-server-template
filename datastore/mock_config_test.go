package datastore

import (
	"github.com/stretchr/testify/mock"
)

// MockConfig is a mock implementation of the interfaces.IConfig interface.
type MockConfig struct {
	mock.Mock
}

func (m *MockConfig) GetDatabasePath() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfig) GetServerPort() uint16 {
	args := m.Called()
	return uint16(args.Int(0))
}

func (m *MockConfig) GetServerAddress() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfig) GetServerCert() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfig) GetServerKey() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfig) GetServerCA() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfig) IsTLSEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}
