package filedb

import (
	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/service/filedb/res"
)

func Initialize() error {
	err := db.Synchronize(
		new(res.File),
	)
	return err
}
