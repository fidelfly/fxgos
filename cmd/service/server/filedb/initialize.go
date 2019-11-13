package filedb

import (
	"github.com/fidelfly/fxgos/cmd/service/api/filedb"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/gostool/db"
)

func Initialize() error {
	err := db.Synchronize(
		new(res.File),
	)

	filedb.RegisterServer(&server{})
	return err
}
