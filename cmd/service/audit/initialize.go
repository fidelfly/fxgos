package audit

import (
	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/service/audit/res"
)

func Initialize() error {
	err := db.Synchronize(
		new(res.Systrail),
	)

	return err
}
