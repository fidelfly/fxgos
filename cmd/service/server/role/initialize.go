package role

import (
	"github.com/fidelfly/fxgos/cmd/service/api/role"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/gostool/db"
)

func Initialize() error {
	err := db.Synchronize(
		new(res.Role),
	)
	role.RegisterServer(&server{})
	return err
}
