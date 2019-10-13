package otk

import (
	"github.com/fidelfly/fxgos/cmd/service/otk/res"
	"github.com/fidelfly/gostool/db"
)

func Initialize() error {
	return db.Synchronize(
		new(res.OneTimeKey),
	)
}
