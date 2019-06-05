package system

import (
	"os"

	"github.com/fidelfly/fxgo"
)

var Configuration = Properties{}
var Runtime = &Configuration.Runtime
var Database = &Configuration.Database
var OAuth2 = &Configuration.Oauth2
var TLS = &Configuration.TLS

func SupportTLS() bool {
	if len(TLS.CertFile) == 0 || len(TLS.KeyFile) == 0 {
		return false
	}

	if _, err := os.Stat(TLS.CertFile); err != nil {
		return false
	}

	if _, err := os.Stat(TLS.KeyFile); err != nil {
		return false
	}
	return true
}

type Properties struct {
	Version  string
	Runtime  RuntimeProperties
	Database DatabaseProperties
	Oauth2   OAuth2Properties
	TLS      TLSConfig
}

type RuntimeProperties struct {
	fxgo.LogConfig
	WebPath string
	Debug   bool
	Port    int64
}

type DatabaseProperties struct {
	Host     string
	Port     string
	Schema   string
	User     string
	Password string
}

type OAuth2Properties struct {
	Client []AuthClient
}

type AuthClient struct {
	ID     string
	Secret string
	Domain string
}

type TLSConfig struct {
	CertFile string
	KeyFile  string
}
