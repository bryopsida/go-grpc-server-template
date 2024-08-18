// main_test.go
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/bryopsida/go-grpc-server-template/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockIConfig is a mock of IConfig interface using testify/mock
type MockIConfig struct {
	mock.Mock
}

// GetDatabasePath implements interfaces.IConfig.
func (m *MockIConfig) GetDatabasePath() string {
	args := m.Called()
	return args.String(0)
}

// GetServerAddress implements interfaces.IConfig.
func (m *MockIConfig) GetServerAddress() string {
	args := m.Called()
	return args.String(0)
}

// GetServerPort implements interfaces.IConfig.
func (m *MockIConfig) GetServerPort() uint16 {
	args := m.Called()
	return uint16(args.Int(0))
}

// IsTLSEnabled implements interfaces.IConfig.
func (m *MockIConfig) IsTLSEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockIConfig) GetServerCert() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockIConfig) GetServerKey() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockIConfig) GetServerCA() string {
	args := m.Called()
	return args.String(0)
}

// MockListener is a mock of net.Listener using testify/mock
type MockListener struct {
	mock.Mock
}

func (m *MockListener) Accept() (net.Conn, error) {
	args := m.Called()
	return args.Get(0).(net.Conn), args.Error(1)
}

func (m *MockListener) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockListener) Addr() net.Addr {
	args := m.Called()
	return args.Get(0).(net.Addr)
}

// MockServer is a mock of grpc.Server using testify/mock
type MockServer struct {
	mock.Mock
}

func (m *MockServer) Serve(lis net.Listener) error {
	args := m.Called(lis)
	return args.Error(0)
}

func (m *MockServer) GracefulStop() {
	m.Called()
}

func createTestCert() (string, string, error) {
	// Generate a new RSA private key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	// Create a certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Org"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Self-sign the certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return "", "", err
	}

	// Encode the private key to PEM format
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})

	// Encode the certificate to PEM format
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	return string(certPEM), string(privPEM), nil
}

func TestBuildTLSCert(t *testing.T) {
	cert, key, err := createTestCert()
	if err != nil {
		t.Fatalf("failed to create test cert: %v", err)
	}
	tests := []struct {
		name    string
		cert    string
		key     string
		wantErr bool
	}{
		{
			name:    "ValidCertAndKey",
			cert:    cert,
			key:     key,
			wantErr: false,
		},
		{
			name:    "InvalidCert",
			cert:    "invalid-cert",
			key:     "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQ==\n-----END PRIVATE KEY-----",
			wantErr: true,
		},
		{
			name:    "InvalidKey",
			cert:    "-----BEGIN CERTIFICATE-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA7Q3Q6Z5j5Q==\n-----END CERTIFICATE-----",
			key:     "invalid-key",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := buildTLSCert(tt.cert, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildTLSCert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildServerCredentials(t *testing.T) {
	cert, key, err := createTestCert()
	if err != nil {
		t.Fatalf("Failed to create test certificate: %v", err)
	}

	tests := []struct {
		name    string
		config  func() interfaces.IConfig
		wantErr bool
	}{
		{
			name: "ValidConfig",
			config: func() interfaces.IConfig {
				mockConfig := new(MockIConfig)
				mockConfig.On("GetServerCert").Return(cert)
				mockConfig.On("GetServerKey").Return(key)
				mockConfig.On("GetServerCA").Return(cert)
				return mockConfig
			},
			wantErr: false,
		},
		{
			name: "MissingCert",
			config: func() interfaces.IConfig {
				mockConfig := new(MockIConfig)
				mockConfig.On("GetServerCert").Return("")
				mockConfig.On("GetServerKey").Return(key)
				mockConfig.On("GetServerCA").Return(cert)
				return mockConfig
			},
			wantErr: true,
		},
		{
			name: "MissingKey",
			config: func() interfaces.IConfig {
				mockConfig := new(MockIConfig)
				mockConfig.On("GetServerCert").Return(cert)
				mockConfig.On("GetServerKey").Return("")
				mockConfig.On("GetServerCA").Return(cert)
				return mockConfig
			},
			wantErr: true,
		},
		{
			name: "MissingCA",
			config: func() interfaces.IConfig {
				mockConfig := new(MockIConfig)
				mockConfig.On("GetServerCert").Return(cert)
				mockConfig.On("GetServerKey").Return(key)
				mockConfig.On("GetServerCA").Return("")
				return mockConfig
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creds, err := buildServerCredentials(tt.config())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, creds)
			}
		})
	}
}

func TestBuildGrpcOptions(t *testing.T) {
	cert, key, err := createTestCert()
	if err != nil {
		t.Fatalf("Failed to create test certificate: %v", err)
	}

	tests := []struct {
		name    string
		config  func() interfaces.IConfig
		wantErr bool
	}{
		{
			name: "TLSEnabled",
			config: func() interfaces.IConfig {
				mockConfig := new(MockIConfig)
				mockConfig.On("IsTLSEnabled").Return(true)
				mockConfig.On("GetServerCert").Return(cert)
				mockConfig.On("GetServerKey").Return(key)
				mockConfig.On("GetServerCA").Return(cert)
				return mockConfig
			},
			wantErr: false,
		},
		{
			name: "TLSDisabled",
			config: func() interfaces.IConfig {
				mockConfig := new(MockIConfig)
				mockConfig.On("IsTLSEnabled").Return(false)
				return mockConfig
			},
			wantErr: false,
		},
		{
			name: "TLSEnabledWithError",
			config: func() interfaces.IConfig {
				mockConfig := new(MockIConfig)
				mockConfig.On("IsTLSEnabled").Return(true)
				mockConfig.On("GetServerCert").Return("")
				mockConfig.On("GetServerKey").Return(key)
				mockConfig.On("GetServerCA").Return(cert)
				return mockConfig
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				assert.Panics(t, func() { buildGrpcOptions(tt.config()) })
			} else {
				options := buildGrpcOptions(tt.config())
				if tt.config().IsTLSEnabled() {
					assert.NotEmpty(t, options)
					assert.IsType(t, grpc.Creds(nil), options[0])
				} else {
					assert.Empty(t, options)
				}
			}
		})
	}
}

func TestBuildGrpcServer(t *testing.T) {
	tests := []struct {
		name    string
		options []grpc.ServerOption
	}{
		{
			name:    "NoOptions",
			options: nil,
		},
		{
			name: "WithOptions",
			options: []grpc.ServerOption{
				grpc.MaxRecvMsgSize(1024 * 1024),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := buildGrpcServer(tt.options)
			assert.NotNil(t, server)
		})
	}
}

func TestServeGrpc(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*MockListener, *MockServer)
		wantErr bool
	}{
		{
			name: "SuccessfulServe",
			setup: func(lis *MockListener, server *MockServer) {
				lis.On("Addr").Return(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080})
				server.On("Serve", lis).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ServeError",
			setup: func(lis *MockListener, server *MockServer) {
				lis.On("Addr").Return(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080})
				server.On("Serve", lis).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLis := new(MockListener)
			mockServer := new(MockServer)
			tt.setup(mockLis, mockServer)
			grpcServer := interfaces.GrpcServer(mockServer)
			if tt.wantErr {
				assert.Panics(t, func() { serveGrpc(grpcServer, mockLis) })
			} else {
				assert.NotPanics(t, func() { serveGrpc(grpcServer, mockLis) })
			}

			mockLis.AssertExpectations(t)
			mockServer.AssertExpectations(t)
		})
	}
}

func TestRunGrpc(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*MockListener, *MockServer)
		wantErr bool
	}{
		{
			name: "SuccessfulRun",
			setup: func(lis *MockListener, server *MockServer) {
				lis.On("Addr").Return(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080})
				server.On("Serve", lis).Return(nil)
				server.On("GracefulStop").Return()
			},
			wantErr: false,
		},
		{
			name: "ContextCancellation",
			setup: func(lis *MockListener, server *MockServer) {
				lis.On("Addr").Return(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080})
				server.On("Serve", lis).Return(nil)
				server.On("GracefulStop").Return()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLis := new(MockListener)
			mockServer := new(MockServer)
			tt.setup(mockLis, mockServer)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go runGrpc(ctx, mockServer, mockLis)

			time.Sleep(1 * time.Second) // Simulate some running time
			cancel()                    // Cancel the context to stop the server

			time.Sleep(2 * time.Second) // Give some time for the server to stop

			mockLis.AssertExpectations(t)
			mockServer.AssertExpectations(t)
		})
	}
}
