package databroker

import (
	"crypto/tls"
	"encoding/base64"
	"time"

	"github.com/pomerium/pomerium/internal/log"
	"github.com/pomerium/pomerium/pkg/cryptutil"
)

var (
	// DefaultDeletePermanentlyAfter is the default amount of time to wait before deleting
	// a record permanently.
	DefaultDeletePermanentlyAfter = time.Hour
	// DefaultStorageType is the default storage type that Server use
	DefaultStorageType = "memory"
	// DefaultGetAllPageSize is the default page size for GetAll calls.
	DefaultGetAllPageSize = 50
)

type serverConfig struct {
	installationID          string
	deletePermanentlyAfter  time.Duration
	secret                  []byte
	storageType             string
	storageConnectionString string
	storageCAFile           string
	storageCertSkipVerify   bool
	storageCertificate      *tls.Certificate
	getAllPageSize          int
}

func newServerConfig(options ...ServerOption) *serverConfig {
	cfg := new(serverConfig)
	WithDeletePermanentlyAfter(DefaultDeletePermanentlyAfter)(cfg)
	WithStorageType(DefaultStorageType)(cfg)
	WithGetAllPageSize(DefaultGetAllPageSize)(cfg)
	for _, option := range options {
		option(cfg)
	}
	return cfg
}

// A ServerOption customizes the server.
type ServerOption func(*serverConfig)

// WithDeletePermanentlyAfter sets the deletePermanentlyAfter duration.
// If a record is deleted via Delete, it will be permanently deleted after
// the given duration.
func WithDeletePermanentlyAfter(dur time.Duration) ServerOption {
	return func(cfg *serverConfig) {
		cfg.deletePermanentlyAfter = dur
	}
}

// WithGetAllPageSize sets the page size for GetAll calls.
func WithGetAllPageSize(pageSize int) ServerOption {
	return func(cfg *serverConfig) {
		cfg.getAllPageSize = pageSize
	}
}

// WithInstallationID sets the installation id in the config.
func WithInstallationID(installationID string) ServerOption {
	return func(cfg *serverConfig) {
		cfg.installationID = installationID
	}
}

// WithSharedKey sets the secret in the config.
func WithSharedKey(sharedKey string) ServerOption {
	return func(cfg *serverConfig) {
		key, err := base64.StdEncoding.DecodeString(sharedKey)
		if err != nil || len(key) != cryptutil.DefaultKeySize {
			log.Error().Err(err).Msgf("shared key is required and must be %d bytes long", cryptutil.DefaultKeySize)
			return
		}
		cfg.secret = key
	}
}

// WithStorageType sets the storage type.
func WithStorageType(typ string) ServerOption {
	return func(cfg *serverConfig) {
		cfg.storageType = typ
	}
}

// WithStorageConnectionString sets the DSN for storage.
func WithStorageConnectionString(connStr string) ServerOption {
	return func(cfg *serverConfig) {
		cfg.storageConnectionString = connStr
	}
}

// WithStorageCAFile sets the CA file in the config.
func WithStorageCAFile(filePath string) ServerOption {
	return func(cfg *serverConfig) {
		cfg.storageCAFile = filePath
	}
}

// WithStorageCertSkipVerify sets the storageCertSkipVerify in the config.
func WithStorageCertSkipVerify(storageCertSkipVerify bool) ServerOption {
	return func(cfg *serverConfig) {
		cfg.storageCertSkipVerify = storageCertSkipVerify
	}
}

// WithStorageCertificate sets the storageCertificate in the config.
func WithStorageCertificate(certificate *tls.Certificate) ServerOption {
	return func(cfg *serverConfig) {
		cfg.storageCertificate = certificate
	}
}
