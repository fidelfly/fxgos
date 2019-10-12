package user

import (
	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/service/user/res"
	"github.com/fidelfly/fxgos/cmd/utilities/pub"
)

func Initialize() error {
	err := db.Synchronize(
		new(res.User),
	)

	pub.Subscribe(pub.TopicResource, cacheSubscriber)

	return err
}

const (
	StatusDeleted = iota - 2
	StatusInvalid
	StatusDeactivated
	StatusValid
)
