package user

import (
	"github.com/fidelfly/fxgos/cmd/service/api/user"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/gostool/db"
)

func Initialize() error {
	err := db.Synchronize(
		new(res.User),
	)

	user.RegisterServer(&server{})
	return err
}
