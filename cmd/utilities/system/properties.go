package system

import (
	"os"

	"github.com/fidelfly/fxgo"
	"github.com/fidelfly/fxgo/authx"

	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/pkg/mail"
)

var Configuration = Properties{}
var Runtime = &Configuration.Runtime
var Database = &Configuration.Database
var OAuth2 = &Configuration.Oauth2
var TLS = &Configuration.TLS
var Mail = &Configuration.Mail

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
	Mail     mail.Config
}

type RuntimeProperties struct {
	fxgo.LogConfig
	AssetPath     string
	WebPath       string
	Debug         bool
	Port          int64
	TemporaryPath string
}

type DatabaseProperties struct {
	db.Server
	fxgo.LogConfig
}

type OAuth2Properties struct {
	Client []authx.AuthClient
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
