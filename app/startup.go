package app

import (
	"github.com/lyismydg/fxgos/system"
	"github.com/lyismydg/fxgos/router"
	"net/http"
	"github.com/sirupsen/logrus"
	_ "github.com/lyismydg/fxgos/caches"
)

func StartService() (err error) {
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	// Parse Config File
	err = system.InitConfig()
	if err != nil {
		return
	}

	// Setup logs
	err = system.SetupLog()
	if err != nil {
		return
	}

	// Init Database
	err = system.InitDatabase(*system.Database)
	if err != nil {
		return
	}

	// Setup Router
	myRouter, err := router.SetupRouter()
	if err != nil {
		return
	}

	//start Server
	server := &http.Server{
		Handler: myRouter,
		Addr: ":" + system.Runtime.Port,
	}

	if system.SupportTLS() {
		logrus.Fatal(server.ListenAndServeTLS(system.TLS.CertFile, system.TLS.KeyFile))
	} else {
		logrus.Fatal(server.ListenAndServe())
	}

	return
}
