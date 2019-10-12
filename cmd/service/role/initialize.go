package role

import (
	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/service/role/res"
	"github.com/fidelfly/fxgos/cmd/utilities/pub"
)

func Initialize() error {
	err := db.Synchronize(
		new(res.Role),
	)

	pub.Subscribe(pub.TopicResource, cacheSubscriber)
	return err
}
