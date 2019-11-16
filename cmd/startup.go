package main

import (
	"context"
	"flag"
	"time"

	"github.com/fidelfly/gox/gosrvx"

	"github.com/fidelfly/gox/confx"
	"github.com/fidelfly/gox/logx"
	"github.com/fidelfly/gox/pkg/filex"
	"github.com/go-xorm/xorm"

	_ "github.com/fidelfly/fxgos/cmd/obsolete/caches"
	"github.com/fidelfly/fxgos/cmd/router"
	"github.com/fidelfly/fxgos/cmd/service"
	"github.com/fidelfly/fxgos/cmd/service/api/user"
	"github.com/fidelfly/fxgos/cmd/service/res"
	auditSrv "github.com/fidelfly/fxgos/cmd/service/server/audit"
	cronSrv "github.com/fidelfly/fxgos/cmd/service/server/cron"
	daSrv "github.com/fidelfly/fxgos/cmd/service/server/da"
	fileSrv "github.com/fidelfly/fxgos/cmd/service/server/filedb"
	iamSrv "github.com/fidelfly/fxgos/cmd/service/server/iam"
	otkSrv "github.com/fidelfly/fxgos/cmd/service/server/otk"
	roleSrv "github.com/fidelfly/fxgos/cmd/service/server/role"
	userSrv "github.com/fidelfly/fxgos/cmd/service/server/user"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
	"github.com/fidelfly/fxgos/cmd/utilities/system"
	"github.com/fidelfly/gostool/db"
	"github.com/fidelfly/gostool/mail"
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

	//start cron jobs
	cronSrv.Start()

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
	err = cronSrv.Initialize()
	if err != nil {
		return
	}
	//init mail alert
	mail.InitMailDelegator(*system.Mail)
	logx.Info("Mail function is initialized.")

	//init file db
	err = fileSrv.Initialize()
	if err != nil {
		return
	}
	logx.Info("File db is initialized.")
	//init iam
	err = iamSrv.Initialize()
	if err != nil {
		return
	}
	logx.Info("Iam function is initialized.")
	//init da
	err = daSrv.Initialize()
	if err != nil {
		return
	}
	logx.Info("Data access function is initialized.")
	//init otk
	err = otkSrv.Initialize()
	if err != nil {
		return
	}
	logx.Info("One-Time-Token function is initialized.")

	//init audit
	err = auditSrv.Initialize()
	if err != nil {
		return
	}

	//init user
	err = userSrv.Initialize()
	if err != nil {
		return
	}
	logx.Info("User function is initialized.")
	//init role
	err = roleSrv.Initialize()
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

	return service.Start()
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
