package otk

import (
	"github.com/fidelfly/fxgos/cmd/service/api/otk"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/gostool/db"
)

func Initialize() error {
	if err := db.Synchronize(
		new(res.OneTimeKey),
	); err != nil {
		return err
	}
	otk.RegisterServer(&server{})
	return nil
}
