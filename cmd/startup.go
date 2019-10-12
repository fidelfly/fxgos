package main

import (
	"context"
	"flag"
	"time"

	"github.com/fidelfly/fxgo/gosrvx"

	// "github.com/fidelfly/fxgo/gosrvx"
	"github.com/fidelfly/fxgo/confx"
	"github.com/fidelfly/fxgo/logx"
	"github.com/fidelfly/fxgo/pkg/filex"
	"github.com/go-xorm/xorm"

	_ "github.com/fidelfly/fxgos/cmd/obsolete/caches"
	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/pkg/mail"
	"github.com/fidelfly/fxgos/cmd/router"
	"github.com/fidelfly/fxgos/cmd/service/audit"
	"github.com/fidelfly/fxgos/cmd/service/filedb"
	"github.com/fidelfly/fxgos/cmd/service/iam"
	"github.com/fidelfly/fxgos/cmd/service/otk"
	"github.com/fidelfly/fxgos/cmd/service/role"
	"github.com/fidelfly/fxgos/cmd/service/user"
	"github.com/fidelfly/fxgos/cmd/service/user/res"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
	"github.com/fidelfly/fxgos/cmd/utilities/system"
)

const defConfigFile = "config.toml"

func StartService() (err error) {
	defer func() {
		if err != nil {
			logx.Error(err)
		}
	}()

	defer func() {
		if err != nil {
			logx.Error(err)
		}
		if err := recover(); err != nil {
			logx.Error("Panic Error: ", err)
		}
	}()

	// Parse Config File
	err = parseConfig()
	if err != nil {
		return
	}
	logx.Info("Configuration is loaded.")

	// Setup logs
	err = setupLog()
	if err != nil {
		return
	}
	logx.Info("Log is setup.")
	// init runtime
	err = initRuntime()
	if err != nil {
		return
	}
	logx.Info("Runtime is initialized")
	// Init Database
	err = initDatabase()
	if err != nil {
		return
	}
	logx.Info("Database is connected.")

	// init function
	err = initFunction()
	if err != nil {
		return
	}
	// init data
	err = initData()
	if err != nil {
		return
	}

	// start Server
	if system.SupportTLS() {
		gosrvx.ListenAndServeTLS(system.TLS.CertFile, system.TLS.KeyFile, router.GetRootRouter(), system.Runtime.Port)
	} else {
		gosrvx.ListenAndServe(router.GetRootRouter(), system.Runtime.Port)
	}

	return nil
}

func parseConfig() error {
	//return  gosrvx.InitTomlConfig("config.toml", &system.Configuration)
	var configFile = ""
	flag.StringVar(&configFile, "config", defConfigFile, "Set Config File")
	flag.Parse()

	return confx.ParseToml(configFile, &system.Configuration)
}

func setupLog() error {
	gosrvx.SetupLogs(&system.Runtime.LogConfig)
	dbLoger = gosrvx.NewLog(&system.Database.LogConfig)
	return nil
}

var dbLoger *logx.Logger

func initDatabase() error {
	err := db.InitEngine(&system.Database.Server, func(engine *xorm.Engine) {
		engine.ShowSQL(true)
		engine.SetConnMaxLifetime(3595 * time.Second)
		if dbLoger != nil {
			engine.SetLogger(&db.Logger{Logger: dbLoger})
		}
		engine.SetLogLevel(db.GetLogLevel(system.Database.LogLevel))
	})
	return err
}

func initRuntime() error {
	return filex.MustFoler(system.Runtime.TemporaryPath)
}

func initFunction() (err error) {
	//init mail alert
	mail.InitMailDelegator(*system.Mail)
	logx.Info("Mail function is initialized.")

	//init file db
	err = filedb.Initialize()
	if err != nil {
		return
	}
	logx.Info("File db is initialized.")
	//init iam
	err = iam.Initialize()
	if err != nil {
		return
	}
	logx.Info("Iam function is initialized.")
	//init otk
	err = otk.Initialize()
	if err != nil {
		return
	}
	logx.Info("One-Time-Token function is initialized.")

	//init audit
	err = audit.Initialize()
	if err != nil {
		return
	}

	//init user
	err = user.Initialize()
	if err != nil {
		return
	}
	logx.Info("User function is initialized.")
	//init role
	err = role.Initialize()
	if err != nil {
		return
	}
	logx.Info("Role function is initialized.")

	//init router
	err = router.Initialize()
	if err != nil {
		return
	}
	logx.Info("Router function is initialized.")
	return err
}

func initData() (err error) {
	err = initSa()
	return
}

func initSa() error {
	if _, err := user.ReadByCode(context.Background(), "sa"); err != nil {
		if err == syserr.ErrNotFound {
			_, err = user.Create(context.Background(), &res.User{
				Code:       "sa",
				Name:       "Super Administrator",
				Email:      "fidel.xu@gmail.com",
				Password:   "valerie",
				Region:     "Shen Zhen",
				Status:     user.StatusValid,
				SuperAdmin: true,
			})
			return err
		}
		return err
	}

	return nil
}
