package audit

import (
	"github.com/fidelfly/fxgos/cmd/service/audit/res"
	"github.com/fidelfly/gostool/db"
)

func Initialize() error {
	err := db.Synchronize(
		new(res.Systrail),
	)

	return err
}
