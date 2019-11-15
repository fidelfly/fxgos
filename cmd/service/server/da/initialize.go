package da

import (
	"github.com/fidelfly/fxgos/cmd/service/api/da"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/gostool/db"
)

func Initialize() error {
	err := db.Synchronize(
		new(res.SecurityGroup),
		new(res.UserSg),
		new(res.ResourceSg),
	)

	initSgCache()
	initDaModel()

	da.RegisterServer(&server{})

	return err
}
