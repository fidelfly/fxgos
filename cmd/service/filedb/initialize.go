package filedb

import (
	"github.com/fidelfly/fxgos/cmd/service/filedb/res"
	"github.com/fidelfly/gostool/db"
)

func Initialize() error {
	err := db.Synchronize(
		new(res.File),
	)
	return err
}
