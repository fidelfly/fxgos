package da

import (
	"github.com/fidelfly/fxgos/cmd/service/api/da"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/gostool/db"
)

func (s server) Start() error {
	initSgCache()
	initDaModel()
	return nil
}

func Initialize() error {
	err := db.Synchronize(
		new(res.SecurityGroup),
		new(res.UserSg),
		new(res.ResourceSg),
	)
	da.RegisterServer(&server{})

	return err
}
