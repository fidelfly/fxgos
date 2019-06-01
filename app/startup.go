package app

import (
	"github.com/fidelfly/fxgo"
	"github.com/sirupsen/logrus"

	_ "github.com/fidelfly/fxgos/caches"
	"github.com/fidelfly/fxgos/system"
)

func StartService() (err error) {
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	// Parse Config File
	err = fxgo.InitTomlConfig("config.toml", &system.Configuration)
	if err != nil {
		return
	}

	// Setup logs
	fxgo.SetupLogs(&system.Runtime.LogConfig)

	// Init Database
	err = initDatabase(*system.Database)
	if err != nil {
		return
	}

	// Setup Router
	myRouter, err := setupRouter()
	if err != nil {
		return
	}

	//start Server
	if system.SupportTLS() {
		fxgo.ListenAndServeTLS(system.TLS.CertFile, system.TLS.KeyFile, myRouter, system.Runtime.Port)
	} else {
		fxgo.ListenAndServe(myRouter, system.Runtime.Port)
	}

	return
}
