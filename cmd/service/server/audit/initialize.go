package audit

import (
	"github.com/fidelfly/fxgos/cmd/service/api/audit"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/gostool/db"
)

func Initialize() error {
	err := db.Synchronize(
		new(res.Systrail),
	)
	audit.RegisterServer(&server{})
	return err
}
