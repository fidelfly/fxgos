package otk

import (
	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/service/otk/res"
)

func Initialize() error {
	return db.Synchronize(
		new(res.OneTimeKey),
	)
}
